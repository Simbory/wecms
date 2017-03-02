package controllers

import (
	"github.com/Simbory/wemvc"
	"github.com/Simbory/wecms"
	"strings"
)

type TemplateController struct {
	wemvc.Controller
}

func getUser() *wecms.User {
	return &wecms.User{
		Id: wecms.NewID(),
		Email: "simbory@sina.cn",
		UserName: "Simbory.Lu",
		Password: "",
		FullName: "Simbory Lu",
		FirstName: "Simbory",
		LastName: "Lu",
		Roles: []wecms.RoleType{wecms.Administrator},
	}
}

func (ctrl TemplateController) GetChildren() interface{} {
	rep := wecms.GetRepository("master")
	if rep == nil {
		return ctrl.NotFound()
	}
	var id = wecms.ID(strings.ToLower(ctrl.RouteData()["id"]))
	children,err := rep.Editing(getUser()).ChildTemplateEntries(id)
	if err != nil {
		return ctrl.NotFound()
	}
	ctrl.ViewData["Children"] = children
	return ctrl.View()
}