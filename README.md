# Obslytics (Observability Analytics)

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/gmfrasca/obslytics)
[![Latest Release](https://img.shields.io/github/release/thanos-community/obslytics.svg?style=flat-square)](https://github.com/gmfrasca/obslytics/releases/latest)
[![CI](https://github.com/gmfrasca/obslytics/workflows/go/badge.svg)](https://github.com/gmfrasca/obslytics/actions?query=workflow%3Ago)
[![Go Report Card](https://goreportcard.com/badge/github.com/gmfrasca/obslytics)](https://goreportcard.com/report/github.com/gmfrasca/obslytics)
[![Slack](https://img.shields.io/badge/join%20slack-%23analytics-brightgreen.svg)](https://slack.cncf.io/)

## Usage

[embedmd]:# (obslytics-help.txt $)
```$
usage: obslytics [<flags>] <command> [<args> ...]

Integrate Observability data into your Analytics pipelines

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and
                           --help-man).
      --version            Show application version.
      --log.level=info     Log filtering level.
      --log.format=logfmt  Log format to use.

Commands:
  help [<command>...]
    Show help.

  export --match=MATCH --min-time=MIN-TIME --max-time=MAX-TIME --resolution=RESOLUTION [<flags>]
    Export observability series data into popular analytics formats.


```

## Contributing

Any contributions are welcome! Just use GitHub Issues and Pull Requests as usual.
We follow [Thanos Go coding style](https://thanos.io/contributing/coding-style-guide.md/) guide.

## Initial Author

[@bwplotka](https://bwplotka.dev)
