package wecms

type Property interface {
	Name() string
	DisplayTitle() string
	Type() string
	Mandatory() bool
	Validation() string
	DefaultValue() string
	Template() *Template
	Section() Section
	Value(item *Item) string
}

type tempProperty struct {
	PName         string `bson:"name" json:"name" xml:"name" field:"name"`
	PDisplayTitle string `bson:"display_title" json:"display_title" xml:"display_title" field:"display_title"`
	PType         string `bson:"type" json:"type" xml:"type" field:"type"`
	PMandatory    bool   `bson:"mandatory" json:"mandatory" xml:"mandatory" field:"mandatory"`
	PValidation   string `bson:"validation" json:"validation" xml:"validation" field:"validation"`
	PDefaultValue string `bson:"default_value" json:"default_value" xml:"default_value" field:"default_value"`
	PSection      string `bson:"section" json:"section" xml:"section" field:"section"`

	template *Template
}

func (prop *tempProperty) Name() string {
	return prop.PName
}

func (prop *tempProperty) DisplayTitle() string {
	return prop.PDisplayTitle
}

func (prop *tempProperty) Type() string {
	return prop.PType
}

func (prop *tempProperty) Mandatory() bool {
	return prop.PMandatory
}

func (prop *tempProperty) Validation() string {
	return prop.PValidation
}

func (prop *tempProperty) DefaultValue() string {
	return prop.PDefaultValue
}

func (prop *tempProperty) Template() *Template {
	return prop.template
}

/*
func (prop *tempProperty) Section() Section {
	return prop.Template().GetSection(prop.PSection)
}
*/

func (prop *tempProperty) Value(item *Item) string {
	if item != nil && prop.Template().Id.Eq(item.Template().Id) {
		if len(item.Values) > 0 {
			for _, value := range item.Values {
				if value.FieldName == prop.Name() {
					return value.Value
				}
			}
		}
		return prop.DefaultValue()
	} else {
		return ""
	}
}
