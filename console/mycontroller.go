package console

import (
	"encoding/csv"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"kevin/sdweb/core"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type MainController struct {
	beego.Controller
	Logger *logrus.Logger
}

func (c *MainController) Prepare() {
	// 初始化 Logrus 日志对象
	c.Logger = core.LogInfoInit()
}

// Index /* 欢迎页 */
func (c *MainController) Index() {
	c.TplName = "index.html"
}

func (c *MainController) InsertData() {
	core.InsertData()
	c.Data["Success"] = "数据插入成功"
	c.TplName = "success.html"
}

// Restart /* 重启脚本 */
func (c *MainController) Restart0() {

	cmd := exec.Command("/bin/bash", os.Getenv("bashPath"))
	outList, err := core.RunCmd(cmd, c.Logger)
	if err != nil {
		c.Logger.Error(err)
	}
	c.Logger.Info(outList)
	c.Data["outList"] = outList
	// 渲染输出结果
	c.TplName = "restart.html"
}

var upgrader = websocket.Upgrader{}
var wsBroadcast = make(chan string)

func (c *MainController) WsHandler() {
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		c.Logger.Errorf("%s: %v", "升级到 WebSocket 连接失败", err)
		return
	}

	// 监听新消息并推送到客户端
	go func() {
		for msg := range wsBroadcast {

			err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				c.Logger.Errorf("向客户端推送新消息失败：%v", err)
				break
			}
		}
		ws.Close()
	}()

	//用于监听客户端的消息 暂时没用
	//go func() {
	//	for {
	//		_, msg, err := ws.ReadMessage()
	//		if err != nil {
	//			c.Logger.Errorf("接收客户端回复的消息失败：%v", err)
	//			break
	//		}
	//		// 收到客户端消息后，可以在这里处理
	//		c.Logger.Info(msg)
	//	}
	//}()

}

// Restart /* 重启脚本 */
func (c *MainController) Restart() {
	// 创建通道
	wsBroadcast = make(chan string)
	outputCh := make(chan string)

	// 启动子协程监听命令输出
	go func() {
		cmd := exec.Command("/bin/bash", os.Getenv("bashPath"))
		err := core.RunCmdCh(cmd, c.Logger, outputCh)
		handleError(c.Logger, err, "命令执行失败")
		close(outputCh) // 关闭通道，表示命令已经执行完毕
	}()

	// 启动 WebSocket 服务器并监听通道消息
	go c.WsHandler()

	// 推送输出结果到 WebSocket 服务器
	go func() {
		for output := range outputCh {
			wsBroadcast <- output
		}
		close(wsBroadcast)
	}()
	//wsIp := "ws://" + os.Getenv("appPort") + "/sdweb/ws"
	//c.Logger.Info("wsip: " + wsIp)
	//c.Data["WsIp"] = wsIp
	c.TplName = "restart.html"
}

// SystemInfo /* 系统情况页,每次点击刷新 */
func (c *MainController) SystemInfo() {
	sys := core.GetSystemInfo()
	c.Data["Disk"] = sys.DiskStorage
	c.Data["Memory"] = sys.MemorySystem
	c.Data["Cpu"] = sys.CpuSystem
	c.TplName = "sys.html"
}

// LoraInfo /* Lora页 */
func (c *MainController) LoraInfo() {
	c.Data["Type"] = "Lora模型页"
	files := core.GetFiles("lora", c.Logger)
	c.Data["Files"] = files
	c.TplName = "table.html"
}

// SdBaseInfo /* SD页 */
func (c *MainController) SdBaseInfo() {
	c.Data["Type"] = "SD模型页"
	files := core.GetFiles("sd", c.Logger)
	c.Data["Files"] = files
	c.TplName = "table.html"
}

// RemoveInfo /* 待删除页 */
func (c *MainController) RemoveInfo() {
	c.Data["Type"] = "待删除模型页"
	files := core.GetFiles("del", c.Logger)
	c.Data["Files"] = files
	c.TplName = "table.html"
}

