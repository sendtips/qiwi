/*
Package qiwi implements QIWI API
as a client library.

Behind this library there are two main structures:
Payment which carries all our requests and RSP responses
and Notify which holds inbound requests from RSP.

For payment sessions are CardRequest, ApplePay and GooglePay methods of Payments are available.

Example to process ApplePay payment

    pay := qiwi.New("OrderID", "SiteID", "TOKEN", "http://example.com/qiwi-api")

    err := pay.ApplePay(context.TODO(), 300, "ApplePayTokenString") // Request a session for 3.00RUB
    if err != nil {
        fmt.Printf("Error occurred: %s", err)
    }

You may pass hook payload to NewNotify function,
or use NotifyParseHTTPRequest which works directly for http.Request
the Notify object will be returned with the payment status.

Example of receiving Notify from incoming http.Request

    //..
    notify, err := qiwi.NotifyParseHTTPRequest("YOUTOKEN", w.ResponseWriter, r.Request)
    if err != nil {
        fmt.Println("[QIWI] Error while loading incoming notification:", err)
    }

Or you may process received data by yourself and pass the payload to NewNotify

    //..
    notify, err := qiwi.NewNotify("YOUTOKEN", "SIGNATURE", payload)
    if err != nil {
        fmt.Println("[QIWI] Error while parsing notification:", err)
    }


QIWI uses ISO8601 (2021-07-29T16:30:00+03:00) time format,
so use a build-in QIWITime custom time type */
package qiwi
