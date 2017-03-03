package wecms

import (
	"github.com/Simbory/wemvc"
	wecmsCtrls "github.com/Simbory/wecms/ui/areas/wecms/controllers"
	"github.com/Simbory/wecms/ui/areas/wecms/filters"
)

func init() {
	area,err := wemvc.NewArea("areas/wecms", "wecms")
	if err != nil {
		return
	}
	area.SetPathFilter("/", filters.LoginFilter)
	area.Route("/account/<action>", wecmsCtrls.AccountController{}, "")
	area.Route("/template/<action>/<id=>", wecmsCtrls.TemplateController{}, "")
}