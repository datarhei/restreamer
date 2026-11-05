// Copyright 2020 Matthew Holt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package acmez

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/mholt/acmez/v2/acme"
	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
	"golang.org/x/net/idna"
)

// NewCSR creates and signs a Certificate Signing Request (CSR) for the given subject
// identifiers (SANs) with the private key.
//
// If you need extensions or other customizations, this function is too opinionated.
// Instead, create a new(x509.CertificateRequest), then fill out the relevant fields
// (SANs, extensions, etc.), send it into x509.CreateCertificateRequest(), then pass
// that result into x509.ParseCertificateRequest() to get the final, parsed CSR. We
// chose this API to offer the most common convenience functions, but also to give
// users advanced flexibility when needed, all while reducing allocations from
// encoding & decoding each CSR and minimizing having to pass the private key around.
//
// Supported SAN types are IPs, email addresses, URIs, and DNS names.
//
// EXPERIMENTAL: This API is subject to change or removal without a major version bump.
func NewCSR(privateKey crypto.Signer, sans []string) (*x509.CertificateRequest, error) {
	if len(sans) == 0 {
		return nil, fmt.Errorf("no SANs provided: %v", sans)
	}

	csrTemplate := new(x509.CertificateRequest)
	for _, name := range sans {
		if ip := net.ParseIP(name); ip != nil {
			csrTemplate.IPAddresses = append(csrTemplate.IPAddresses, ip)
		} else if strings.Contains(name, "@") {
			csrTemplate.EmailAddresses = append(csrTemplate.EmailAddresses, name)
		} else if u, err := url.Parse(name); err == nil && strings.Contains(name, "/") {
			csrTemplate.URIs = append(csrTemplate.URIs, u)
		} else {
			// "The domain name MUST be encoded in the form in which it would appear
			// in a certificate.  That is, it MUST be encoded according to the rules
			// in Section 7 of [RFC5280]." §7.1.4
			normalizedName, err := idna.ToASCII(name)
			if err != nil {
				return nil, fmt.Errorf("converting identifier '%s' to ASCII: %v", name, err)
			}
			csrTemplate.DNSNames = append(csrTemplate.DNSNames, normalizedName)
		}
	}

	// to properly fill out the CSR, we need to create it, then parse it
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privateKey)
	if err != nil {
		return nil, fmt.Errorf("generating CSR: %v", err)
	}
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		return nil, fmt.Errorf("parsing generated CSR: %v", err)
	}

	return csr, nil
}

// OrderParameters contains high-level input parameters for ACME transactions,
// the state of which are represented by Order objects. This type is used as a
// convenient high-level way to convey alk the configuration needed to obtain a
// certificate (except the private key, which is provided separately to prevent
// inadvertent exposure of secret material) through ACME in one consolidated value.
//
// Account, Identifiers, and CSR fields are REQUIRED.
type OrderParameters struct {
	// The ACME account with which to perform certificate operations.
	// It should already be registered with the server and have a
	// "valid" status.
	Account acme.Account

	// The list of identifiers for which to issue the certificate.
	// Identifiers may become Subject Alternate Names (SANs) in the
	// certificate. This slice must be consistent with the SANs
	// listed in the CSR. The OrderFromCSR() function can be
	// called to ensure consistency in most cases.
	//
	// Supported identifier types are currently: dns, ip,
	// permanent-identifier, and hardware-module.
	Identifiers []acme.Identifier

	// CSR is a type that can provide the Certificate Signing
	// Request, which is needed when finalizing the ACME order.
	// It is invoked after challenges have completed and before
	// finalization.
	CSR CSRSource

	// Optionally customize the lifetime of the certificate by
	// specifying the NotBefore and/or NotAfter dates for the
	// certificate. Not all CAs support this. Check your CA's
	// ACME service documentation.
	NotBefore, NotAfter time.Time

	// Set this to the old certificate if a certificate is being renewed.
	//
	// DRAFT: EXPERIMENTAL ARI DRAFT SPEC. Subject to change/removal.
	Replaces *x509.Certificate
}

// OrderParametersFromCSR makes a valid OrderParameters from the given CSR.
// If necessary, the returned parameters may be further customized before using.
//
// EXPERIMENTAL: This API is subject to change or removal without a major version bump.
func OrderParametersFromCSR(account acme.Account, csr *x509.CertificateRequest) (OrderParameters, error) {
	ids, err := createIdentifiersUsingCSR(csr)
	if err != nil {
		return OrderParameters{}, err
	}
	if len(ids) == 0 {
		return OrderParameters{}, errors.New("no subjects found in CSR")
	}
	return OrderParameters{
		Account:     account,
		Identifiers: ids,
		CSR:         StaticCSR(csr),
	}, nil
}

// CSRSource is an interface that provides users of this
// package the ability to provide a CSR as part of the
// ACME flow. This allows the final CSR to be provided
// just before the Order is finalized, which is useful
// for certain challenge types (e.g. device-attest-01,
// where the key used for signing the CSR doesn't exist
// until the challenge has been validated).
//
// EXPERIMENTAL: Subject to change (though unlikely, and nothing major).
type CSRSource interface {
	// CSR returns a Certificate Signing Request that will be
	// given to the ACME server. This function is called after
	// an ACME challenge completion and before order finalization.
	//
	// The returned CSR must have the Raw field populated with the
	// DER-encoded certificate request signed by the private key.
	// Typically this involves creating a template CSR, then calling
	// x509.CreateCertificateRequest(), then x509.ParseCertificateRequest()
	// on the output. That should return a valid CSR. The NewCSR()
	// function in this package does this for you, but if you need more
	// control you should make it yourself.
	//
	// The Subject CommonName field is NOT considered.
	CSR(context.Context, []acme.Identifier) (*x509.CertificateRequest, error)
}

