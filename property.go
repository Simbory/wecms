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
}

type property struct {
	PName         string `bson:"name" json:"name" xml:"name" field:"name"`
	PDisplayTitle string `bson:"display_title" json:"display_title" xml:"display_title" field:"display_title"`
	PType         string `bson:"type" json:"type" xml:"type" field:"type"`
	PMandatory    bool   `bson:"mandatory" json:"mandatory" xml:"mandatory" field:"mandatory"`
	PValidation   string `bson:"validation" json:"validation" xml:"validation" field:"validation"`
	PDefaultValue string `bson:"default_value" json:"default_value" xml:"default_value" field:"default_value"`
	PSection      string `bson:"section" json:"section" xml:"section" field:"section"`

	template *Template
}

func (prop *property) Name() string {
	return prop.PName
}

func (prop *property) DisplayTitle() string {
	return prop.PDisplayTitle
}

func (prop *property) Type() string {
	return prop.PType
}

func (prop *property) Mandatory() bool {
	return prop.PMandatory
}

func (prop *property) Validation() string {
	return prop.PValidation
}

func (prop *property) DefaultValue() string {
	return prop.PDefaultValue
}

func (prop *property) Template() *Template {
	return prop.template
}

func (prop *property) Section() Section {
	return prop.Template().GetSection(prop.PSection)
}