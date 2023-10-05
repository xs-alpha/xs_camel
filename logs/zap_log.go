package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logFile = "./logs/running.log"

func SetupLogger() {
	fmt.Println("开启日志")
	// 检查 ./logs 文件夹是否存在，如果不存在则创建它
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		err := os.MkdirAll("./logs", os.ModePerm)
		if err != nil {
			log.Println("创建日志文件夹失败:", err)
			return
		}
	}
	logFileLocation, _ := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(logFileLocation)
}

func CloseLogger() {
	fmt.Println("关闭日志")
	// 将日志输出设置为 os.Stderr，这会关闭日志输出
	log.SetOutput(os.Stderr)
}

func MonitorFileSize(maxSize int64, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			log.Println("开启日志文件大小检测--")
			fileInfo, err := os.Stat(logFile)
			if err != nil {
				log.Println("Error getting file info:", err)
				return
			}

			fileSize := fileInfo.Size()
			if fileSize > maxSize {
				log.Printf("File size exceeded %d bytes. Deleting %s...\n", maxSize, logFile)

				// 删除文件
				err := os.Remove(logFile)
				if err != nil {
					log.Println("Error deleting file:", err)
					return
				}

				log.Println("File deleted.")
				return // 退出 Goroutine
			} else {
				log.Printf("当前大小：%d, 最大大小:%d, 未超过限制", fileSize, maxSize)
			}
		}
	}
}
