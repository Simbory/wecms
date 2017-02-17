package wecms

import "time"

type TemplateEntry struct {
	Id         ID `bson:"_id"`
	Name       string
	Type       string
	Container  ID
	CreateTime time.Time
	UpdateTime time.Time
	CreatedBy  string
	UpdatedBy  string
}