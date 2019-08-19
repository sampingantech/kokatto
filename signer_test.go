package kokatto

import "testing"

func TestSign(t *testing.T) {
	type args struct {
		req    Signer
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal case",
			args: args{
				req: &OTPRequest{
					ClientID:    "8080",
					AppType:     "OPM",
					MediaType:   "VOI",
					PhoneNumber: "08000000001",
					OTPMCode:    "1234",
					Timestamp:   "2019-08-19T16:59:01+0700",
				},
				secret: []byte("25e5d4b904d8aa3144df3c5e15d6a60fc997f3101888eafe0fa9615726641123"),
			},
			want: "9B2D51B6F64DC2C314D98E762461740567265DEF33493AC087B8EF020AD6316F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Sign(tt.args.req, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.args.req.Validate(tt.want) {
				t.Errorf("failed to validate, signature = %v", tt.want)
			}
		})
	}
}
