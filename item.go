package wecms

type ItemValue struct {
	FieldName string
	Value     string
}

type Item struct {
	Id         ID `bson:"_id"`
	Name       string
	TemplateId ID
	ParentId   ID
	Values     []*ItemValue

	currentRep *Repository
	children   []*Item
}

func (item *Item) tryLoadChildren() {
	if item.children != nil {
		return
	}
	childItems,err := item.currentRep.getChildItems(item.Id)
	if err != nil {
		return
	}
	if childItems == nil {
		childItems = make([]*Item, 0)
	}
	item.children = childItems
}

func (item *Item) CurrentRepository() *Repository {
	return item.currentRep
}

// Children get the child items
func (item *Item) Children() []*Item {
	item.tryLoadChildren()
	return item.children
}

// Parent get the parent item of current item
func (item *Item) Parent() *Item {
	if item.Id.Eq(RootID) {
		return nil
	}
	p,err := item.currentRep.GetItem(item.ParentId)
	if err != nil {
		return nil
	}
	return p
}

// Template get the template of current item
func (item *Item) Template() *Template {
	return item.currentRep.GetTemplate(item.TemplateId)
}

// Value get the value by field name
func (item *Item) Value(fieldName string) string {
	template := item.Template()
	if template == nil {
		return ""
	}
	field := template.GetField(fieldName)
	if field == nil {
		return ""
	}
	var value string
	var got bool
	if len(item.Values) > 0 {
		for _, valueSetting := range item.Values {
			if valueSetting.FieldName == fieldName {
				got = true
				value = valueSetting.Value
				break
			}
		}
	}
	if !got {
		return field.DefaultValue
	}
	return value
}