package fod

type Filter interface {
	// filter an entry
	Filter([]*Entry) []*Entry
}
