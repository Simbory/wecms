package wecms

import "time"

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
	Sections []*section
	Properties []*property

	rep       *Repository
	ancestors []*Template
	bases     []*Template
}

func (t *Template) CurrentRepository() *Repository {
	return t.rep
}

func (t *Template) AncestorTemplates() []*Template {
	if t.ancestors != nil {
		return t.ancestors
	}
	if len(t.Bases) == 0 {
		t.ancestors = []*Template{}
	} else {
		var bases = map[ID]*Template{}
		for _, baseId := range t.Bases {
			// check and ensure the base template ID is not the same as current template id
			if baseId.Eq(t.Id) {
				continue
			}
			// check if the base template ID is already exist
			if _,ok := bases[baseId]; ok {
				continue
			}
			// Get base template by id
			baseTemp,err := t.CurrentRepository().getTemplate(baseId)
			if err != nil {
				continue
			}
			basesOfBase := baseTemp.BaseTemplates()
			var circularDependency = false
			if len(basesOfBase) > 0 {
				for _, baseBase := range basesOfBase {
					// check circular template dependency
					if baseBase.Id.Eq(t.Id) {
						circularDependency = true
						break
					}
					if _,ok := bases[baseBase.Id]; !ok {
						bases[baseBase.Id] =  baseBase
					}
				}
			}
			if !circularDependency {
				bases[baseTemp.Id] = baseTemp
			}
		}
		if len(bases) > 0 {
			var arr = make([]*Template, len(bases))
			for _, base := range bases {
				arr = append(arr, base)
			}
			t.ancestors = arr
		} else {
			t.ancestors = []*Template{}
		}
	}
	return t.ancestors
}

func (t *Template) BaseTemplates() []*Template {
	if t.bases != nil {
		return t.bases
	}
	ancestors := t.AncestorTemplates()
	if len(ancestors) == 0 {
		t.bases = []*Template{}
	} else {
		bases := []*Template{}
		for _, ancestor := range ancestors {
			for _, base := range t.Bases {
				if ancestor.Id.Eq(base) {
					bases = append(bases, ancestor)
				}
			}
		}
		t.bases = bases
	}
	return t.bases
}

func (t *Template) GetSection(name string) Section {
	if len(t.Sections) > 0 {
		for _, section := range t.Sections {
			if section.Name() == name {
				section.template = t
				return section
			}
		}
	}
	ancestors := t.AncestorTemplates()
	if len(ancestors) > 0 {
		for _, ancestor := range ancestors {
			if len(ancestor.Sections) > 0 {
				for _, section := range ancestor.Sections {
					if section.Name() == name {
						section.template = t
						return section
					}
				}
			}
		}
	}
	return nil
}

func (t *Template) AllSections() map[string]Section {
	var sections = map[string]Section{}
	if len(t.Sections) > 0 {
		for _, section := range t.Sections {
			if _,ok := sections[section.Name()]; !ok {
				section.template = t
				sections[section.Name()] = section
			}
		}
	}
	ancestors := t.AncestorTemplates()
	if len(ancestors) > 0 {
		for _, ancestor := range ancestors {
			if len(ancestor.Sections) > 0 {
				for _, section := range ancestor.Sections {
					if _,ok := sections[section.Name()]; !ok {
						section.template = t
						sections[section.Name()] = section
					}
				}
			}
		}
	}
	return sections
}

func (t *Template) GetProperty(name string) Property {
	if len(t.Properties) > 0 {
		for _, property := range t.Properties {
			if property.Name() == name {
				return property
			}
		}
	}
	ancestors := t.AncestorTemplates()
	if len(ancestors) > 0 {
		for _, ancestor := range ancestors {
			if len(ancestor.Properties) > 0 {
				for _, prop := range ancestor.Properties {
					if prop.Name() == name {
						return prop
					}
				}
			}
		}
	}
	return nil
}

func (t *Template) AllProperties() map[string]Property {
	var properties = map[string]Property{}
	if len(t.Properties) > 0 {
		for _, prop := range t.Properties {
			if _,ok := properties[prop.Name()]; !ok {
				prop.template = t
				properties[prop.Name()] = prop
			}
		}
	}
	ancestors := t.AncestorTemplates()
	if len(ancestors) > 0 {
		for _, ancestor := range ancestors {
			if len(ancestor.Properties) > 0 {
				for _, section := range ancestor.Properties {
					if _,ok := properties[section.Name()]; !ok {
						section.template = t
						properties[section.Name()] = section
					}
				}
			}
		}
	}
	return properties
}