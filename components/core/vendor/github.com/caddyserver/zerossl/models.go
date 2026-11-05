package zerossl

import "fmt"

type APIError struct {
	Success   anyBool `json:"success"`
	ErrorInfo struct {
		Code int    `json:"code"`
		Type string `json:"type"`

		// for domain verification only; each domain is grouped into its
		// www and non-www variant for CNAME validation, or its URL
		// for HTTP validation
		Details map[string]map[string]ValidationError `json:"details"`
	} `json:"error"`
}

func (ae APIError) Error() string {
	if ae.ErrorInfo.Code == 0 && ae.ErrorInfo.Type == "" && len(ae.ErrorInfo.Details) == 0 {
		return "<missing error info>"
	}
	return fmt.Sprintf("API error %d: %s (details=%v)",
		ae.ErrorInfo.Code, ae.ErrorInfo.Type, ae.ErrorInfo.Details)
}

type ValidationError struct {
	CNAMEValidationError
	HTTPValidationError
}

type CNAMEValidationError struct {
	CNAMEFound    int    `json:"cname_found"`
	RecordCorrect int    `json:"record_correct"`
	TargetHost    string `json:"target_host"`
	TargetRecord  string `json:"target_record"`
	ActualRecord  string `json:"actual_record"`
}

type HTTPValidationError struct {
	FileFound int    `json:"file_found"`
	Error     bool   `json:"error"`
	ErrorSlug string `json:"error_slug"`
	ErrorInfo string `json:"error_info"`
}

type CertificateObject struct {
	ID                string  `json:"id"` // "certificate hash"
	Type              string  `json:"type"`
	CommonName        string  `json:"common_name"`
	AdditionalDomains string  `json:"additional_domains"`
	Created           string  `json:"created"`
	Expires           string  `json:"expires"`
	Status            string  `json:"status"`
	ValidationType    *string `json:"validation_type,omitempty"`
	ValidationEmails  *string `json:"validation_emails,omitempty"`
	ReplacementFor    string  `json:"replacement_for,omitempty"`
	FingerprintSHA1   *string `json:"fingerprint_sha1"`
	BrandValidation   any     `json:"brand_validation"`
	Validation        *struct {
		EmailValidation map[string][]string         `json:"email_validation,omitempty"`
		OtherMethods    map[string]ValidationObject `json:"other_methods,omitempty"`
	} `json:"validation,omitempty"`
}

type ValidationObject struct {
	FileValidationURLHTTP  string   `json:"file_validation_url_http"`
	FileValidationURLHTTPS string   `json:"file_validation_url_https"`
	FileValidationContent  []string `json:"file_validation_content"`
	CnameValidationP1      string   `json:"cname_validation_p1"`
	CnameValidationP2      string   `json:"cname_validation_p2"`
}

type CertificateBundle struct {
	CertificateCrt string `json:"certificate.crt"`
	CABundleCrt    string `json:"ca_bundle.crt"`
}

type CertificateList struct {
	TotalCount     int                 `json:"total_count"`
	ResultCount    int                 `json:"result_count"`
	Page           string              `json:"page"` // don't ask me why this is a string
	Limit          int                 `json:"limit"`
	ACMEUsageLevel string              `json:"acmeUsageLevel"`
	ACMELocked     bool                `json:"acmeLocked"`
	Results        []CertificateObject `json:"results"`
}

type ValidationStatus struct {
	ValidationCompleted int `json:"validation_completed"`
	Details             map[string]struct {
		Method string `json:"method"`
		Status string `json:"status"`
	} `json:"details"`
}
