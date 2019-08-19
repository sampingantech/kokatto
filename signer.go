package kokatto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

// Signer interface for signed request
type Signer interface {
	QueryString() (string, error)
	SetSignature(string)
	Validate(string) bool
}

// Sign request using HMAC-SHA256 private key
// signature = HMAC-SHA256(MD5(QueryString), secret)
func Sign(req Signer, secret []byte) error {
	if len(secret) == 0 {
		return errors.New("invalid key")
	}
	// md5 query string
	hasher := md5.New()
	query, err := req.QueryString()
	if err != nil {
		return err
	}
	_, err = hasher.Write([]byte(query))
	if err != nil {
		return errors.Wrap(err, "err write md5")
	}
	queryHash := hex.EncodeToString(hasher.Sum(nil))

	// hmac query string
	h := hmac.New(sha256.New, secret)
	_, err = h.Write([]byte(queryHash))
	if err != nil {
		return errors.Wrap(err, "err write hmac")
	}

	req.SetSignature(strings.ToUpper(hex.EncodeToString(h.Sum(nil))))
	return nil
}
