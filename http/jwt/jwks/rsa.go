package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
)

// RSA parses a JSONKey and turns it into an RSA public key.
func (j *jwkImpl) rsa() (publicKey *rsa.PublicKey, err error) {

	// Check if the key has already been computed.
	if j.precomputed != nil {
		var ok bool
		if publicKey, ok = j.precomputed.(*rsa.PublicKey); ok {
			return publicKey, nil
		}
	}

	// Confirm everything needed is present.
	if j.key.Exponent == "" || j.key.Modulus == "" {
		return nil, fmt.Errorf("%w: rsa", ErrMissingAssets)
	}

	// Decode the exponent from Base64.
	//
	// According to RFC 7518, this is a Base64 URL unsigned integer.
	// https://tools.ietf.org/html/rfc7518#section-6.3
	var exponent []byte
	if exponent, err = base64.RawURLEncoding.DecodeString(j.key.Exponent); err != nil {
		return nil, err
	}

	// Decode the modulus from Base64.
	var modulus []byte
	if modulus, err = base64.RawURLEncoding.DecodeString(j.key.Modulus); err != nil {
		return nil, err
	}

	/*
		cert := "-----BEGIN CERTIFICATE-----\n" + k.X5c[0] + "\n-----END CERTIFICATE-----"

		publicKey, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			continue
		}
	*/

	// Create the RSA public key.
	publicKey = &rsa.PublicKey{}

	// Turn the exponent into an integer.
	//
	// According to RFC 7517, these numbers are in big-endian format.
	// https://tools.ietf.org/html/rfc7517#appendix-A.1
	publicKey.E = int(big.NewInt(0).SetBytes(exponent).Uint64())

	// Turn the modulus into a *big.Int.
	publicKey.N = big.NewInt(0).SetBytes(modulus)

	// Keep the public key so it won't have to be computed every time.
	j.precomputed = publicKey

	return publicKey, nil
}
