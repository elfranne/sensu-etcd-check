[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/elfranne/sensu-etcd-check)
![Go Test](https://github.com/elfranne/sensu-etcd-check/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/elfranne/sensu-etcd-check/workflows/goreleaser/badge.svg)

# sensu-etcd-check

## Table of Contents

- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The sensu-etcd-check is a [Sensu Check][6] that connects to one or more [etcd][11] cluster
endpoints and reports the database size reported by `etcdctl endpoint status`. The check goes
critical when the database size exceeds a configurable threshold, which is useful for catching
etcd clusters approaching their configured storage quota before they hit it and go read-only.

## Usage examples

```text
Usage:
  check-etcd [flags]

Flags:
  -c, --cert-file string       Path to the cert
  -h, --help                   help for check-etcd
  -k, --key-file string        Path to the key
  -s, --size int               Maximum database Size (default 1500000000)
  -t, --timeout int            Request timeout (default 5)
      --trusted-ca-file string Path to the CA file
  -u, --url strings            Url of etcd instance(s) (default [http://127.0.0.1:2379])
```

- `--url` accepts a comma-separated list of endpoints; the check queries the status of the
  first URL in the list.
- `--size` is the maximum acceptable `DbSize` in bytes before the check returns critical.
  Etcd's default storage quota is 2GB, so the default threshold here is set to 1.5GB to leave
  room to react.
- `--cert-file`, `--key-file`, and `--trusted-ca-file` must all be provided together to enable
  mutual TLS when connecting to etcd.

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add elfranne/sensu-etcd-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][8].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-etcd-check
  namespace: default
spec:
  command: sensu-etcd-check --url http://127.0.0.1:2379 --size 1500000000
  subscriptions:
  - etcd
  runtime_assets:
  - elfranne/sensu-etcd-check
```

A TLS-enabled example, for clusters using client certificate authentication:

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-etcd-check
  namespace: default
spec:
  command: >-
    sensu-etcd-check
    --url https://127.0.0.1:2379
    --cert-file /etc/etcd/ssl/etcd-client.crt
    --key-file /etc/etcd/ssl/etcd-client.key
    --trusted-ca-file /etc/etcd/ssl/ca.crt
  subscriptions:
  - etcd
  runtime_assets:
  - elfranne/sensu-etcd-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-etcd-check repository:

```
go build
```

## Additional notes

The check only queries the status of the first endpoint in `--url`; additional endpoints are
passed through to the etcd client for connection purposes but are not each individually checked.

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[8]: https://bonsai.sensu.io/assets/elfranne/sensu-etcd-check
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[11]: https://etcd.io/
