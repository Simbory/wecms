package wecms

type Section interface {
	Name() string
	SortOrder() int
	Owner() *Template
	Template() *Template
}