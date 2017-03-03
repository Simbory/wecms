package wecms

import (
	"errors"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type RepEditing struct {
	currentRep  *Repository
	currentUser string
}

func (editing *RepEditing) saveTemplate(t *Template) error {
	if t == nil {
		return errParamNil("t")
	}
	if len(t.Name) == 0 {
		return errors.New("The name of the template cannot be empty.")
	}
	if len(t.Id) < 1 {
		t.Id = NewID()
	}
	t.Type = "Template"
	t.UpdateTime = time.Now()
	t.UpdatedBy = editing.currentUser
	if len(t.Container) == 0 {
		t.Container = RootID
	}
	session := editing.currentRep.getSession()
	if session == nil {
		return errSessionNil(editing.currentRep.name)
	}
	defer session.Close()

	coll := session.DB(editing.currentRep.dbName).C("templates")

	count,err := coll.FindId(t.Id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		err = coll.UpdateId(t.Id, t)
		if err != nil {
			return err
		}
	} else {
		t.CreatedBy = editing.currentUser
		t.CreateTime = t.UpdateTime
		err = coll.Insert(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (editing *RepEditing) SaveTemplate(t *Template) error {
	err := editing.saveTemplate(t)
	if err != nil {
		return err
	}
	editing.currentRep.templateCache[t.Id] = t
	return nil
}

func (editing *RepEditing) SaveTemplateEntry(entry *TemplateEntry) error {
	if entry == nil {
		return errParamNil("entry")
	}
	session := editing.currentRep.getSession()
	if session == nil {
		return errSessionNil(editing.currentRep.name)
	}
	defer session.Close()

	coll := session.DB(editing.currentRep.dbName).C("templates")
	if len(entry.Id) == 0 {
		entry.Id = NewID()
		err := coll.Insert(entry)
		if err != nil {
			return err
		}
		return nil
	} else {
		count,err := coll.FindId(entry.Id).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			err = coll.UpdateId(entry.Id, entry)
		} else {
			err = coll.Insert(entry)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (editing *RepEditing) GetTemplateEntry(id ID) (*TemplateEntry, error) {
	session := editing.currentRep.getSession()
	if session == nil {
		return nil, errSessionNil(editing.currentRep.name)
	}
	defer session.Close()
	coll := session.DB(editing.currentRep.dbName).C("templates")
	var entry TemplateEntry
	err := coll.FindId(id).One(&entry)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil,nil
		}
		return nil,err
	}
	return &entry,nil
}

// ChildTemplateEntries get the child template entries
func (editing *RepEditing) ChildTemplateEntries(parentId ID) ([]*TemplateEntry, error) {
	session := editing.currentRep.getSession()
	if session == nil {
		return nil, errSessionNil(editing.currentRep.name)
	}
	defer session.Close()
	coll := session.DB(editing.currentRep.dbName).C("templates")
	var entries []*TemplateEntry
	err := coll.Find(bson.M{"container": parentId}).All(&entries)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil,nil
		}
		return nil, err
	}
	return entries, nil
}

// saveItem save the item data to mongodb
func (editing *RepEditing) saveItem(item *Item) error {
	if item == nil {
		return errParamNil("item")
	}
	if len(item.Name) == 0 {
		return errors.New("the name of the item cannot be empty")
	}
	if len(item.TemplateId) == 0 {
		return errors.New("the template ID of cannot be empty")
	}
	if t := editing.currentRep.GetTemplate(item.TemplateId); t == nil {
		return fmt.Errorf("invalid template ID: %s", string(item.TemplateId))
	}
	if len(item.ParentId) == 0 {
		return errors.New("The parent ID cannot be empty")
	}
	parent,err := editing.currentRep.getItem(item.ParentId)
	if err != nil {
		return err
	}
	if parent == nil {
		return fmt.Errorf("Invalid parent ID, the parent item cannot be found: %s", string(item.ParentId))
	}
	if len(item.Id) == 0 {
		item.Id = NewID()
	}
	session := editing.currentRep.getSession()
	if session == nil {
		return errSessionNil(editing.currentRep.name)
	}
	defer session.Close()

	coll := session.DB(editing.currentRep.dbName).C("items")
	count,err := coll.FindId(item.Id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		coll.UpdateId(item.Id, item)
	} else {
		coll.Insert(item)
	}
	return nil
}

func (editing *RepEditing) SaveItem(item *Item) error {
	err := editing.saveItem(item)
	if err != nil {
		return err
	}
	item.currentRep = editing.currentRep
	if editing.currentRep.itemCache == nil {
		editing.currentRep.itemCache = make(repCache, editing.currentRep.itemCacheSize)
	}
	editing.currentRep.itemCache[item.Id] = item
	parent,_ := editing.currentRep.getItem(item.ParentId)
	parent.children = nil
	return nil
}

func (editing *RepEditing) MoveItem(item *Item, newParent ID) error {
	if item == nil {
		return errParamNil("item")
	}
	newParentItem,err := editing.currentRep.GetItem(newParent)
	if err != nil {
		return err
	}
	if newParentItem == nil {
		return fmt.Errorf("Invalid parent ID: %s", string(newParent))
	}
	oldParent,_ := editing.currentRep.GetItem(item.ParentId)
	item.ParentId = newParent
	editing.saveItem(item)
	if oldParent != nil {
		oldParent.children = nil
	}
	newParentItem.children = nil
	return nil
}