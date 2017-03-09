package wecms

type Section interface {
	Name() string
	SortOrder() int
	Template() *Template
}

type section struct {
	PName string `bson:"name" json:"name" xml:"name" field:"name"`
	POrder int `bson:"order" json:"order" xml:"order" field:"order"`

	template *Template
}

func (s *section) Name() string {
	return s.PName
}

func (s *section) SortOrder() int {
	return s.POrder
}

func (s *section) Template() *Template {
	return s.template
}
