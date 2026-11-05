package jwks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"
	"math/big"
)

const (
	// p256 represents a cryptographic elliptical curve type.
	p256 = "P-256"

	// p384 represents a cryptographic elliptical curve type.
	p384 = "P-384"

	// p521 represents a cryptographic elliptical curve type.
	p521 = "P-521"
)

// ECDSA parses a JSONKey and turns it into an ECDSA public key.
func (j *jwkImpl) ecdsa() (publicKey *ecdsa.PublicKey, err error) {
	// Check if the key has already been computed.
	if j.precomputed != nil {
		var ok bool
		if publicKey, ok = j.precomputed.(*ecdsa.PublicKey); ok {
			return publicKey, nil
		}
	}

	// Confirm everything needed is present.
	if j.key.X == "" || j.key.Y == "" || j.key.Curve == "" {
		return nil, fmt.Errorf("%w: ecdsa", ErrMissingAssets)
	}

	// Decode the X coordinate from Base64.
	//
	// According to RFC 7518, this is a Base64 URL unsigned integer.
	// https://tools.ietf.org/html/rfc7518#section-6.3
	var xCoordinate []byte
	if xCoordinate, err = base64.RawURLEncoding.DecodeString(j.key.X); err != nil {
		return nil, err
	}

	// Decode the Y coordinate from Base64.
	var yCoordinate []byte
	if yCoordinate, err = base64.RawURLEncoding.DecodeString(j.key.Y); err != nil {
		return nil, err
	}

	// Create the ECDSA public key.
	publicKey = &ecdsa.PublicKey{}

	// Set the curve type.
	var curve elliptic.Curve
	switch j.key.Curve {
	case p256:
		curve = elliptic.P256()
	case p384:
		curve = elliptic.P384()
	case p521:
		curve = elliptic.P521()
	}
	publicKey.Curve = curve

	// Turn the X coordinate into *big.Int.
	//
	// According to RFC 7517, these numbers are in big-endian format.
	// https://tools.ietf.org/html/rfc7517#appendix-A.1
	publicKey.X = big.NewInt(0).SetBytes(xCoordinate)

	// Turn the Y coordinate into a *big.Int.
	publicKey.Y = big.NewInt(0).SetBytes(yCoordinate)

	// Keep the public key so it won't have to be computed every time.
	j.precomputed = publicKey

	return publicKey, nil
}
