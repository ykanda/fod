package fod

// Filter : filter interface
type Filter interface {
	// filter an entry
	filter([]*Entry) []*Entry
}
