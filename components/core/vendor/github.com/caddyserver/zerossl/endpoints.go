package zerossl

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CreateCertificate creates a certificate. After creating a certificate, its identifiers must be verified before
// the certificate can be downloaded. The CSR must have been fully created using x509.CreateCertificateRequest
// (its Raw field must be filled out).
func (c Client) CreateCertificate(ctx context.Context, csr *x509.CertificateRequest, validityDays int) (CertificateObject, error) {
	payload := struct {
		CertificateDomains        string `json:"certificate_domains"`
		CertificateCSR            string `json:"certificate_csr"`
		CertificateValidityDays   int    `json:"certificate_validity_days,omitempty"`
		StrictDomains             int    `json:"strict_domains,omitempty"`
		ReplacementForCertificate string `json:"replacement_for_certificate,omitempty"`
	}{
		CertificateDomains:      strings.Join(identifiersFromCSR(csr), ","),
		CertificateCSR:          csr2pem(csr.Raw),
		CertificateValidityDays: validityDays,
		StrictDomains:           1,
	}

	var result CertificateObject
	if err := c.httpPost(ctx, "/certificates", nil, payload, &result); err != nil {
		return CertificateObject{}, err
	}

	return result, nil
}

// VerifyIdentifiers tells ZeroSSL that you are ready to prove control over your domain/IP using the method specified.
// The credentials from CreateCertificate must be used to verify identifiers. At least one email is required if using
// email verification method.
func (c Client) VerifyIdentifiers(ctx context.Context, certificateID string, method VerificationMethod, emails []string) (CertificateObject, error) {
	payload := struct {
		ValidationMethod VerificationMethod `json:"validation_method"`
		ValidationEmail  string             `json:"validation_email,omitempty"`
	}{
		ValidationMethod: method,
	}
	if method == EmailVerification && len(emails) > 0 {
		payload.ValidationEmail = strings.Join(emails, ",")
	}

	endpoint := fmt.Sprintf("/certificates/%s/challenges", url.QueryEscape(certificateID))

	var result CertificateObject
	if err := c.httpPost(ctx, endpoint, nil, payload, &result); err != nil {
		return CertificateObject{}, err
	}

	return result, nil
}

// DownloadCertificateFile writes the certificate bundle as a zip file to the provided output writer.
func (c Client) DownloadCertificateFile(ctx context.Context, certificateID string, includeCrossSigned bool, output io.Writer) error {
	endpoint := fmt.Sprintf("/certificates/%s/download", url.QueryEscape(certificateID))

	qs := url.Values{}
	if includeCrossSigned {
		qs.Set("include_cross_signed", "1")
	}

	url := c.url(endpoint, qs)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient().Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: HTTP %d", resp.StatusCode)
	}

	if _, err := io.Copy(output, resp.Body); err != nil {
		return err
	}

	return nil
}

func (c Client) DownloadCertificate(ctx context.Context, certificateID string, includeCrossSigned bool) (CertificateBundle, error) {
	endpoint := fmt.Sprintf("/certificates/%s/download/return", url.QueryEscape(certificateID))

	qs := url.Values{}
	if includeCrossSigned {
		qs.Set("include_cross_signed", "1")
	}

	var result CertificateBundle
	if err := c.httpGet(ctx, endpoint, qs, &result); err != nil {
		return CertificateBundle{}, err
	}

	return result, nil
}

func (c Client) GetCertificate(ctx context.Context, certificateID string) (CertificateObject, error) {
	endpoint := fmt.Sprintf("/certificates/%s", url.QueryEscape(certificateID))

	var result CertificateObject
	if err := c.httpGet(ctx, endpoint, nil, &result); err != nil {
		return CertificateObject{}, err
	}

	return result, nil
}

// ListCertificateParameters specifies how to search or list certificates on the account.
// An empty set of parameters will return no results.
type ListCertificatesParameters struct {
	// Return certificates with this status.
	Status string

	// Return these types of certificates.
	Type string

	// The CommonName or SAN.
	Search string

	// The page number. Default: 1
	Page int

	// How many per page. Default: 100
	Limit int
}

