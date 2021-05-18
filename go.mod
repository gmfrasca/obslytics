module github.com/thanos-community/obslytics

go 1.14

require (
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/OneOfOne/xxhash v1.2.6 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/cortexproject/cortex v1.8.1-0.20210422151339-cf1c444e0905
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20191106031601-ce3c9ade29de // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-ieproxy v0.0.0-20191113090002-7c0f6868bffe // indirect
	github.com/oklog/run v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.21.0
	github.com/prometheus/prometheus v1.8.2-0.20210421143221-52df5ef7a3be
	github.com/sercand/kuberesolver v2.4.0+incompatible // indirect
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/thanos-io/thanos v0.20.1
	github.com/xitongsys/parquet-go v1.5.2
	github.com/xitongsys/parquet-go-source v0.0.0-20200817004010-026bad9b25d0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.uber.org/automaxprocs v1.3.0
	google.golang.org/grpc v1.36.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/utils v0.0.0-20200731180307-f00132d28269 // indirect
)

// Compatibility constraints
replace (
	// Using a 3rd-party branch for custom dialer - see https://github.com/bradfitz/gomemcache/pull/86
	github.com/bradfitz/gomemcache => github.com/themihai/gomemcache v0.0.0-20180902122335-24332e2d58ab
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	github.com/prometheus/common => github.com/prometheus/common v0.13.0
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20210421143221-52df5ef7a3be
	github.com/sercand/kuberesolver => github.com/sercand/kuberesolver v2.4.0+incompatible
	github.com/thanos-io/thanos => github.com/thanos-io/thanos v0.20.1
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	k8s.io/client-go => k8s.io/client-go v0.18.3
	k8s.io/klog => k8s.io/klog v0.3.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/api => k8s.io/api v0.18.19
)