// AddModel /* 增加模型 */
func (c *MainController) AddModel() {
	c.TplName = "index.html"
}

// RemoveModel /* 移除模型但不删 */
func (c *MainController) RemoveModel() {
	id := c.GetString("id")
	core.RemoveFileById(id, c.Logger)
	c.Data["id"] = id
	c.TplName = "index.html"
}

// BackModel /* 移回模型 */
func (c *MainController) BackModel() {
	c.Data["file_name"] = c.GetString("file_name")
	c.TplName = "index.html"
}

// DeleteModel /* 删除模型-真删 */
func (c *MainController) DeleteModel() {
	c.Data["file_name"] = c.GetString("file_name")

	c.TplName = "index.html"
}

func (c *MainController) UploadTmp() {
	fileType := c.GetString("type")
	var uploadPath string
	if fileType == "lora" {
		uploadPath = os.Getenv("loraPath")
	} else {
		uploadPath = os.Getenv("sdPath")
	}

	// 获取上传文件数据
	file, header, err := c.GetFile("file")
	fileName := header.Filename // 获取上传文件的文件名
	fileSize := header.Size     // 获取上传文件的文件大小，单位为字节
	if err != nil {
		// 上传失败
		c.Logger.Error("上传失败:", err)
		return
	}
	defer file.Close()

	// 保存文件
	filePath := path.Join(uploadPath, fileName)
	err = c.SaveToFile("file", filePath)
	if err != nil {
		// 上传失败
		c.Logger.Error("上传失败:", fileName, err)
		return
	}

	//更新文件到data.csv
	c.Data["Success"] = AddOneData(fileName, fileType, fileSize, c.Logger)
	c.TplName = "success.html" // 上传成功

}

func (c *MainController) UploadFile() {
	// 获取表单数据
	//fileType := c.GetString("type")
	f, h, err := c.GetFile("chunk")
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Ctx.Output.Body([]byte("获取上传文件失败"))
		return
	}
	defer f.Close()

	chunkIndex, _ := strconv.Atoi(c.GetString("chunkIndex"))
	totalChunks, _ := strconv.Atoi(c.GetString("totalChunks"))
	fileName := h.Filename

	tmpPath := path.Join(os.Getenv("tmpPath"), fileName)

	// 创建目录
	err = os.MkdirAll(tmpPath+fileName, os.ModePerm)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Ctx.Output.Body([]byte("创建目录失败"))
		return
	}

	// 保存分片文件
	savePath := path.Join(tmpPath, fileName, strconv.Itoa(chunkIndex))
	err = c.SaveToFile("chunk", savePath)
	if err != nil {
		c.Ctx.Output.Status = 500
		c.Ctx.Output.Body([]byte("保存分片文件失败"))
		return
	}

	// 检查是否已上传所有文件分片
	if chunkIndex == totalChunks-1 {
		err = mergeFile(fileName, totalChunks, tmpPath, os.Getenv("uploadPath"))
		if err != nil {
			c.Ctx.Output.Status = 500
			c.Ctx.Output.Body([]byte("合并文件失败"))
			return
		}
	}

	//更新文件到data.csv
	//AddOneData(fileName, fileType, h.Size)

	// 返回响应
	c.Ctx.Output.Status = 200
	c.Ctx.Output.Body([]byte("上传文件成功"))
}

