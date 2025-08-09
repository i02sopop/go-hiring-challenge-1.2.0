// Package filter defines the filter model to fetch data by the filter.
package filter

// OperationType defines the filter operation type for comparison (equal, less than,
// great than, ...).
type OperationType string

const (
	// Equal comparison.
	Equal OperationType = "="
	// LessThan comparison.
	LessThan OperationType = "<"
)

// Filter defiles a filter for a search, so we can rule out result that doesn't comply
// with the filter.
type Filter struct {
	Key       string
	Value     string
	Operation OperationType
}
