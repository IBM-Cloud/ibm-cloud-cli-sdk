package models

// ByLastIndex sorts PaginationURLs by LastIndex
type ByLastIndex []PaginationURL

func (a ByLastIndex) Len() int { return len(a) }

func (a ByLastIndex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByLastIndex) Less(i, j int) bool { return a[i].LastIndex < a[j].LastIndex }

type PaginationURL struct {
	LastIndex int
	NextURL   string
}
