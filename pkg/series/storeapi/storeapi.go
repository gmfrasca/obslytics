package storeapi

import (
	"context"
	"io"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/storage"
	"github.com/gmfrasca/obslytics/pkg/series"
	"github.com/thanos-io/thanos/pkg/store/labelpb"
	"github.com/thanos-io/thanos/pkg/store/storepb"
	tracing "github.com/thanos-io/thanos/pkg/tracing/client"
	"google.golang.org/grpc"

	"math"
	"github.com/go-kit/kit/log/level"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/credentials"
	thanostracing "github.com/thanos-io/thanos/pkg/tracing"

)

func newCustomClientConfig(logger log.Logger, cert, key, caCert, serverName string, insecureSkipVerify bool) (*tls.Config, error) {
	var certPool *x509.CertPool
	if caCert != "" {
		caPEM, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, errors.Wrap(err, "reading client CA")
		}

		certPool = x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caPEM) {
			return nil, errors.Wrap(err, "building client CA")
		}
		level.Info(logger).Log("msg", "TLS client using provided certificate pool")
	} else {
		var err error
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, errors.Wrap(err, "reading system certificate pool")
		}
		level.Info(logger).Log("msg", "TLS client using system certificate pool")
	}

	tlsCfg := &tls.Config{
		RootCAs: certPool,
		InsecureSkipVerify: insecureSkipVerify,
	}

	if (key != "") != (cert != "") {
		return nil, errors.New("both client key and certificate must be provided")
	}

	if cert != "" {
		cert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, errors.Wrap(err, "client credentials")
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
		level.Info(logger).Log("msg", "TLS client authentication enabled")
	}
	return tlsCfg, nil
}


// StoreClientGRPCOpts creates gRPC dial options for connecting to a store client.
func InsecureClient(logger log.Logger, reg *prometheus.Registry, tracer opentracing.Tracer, secure bool, cert, key, caCert, serverName string) ([]grpc.DialOption, error) {
	grpcMets := grpc_prometheus.NewClientMetrics()
	grpcMets.EnableClientHandlingTimeHistogram(
		grpc_prometheus.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	)
	dialOpts := []grpc.DialOption{
		// We want to make sure that we can receive huge gRPC messages from storeAPI.
		// On TCP level we can be fine, but the gRPC overhead for huge messages could be significant.
		// Current limit is ~2GB.
		// TODO(bplotka): Split sent chunks on store node per max 4MB chunks if needed.
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpcMets.UnaryClientInterceptor(),
				thanostracing.UnaryClientInterceptor(tracer),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpcMets.StreamClientInterceptor(),
				thanostracing.StreamClientInterceptor(tracer),
			),
		),
	}
	if reg != nil {
		reg.MustRegister(grpcMets)
	}



	level.Info(logger).Log("msg", "enabling client to server TLS")

	tlsCfg, err := newCustomClientConfig(logger, cert, key, caCert, serverName, !secure)
	if err != nil {
		return nil, err
	}

	return append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))), nil
}


// Series implements input.Reader.
type Series struct {
	logger log.Logger
	conf   series.Config
}

func NewSeries(logger log.Logger, conf series.Config) (Series, error) {
	return Series{logger: logger, conf: conf}, nil
}

func (i Series) Read(ctx context.Context, params series.Params) (series.Set, error) {
	dialOpts, err := InsecureClient(i.logger, nil, tracing.NoopTracer(),
		!i.conf.TLSConfig.InsecureSkipVerify,
		i.conf.TLSConfig.CertFile,
		i.conf.TLSConfig.KeyFile,
		i.conf.TLSConfig.CAFile,
		i.conf.TLSConfig.ServerName)

	if err != nil {
		return nil, errors.Wrap(err, "error initializing GRPC options")
	}

	conn, err := grpc.DialContext(ctx, i.conf.Endpoint, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing GRPC dial context")
	}


	matchers, err := storepb.TranslatePromMatchers(params.Matchers...)
	if err != nil {
		return nil, err
	}

	client := storepb.NewStoreClient(conn)
	seriesClient, err := client.Series(ctx, &storepb.SeriesRequest{
		MinTime:                 timestamp.FromTime(params.MinTime),
		MaxTime:                 timestamp.FromTime(params.MaxTime),
		Matchers:                matchers,
		PartialResponseStrategy: storepb.PartialResponseStrategy_ABORT,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "storepb.Series against %v", i.conf.Endpoint)
	}

	return &iterator{
		ctx:    ctx,
		conn:   conn,
		client: seriesClient,
		mint:   timestamp.FromTime(params.MinTime),
		maxt:   timestamp.FromTime(params.MaxTime),
	}, nil
}

// iterator implements input.Set.
type iterator struct {
	ctx           context.Context
	conn          *grpc.ClientConn
	client        storepb.Store_SeriesClient
	currentSeries *storepb.Series

	mint, maxt int64

	err error
}

func (i *iterator) Next() bool {
	seriesResp, err := i.client.Recv()
	if err == io.EOF {
		return false
	}
	if err != nil {
		i.err = err
		return false
	}

	i.currentSeries = seriesResp.GetSeries()
	return true
}

func (i *iterator) At() storage.Series {
	// We support only raw data for now.
	return newChunkSeries(
		labelpb.LabelsToPromLabels(i.currentSeries.Labels),
		i.currentSeries.Chunks,
		i.mint, i.maxt,
		[]storepb.Aggr{storepb.Aggr_COUNT, storepb.Aggr_SUM},
	)
}

func (i *iterator) Warnings() storage.Warnings { return nil }

func (i *iterator) Err() error {
	return i.err
}

func (i *iterator) Close() error {
	if err := i.client.CloseSend(); err != nil {
		return err
	}

	return i.conn.Close()
}
