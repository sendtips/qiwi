/*
Package qiwi implements QIWI API
As a client library.

Example to process ApplePay payment

    pay := qiwi.New("OrderID", "SiteID", "TOKEN", "http://example.com/qiwi-api")

    err := pay.ApplePay(context.TODO(), 300, "ApplePayTokenString") // Request a session for 3.00RUB
    if err != nil {
        fmt.Printf("Error occurred: %v", err)
    }

If error is nil, you need to wait for a web hook.
You should pass hook payload to NewNotify function,
the Notify object will be returned with payment status

QIWI uses ISO8601 time format, unusual in Go
for that we implement custom time parsers for JSON
example: 2021-07-29T16:30:00+03:00 */
package qiwi
