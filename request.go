package kokatto

import (
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

// Request helper for kokatto
type Request struct {
	Signature string `url:"signature,omitempty"`
}

// QueryString return object as query string
func (r *Request) QueryString(obj interface{}) (string, error) {
	v, err := query.Values(obj)
	if err != nil {
		return "", errors.Wrap(err, "QueryString err")
	}
	return v.Encode(), nil
}

// SetSignature will set request signature
func (r *Request) SetSignature(signature string) {
	r.Signature = signature
}

// Validate will return true if signature matched
func (r *Request) Validate(signature string) bool {
	return r.Signature == signature
}

// OTPRequest attribute
type OTPRequest struct {
	Request
	ClientID    string `url:"clientId"`
	AppType     string `url:"appType"`
	MediaType   string `url:"mediaType"`
	PhoneNumber string `url:"phoneNumber"`
	OTPMCode    string `url:"otpmCode,omitempty"`
	Timestamp   string `url:"timestamp"`
}

// QueryString return OTPRequest object as query string
func (o *OTPRequest) QueryString() (string, error) {
	return o.Request.QueryString(o)
}

// DeliveryStatusRequest attribute
type DeliveryStatusRequest struct {
	Request
	ClientID  string `url:"clientId"`
	AppType   string `url:"appType"`
	RequestID string `url:"requestId"`
	Timestamp string `url:"timestamp"`
}

// QueryString return DeliveryStatusRequest object as query string
func (d *DeliveryStatusRequest) QueryString() (string, error) {
	return d.Request.QueryString(d)
}
