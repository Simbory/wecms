package wecms

import (
	"github.com/Simbory/wemvc"
	wecmsCtrls "github.com/Simbory/wecms/ui/areas/wecms/controllers"
)

func init() {
	area,err := wemvc.NewArea("areas/wecms", "wecms")
	if err != nil {
		return
	}
	area.Route("/template/<action>/<id=>", wecmsCtrls.TemplateController{}, "")
}