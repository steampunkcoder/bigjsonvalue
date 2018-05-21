# bigjsonvalue

[![Build status][travis-image]][travis-url] [![License][license-image]][license-url] [![GoDoc][godoc-img]][godoc-url]

> `bigjsonvalue` replaces `interface{}` for decoding unknown JSON values

`BigJSONValue` and `NatJSONValue` are wrappers around the `interface{}` type, to force
[`json.Unmarshal()`](https://golang.org/pkg/encoding/json/#Unmarshal)
to decode integer values as integers instead of `float64`.
Main reason for this is `float64` doesn't have enough precision
to store exact values of large `int64` or `uint64` values.
Instead of trying to unmarshal unknown JSON values into an `interface{}`,
unmarshal into a `BigJSONValue` or `NatJSONValue` instead.

### Example Usage

The impetus for `bigjsonvalue` is decoding the JSON encoded output from
[`wal2json`](https://github.com/eulerto/wal2json), which is an output plugin for
[PostgreSQL Logical-Decoding](https://www.postgresql.org/docs/current/static/logicaldecoding.html).
This example uses the [`pgx`](https://github.com/jackc/pgx) library to
connect to the [`wal2json`](https://github.com/eulerto/wal2json) output of
[PostgreSQL Logical-Decoding](https://www.postgresql.org/docs/current/static/logicaldecoding.html).

```golang
import (
        pgx "github.com/jackc/pgx"
        bjv "github.com/steampunkcoder/bigjsonvalue"
)

// WalOldKeys models the "oldkeys" map in a WAL change JSON
type WalOldKeys struct {
        KeyNames  []string           `json:"keynames"`
        KeyTypes  []string           `json:"keytypes"`
        KeyValues []bjv.BigJSONValue `json:"keyvalues"`
}

// WalChangeRec models a single change record in a WAL change JSON
type WalChangeRec struct {
        Kind         string             `json:"kind"`
        Schema       string             `json:"schema"`
        Table        string             `json:"table"`
        ColumnNames  []string           `json:"columnnames"`
        ColumnTypes  []string           `json:"columntypes"`
        ColumnValues []bjv.BigJSONValue `json:"columnvalues"`
        OldKeys      WalOldKeys         `json:"oldkeys"`
}

// WalChangeTx models an entire change list transaction in a WAL change JSON
// obtained from a pgx.WalMessage notification
type WalChangeTx struct {
        Changes []WalChangeRec `json:"change"`
}

func WalListenLoop() {
        rConn, _ := pgx.ReplicationConnect(...)
        rConn.CreateReplicationSlot("myslot", "wal2json")
        rConn.StartReplication("myslot", 0, -1, ...)
        for {
                rMsg, _ := rConn.WaitForReplicationMessage(...)
                if rMsg.WalMessage != nil {
                        var chgTx WalChangeTx
                        json.Unmarshal(rMsg.WalMessage.WalData, &chgTx)
                } else if rMsg.ServerHeartbeat != nil {
                        sMsg, _ := pgx.NewStandbyStatus(...)
                        rConn.SendStandbyStatus(sMsg)
                }
        }
}

```

### API Docs
[![GoDoc][godoc-img]][godoc-url]

### License
Released under MIT License [![License][license-image]][license-url]

[godoc-url]: https://godoc.org/github.com/steampunkcoder/bigjsonvalue
[godoc-img]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square
[license-url]: LICENSE
[travis-image]: https://img.shields.io/travis/steampunkcoder/bigjsonvalue.svg?style=flat-square
[travis-url]: https://travis-ci.org/steampunkcoder/bigjsonvalue

