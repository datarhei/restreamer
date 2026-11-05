package copy

import "github.com/datarhei/core/v16/config/value"

func StringMap(src map[string]string) map[string]string {
	dst := make(map[string]string)

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func TenantSlice(src []value.Auth0Tenant) []value.Auth0Tenant {
	dst := Slice(src)

	for i, t := range src {
		dst[i].Users = Slice(t.Users)
	}

	return dst
}

func Slice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)

	return dst
}
