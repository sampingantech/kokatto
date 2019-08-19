# kokatto

Unofficial [Kokatto](https://www.kokatto.com) Go Client

[![go report card](https://goreportcard.com/badge/github.com/sampingantech/kokatto "go report card")](https://goreportcard.com/report/github.com/sampingantech/kokatto)
[![GoDoc](https://godoc.org/github.com/sampingantech/kokatto?status.svg)](https://godoc.org/github.com/sampingantech/kokatto)

## Example

Setup client

```
client, err  := kokatto.NewClient(kokatto.Params{
    ClientID:   "kokatto-client-id",
    PrivateKey: "kokatto-private-key",
})
```

Sending OTP, OTP generated by Kokatto
```
rsp, err := client.SendOTP("080000xxxx")
```

Sending OTP, OTP specified by client
```
rsp, err := client.SendOTP("080000xxxx", kokatto.WithOTP("1234"))
```

Get delivery status
```
status, err := client.DeliveryStatus(rsp.RequestID)
```
