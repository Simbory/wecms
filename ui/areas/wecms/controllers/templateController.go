package controllers

import (
	"github.com/Simbory/wemvc"
	"github.com/Simbory/wecms"
	"strings"
	"time"
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

func (ctrl TemplateController) GetNew() interface{} {
	var parent = ctrl.Request().URL.Query().Get("parent")
	if len(parent) == 0 {
		parent = string(wecms.RootID)
	}
	if parent == string(wecms.RootID) {
		ctrl.ViewData["parent"] = parent
	} else {
		rep := wecms.GetRepository("master")
		if rep == nil {
			return ctrl.NotFound()
		}
		pEntry,err := rep.Editing(getUser()).GetTemplateEntry(wecms.ID(parent))
		if err != nil {
			panic(err)
		}
		if pEntry == nil {
			return ctrl.NotFound()
		}
		ctrl.ViewData["parent"] = parent
	}
	return ctrl.View()
}

func (ctrl TemplateController) PostNew() interface{} {
	var entry = wecms.TemplateEntry{
		Name: ctrl.Request().Form.Get("Name"),
		Type: ctrl.Request().Form.Get("Type"),
		Container: wecms.ID(ctrl.Request().Form.Get("Container")),
		CreatedBy: ctrl.Request().Form.Get("CreatedBy"),
		UpdatedBy: ctrl.Request().Form.Get("UpdatedBy"),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	rep := wecms.GetRepository("master")
	if rep == nil {
		return ctrl.NotFound()
	}
	err := rep.Editing(getUser()).SaveTemplateEntry(&entry)
	if err != nil {
		panic(err)
	}
	return ctrl.PlainText("OK")
}