// 合并文件
func mergeFile(fileName string, totalChunks int, tmpPath string, uploadPath string) error {
	// 创建目标文件
	f, err := os.Create(path.Join(uploadPath, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	// 逐一读取分片文件内容并写入目标文件
	for i := 0; i < totalChunks; i++ {
		chunkPath := "./uploads/" + fileName + "/" + strconv.Itoa(i)
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		_, err = io.Copy(f, chunkFile)
		if err != nil {
			return err
		}
		err = os.Remove(chunkPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddOneData(fileName string, fileType string, fileSize int64, log *logrus.Logger) string {
	dataPath := "./conf/data.csv"
	size := float64(fileSize / 1024 / 1024) // 转换为 MB
	// 打开 data.csv 文件
	dataFile, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("failed to open data.csv file: %v", err)
		return "上传失败,读取data.csv文件路径失败"
	}
	defer dataFile.Close()

	// 读取 data.csv 文件以获取最新的序号
	data, err := os.Open(dataPath)
	if err != nil {
		log.Errorf("failed to open data.csv file: %v", err)
		return "上传失败,打开data.csv文件失败"
	}
	defer data.Close()
	reader := csv.NewReader(data)
	records, err := reader.ReadAll()
	if err != nil {
		log.Errorf("failed to read data.csv file: %v", err)
		return "上传失败,读取data.csv数据失败"
	}
	var id int
	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		id, _ = strconv.Atoi(lastRecord[0])
	}
	id++ // 增加序号以创建新记录

	//写入文件
	err = core.AddFileInfo(id, fileName, fileType, size, dataFile, log)
	if err != nil {
		log.Errorf("failed to write data.csv file: %v", err)
		return "上传失败,写入data.csv数据失败"
	}
	return "上传成功"
}

// 在服务器端，我们需要接收并拼接所有块，形成原始文件。示例代码如下： -老
func (c *MainController) Upload() {
	// 获取上传文件块的索引和总块数
	blockIndex, _ := c.GetInt("blockIndex")
	blockCount, _ := c.GetInt("blockCount")

	// 获取上传文件块数据
	file, header, err := c.GetFile("file")
	if err != nil {
		// 上传失败
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer file.Close()

	// 保存上传文件块到服务器临时文件夹
	blockPath := path.Join("tmp", header.Filename+".part"+strconv.Itoa(blockIndex))
	err = c.SaveToFile("file", blockPath)
	if err != nil {
		// 上传失败
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 判断是否已上传所有块
	if blockIndex == blockCount-1 {
		// 已上传所有块，将所有块拼接成原始文件
		filePath := path.Join("upload", header.Filename)
		f, err := os.Create(filePath)
		if err != nil {
			// 上传失败
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": err.Error(),
			}
			c.ServeJSON()
			return
		}
		defer f.Close()

		// 拼接上传文件块，形成原始文件
		for i := 0; i < blockCount; i++ {
			blockPath := path.Join("tmp", header.Filename+".part"+strconv.Itoa(i))
			block, err := os.Open(blockPath)
			if err != nil {
				// 上传失败
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"message": err.Error(),
				}
				c.ServeJSON()
				return
			}
			defer block.Close()

			_, err = io.Copy(f, block)
			if err != nil {
				// 上传失败
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"message": err.Error(),
				}
				c.ServeJSON()
				return
			}

			// 删除上传文件块
			os.Remove(blockPath)
		}

		// 删除临时文件夹
		os.Remove(path.Join("tmp", header.Filename))

		// 上传成功
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": "Upload completed",
		}
		c.ServeJSON()
		return
	}

	// 上传成功，等待上传下一个块
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Block uploaded",
	}
	c.ServeJSON()
}

// UploadModel /* 上传模型-老 */
func (c *MainController) UploadModel() {

	// 获取表单数据
	fileType := c.GetString("type")
	//file, header, err := c.GetFile("file")
	//if err != nil {
	//	// 处理上传文件失败的情况
	//	c.Logger.Error(err)
	//	return
	//}
	c.Logger.Info(fileType)

	// 处理上传文件成功的情况
	// ...

	c.TplName = "upload.html"
}

// 封装错误处理函数
func handleError(log *logrus.Logger, err error, message string) {
	if err != nil {
		log.Errorf("%s: %v", message, err)
	}
}

//func (c *MainController) Get() {
//	c.Data["Website"] = "beego.me"
//	c.Data["Email"] = "astaxie@gmail.com"
//	c.TplName = "index.html"
//}
//
//func (c *MainController) Post() {
//	name := c.GetString("name")
//	if name == "" {
//		name = "旅行者"
//	}
//	c.Data["Name"] = name
//	c.TplName = "result.tpl"
//}
