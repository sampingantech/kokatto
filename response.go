package kokatto

// OTPResponse attribute
type OTPResponse struct {
	Status       string `json:"status"`
	StatusCode   string `json:"statusCode"`
	RequestID    string `json:"requestId"`
	OTPM         string `json:"otpm"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage"`
}

// DeliveryData attribute
type DeliveryData struct {
	RequestID            string `json:"requestId"`
	DestinationAddress   string `json:"destinationAddress"`
	OTPCode              string `json:"otpCode"`
	OTPStatus            string `json:"otpStatus"`
	OTPStatusDescription string `json:"otpStatusDescription"`
}

// DeliveryStatusResponse attribute
type DeliveryStatusResponse struct {
	Status       string       `json:"status"`
	StatusCode   string       `json:"statusCode"`
	Message      string       `json:"message"`
	ErrorMessage string       `json:"errorMessage"`
	Data         DeliveryData `json:"data"`
}
