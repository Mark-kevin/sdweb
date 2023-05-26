package core

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	cpu "github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GetSystemInfo0() SystemMeta {
	var systemMeta SystemMeta
	memorySystem := make(map[string]string)
	cpuSystem := make(map[string]string)
	diskStorage := make(map[string]string)
	cpuUsage := make(map[string]string)
	loadAverage := make(map[string]string)
	memorySystem["Info"] = "memorySystem信息"
	cpuSystem["Info"] = "cpuSystem信息"
	diskStorage["Info"] = "diskStorage信息"
	cpuUsage["Info"] = "cpuUsage信息"
	loadAverage["Info"] = "loadAverage信息"

	systemMeta.MemorySystem = memorySystem
	systemMeta.CpuSystem = cpuSystem
	systemMeta.CpuUsage = cpuUsage
	systemMeta.DiskStorage = diskStorage
	systemMeta.LoadAverage = loadAverage
	return systemMeta
}

func GetSystemInfo() SystemMeta {
	var systemMeta SystemMeta
	memorySystem := make(map[string]string)
	cpuSystem := make(map[string]string)
	diskStorage := make(map[string]string)
	cpuUsage := make(map[string]string)
	loadAverage := make(map[string]string)

	cpuSystems := make([]map[string]string, 0)
	diskStorages := make([]map[string]string, 0)
	cpuUsages := make(map[string]string)

	//内存
	v, _ := mem.VirtualMemory()
	memorySystem["Info"] = v.String()
	fmt.Println(memorySystem["Info"])
	//memorySystem["Total"] = string(v.Total)
	//memorySystem["Free"] = string(v.Free)
	//memorySystem["UsedPercent"] = fmt.Sprintf("%.2f", v.UsedPercent)
	//memorySystem["Total"] = string(v.Total)

	//cpu
	// 将 CPU 信息存储在 map 中
	cpuInfos, _ := cpu.Info()
	for _, info := range cpuInfos {
		cpuInfo := make(map[string]string)
		cpuInfo["ModelName"] = info.ModelName
		cpuInfo["Cores"] = fmt.Sprintf("%d", info.Cores)
		cpuInfo["PhysicalID"] = info.PhysicalID
		cpuInfo["CoreID"] = info.CoreID
		cpuSystems = append(cpuSystems, cpuInfo)
	}
	cpuSystemString, _ := json.Marshal(cpuSystems)
	cpuSystem["Info"] = string(cpuSystemString)
	fmt.Println(cpuSystem["Info"])

	// CPU使用率
	cpuPercent, _ := cpu.Percent(time.Second, false)
	// 将每个 CPU 的使用率存储在 map 中
	count := len(cpuPercent)
	cpuUsages["Count"] = fmt.Sprintf("%d", count)
	for i, percent := range cpuPercent {
		cpuUsages[fmt.Sprintf("CPU%d", i)] = fmt.Sprintf("%.2f%%", percent)
	}
	cpuUsageString, _ := json.Marshal(cpuUsages)
	cpuUsage["Info"] = string(cpuUsageString)
	fmt.Println(cpuUsage["Info"])

	//磁盘disk
	partitions, _ := disk.Partitions(true)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Println(err)
			continue
		}
		diskInfo := make(map[string]string)
		diskInfo["Filesystem"] = usage.Fstype
		diskInfo["Size"] = fmt.Sprintf("%.2f GB", float64(usage.Total)/(1024*1024*1024))
		diskInfo["Used"] = fmt.Sprintf("%.2f GB", float64(usage.Used)/(1024*1024*1024))
		diskInfo["Free"] = fmt.Sprintf("%.2f GB", float64(usage.Free)/(1024*1024*1024))
		diskInfo["Usage"] = fmt.Sprintf("%.2f%%", usage.UsedPercent)
		diskStorages = append(diskStorages, diskInfo)
	}
	diskStorageString, _ := json.Marshal(diskStorages)
	diskStorage["Info"] = string(diskStorageString)
	fmt.Println(diskStorage["Info"])

	//load average
	stat, _ := load.Avg()
	jsonstat, _ := json.Marshal(stat)
	loadAverage["Info"] = string(jsonstat)
	fmt.Println(diskStorage["Info"])

	systemMeta.MemorySystem = memorySystem
	systemMeta.CpuSystem = cpuSystem
	systemMeta.CpuUsage = cpuUsage
	systemMeta.DiskStorage = diskStorage
	systemMeta.LoadAverage = loadAverage
	return systemMeta
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
			// 取第四列的值，并判断是否符合条件
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

func ReadFiles(path string, log *logrus.Logger) [][]string {
	file, err := os.Open(path)
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
	return lines
}
func ReWriteFile(path string, lines [][]string, log *logrus.Logger) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	// 创建一个 CSV Writer
	writer := csv.NewWriter(file)
	writer.Comma = ',' // 设置分隔符

	// 将 CSV 数据写入文件
	err = writer.WriteAll(lines)
	if err != nil {
		log.Error("报错::", err)
		return
	}

	// 刷新缓冲区
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Error("报错::", err)
		return
	}
}

func RemoveFile(id string, fType string, log *logrus.Logger) {
	lines := ReadFiles("./conf/data.csv", log)
	var newLines [][]string
	for _, line := range lines {
		// 取第1列和第3列的值，并判断是否符合条件
		if line[0] == id && line[2] == fType {
			line[4] = "1"
		}
		newLines = append(newLines, line)
	}
	ReWriteFile("./conf/data.csv", newLines, log)
}
func BackFile(id string, fType string, log *logrus.Logger) {
	lines := ReadFiles("./conf/data.csv", log)
	var newLines [][]string
	for _, line := range lines {
		// 取第1列和第3列的值，并判断是否符合条件
		if line[0] == id && line[2] == fType {
			line[4] = "0"
		}
		newLines = append(newLines, line)
	}
	ReWriteFile("./conf/data.csv", newLines, log)
}
func DeleteFile(id string, fType string, log *logrus.Logger) {
	lines := ReadFiles("./conf/data.csv", log)
	var newLines [][]string
	for _, line := range lines {
		// 取第1列和第3列的值，并判断是否符合条件
		if line[0] == id && line[2] == fType {
			//删除文件
			fPath := ""
			if fType == "sd" {
				fPath = os.Getenv("sdPath")
			} else if fType == "lora" {
				fPath = os.Getenv("loraPath")
			} else {
				fmt.Println("删除文件类型错误", fType, id)
				continue
			}
			//删除目录的对应文件
			err := os.Remove(filepath.Join(fPath, line[1]))
			if err != nil {
				log.Error(err)
			}
			continue
		}
		newLines = append(newLines, line)
	}
	ReWriteFile("./conf/data.csv", newLines, log)
}

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
