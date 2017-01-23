package wecms

type ItemValue struct {
	FieldName string
	Value     string
}

type Item struct {
	Id         ID `bson:"__id"`
	Name       string
	TemplateId ID
	ParentId   ID
	Values     []ItemValue

	currentRep *Repository
}

func (item *Item) CurrentRepository() *Repository {
	return item.currentRep
}