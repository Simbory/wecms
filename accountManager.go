package wecms

import (
	"errors"
	"crypto/md5"
	"io"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type AccountManager struct {
	currentRep *Repository
	manager    *User
}

func (accMng *AccountManager) NewUser(email, userName, pwd, fullName, firstName, lastName string, roles []RoleType) error {
	if accMng.manager == nil {
		return errors.New("Invalid manager： the manager cannot be nil")
	}
	if !accMng.manager.CanManage() {
		return errors.New("The manager has no rights to create a new user")
	}
	if len(email) == 0 {
		return errParamEmpty("email")
	}
	if len(userName) == 0 {
		return errParamEmpty("userName")
	}
	if len(pwd) == 0 {
		return errParamEmpty("pwd")
	}
	if len(fullName) == 0 {
		return errParamEmpty("fullName")
	}
	if len(firstName) == 0 {
		return errParamEmpty("firstName")
	}
	if len(lastName) == 0 {
		return errParamEmpty("lastName")
	}
	var user = &User{
		Id: NewID(),
		Email: email,
		UserName: userName,
		Password: accMng.encryptPad(pwd),
		FullName: fullName,
		FirstName: firstName,
		LastName: lastName,
		Roles: roles,
	}
	tempUser,err := accMng.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if tempUser != nil {
		return fmt.Errorf("The email is already used: %s", user.Email)
	}

	session := accMng.currentRep.getSession()
	if session == nil {
		return errSessionNil(accMng.currentRep.dbName)
	}
	defer session.Close()

	db := session.DB(accMng.currentRep.dbName)
	coll := db.C("users")
	return coll.Insert(user)
}

func (accMng *AccountManager) encryptPad(pwd string) string {
	if len(pwd) == 0 {
		return ""
	}
	md5Hash := md5.New()
	io.WriteString(md5Hash, pwd)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func (accMng *AccountManager) GetUserByEmail(email string) (*User,error) {
	session := accMng.currentRep.getSession()
	if session == nil {
		return nil, errSessionNil(accMng.currentRep.dbName)
	}
	defer session.Close()
	db := session.DB(accMng.currentRep.dbName)
	coll := db.C("users")

	var user *User
	err := coll.Find(bson.M{"email": email}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}