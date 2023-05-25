package main

import (
	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"kevin/sdweb/console"
	"kevin/sdweb/core"
)

type Meta struct {
	Syslog   *logrus.Logger
	YamlConf *core.YamlConf
}

func StartInit() *Meta {
	//启用日志引擎logrus
	log := core.LogInfoInit()
	log.Info("日志引擎启用成功...")

	//启用配置引擎yaml读取配置参数
	conf := core.YamlInit()
	log.Info("conf.yaml参数读取成功...")

	//配置全局变量
	core.ConstantsInit(*conf, log)

	//启用web引擎 beego
	console.WebInit()
	beego.Run(conf.App.AppPort)
	return &Meta{Syslog: log, YamlConf: conf}
}

func main() {
	StartInit()

}
