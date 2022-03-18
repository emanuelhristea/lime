package license

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	"golang.org/x/crypto/ed25519"
)

var (
	// ErrInvalidSignature is a ...
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrMalformedLicense is a ...
	ErrMalformedLicense = errors.New("malformed license")

	// generate new ed25519 key and replaces !!!
	privateKey = []byte("M7JsfVjXCj/60wflqhWisMh0tzHC7ozEB55EjsOT8ZXEgIn1/wXJhpPV47NDLrsuIc6gdcQcesQmyk2OBMmsqw==")
	publicKey  = []byte("xICJ9f8FyYaT1eOzQy67LiHOoHXEHHrEJspNjgTJrKs=")
)

// Role is
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
	GuestRole Role = "guest"
)

// License is a ...
type License struct {
	Iss string          `json:"iss"` // Issued By
	Cus string          `json:"cus"` // Customer ID
	Sub uint64          `json:"sub"` // Subscriber ID
	Typ string          `json:"typ"` // License Type
	Lim Limits          `json:"lim"` // License Limit (e.g. Site)
	Iat time.Time       `json:"iat"` // Issued At
	Exp time.Time       `json:"exp"` // Expires At
	Dat json.RawMessage `json:"dat"` // Metadata
}

type Key struct {
	Key   string `json:"key,omitempty"`
	Hash  string `json:"hash,omitempty"`
	Mac   string `json:"mac,omitempty"`
	Valid bool   `json:"valid,omitempty"`
}

// Limits is a ...
type Limits struct {
	Tandem  bool `json:"tandem"`
	Triaxis bool `json:"triaxis"`
	Robots  bool `json:"robots"`
	Period  int  `json:"expiry"`
	Devices int  `json:"devices"`
}

type Subscription struct {
	Plan       string `json:"plan"`                  // Subscription plan
	PurchaseID string `json:"purchase_id"`           // transaction id
	Limits     Limits `json:"limits"`                // License Limit (e.g. Site)
	InUse      int    `json:"in_use"`                //
	LicenseKey Key    `json:"license_key,omitempty"` //
	Role       Role   `json:"role"`                  //
	Status     bool   `json:"status"`                //
	ExpiresIn  int    `json:"expires_in"`
}

// Expired is a ...
func (l *License) Expired() bool {
	return !l.Exp.IsZero() && time.Now().After(l.Exp)
}

// Encode is a ...
func (l *License) Encode(privateKey ed25519.PrivateKey) ([]byte, error) {
	msg, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	sig := ed25519.Sign(privateKey, msg)
	buf := new(bytes.Buffer)
	buf.Write(sig)
	buf.Write(msg)

	block := &pem.Block{
		Type:  "LICENSE KEY",
		Bytes: buf.Bytes(),
	}
	return pem.EncodeToMemory(block), nil
}

// Decode is a ...
func Decode(data []byte, publicKey ed25519.PublicKey) (*License, error) {
	block, _ := pem.Decode(data)
	if block == nil || len(block.Bytes) < ed25519.SignatureSize {
		return nil, ErrMalformedLicense
	}

	sig := block.Bytes[:ed25519.SignatureSize]
	msg := block.Bytes[ed25519.SignatureSize:]

	verified := ed25519.Verify(publicKey, msg, sig)
	if !verified {
		return nil, ErrInvalidSignature
	}
	out := new(License)
	err := json.Unmarshal(msg, out)
	return out, err
}

// GetPrivateKey is a ...
func GetPrivateKey() ed25519.PrivateKey {
	key, _ := DecodePrivateKey(privateKey)
	return key
}

// GetPublicKey is a ...
func GetPublicKey() ed25519.PublicKey {
	key, _ := DecodePublicKey(publicKey)
	return key
}
