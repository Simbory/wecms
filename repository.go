package wecms

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	name          string
	conn          string
	dbName        string
	tempCacheSize int
	itemCacheSize int
	session       *mgo.Session
	templateCache repCache
	itemCache     repCache
}

// getSession clone a new mgo session from the main session
func (rep *Repository) getSession() *mgo.Session {
	if rep.session == nil {
		s,err := mgo.Dial(rep.conn)
		if err != nil {
			return nil
		}
		rep.session = s
	}
	return rep.session.Clone()
}

// getTemplate get template from database by template ID
func (rep *Repository) getTemplate(id ID) (*Template, error) {
	session := rep.getSession()
	if session == nil {
		return nil, errors.New("the data session of this repository is nil")
	}
	defer session.Close()

	coll := session.DB(rep.dbName).C("templates")
	var t *Template
	err := coll.Find(bson.M{"_id": id, "type":"Template"}).One(t)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

// GetTemplate get the template from repository cache. If the template cannot be found in cache, try to get the template
// from the database.
func (rep *Repository) GetTemplate(templateId ID) *Template {
	if rep.templateCache == nil {
		rep.templateCache = make(repCache, rep.tempCacheSize)
	}

	if temp, ok := rep.templateCache[templateId]; ok {
		return temp.(*Template)
	} else {
		temp, err := rep.getTemplate(templateId)
		if err != nil {
			return nil
		} else {
			rep.templateCache[templateId] = temp
		}
		return temp
	}
}

// getItem get item from database by item ID
func (rep *Repository) getItem(id ID) (*Item, error) {
	session := rep.getSession()
	if session == nil {
		return nil, errors.New("the data session of this repository is nil")
	}
	defer session.Close()

	coll := session.DB(rep.dbName).C("items")
	var item *Item
	err := coll.FindId(id).One(item)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (rep *Repository) GetItem(itemId ID) (*Item,error) {
	if rep.itemCache == nil {
		rep.itemCache = make(repCache, rep.itemCacheSize)
	}
	if item,ok := rep.itemCache[itemId]; ok {
		return item.(*Item), nil
	} else {
		item,err := rep.getItem(itemId)
		if err != nil {
			return nil,err
		} else if item != nil {
			item.currentRep = rep
			rep.itemCache[itemId] = item
		}
		return item, nil
	}
}

func (rep *Repository) ClearCache() {
	if rep.itemCache != nil {
		rep.itemCache.clear()
	}
	if rep.templateCache != nil {
		rep.templateCache.clear()
	}
}

func (rep *Repository) getChildItems(parentId ID) ([]*Item, error) {
	session := rep.getSession()
	if session == nil {
		return nil, errors.New("the data session of this repository is nil")
	}
	defer session.Close()

	coll := session.DB(rep.dbName).C("items")
	items := []Item{}
	err := coll.Find(bson.M{"parentid": parentId}).All(&items)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		var results []*Item
		for _, item := range items {
			if len(item.Id) == 0 {
				continue
			}
			newItem,err := rep.GetItem(item.Id)
			if newItem == nil || err != nil {
				continue
			}
			results = append(results, newItem)
		}
		return results,nil
	} else {
		return nil, nil
	}
}

func (rep *Repository) Editing(user *User) *RepEditing {
	if user.CanDev() {
		return &RepEditing{rep, user.UserName}
	}
	return nil
}

var reps map[string]*Repository

// RegRepository register a new data repository to the repository list
func RegRepository(newRep *Repository) error {
	assertNotNil(newRep, "newRep")
	if len(newRep.name) == 0 || newRep.session == nil {
		return errors.New("Invalid repository: the name is empty or the session is nil")
	}
	if reps == nil {
		reps = make(map[string]*Repository, 3)
	}
	if _,ok := reps[newRep.name]; ok {
		return fmt.Errorf("Duplicated repository name '%s': the repository '%s' is already exist.", newRep.name, newRep.name)
	} else {
		reps[newRep.name] = newRep
	}
	return nil
}

// GetRepository get the data repository by name
func GetRepository(name string) *Repository {
	if len(reps) == 0 {
		return nil
	}
	if rep,ok := reps[name]; ok {
		return rep
	} else {
		return nil
	}
}

// NewRepository Create a new repository
func NewRepository(name, conn, dbName string, tempCacheSize, itemCacheSize int) (*Repository, error) {
	session, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	rep := &Repository{
		name:          name,
		conn:          conn,
		dbName:        dbName,
		session:       session,
		tempCacheSize: tempCacheSize,
		itemCacheSize: itemCacheSize,
		templateCache: make(repCache, tempCacheSize),
		itemCache:     make(repCache, itemCacheSize),
	}
	return rep, nil
}