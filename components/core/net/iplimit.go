package net

import (
	"fmt"
	"net"
	"strings"
)

// The IPLimiter interface allows to check whether a certain IP
// is allowed.
type IPLimiter interface {
	// Tests whether the IP is allowed in respect to the underlying implementation
	IsAllowed(ip string) bool
}

// IPLimit implements the IPLimiter interface by having an allow and block list
// of CIDR ranges.
type iplimit struct {
	// Array of allowed IP ranges
	allowlist []*net.IPNet

	// Array of blocked IP ranges
	blocklist []*net.IPNet
}

// NewIPLimiter creates a new IPLimiter with the given IP ranges for the
// allowed and blocked IPs. Empty strings are ignored. Returns an error
// if an invalid IP range has been found.
func NewIPLimiter(blocklist, allowlist []string) (IPLimiter, error) {
	ipl := &iplimit{}

	for _, ipblock := range blocklist {
		ipblock = strings.TrimSpace(ipblock)
		if len(ipblock) == 0 {
			continue
		}

		_, cidr, err := net.ParseCIDR(ipblock)
		if err != nil {
			return nil, fmt.Errorf("the IP block %s in the block list is invalid", ipblock)
		}

		ipl.blocklist = append(ipl.blocklist, cidr)
	}

	for _, ipblock := range allowlist {
		ipblock = strings.TrimSpace(ipblock)
		if len(ipblock) == 0 {
			continue
		}

		_, cidr, err := net.ParseCIDR(ipblock)
		if err != nil {
			return nil, fmt.Errorf("the IP block %s in the allow list is invalid", ipblock)
		}

		ipl.allowlist = append(ipl.allowlist, cidr)
	}

	return ipl, nil
}

// IsAllowed checks whether the provided IP is allowed according
// to the IP ranges in the allow- and blocklists.
func (ipl *iplimit) IsAllowed(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, r := range ipl.blocklist {
		if r.Contains(parsedIP) {
			return false
		}
	}

	if len(ipl.allowlist) == 0 {
		return true
	}

	for _, r := range ipl.allowlist {
		if r.Contains(parsedIP) {
			return true
		}
	}

	return false
}

type nulliplimiter struct{}

func NewNullIPLimiter() IPLimiter {
	return &nulliplimiter{}
}

func (ipl *nulliplimiter) IsAllowed(ip string) bool {
	return true
}
