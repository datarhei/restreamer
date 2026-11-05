// Package zerossl implements the ZeroSSL REST API.
// See the API documentation on the ZeroSSL website: https://zerossl.com/documentation/api/
package zerossl

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// The base URL to the ZeroSSL API.
const BaseURL = "https://api.zerossl.com"

// ListAllCertificates returns parameters that lists all the certificates on the account;
// be sure to set Page and Limit if paginating.
func ListAllCertificates() ListCertificatesParameters {
	return ListCertificatesParameters{
		Status: "draft,pending_validation,issued,cancelled,revoked,expired",
	}
}

func identifiersFromCSR(csr *x509.CertificateRequest) []string {
	var identifiers []string
	if csr.Subject.CommonName != "" {
		// deprecated for like 20 years, but oh well
		identifiers = append(identifiers, csr.Subject.CommonName)
	}
	identifiers = append(identifiers, csr.DNSNames...)
	identifiers = append(identifiers, csr.EmailAddresses...)
	for _, ip := range csr.IPAddresses {
		identifiers = append(identifiers, ip.String())
	}
	for _, uri := range csr.URIs {
		identifiers = append(identifiers, uri.String())
	}
	return identifiers
}

func csr2pem(csrASN1DER []byte) string {
	return fmt.Sprintf("-----BEGIN CERTIFICATE REQUEST-----\n%s\n-----END CERTIFICATE REQUEST-----",
		base64.StdEncoding.EncodeToString(csrASN1DER))
}

// VerificationMethod represents a way of verifying identifiers with ZeroSSL.
type VerificationMethod string

// Verification methods.
const (
	EmailVerification VerificationMethod = "EMAIL"
	CNAMEVerification VerificationMethod = "CNAME_CSR_HASH"
	HTTPVerification  VerificationMethod = "HTTP_CSR_HASH"
	HTTPSVerification VerificationMethod = "HTTPS_CSR_HASH"
)

// RevocationReason represents various reasons for revoking a certificate.
type RevocationReason string

const (
	UnspecifiedReason    RevocationReason = "unspecified"          // default
	KeyCompromise        RevocationReason = "keyCompromise"        // lost control of private key
	AffiliationChanged   RevocationReason = "affiliationChanged"   // identify information changed
	Superseded           RevocationReason = "Superseded"           // certificate replaced -- do not revoke for this reason, however
	CessationOfOperation RevocationReason = "cessationOfOperation" // domains are no longer in use
)
