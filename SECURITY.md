# Security Policy

## Supported versions

Pre-1.0, only the latest tagged release receives fixes.

## Reporting a vulnerability

Please report suspected vulnerabilities privately via
[GitHub Security Advisories](https://github.com/dimchansky/h3-go/security/advisories/new)
rather than opening a public issue. Reports are handled on a best-effort
basis; this is a volunteer-maintained project with no response-time SLA.

## Scope notes

The production library is pure, safe Go: no cgo, no `unsafe`, no
dependencies (enforced by CI gates). The cgo parity harness and the
`interop/uberdiff` module compile third-party C code, but only behind
opt-in build tags in test infrastructure — they are never part of a
consumer's build.
