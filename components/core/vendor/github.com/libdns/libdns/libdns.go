// Package libdns defines core interfaces that should be implemented by DNS
// provider clients. They are small and idiomatic Go interfaces with
// well-defined semantics.
//
// Records are described independently of any particular zone, a convention
// that grants Record structs portability across zones. As such, record names
// are partially qualified, i.e. relative to the zone. For example, an A
// record called "sub" in zone "example.com." represents a fully-qualified
// domain name (FQDN) of "sub.example.com.". Implementations should expect
// that input records conform to this standard, while also ensuring that
// output records do; adjustments to record names may need to be made before
// or after provider API calls, for example, to maintain consistency with
// all other libdns packages. Helper functions are available in this package
// to convert between relative and absolute names.
//
// Although zone names are a required input, libdns does not coerce any
// particular representation of DNS zones; only records. Since zone name and
// records are separate inputs in libdns interfaces, it is up to the caller
// to pair a zone's name with its records in a way that works for them.
//
// All interface implementations must be safe for concurrent/parallel use,
// meaning 1) no data races, and 2) simultaneous method calls must result
// in either both their expected outcomes or an error.
//
// For example, if AppendRecords() is called at the same time and two API
// requests are made to the provider at the same time, the result of both
// requests must be visible after they both complete; if the provider does
// not synchronize the writing of the zone file and one request overwrites
// the other, then the client implementation must take care to synchronize
// on behalf of the incompetent provider. This synchronization need not be
// global; for example: the scope of synchronization might only need to be
// within the same zone, allowing multiple requests at once as long as all
// of them are for different zones. (Exact logic depends on the provider.)
package libdns

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RecordGetter can get records from a DNS zone.
type RecordGetter interface {
	// GetRecords returns all the records in the DNS zone.
	//
	// Implementations must honor context cancellation and be safe for
	// concurrent use.
	GetRecords(ctx context.Context, zone string) ([]Record, error)
}

// RecordAppender can non-destructively add new records to a DNS zone.
type RecordAppender interface {
	// AppendRecords creates the requested records in the given zone
	// and returns the populated records that were created. It never
	// changes existing records.
	//
	// Implementations must honor context cancellation and be safe for
	// concurrent use.
	AppendRecords(ctx context.Context, zone string, recs []Record) ([]Record, error)
}

// RecordSetter can set new or update existing records in a DNS zone.
type RecordSetter interface {
	// SetRecords updates the zone so that the records described in the
	// input are reflected in the output. It may create or overwrite
	// records or -- depending on the record type -- delete records to
	// maintain parity with the input. No other records are affected.
	// It returns the records which were set.
	//
	// Records that have an ID associating it with a particular resource
	// on the provider will be directly replaced. If no ID is given, this
	// method may use what information is given to do lookups and will
	// ensure that only necessary changes are made to the zone.
	//
	// Implementations must honor context cancellation and be safe for
	// concurrent use.
	SetRecords(ctx context.Context, zone string, recs []Record) ([]Record, error)
}

// RecordDeleter can delete records from a DNS zone.
type RecordDeleter interface {
	// DeleteRecords deletes the given records from the zone if they exist.
	// It returns the records that were deleted.
	//
	// Records that have an ID to associate it with a particular resource on
	// the provider will be directly deleted. If no ID is given, this method
	// may use what information is given to do lookups and delete only
	// matching records.
	//
	// Implementations must honor context cancellation and be safe for
	// concurrent use.
	DeleteRecords(ctx context.Context, zone string, recs []Record) ([]Record, error)
}

// ZoneLister can list available DNS zones.
type ZoneLister interface {
	// ListZones returns the list of available DNS zones for use by
	// other libdns methods.
	//
	// Implementations must honor context cancellation and be safe for
	// concurrent use.
	ListZones(ctx context.Context) ([]Zone, error)
}

// Record is a generalized representation of a DNS record.
//
// The values of this struct should be free of zone-file-specific syntax,
// except if this struct's fields do not sufficiently represent all the
// fields of a certain record type; in that case, the remaining data for
// which there are not specific fields should be stored in the Value as
// it appears in the zone file.
type Record struct {
	// provider-specific metadata
	ID string

	// general record fields
	Type  string
	Name  string // partially-qualified (relative to zone)
	Value string
	TTL   time.Duration

	// type-dependent record fields
	Priority uint // HTTPS, MX, SRV, and URI records
	Weight   uint // SRV and URI records
}

// Zone is a generalized representation of a DNS zone.
type Zone struct {
	Name string
}

// ToSRV parses the record into a SRV struct with fully-parsed, literal values.
//
// EXPERIMENTAL; subject to change or removal.
func (r Record) ToSRV() (SRV, error) {
	if r.Type != "SRV" {
		return SRV{}, fmt.Errorf("record type not SRV: %s", r.Type)
	}

	fields := strings.Fields(r.Value)
	if len(fields) != 2 {
		return SRV{}, fmt.Errorf("malformed SRV value; expected: '<port> <target>'")
	}

	port, err := strconv.Atoi(fields[0])
	if err != nil {
		return SRV{}, fmt.Errorf("invalid port %s: %v", fields[0], err)
	}
	if port < 0 {
		return SRV{}, fmt.Errorf("port cannot be < 0: %d", port)
	}

	parts := strings.SplitN(r.Name, ".", 3)
	if len(parts) < 3 {
		return SRV{}, fmt.Errorf("name %v does not contain enough fields; expected format: '_service._proto.name'", r.Name)
	}

	return SRV{
		Service:  strings.TrimPrefix(parts[0], "_"),
		Proto:    strings.TrimPrefix(parts[1], "_"),
		Name:     parts[2],
		Priority: r.Priority,
		Weight:   r.Weight,
		Port:     uint(port),
		Target:   fields[1],
	}, nil
}

// SRV contains all the parsed data of an SRV record.
//
// EXPERIMENTAL; subject to change or removal.
type SRV struct {
	Service  string // no leading "_"
	Proto    string // no leading "_"
	Name     string
	Priority uint
	Weight   uint
	Port     uint
	Target   string
}

// ToRecord converts the parsed SRV data to a Record struct.
//
// EXPERIMENTAL; subject to change or removal.
func (s SRV) ToRecord() Record {
	return Record{
		Type:     "SRV",
		Name:     fmt.Sprintf("_%s._%s.%s", s.Service, s.Proto, s.Name),
		Priority: s.Priority,
		Weight:   s.Weight,
		Value:    fmt.Sprintf("%d %s", s.Port, s.Target),
	}
}

// RelativeName makes fqdn relative to zone. For example, for a FQDN of
// "sub.example.com" and a zone of "example.com", it outputs "sub".
//
// If fqdn cannot be expressed relative to zone, the input fqdn is returned.
func RelativeName(fqdn, zone string) string {
	// liberally ignore trailing dots on both fqdn and zone, because
	// the relative name won't have a trailing dot anyway; I assume
	// this won't be problematic...?
	// (initially implemented because Cloudflare returns "fully-
	// qualified" domains in their records without a trailing dot,
	// but the input zone typically has a trailing dot)
	return strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(fqdn, "."), strings.TrimSuffix(zone, ".")), ".")
}

// AbsoluteName makes name into a fully-qualified domain name (FQDN) by
// prepending it to zone and tidying up the dots. For example, an input
// of name "sub" and zone "example.com." will return "sub.example.com.".
func AbsoluteName(name, zone string) string {
	if zone == "" {
		return strings.Trim(name, ".")
	}
	if name == "" || name == "@" {
		return zone
	}
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	return name + zone
}
