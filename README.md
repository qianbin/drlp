# drlp

[![Build Status](https://github.com/qianbin/drlp/workflows/Go/badge.svg)](https://github.com/qianbin/drlp/actions)
[![GoDoc](https://godoc.org/github.com/qianbin/drlp?status.svg)](http://godoc.org/github.com/qianbin/drlp)
[![Go Report](https://goreportcard.com/badge/github.com/qianbin/drlp)](https://goreportcard.com/report/github.com/qianbin/drlp)

Short for Direct-RLP: A fast in-place RLP encoder


### Installation

It requires Go 1.19 or newer.

```bash
go get github.com/qianbin/drlp
```

### Usage

Number and string 
```go
var buf drlp.Buffer

buf.PutUint(10)
buf.PutString([]byte("hello drlp"))

fmt.Printf("%x\n", buf)
// 0a8a68656c6c6f2064726c70
```

List
```go
var buf drlp.Buffer

buf.PutString([]byte("this is a list"))

li := buf.List()
buf.PutString([]byte("lis content"))
encodedList := li.End()

fmt.Printf("%x\n", encodedList)
// cc8b6c697320636f6e74656e74
fmt.Printf("%x\n", buf)
// 8e746869732069732061206c697374cc8b6c697320636f6e74656e74
```

## License

The MIT License