func (c Client) ListCertificates(ctx context.Context, params ListCertificatesParameters) (CertificateList, error) {
	qs := url.Values{}
	if params.Status != "" {
		qs.Set("certificate_status", params.Status)
	}
	if params.Type != "" {
		qs.Set("certificate_type", params.Type)
	}
	if params.Search != "" {
		qs.Set("search", params.Search)
	}
	if params.Limit != 0 {
		qs.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Page != 0 {
		qs.Set("page", strconv.Itoa(params.Page))
	}

	var result CertificateList
	if err := c.httpGet(ctx, "/certificates", qs, &result); err != nil {
		return CertificateList{}, err
	}

	return result, nil
}

func (c Client) VerificationStatus(ctx context.Context, certificateID string) (ValidationStatus, error) {
	endpoint := fmt.Sprintf("/certificates/%s/status", url.QueryEscape(certificateID))

	var result ValidationStatus
	if err := c.httpGet(ctx, endpoint, nil, &result); err != nil {
		return ValidationStatus{}, err
	}

	return result, nil
}

func (c Client) ResendVerificationEmail(ctx context.Context, certificateID string) error {
	endpoint := fmt.Sprintf("/certificates/%s/challenges/email", url.QueryEscape(certificateID))

	var result struct {
		Success anyBool `json:"success"`
	}
	if err := c.httpGet(ctx, endpoint, nil, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("got %v without any error status", result)
	}

	return nil
}

// Only revoke a certificate if the private key is compromised, the certificate was a mistake, or
// the identifiers are no longer in use. Do not revoke a certificate when renewing it.
func (c Client) RevokeCertificate(ctx context.Context, certificateID string, reason RevocationReason) error {
	endpoint := fmt.Sprintf("/certificates/%s/revoke", url.QueryEscape(certificateID))

	qs := url.Values{"reason": []string{string(reason)}}

	var result struct {
		Success anyBool `json:"success"`
	}
	if err := c.httpGet(ctx, endpoint, qs, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("got %v without any error status", result)
	}

	return nil
}

// CancelCertificate cancels a certificate that has not been issued yet (is in draft or pending_validation state).
func (c Client) CancelCertificate(ctx context.Context, certificateID string) error {
	endpoint := fmt.Sprintf("/certificates/%s/cancel", url.QueryEscape(certificateID))

	var result struct {
		Success anyBool `json:"success"`
	}
	if err := c.httpPost(ctx, endpoint, nil, nil, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("got %v without any error status", result)
	}

	return nil
}

// ValidateCSR sends the CSR to ZeroSSL for validation. Pass in the ASN.1 DER-encoded bytes;
// this is found in x509.CertificateRequest.Raw after calling x5p9.CreateCertificateRequest.
func (c Client) ValidateCSR(ctx context.Context, csrASN1DER []byte) error {
	payload := struct {
		CSR string `json:"csr"`
	}{
		CSR: csr2pem(csrASN1DER),
	}

	var result struct {
		Valid bool `json:"valid"`
		Error any  `json:"error"`
	}
	if err := c.httpPost(ctx, "/validation/csr", nil, payload, &result); err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("invalid CSR: %v", result.Error)
	}
	return nil
}

func (c Client) GenerateEABCredentials(ctx context.Context) (keyID, hmacKey string, err error) {
	var result struct {
		APIError
		EABKID     string `json:"eab_kid"`
		EABHMACKey string `json:"eab_hmac_key"`
	}
	err = c.httpPost(ctx, "/acme/eab-credentials", nil, nil, &result)
	if err != nil {
		return
	}
	if !result.Success {
		err = fmt.Errorf("failed to create EAB credentials: %v", result.APIError)
	}
	return result.EABKID, result.EABHMACKey, err
}
