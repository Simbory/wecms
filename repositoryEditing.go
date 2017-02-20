package wecms

import (
	"errors"
	"time"
)

type RepositoryEditing struct {
	currentRep  *Repository
	currentUser string
}

func (editing *RepositoryEditing) saveTemplate(t *Template) error {
	if t == nil {
		return errParamNil("t")
	}
	if len(t.Id) < 1 {
		t.Id = NewID()
	}
	if len(t.Name) == 0 {
		return errors.New("The name of the template cannot be empty.")
	}
	t.Type = "Template"
	t.UpdateTime = time.Now()
	t.UpdatedBy = editing.currentUser
	if len(t.Container) == 0 {
		t.Container = RootID
	}
	session := editing.currentRep.getSession()
	if session == nil {
		return errSessionNil(editing.currentRep.dbName)
	}
	defer session.Close()

	db := session.DB(editing.currentRep.dbName)
	coll := db.C("templates")

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

func (editing *RepositoryEditing) SaveTemplate(t *Template) error {
	err := editing.saveTemplate(t)
	if err != nil {
		return err
	}
	editing.currentRep.templateCache[t.Id] = t
	return nil
}
