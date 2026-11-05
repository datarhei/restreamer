package url

import (
	"net"
	"net/url"
	"regexp"
	"strings"
)

type URL struct {
	Scheme      string
	Opaque      string        // encoded opaque data
	User        *url.Userinfo // username and password information
	Host        string        // host or host:port
	RawPath     string        // path (relative paths may omit leading slash)
	RawQuery    string        // encoded query values, without '?'
	RawFragment string        // fragment for references, without '#'
}

func (u *URL) Hostname() string {
	if !strings.Contains(u.Host, ":") {
		return u.Host
	}

	hostname, _, _ := net.SplitHostPort(u.Host)

	return hostname
}

func (u *URL) Port() string {
	if !strings.Contains(u.Host, ":") {
		return ""
	}

	_, port, _ := net.SplitHostPort(u.Host)

	return port
}

var reScheme = regexp.MustCompile(`(?i)^([a-z][a-z0-9.+-:]*):/{1,3}`)

// Validate checks whether the given address is a valid URL, based on the
// relaxed version of Parse in this package.
func Validate(address string) error {
	_, err := Parse(address)

	return err
}

// Parse parses an URL into its components. It is a more relaxed version of
// url.Parse as it's not checking the escaping of the path, query, and fragment.
func Parse(address string) (*URL, error) {
	address, frag, _ := strings.Cut(address, "#")

	u := &URL{
		RawFragment: frag,
	}

	matches := reScheme.FindStringSubmatch(address)
	if matches != nil {
		u.Scheme = matches[1]
		address = strings.Replace(address, u.Scheme+":", "", 1)
	}

	address, query, _ := strings.Cut(address, "?")
	u.RawQuery = query

	if strings.HasPrefix(address, "///") {
		u.RawPath = strings.TrimPrefix(address, "//")
		return u, nil
	}

	if strings.HasPrefix(address, "//") {
		host, path, _ := strings.Cut(address[2:], "/")
		u.RawPath = "/" + path

		parsedHost, err := url.Parse("//" + host)
		if err != nil {
			return nil, err
		}

		u.User = parsedHost.User
		u.Host = parsedHost.Host

		return u, nil
	}

	if strings.HasPrefix(address, "/") {
		u.RawPath = address

		return u, nil
	}

	scheme, address, _ := strings.Cut(address, ":")

	u.Scheme = scheme
	u.Opaque = address

	return u, nil
}

// HasScheme returns whether the address has an URL scheme prefix
func HasScheme(address string) bool {
	if ok := reScheme.MatchString(address); ok {
		return true
	}

	return false
}

// Lookup returns the first resolved IP address of the host in the address.
// If the address doesn't have an URL scheme prefix, an empty string and no
// error is returned. If the address is an URL and the lookup fails, an
// error is returned.
func Lookup(address string) (string, error) {
	if !HasScheme(address) {
		return "", nil
	}

	u, err := Parse(address)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	if len(host) == 0 {
		return "", nil
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}

	return addrs[0], nil
}
