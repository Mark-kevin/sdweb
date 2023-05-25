package core

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type YamlConf struct {
	App struct {
		AppName string `yaml:"app_name"`
		AppPort string `yaml:"app_port"`
	} `yaml:"app"`
	Path struct {
		StableDiffusion string `yaml:"stable_diffusion"`
		LoraPath        string `yaml:"lora_path"`
		SdBasePath      string `yaml:"sd_base_path"`
		removePath      string `yaml:"remove_path"`
		tmpPath         string `yaml:"tmp_path"`
	} `yaml:"path"`
	Bash struct {
		Path string `yaml:"path"`
		File string `yaml:"file"`
	} `yaml:"bash"`
}

func YamlInit() *YamlConf {
	// 读取配置文件
	data, err := os.ReadFile("./conf/conf.yml")
	if err != nil {
		panic("failed to read config file")
	}

	// 解析配置文件
	var config YamlConf
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("failed to parse config file")
	}

	return &config
}

func ConstantsInit(config YamlConf, log *logrus.Logger) {
	path := config.Path.StableDiffusion
	loraPath := filepath.Join(path, config.Path.LoraPath)
	sdPath := filepath.Join(path, config.Path.SdBasePath)
	delPath := filepath.Join(path, config.Path.removePath)
	bashPath := filepath.Join(config.Bash.Path, config.Bash.File)
	tmpPath := filepath.Join(config.Path.tmpPath, config.Bash.File)
	appPort := config.App.AppPort
	// 设置全局变量
	os.Setenv("loraPath", loraPath)
	os.Setenv("sdPath", sdPath)
	os.Setenv("delPath", delPath)
	os.Setenv("bashPath", bashPath)
	os.Setenv("tmpPath", tmpPath)
	os.Setenv("appPort", appPort)
	log.Info("设置全局变量:" + "loraPath:" + loraPath + ", sdPath:" + sdPath + ", delPath:" + delPath)
	log.Info("设置全局变量:"+"bashPath:"+bashPath, ", tmpPath:"+tmpPath+", appPort:"+appPort)
}
