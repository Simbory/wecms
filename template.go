package wecms

import "time"

type TemplateField struct {
	Id              ID
	Name            string
	DisplayTitle    string
	FieldType       string
	Mandatory       bool
	ValidationRegex string
	DefaultValue    string

	section *TemplateSection
}

type TemplateSection struct {
	Name   string
	Fields []*TemplateField
}

func (s *TemplateSection) GetField(name string) *TemplateField {
	if len(s.Fields) == 0 {
		return nil
	}
	for _, f := range s.Fields {
		if f.Name == name {
			f.section = s
			return f
		}
	}
	return nil
}

type Template struct {
	Id         ID `bson:"_id"`
	Name       string
	Type       string
	Container  ID
	CreateTime time.Time
	UpdateTime time.Time
	CreatedBy  string
	UpdatedBy  string
	Bases      []ID
	Sections   []*TemplateSection
}

func (t *Template) GetSection(name string) *TemplateSection {
	if len(t.Sections) == 0 {
		return nil
	}
	for _, section := range t.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (t *Template) GetField(name string) *TemplateField {
	if len(t.Sections) == 0 {
		return nil
	}
	for _, section := range t.Sections {
		if len(section.Fields) == 0 {
			continue
		}
		if f := section.GetField(name); f != nil {
			return f
		}
	}
	return nil
}