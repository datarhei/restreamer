// Copyright 2023 Matthew Holt
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
	"fmt"
	"strings"

	"github.com/mholt/acmez/v2/acme"
)

// MailReplyChallengeResponse builds an email response body including headers to reply to the
// email-reply-00 challenge email. This function only builds the email body; sending the
// message has to be performed by the caller of this function. The mailSubject and
// messageId come from the challenge mail, and if there is no reply-to header in the
// challenge email, the replyTo parameter should be empty.
func MailReplyChallengeResponse(c acme.Challenge, mailSubject string, messageId string, replyTo string) (string, error) {
	if replyTo == "" {
		replyTo = c.From
	}
	tokenPart1 := strings.TrimPrefix(mailSubject, "ACME: ")
	keyAuth, err := c.MailReply00KeyAuthorization(tokenPart1)
	if err != nil {
		return "", fmt.Errorf("failed creating key authorization: %w", err)
	}
	msg := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"In-Reply-To: %s\r\n"+
		"Subject: RE: ACME: %s\r\n"+
		"Content-Type: text/plain\r\n"+
		"\r\n"+
		"-----BEGIN ACME RESPONSE-----\r\n"+
		"%s\r\n"+
		"-----END ACME RESPONSE-----\r\n", replyTo, c.Identifier.Value, messageId, tokenPart1, keyAuth)
	return msg, nil
}
