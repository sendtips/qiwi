/*
Package qiwi implements QIWI API
as a client library.

Behind this library there are two main structures:
Payment which carries all our requests and RSP responses
and Notify which holds inbound requests from RSP.

Example to process ApplePay payment

    pay := qiwi.New("OrderID", "SiteID", "TOKEN", "http://example.com/qiwi-api")

    err := pay.ApplePay(context.TODO(), 300, "ApplePayTokenString") // Request a session for 3.00RUB
    if err != nil {
        fmt.Printf("Error occurred: %v", err)
    }

You may pass hook payload to NewNotify function,
or use NotifyParseHTTPRequest which works directly for http.Request
the Notify object will be returned with the payment status.

QIWI uses ISO8601 time format, 
so we implement a custom time parser
example: 2021-07-29T16:30:00+03:00 */
package qiwi
