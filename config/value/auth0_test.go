package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuth0Value(t *testing.T) {
	tenants := []Auth0Tenant{}

	v := NewTenantList(&tenants, nil, " ")
	require.Equal(t, "(empty)", v.String())

	v.Set("auth0://clientid@domain?aud=audience&user=user1&user=user2 auth0://domain2?aud=audience2&user=user3")
	require.Equal(t, []Auth0Tenant{
		{
			Domain:   "domain",
			ClientID: "clientid",
			Audience: "audience",
			Users:    []string{"user1", "user2"},
		},
		{
			Domain:   "domain2",
			Audience: "audience2",
			Users:    []string{"user3"},
		},
	}, tenants)
	require.Equal(t, "auth0://clientid@domain?aud=audience&user=user1&user=user2 auth0://domain2?aud=audience2&user=user3", v.String())
	require.NoError(t, v.Validate())

	v.Set("eyJkb21haW4iOiJkYXRhcmhlaS5ldS5hdXRoMC5jb20iLCJhdWRpZW5jZSI6Imh0dHBzOi8vZGF0YXJoZWkuY29tL2NvcmUiLCJ1c2VycyI6WyJhdXRoMHx4eHgiXX0=")
	require.Equal(t, []Auth0Tenant{
		{
			Domain:   "datarhei.eu.auth0.com",
			ClientID: "",
			Audience: "https://datarhei.com/core",
			Users:    []string{"auth0|xxx"},
		},
	}, tenants)
	require.Equal(t, "auth0://datarhei.eu.auth0.com?aud=https%3A%2F%2Fdatarhei.com%2Fcore&user=auth0%7Cxxx", v.String())
	require.NoError(t, v.Validate())
}
