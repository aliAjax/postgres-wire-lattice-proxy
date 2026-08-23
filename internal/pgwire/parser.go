package pgwire

import (
	"errors"
	"strings"
)

type QueryKind string

const (
	ReadOnly QueryKind = "read"
	Write    QueryKind = "write"
	Txn      QueryKind = "txn"
	Unknown  QueryKind = "unknown"
	Locking  QueryKind = "locking"
)

type QueryInfo struct {
	Raw, Normalized, Fingerprint string
	Kind                         QueryKind
	Pins                         bool
}

func Classify(sql string) QueryInfo {
	n := strings.TrimSpace(sql)
	low := strings.ToLower(n)
	i := QueryInfo{Raw: sql, Normalized: normalize(low), Fingerprint: fingerprint(low)}
	switch {
	case strings.HasPrefix(low, "select") || strings.HasPrefix(low, "show") || strings.HasPrefix(low, "explain"):
		i.Kind = ReadOnly
	case strings.HasPrefix(low, "insert") || strings.HasPrefix(low, "update") || strings.HasPrefix(low, "delete") || strings.HasPrefix(low, "create") || strings.HasPrefix(low, "alter") || strings.HasPrefix(low, "drop"):
		i.Kind = Write
	case strings.HasPrefix(low, "begin") || strings.HasPrefix(low, "start transaction") || strings.HasPrefix(low, "commit") || strings.HasPrefix(low, "rollback") || strings.HasPrefix(low, "savepoint"):
		i.Kind = Txn
	case strings.HasPrefix(low, "set ") || strings.HasPrefix(low, "listen") || strings.Contains(low, "temporary") || strings.Contains(low, "declare ") || strings.Contains(low, "advisory_lock"):
		i.Pins = true
		i.Kind = Unknown
	default:
		i.Kind = Unknown
	}
	if strings.Contains(low, "for update") || strings.Contains(low, "for share") || strings.Contains(low, "nextval(") || strings.Contains(low, "random(") {
		i.Kind = Locking
	}
	if strings.Contains(low, "; select") || strings.Contains(low, ";update") {
		i.Kind = Unknown
	}
	return i
}
func normalize(s string) string { p := strings.Fields(s); return strings.Join(p, " ") }
func fingerprint(s string) string {
	b := []byte(normalize(s))
	for i, c := range b {
		if c >= '0' && c <= '9' {
			b[i] = '?'
		}
	}
	return string(b)
}
func ValidateQuery(q string) error {
	if len(q) > 1<<20 {
		return errors.New("query too large")
	}
	if strings.IndexByte(q, 0) >= 0 {
		return errors.New("nul byte in query")
	}
	return nil
}
