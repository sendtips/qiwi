package qiwi

// {
//        "payment":{
//           "paymentid":"4504751",
//           "tokendata":{
//              "paymenttoken":"4cc975be-483f-8d29-2b7de3e60c2f",
//              "expireddate":"2021-12-31t00:00:00+03:00"
//           },
//           "type":"payment",
//           "createddatetime":"2019-10-08t11:31:37+03:00",
//           "status":{
//              "value":"success",
//              "changeddatetime":"2019-10-08t11:31:37+03:00"
//           },
//           "amount":{
//              "value":2211.24,
//              "currency":"RUB"
//           },
//           "paymentMethod":{
//              "type":"CARD",
//              "maskedPan":"220024/*/*/*/*/*/*5036",
//              "rrn":null,
//              "authCode":null,
//              "type":"CARD"
//           },
//           "paymentCardInfo": {
//              "issuingCountry": "810",
//              "issuingBank": "QiwiBank",
//              "paymentSystem": "VISA",
//              "fundingSource": "CREDIT",
//              "paymentSystemProduct": "P|Visa Gold"
//           },
//           "customer":{
//              "ip":"79.142.20.248",
//              "account":"token32",
//              "phone":"0"
//           },
//           "billId":"testing122",
//           "customFields":{},
//           "flags":[
//              "SALE"
//           ]
//        },
//        "type":"PAYMENT",
//        "version":"1"
//     }

// payment	Payment information	Object	Always
// type	Operation type	String(200)	Always
// paymentId	Payment operation unique identifier in RSP's system	String(200)	Always
// createdDateTime	System date of the operation creation	URL-encoded string
// YYYY-MM-DDThh:mm:ss	Always
// amount	Object	Operation amount data	Always
// value	Operation amount rounded down to two decimals	Number(6.2)	Always
// currency	Operation currency (Code: Alpha-3 ISO 4217)	String(3)	Always
// billId	Corresponding invoice ID for the operation	String(200)	Always
// status	Operation status data	Object	Always
// value	Operation status value	String	Always
// changedDatetime	Date of operation status update	URL-encoded string
// YYYY-MM-DDThh:mm:ssZ	Always
// reasonCode	Rejection reason code	String(200)	Always
// reasonMessage	Rejection reason description	String(200)	Always
// errorCode	Error code	Number	Always
// paymentMethod	Payment method data	Object	Always
// type	Payment method type	String
// paymentToken	Card payment token	String	When payment token is used for the operation
// maskedPan	Masked card PAN	String	Always
// rrn	Payment RRN (ISO 8583)	Number	Always
// authCode	Payment Auth code	Number	Always
// paymentCardInfo	Card information	Object	Always
// issuingCountry	Issuer country code	String(3)	Always
// issuingBank	Issuer name	String	Always
// paymentSystem	Card's payment system	String	Always
// fundingSource	Card's type (debit/credit/..)	String	Always
// paymentSystemProduct	Card's category	String	Always
// customer	Customer data	Object	Always
// phone	Customer phone number	String	Always
// email	Customer e-mail	String	Always
// account	Customer ID in RSP system	String	Always
// ip	Customer IP address	String	Always
// country	Customer country from address string	String	Always
// customFields	Fields with additional information for the operation	Object	Always
// cf1	Extra field with some information to operation data	String(256)	Always
// cf2	Extra field with some information to operation data	String(256)	Always
// cf3	Extra field with some information to operation data	String(256)	Always
// cf4	Extra field with some information to operation data	String(256)	Always
// cf5	Extra field with some information to operation data	String(256)	Always
// FROM_MERCHANT_CONTRACT_ID	Contract ID between the merchant and the customer	String(256)	Always
// FROM_MERCHANT_BOOKING_NUMBER	Booking number for the customer	String(256)	Always
// FROM_MERCHANT_PHONE	Phone number of the customer	String(256)	Always
// FROM_MERCHANT_FULL_NAME	Full name of the customer	String(256)	Always
// flags	Additional API commands	Array of Strings. Possible elements — SALE / REVERSAL	Always
// tokenData	Payment token data	Object	When payment token issue was requested
// paymentToken	Card payment token	String	When payment token issue was requested
// expiredDate	Payment token expiration date. ISO-8601 Date format:
// YYYY-MM-DDThh:mm:ss±hh:mm	String	When payment token issue was requested
// paymentSplits	Split payments description	Array(Objects)	For split payments
// type	Data type. Always MERCHANT_DETAILS	String	For split payments
// siteUid	Merchant ID	String	For split payments
// splitAmount	Merchant reimbursement	Object	For split payments
// value	Amount of reimbursement	Number	For split payments
// currency	Text code of reimbursement currency, by ISO	String(3)	For split payments
// splitCommissions	Commission data	Object	For split payments
// merchantCms	Commission from merchant	Object	For split payments
// value	Amount of commission	Number	For split payments
// currency	Text code of commission currency, by ISO	String(3)	For split payments
// orderId	Order number	String	For split payments
// comment	Comment for the order	String	For split payments
// type	Notification type (PAYMENT)	String(200)	Always
// version	Notification version	String	Always
