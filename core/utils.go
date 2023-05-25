package core

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

func GetSystemInfo() SystemMeta {
	var SystemMeta SystemMeta

	return SystemMeta
}

func GetFiles(fileType string, log *logrus.Logger) []FileMeta {
	var files []FileMeta
	file, err := os.Open("./conf/data.csv")
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	// 创建一个 CSV Reader
	reader := csv.NewReader(file)
	reader.Comma = ',' // 设置分隔符

	// 读取 CSV 文件的每一行
	lines, err := reader.ReadAll()
	if err != nil {
		log.Error(err)
	}

	//非del 和 del处理不同
	if fileType != "del" {
		for _, line := range lines {
			// 取第三列和第五列的值，并判断是否符合条件
			if line[2] == fileType && line[4] == "0" {
				file := FileMeta{
					Id:    line[0],
					Name:  line[1],
					Type:  line[2],
					Size:  line[3],
					IsDel: line[4],
				}
				files = append(files, file)
			}
		}
	} else {
		for _, line := range lines {
			// 取第三列和第五列的值，并判断是否符合条件
			if line[4] == "1" {
				file := FileMeta{
					Id:    line[0],
					Name:  line[1],
					Type:  line[2],
					Size:  line[3],
					IsDel: line[4],
				}
				files = append(files, file)
			}
		}
	}
	return files
}

func RemoveFileById(id string, log *logrus.Logger) {}

func DeleteFileById(id string, log *logrus.Logger) {}

func RunCmd(cmd *exec.Cmd, log *logrus.Logger) ([]string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("无法获取标准输出：%v", err)
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	log.Info("stderr: ", stderr)
	if err != nil {
		log.Errorf("无法获取标准错误输出：%v", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("无法启动命令：%v", err)
		return nil, err
	}

	var outList []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		outList = append(outList, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		log.Errorf("命令运行出错：%v", err)
		return outList, err
	}

	return outList, nil
}

// RunCmdCh 执行cmd命令
func RunCmdCh(cmd *exec.Cmd, log *logrus.Logger, outputCh chan string) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("无法获取标准输出：%v", err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	log.Info("CMD错误监控 stderr管道指针(/u0026=$): ", stderr)
	if err != nil {
		log.Errorf("无法获取标准错误输出：%v", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("无法启动命令：%v", err)
		return err
	}

	// 读取命令输出，实时发送到通道
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		outputCh <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("读取命令输出失败：%v", err)
		return err
	}
	// 等待命令执行完毕
	err = cmd.Wait()
	if err != nil {
		log.Errorf("命令执行错误：%v", err)
		return err
	}
	return nil
}

func InsertData() error {
	sdPath := os.Getenv("sdPath")
	loraPath := os.Getenv("loraPath")
	dataPath := "./conf/data.csv"

	// 打开 data.csv 文件
	dataFile, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open data.csv file: %v", err)
	}
	defer dataFile.Close()

	// 读取 sd 目录下的文件信息，并将其写入 data.csv 文件
	err = writeFilesInfo(sdPath, "sd", dataFile)
	if err != nil {
		return fmt.Errorf("failed to write files info in sd path: %v", err)
	}

	// 读取 lora 目录下的文件信息，并将其写入 data.csv 文件
	err = writeFilesInfo(loraPath, "lora", dataFile)
	if err != nil {
		return fmt.Errorf("failed to write files info in lora path: %v", err)
	}

	return nil
}

// 写入目录 path 中的文件信息到 f 中
func AddFileInfo(id int, fileName string, fileType string, fileSize float64, f *os.File, log *logrus.Logger) error {
	// 获取目录中的文件列表
	_, err := fmt.Fprintf(f, "%d,%s,%s,%.2f MB,0\n", id, fileName, fileType, fileSize)
	if err != nil {
		log.Errorf("failed to write files info: %v", err)
		return err
	}
	return nil
}

// 写入目录 path 中的文件信息到 f 中
func writeFilesInfo(path, dir string, f *os.File) error {
	// 获取目录中的文件列表
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", path, err)
	}

	for i, file := range files {
		// 写入文件信息
		fileName := file.Name()
		fileSize := float64(file.Size()) / 1024 / 1024 // 转换为 MB
		_, err := fmt.Fprintf(f, "%d,%s,%s,%.2f MB,0\n", i+1, fileName, dir, fileSize)
		if err != nil {
			return fmt.Errorf("failed to write files info: %v", err)
		}
	}

	return nil
}

func getSDStatus() bool {
	return true
}

func main() {

}