// StaticCSR returns a CSRSource that simply returns the input CSR.
func StaticCSR(csr *x509.CertificateRequest) CSRSource { return staticCSR{csr} }

// staticCSR is a CSRSource that returns an existing CSR.
type staticCSR struct{ *x509.CertificateRequest }

// CSR returns the associated CSR.
func (cs staticCSR) CSR(_ context.Context, _ []acme.Identifier) (*x509.CertificateRequest, error) {
	return cs.CertificateRequest, nil
}

// Interface guard
var _ CSRSource = (*staticCSR)(nil)

var (
	oidExtensionSubjectAltName = []int{2, 5, 29, 17}
	oidPermanentIdentifier     = []int{1, 3, 6, 1, 5, 5, 7, 8, 3}
	oidHardwareModuleName      = []int{1, 3, 6, 1, 5, 5, 7, 8, 4}
)

// RFC 5280 - https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.6
//
//	OtherName ::= SEQUENCE {
//	  type-id    OBJECT IDENTIFIER,
//	  value      [0] EXPLICIT ANY DEFINED BY type-id }
type otherName struct {
	TypeID asn1.ObjectIdentifier
	Value  asn1.RawValue
}

// permanentIdentifier is defined in RFC 4043 as an optional feature that can be
// used by a CA to indicate that two or more certificates relate to the same
// entity.
//
// The OID defined for this SAN is "1.3.6.1.5.5.7.8.3".
//
// See https://www.rfc-editor.org/rfc/rfc4043
//
//	PermanentIdentifier ::= SEQUENCE {
//	  identifierValue    UTF8String OPTIONAL,
//	  assigner           OBJECT IDENTIFIER OPTIONAL
//	}
type permanentIdentifier struct {
	IdentifierValue string                `asn1:"utf8,optional"`
	Assigner        asn1.ObjectIdentifier `asn1:"optional"`
}

// hardwareModuleName is defined in RFC 4108 as an optional feature that can be
// used to identify a hardware module.
//
// The OID defined for this SAN is "1.3.6.1.5.5.7.8.4".
//
// See https://www.rfc-editor.org/rfc/rfc4108#section-5
//
//	HardwareModuleName ::= SEQUENCE {
//	  hwType OBJECT IDENTIFIER,
//	  hwSerialNum OCTET STRING
//	}
type hardwareModuleName struct {
	Type         asn1.ObjectIdentifier
	SerialNumber []byte `asn1:"tag:4"`
}

func forEachSAN(der cryptobyte.String, callback func(tag int, data []byte) error) error {
	if !der.ReadASN1(&der, cryptobyte_asn1.SEQUENCE) {
		return errors.New("invalid subject alternative name extension")
	}
	for !der.Empty() {
		var san cryptobyte.String
		var tag cryptobyte_asn1.Tag
		if !der.ReadAnyASN1Element(&san, &tag) {
			return errors.New("invalid subject alternative name extension")
		}
		if err := callback(int(tag^0x80), san); err != nil {
			return err
		}
	}

	return nil
}

// createIdentifiersUsingCSR extracts the list of ACME identifiers from the
// given Certificate Signing Request.
func createIdentifiersUsingCSR(csr *x509.CertificateRequest) ([]acme.Identifier, error) {
	var ids []acme.Identifier
	for _, name := range csr.DNSNames {
		ids = append(ids, acme.Identifier{
			Type:  "dns", // RFC 8555 §9.7.7
			Value: name,
		})
	}
	for _, ip := range csr.IPAddresses {
		ids = append(ids, acme.Identifier{
			Type:  "ip", // RFC 8738
			Value: ip.String(),
		})
	}
	for _, email := range csr.EmailAddresses {
		ids = append(ids, acme.Identifier{
			Type:  "email", // RFC 8823
			Value: email,
		})
	}

	// Extract permanent identifiers and hardware module values.
	// This block will ignore errors.
	for _, ext := range csr.Extensions {
		if ext.Id.Equal(oidExtensionSubjectAltName) {
			err := forEachSAN(ext.Value, func(tag int, data []byte) error {
				var on otherName
				if rest, err := asn1.UnmarshalWithParams(data, &on, "tag:0"); err != nil || len(rest) > 0 {
					return nil
				}

				switch {
				case on.TypeID.Equal(oidPermanentIdentifier):
					var pi permanentIdentifier
					if _, err := asn1.Unmarshal(on.Value.Bytes, &pi); err == nil {
						ids = append(ids, acme.Identifier{
							Type:  "permanent-identifier", // draft-acme-device-attest-00 §3
							Value: pi.IdentifierValue,
						})
					}
				case on.TypeID.Equal(oidHardwareModuleName):
					var hmn hardwareModuleName
					if _, err := asn1.Unmarshal(on.Value.Bytes, &hmn); err == nil {
						ids = append(ids, acme.Identifier{
							Type:  "hardware-module", // draft-acme-device-attest-00 §4
							Value: string(hmn.SerialNumber),
						})
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return ids, nil
}
