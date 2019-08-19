package kokatto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// TODO(sakti): change scheme to https if supported
	createURL         = "http://otpm.kokatto.com/otpm/create"
	deliveryStatusURL = "http://otpm.kokatto.com/otpm/status"
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
func NewClient(p Params) (*Kokatto, error) {
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

	return &Kokatto{client: c, params: p, secretKey: []byte(p.PrivateKey), tz: tz}, nil
}

// Kokatto client
type Kokatto struct {
	client    *http.Client
	params    Params
	secretKey []byte
	tz        *time.Location
}

func (k *Kokatto) requestOTP() *OTPRequest {
	return &OTPRequest{
		ClientID:  k.params.ClientID,
		AppType:   k.params.AppType,
		MediaType: k.params.MediaType,
		Timestamp: time.Now().In(k.tz).Format(timeFormat),
	}
}

func (k *Kokatto) requestDeliveryStatus() *DeliveryStatusRequest {
	return &DeliveryStatusRequest{
		ClientID:  k.params.ClientID,
		AppType:   k.params.AppType,
		Timestamp: time.Now().In(k.tz).Format(timeFormat),
	}
}

// SendOTP will call kokatto request OTP
func (k *Kokatto) SendOTP(phoneNumber string, opts ...Option) (rsp OTPResponse, err error) {
	o := evaluateOptions(opts)
	req := k.requestOTP()

	req.PhoneNumber = phoneNumber
	if o.otpCode != "" {
		req.OTPMCode = o.otpCode
	}

	err = Sign(req, k.secretKey)
	if err != nil {
		return rsp, err
	}
	query, err := req.QueryString()
	if err != nil {
		return rsp, err
	}

	result, err := k.client.Post(createURL+"?"+query, "application/json", nil)
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
func (k *Kokatto) DeliveryStatus(reqID string) (rsp DeliveryStatusResponse, err error) {
	req := k.requestDeliveryStatus()
	req.RequestID = reqID

	err = Sign(req, k.secretKey)
	if err != nil {
		return rsp, err
	}
	query, err := req.QueryString()
	if err != nil {
		return rsp, err
	}
	fmt.Println("[*] ", deliveryStatusURL+"?"+query)
	result, err := k.client.Get(deliveryStatusURL + "?" + query)
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
			fmt.Println("[*] error disini")
			fmt.Println(string(rspBody))
			return rsp, err
		}
		return rsp, &errRsp
	}
	err = json.Unmarshal(rspBody, &rsp)
	return rsp, err
}
