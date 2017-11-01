# uid

A package for generating unique identifiers, loosely based on MongoDB's [ObjectId](https://docs.mongodb.com/manual/reference/method/ObjectId/).

Identifiers are 12-bytes and roughly time sortable.

> * a 4-byte value representing the seconds since the Unix epoch,
> * a 3-byte machine identifier,
> * a 2-byte process id, and
> * a 3-byte counter, starting with a random value.

## Documentation

Package documentation available on GoDoc: uid: [github.com/billglover/uid](https://godoc.org/github.com/billglover/uid)

## Benchmarks

```plain
goos: darwin
goarch: amd64
pkg: github.com/billglover/uid
BenchmarkUID-4            500000              2313 ns/op
BenchmarkUIDString-4      500000              2307 ns/op
```
