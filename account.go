package wecms

type RoleType int

const (
	Anonymous RoleType = iota
	Editor
	Publisher
	Manager
	Developer
	Administrator
)

type User struct {
	Id        ID `bson:"_id"`
	Email     string
	FullName  string
	FirstName string
	LastName  string
	Roles     []RoleType
}

func (user *User) IsAnonymous() bool {
	if len(user.Roles) == 0 {
		return true
	}
	for _, role := range user.Roles {
		if role != Anonymous {
			return false
		}
	}
	return true
}

func (user *User) CanDev() bool {
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			if role == Administrator || role == Developer {
				return true
			}
		}
	}
	return false
}

func (user *User) CanEdit() bool {
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			if role == Administrator || role == Developer || role == Editor {
				return true
			}
		}
	}
	return false
}

func (user *User) CanManage() bool {
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			if role == Administrator || role == Manager {
				return true
			}
		}
	}
	return false
}

func (user *User) CanPublish() bool {
	if len(user.Roles) > 0 {
		for _, role := range user.Roles {
			if role == Administrator || role == Developer || role == Publisher {
				return true
			}
		}
	}
	return false
}