package wecms

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type repCache map[ID]interface{}

type Repository struct {
	name    string
	conn    string
	session *mgo.Session

	templateCache repCache
	itemCache     repCache
}

// getSession clone a new mgo session from the main session
func (rep *Repository) getSession() *mgo.Session {
	return rep.session.Clone()
}

var reps map[string]*Repository

// RegRepository register a new data repository to the repository list
func RegRepository(name string, newRep *Repository) error {
	assertNotEmpty(name, "name")
	assertNotNil(newRep, "newRep")
	if reps == nil {
		reps = make(map[string]*Repository, 3)
	}
	if _,ok := reps[name]; ok {
		return fmt.Errorf("Duplicated repository name '%s': the repository '%s' is already exist.", name, name)
	} else {
		reps[name] = newRep
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
func NewRepository(name string, conn string, tempCacheSize, itemCacheSize int) (*Repository, error) {
	session, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	rep := &Repository{
		name:          name,
		conn:          conn,
		session:       session,
		templateCache: make(repCache, tempCacheSize),
		itemCache:     make(repCache, itemCacheSize),
	}
	return rep, nil
}