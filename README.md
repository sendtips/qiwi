# QIWI API Go client library

[![GitHub Actions](https://github.com/sendtips/qiwi/workflows/Go/badge.svg)](https://github.com/sendtips/qiwi/actions?workflow=Go)
[![GoDoc](https://godoc.org/github.com/sendtips/qiwi?status.svg)](https://godoc.org/github.com/sendtips/qiwi)
[![codecov](https://codecov.io/gh/sendtips/qiwi/branch/master/graph/badge.svg)](https://codecov.io/gh/sendtips/qiwi)
[![Go Report Card](https://goreportcard.com/badge/github.com/sendtips/qiwi)](https://goreportcard.com/report/github.com/sendtips/qiwi)
[![Sourcegraph](https://sourcegraph.com/github.com/sendtips/qiwi/-/badge.svg)](https://sourcegraph.com/github.com/sendtips/qiwi?badge)
[![sendtips](https://img.shields.io/badge/🍩_Sendtips-@awsom82-black?labelColor=3298dc)](https://sendtips.ru/pay/E2ZfzjVE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsendtips%2Fqiwi.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsendtips%2Fqiwi?ref=badge_shield)


A Go client library to work with [QIWI API](https://developer.qiwi.com/en/).

## Install
Install by import `github.com/sendtips/qiwi` or via `go get github.com/sendtips/qiwi`

The library support `go1.14` and newer.

## Library status
Library in early development, we *not recommend use it on production* till it reach `v1`.

## Tests
Run tests using `go test`

## Example
To obtain a payment session on QIWI website you need to create a new qiwi object via `qiwi.New()` and call its `CardRequest()` method.

```go
package main

import (
    "fmt"
    "context"
    
    "github.com/sendtips/qiwi"
)

func main() {
    pay := qiwi.New("OrderID", "SiteID", "TOKEN", "http://example.com/qiwi-api")

    err := pay.CardRequest(context.TODO(), "PublicKey", 300) // Request a session for 3.00RUB
    if err != nil {
        fmt.Printf("Error occurred: %v", err)
    }

    fmt.Printf("%s", pay.PayURL) // Payment session url on QIWI website
}
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsendtips%2Fqiwi.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsendtips%2Fqiwi?ref=badge_large)