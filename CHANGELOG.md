# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.8] - 2026-09-09

### Changed
- Database size is now checked on every endpoint given with `--url`, not just
  the first one. A multi-endpoint check that previously only inspected the
  first endpoint may now report CRITICAL for the others.
- Argument validation returns errors instead of printing them, so failures
  surface as UNKNOWN rather than CRITICAL. Missing certificate, key or CA
  files now exit 3 instead of 2.
- `--cert-file`, `--key-file` and `--trusted-ca-file` must now be supplied
  together; setting only some of them is an error instead of being partially
  applied.
- TLS configuration errors are reported instead of ignored, and no empty TLS
  config is passed when the TLS flags are unset.
- Output now names the endpoint each result belongs to.

### Added
- Timeout applied to the etcd status request, using the existing `--timeout`
  value.

### Fixed
- Corrected typos in the `--size` and `--url` help text and in the
  "exceeding set limit" message.

## [0.0.1] - 2000-01-01

### Added
- Initial release
