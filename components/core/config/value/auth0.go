package value

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// array of auth0 tenants

type Auth0Tenant struct {
	Domain   string   `json:"domain"`
	Audience string   `json:"audience"`
	ClientID string   `json:"clientid"`
	Users    []string `json:"users"`
}

func (a *Auth0Tenant) String() string {
	u := url.URL{
		Scheme: "auth0",
		Host:   a.Domain,
	}

	if len(a.ClientID) != 0 {
		u.User = url.User(a.ClientID)
	}

	q := url.Values{}
	q.Set("aud", a.Audience)

	for _, user := range a.Users {
		q.Add("user", user)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

type TenantList struct {
	p         *[]Auth0Tenant
	separator string
}

func NewTenantList(p *[]Auth0Tenant, val []Auth0Tenant, separator string) *TenantList {
	v := &TenantList{
		p:         p,
		separator: separator,
	}

	*p = val

	return v
}

// Set allows to set a tenant list in two formats:
// - a separator separated list of bas64 encoded Auth0Tenant JSON objects
// - a separator separated list of Auth0Tenant in URL representation: auth0://[clientid]@[domain]?aud=[audience]&user=...&user=...
func (s *TenantList) Set(val string) error {
	list := []Auth0Tenant{}

	for i, elm := range strings.Split(val, s.separator) {
		t := Auth0Tenant{}

		if strings.HasPrefix(elm, "auth0://") {
			data, err := url.Parse(elm)
			if err != nil {
				return fmt.Errorf("invalid url encoding of tenant %d: %w", i, err)
			}

			t.Domain = data.Host
			t.ClientID = data.User.Username()
			t.Audience = data.Query().Get("aud")
			t.Users = data.Query()["user"]
		} else {
			data, err := base64.StdEncoding.DecodeString(elm)
			if err != nil {
				return fmt.Errorf("invalid base64 encoding of tenant %d: %w", i, err)
			}

			if err := json.Unmarshal(data, &t); err != nil {
				return fmt.Errorf("invalid JSON in tenant %d: %w", i, err)
			}
		}

		list = append(list, t)
	}

	*s.p = list

	return nil
}

func (s *TenantList) String() string {
	if s.IsEmpty() {
		return "(empty)"
	}

	list := []string{}

	for _, t := range *s.p {
		list = append(list, t.String())
	}

	return strings.Join(list, s.separator)
}

func (s *TenantList) Validate() error {
	for i, t := range *s.p {
		if len(t.Domain) == 0 {
			return fmt.Errorf("the domain for tenant %d is missing", i)
		}

		if len(t.Audience) == 0 {
			return fmt.Errorf("the audience for tenant %d is missing", i)
		}
	}

	return nil
}

func (s *TenantList) IsEmpty() bool {
	return len(*s.p) == 0
}
