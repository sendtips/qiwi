/*
Package qiwi implements QIWI API
as a client library.

Behind this library there are two main structures:
Payment which carries all our requests and RSP responses
and Notify which holds inbound requests from RSP.

You should process incoming requests in your own function
and unmarshalling to Notify, the rest will be done by this lib.

Example to process ApplePay payment

    pay := qiwi.New("OrderID", "SiteID", "TOKEN", "http://example.com/qiwi-api")

    err := pay.ApplePay(context.TODO(), 300, "ApplePayTokenString") // Request a session for 3.00RUB
    if err != nil {
        fmt.Printf("Error occurred: %v", err)
    }

If an error is a nil, you need to wait for a webhook.
You should pass hook payload to NewNotify function,
the Notify object will be returned with payment statusЮ

QIWI uses ISO8601 time format, unusual in Go
for that, we implement custom time parsers for JSON
example: 2021-07-29T16:30:00+03:00 */
package qiwi
