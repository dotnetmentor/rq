# rq

Resource Query - for dynamic querying of resources

## General

- A resource can be anything you want (tenants, packages, markets etc.)
- Properties are metdata about a resource (that can be selected but not queried)
- Conditions can be declared to allow filtering of results (by providing parameters as flags)

## Example file

See `rq.yaml` for a basic example.

## Help

Run `rq --help` to get details on how to use the cli.

## Install

### Homebrew tap

```console
brew install dotnetmentor/tap/rq
```

### go install

```console
go install github.com/dotnetmentor/rq@latest
```

### Manual

Download binaries from [release page](https://github.com/dotnetmentor/rq/releases)
