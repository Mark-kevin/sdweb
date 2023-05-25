package console

import (
	"github.com/astaxie/beego"
)

func WebInit() {
	BeegoConf()
	control := &MainController{}
	beego.Router("/sdweb/", control, "get:Index")
	beego.Router("/sdweb/sys", control, "get:SystemInfo")
	beego.Router("/sdweb/lora", control, "get:LoraInfo")
	beego.Router("/sdweb/sd", control, "get:SdBaseInfo")
	beego.Router("/sdweb/del", control, "get:RemoveInfo")

	//操作用

	beego.Router("/sdweb/1", control, "get:InsertData")       //开发 导入数据用
	beego.Router("/sdweb/remove", control, "get:RemoveModel") //
	//beego.Router("/upload", control, "post:UploadModel") // 上传大文件
	//beego.Router("/upload", control, "post:UploadFile") //上传大文件-新
	beego.Router("/sdweb/upload", control, "post:UploadTmp") //上传大文件-临时

	//websocket
	beego.Router("/sdweb/restart", control, "get:Restart")
	beego.Router("/sdweb/cmd_logs", control, "get:GetCmdLogs")
}

func BeegoConf() {
	beego.BConfig.WebConfig.ViewsPath = "templates"
	// 设置上传文件大小限制为 10GB
	beego.BConfig.MaxMemory = 10 * 1024 * 1024 * 1024 // 10GB
	//注册函数
	beego.AddFuncMap("add", add) // 注册 add 函数
}

func add(x, y int) int {
	return x + y
}
