package kokatto

var (
	defaultOptions = &options{}
)

type options struct {
	otpCode string
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option for sauth middleware
type Option func(*options)

// WithOTP will set OTP code to send.
func WithOTP(otpCode string) Option {
	return func(o *options) {
		o.otpCode = otpCode
	}
}
