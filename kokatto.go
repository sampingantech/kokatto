package kokatto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	createURL         = "https://otpm.kokatto.com/otpm/create"
	deliveryStatusURL = "https://otpm-report.kokatto.com/otpm/status"
	defAppType        = "OPM"
	defMediaType      = "VOI"
	tzLocation        = "Asia/Jakarta"
	timeFormat        = "2006-01-02T15:04:05Z0700"
)

// Params for kokatto client
type Params struct {
	ClientID   string
	PrivateKey string
	AppType    string
	MediaType  string
}

// NewClient will return a new kokatto client instance.
func NewClient(p Params) (*Client, error) {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	if p.ClientID == "" {
		return nil, errors.Wrap(ErrMissingParam, "ClientID")
	}
	if p.PrivateKey == "" {
		return nil, errors.Wrap(ErrMissingParam, "PrivateKey")
	}
	if p.AppType == "" {
		p.AppType = defAppType
	}
	if p.MediaType == "" {
		p.MediaType = defMediaType
	}
	tz, err := time.LoadLocation(tzLocation)
	if err != nil {
		return nil, err
	}

	return &Client{client: c, params: p, secretKey: []byte(p.PrivateKey), tz: tz}, nil
}

// Client for kokatto
type Client struct {
	client    *http.Client
	params    Params
	secretKey []byte
	tz        *time.Location
}

func (c *Client) requestOTP() *OTPRequest {
	return &OTPRequest{
		ClientID:  c.params.ClientID,
		AppType:   c.params.AppType,
		MediaType: c.params.MediaType,
		Timestamp: time.Now().In(c.tz).Format(timeFormat),
	}
}

func (c *Client) requestDeliveryStatus() *DeliveryStatusRequest {
	return &DeliveryStatusRequest{
		ClientID:  c.params.ClientID,
		AppType:   c.params.AppType,
		Timestamp: time.Now().In(c.tz).Format(timeFormat),
	}
}

// SendOTP will call kokatto request OTP
func (c *Client) SendOTP(phoneNumber string, opts ...Option) (rsp OTPResponse, err error) {
	o := evaluateOptions(opts)
	req := c.requestOTP()

	req.PhoneNumber = phoneNumber
	if o.otpCode != "" {
		req.OTPMCode = o.otpCode
	}

	err = Sign(req, c.secretKey)
	if err != nil {
		return rsp, err
	}
	query, err := req.QueryString()
	if err != nil {
		return rsp, err
	}

	result, err := c.client.Post(createURL+"?"+query, "application/json", nil)
	if err != nil {
		return rsp, errors.Wrap(err, "SendOTP failed")
	}
	defer result.Body.Close()
	rspBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return rsp, err
	}
	if result.StatusCode != 200 && result.StatusCode != 201 {
		var errRsp Error
		err = json.Unmarshal(rspBody, &errRsp)
		if err != nil {
			return rsp, err
		}
		return rsp, &errRsp
	}
	err = json.Unmarshal(rspBody, &rsp)
	return rsp, err
}

// DeliveryStatus will get delivery status of OTP based on reqID.
func (c *Client) DeliveryStatus(reqID string) (rsp DeliveryStatusResponse, err error) {
	req := c.requestDeliveryStatus()
	req.RequestID = reqID

	err = Sign(req, c.secretKey)
	if err != nil {
		return rsp, err
	}
	query, err := req.QueryString()
	if err != nil {
		return rsp, err
	}
	result, err := c.client.Get(deliveryStatusURL + "?" + query)
	if err != nil {
		return rsp, errors.Wrap(err, "DeliveryStatus failed")
	}
	defer result.Body.Close()
	rspBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return rsp, err
	}
	if result.StatusCode != 200 {
		var errRsp Error
		err = json.Unmarshal(rspBody, &errRsp)
		if err != nil {
			return rsp, err
		}
		return rsp, &errRsp
	}
	err = json.Unmarshal(rspBody, &rsp)
	return rsp, err
}
