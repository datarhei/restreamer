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

package acme

import (
	"context"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ErrUnsupported is used to indicate lack of support by an ACME server.
var ErrUnsupported = fmt.Errorf("unsupported by ACME server")

// RenewalInfo "is a new resource type introduced to ACME protocol.
// This new resource allows clients to query the server for suggestions
// on when they should renew certificates."
//
// ACME Renewal Information (ARI):
// https://www.ietf.org/archive/id/draft-ietf-acme-ari-03.html §4.2
//
// This is a DRAFT specification and the API is subject to change.
type RenewalInfo struct {
	// suggestedWindow (object, required): A JSON object with two keys,
	// "start" and "end", whose values are timestamps, encoded in the
	// format specified in [RFC3339], which bound the window of time
	// in which the CA recommends renewing the certificate.
	SuggestedWindow struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"suggestedWindow"`

	// explanationURL (string, optional): A URL pointing to a page which may
	// explain why the suggested renewal window is what it is. For example,
	// it may be a page explaining the CA's dynamic load-balancing strategy,
	// or a page documenting which certificates are affected by a mass
	// revocation event. Conforming clients SHOULD provide this URL to their
	// operator, if present.
	ExplanationURL string `json:"explanationURL,omitempty"`

	// The following fields are not part of the RenewalInfo object in
	// the ARI spec, but are important for proper conformance to the
	// spec, and are practically useful for implementators:

	// "The unique identifer is constructed by concatenating the
	// base64url-encoding Section 5 of [RFC4648] of the bytes of the
	// keyIdentifier field of certificate's Authority Key Identifier
	// (AKI) Section 4.2.1.1 of [RFC5280] extension, a literal period,
	// and the base64url-encoding of the bytes of the DER encoding of
	// the certificate's Serial Number (without the tag and length bytes).
	// All trailing "=" characters MUST be stripped from both parts of
	// the unique identifier."
	//
	// We generate this once and store it so the certificate does not
	// need to be stored in its decoded form or decoded multiple times.
	UniqueIdentifier string `json:"_uniqueIdentifier,omitempty"`

	// The next poll time based on the Retry-After response header for
	// the benefit of the caller for scheduling renewals. If specified,
	// GetRenewalInfo should not be called again before this time.
	//
	// "The server SHOULD include a Retry-After header indicating the polling
	// interval that the ACME server recommends. Conforming clients SHOULD
	// query the renewalInfo URL again after the Retry-After period has passed,
	// as the server may provide a different suggestedWindow."
	RetryAfter *time.Time `json:"_retryAfter,omitempty"`

	// The client should "select a uniform random time within the suggested
	// window." We select this time when getting the renewal info from the
	// server, though this behavior is ambiguous:
	// https://github.com/aarongable/draft-acme-ari/issues/70
	SelectedTime time.Time `json:"_selectedTime"`
}

// NeedsRefresh returns true if the renewal info needs updating.
// It returns false otherwise, or if the renewal info is empty
// (window is missing), assuming that there is no ARI available.
func (ari RenewalInfo) NeedsRefresh() bool {
	if !ari.HasWindow() {
		return false
	}
	if ari.RetryAfter == nil {
		// TODO: this seems like an unlikely condition, but we could be smart in its absence, like based on the window... play it safe for now though and just always be updating I guess
		return true
	}
	return time.Now().After(*ari.RetryAfter)
}

// HasWindow returns true if this ARI has a window. If not,
// it's likely because ARI is not supported or available.
func (ari RenewalInfo) HasWindow() bool {
	return !ari.SuggestedWindow.Start.IsZero() && !ari.SuggestedWindow.End.IsZero()
}

// SameWindow returns true if this ARI has the same window as the ARI passed in.
// Note that suggested windows can move in either direction, expand, or contract,
// so this method compares both start and end values for exact equality.
func (ari RenewalInfo) SameWindow(other RenewalInfo) bool {
	return ari.SuggestedWindow.Start.Equal(other.SuggestedWindow.Start) &&
		ari.SuggestedWindow.End.Equal(other.SuggestedWindow.End)
}

// GetRenewalInfo returns the ACME Renewal Information (ARI) for the certificate.
// It fills in the Retry-After value, if present, onto the returned struct so
// the caller can poll appropriately. If the ACME server does not support ARI,
// an error wrapping ErrUnsupported will be returned.
func (c *Client) GetRenewalInfo(ctx context.Context, leafCert *x509.Certificate) (RenewalInfo, error) {
	if err := c.provision(ctx); err != nil {
		return RenewalInfo{}, err
	}
	if c.dir.RenewalInfo == "" {
		return RenewalInfo{}, fmt.Errorf("%w: directory does not indicate ARI support (missing renewalInfo)", ErrUnsupported)
	}

	if c.Logger != nil {
		c.Logger.Debug("getting renewal info", zap.Strings("names", leafCert.DNSNames))
	}

	certID, err := ARIUniqueIdentifier(leafCert)
	if err != nil {
		return RenewalInfo{}, err
	}

	var ari RenewalInfo
	resp, err := c.httpReq(ctx, http.MethodGet, c.ariEndpoint(certID), nil, &ari)
	if err != nil {
		return RenewalInfo{}, err
	}
	ari.UniqueIdentifier = certID

	// "The server SHOULD include a Retry-After header indicating the polling
	// interval that the ACME server recommends." draft-ietf-acme-ari-03 §4.2
	raTime, err := retryAfterTime(resp)
	if err != nil && c.Logger != nil {
		c.Logger.Error("invalid Retry-After value", zap.Error(err))
	}
	if !raTime.IsZero() {
		ari.RetryAfter = &raTime
	}

	// "Conforming clients MUST attempt renewal at a time of their choosing
	// based on the suggested renewal window. ... Select a uniform random
	// time within the suggested window." §4.2
	// TODO: It's unclear whether this time should be selected once
	// or every time the client wakes to check ARI (see step 5 of the
	// recommended algorithm); I've inquired here:
	// https://github.com/aarongable/draft-acme-ari/issues/70
	// We add 1 to the start time since we are dealing in seconds for
	// simplicity, but the server may provide sub-second timestamps.
	start, end := ari.SuggestedWindow.Start.Unix()+1, ari.SuggestedWindow.End.Unix()
	ari.SelectedTime = time.Unix(rand.Int63n(end-start)+start, 0).UTC()

	if c.Logger != nil {
		c.Logger.Info("got renewal info",
			zap.Strings("names", leafCert.DNSNames),
			zap.Time("window_start", ari.SuggestedWindow.Start),
			zap.Time("window_end", ari.SuggestedWindow.End),
			zap.Time("selected_time", ari.SelectedTime),
			zap.Timep("recheck_after", ari.RetryAfter),
			zap.String("explanation_url", ari.ExplanationURL),
		)
	}

	return ari, nil
}

// ariEndpoint returns the ARI endpoint URI for certificate with the
// given ARI certificate ID, according to the configured CA's directory.
func (c *Client) ariEndpoint(ariCertID string) string {
	if c.dir.RenewalInfo == "" || ariCertID == "" {
		return ""
	}
	return c.dir.RenewalInfo + "/" + ariCertID
}

// ARIUniqueIdentifier returns the unique identifier for the certificate
// as used by ACME Renewal Information.
// EXPERIMENTAL: ARI is a draft RFC spec: draft-ietf-acme-ari-03
func ARIUniqueIdentifier(leafCert *x509.Certificate) (string, error) {
	if leafCert.SerialNumber == nil {
		return "", fmt.Errorf("no serial number")
	}
	// TODO: Let's Encrypt's reference implementation switched from using
	// SerialNumber.Bytes() to this method, which seems less efficient,
	// but yields the same results !? I asked about it here:
	// https://github.com/letsencrypt/website/issues/1670
	serialDER, err := asn1.Marshal(leafCert.SerialNumber)
	if err != nil {
		return "", err
	}
	if len(serialDER) < 3 {
		return "", fmt.Errorf("serial number DER too short: %d (%x)", len(serialDER), serialDER)
	}
	// skip tag and length; extract only integer bytes
	return base64.RawURLEncoding.EncodeToString(leafCert.AuthorityKeyId) + "." +
		base64.RawURLEncoding.EncodeToString(serialDER[2:]), nil // skip tag and length, just use integer part
}
