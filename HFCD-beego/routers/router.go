package routers

import (
	"HFCD-beego/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/static_detect", &controllers.StaticDetectContronller{}, "*:Get")
	beego.Router("/dynamic_detect", &controllers.DynamicDetectController{}, "*:Get")
	beego.Router("/mix_detect", &controllers.MixDetectController{}, "*:Get")
}
