package kokatto

import (
	"testing"
)

func TestOTPRequest_QueryString(t *testing.T) {
	type fields struct {
		ClientID    string
		AppType     string
		MediaType   string
		PhoneNumber string
		OTPMCode    string
		Signature   string
		Timestamp   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields{
				ClientID:    "8080",
				AppType:     "OPM",
				MediaType:   "VOI",
				PhoneNumber: "08000000001",
				OTPMCode:    "1234",
				Timestamp:   "2019-08-19T16:59:01+0700",
			},
			want: "appType=OPM&clientId=8080&mediaType=VOI&otpmCode=1234&phoneNumber=08000000001&timestamp=2019-08-19T16%3A59%3A01%2B0700",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OTPRequest{
				Request: Request{
					Signature: tt.fields.Signature,
				},
				ClientID:    tt.fields.ClientID,
				AppType:     tt.fields.AppType,
				MediaType:   tt.fields.MediaType,
				PhoneNumber: tt.fields.PhoneNumber,
				OTPMCode:    tt.fields.OTPMCode,
				Timestamp:   tt.fields.Timestamp,
			}
			if got, err := o.QueryString(); got != tt.want {
				if err != nil && !tt.wantErr {
					t.Error(err)
				}
				t.Errorf("OTPRequest.QueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}
