package logs

import (
	"log"
	"os"
	"time"
)

var logFile = "./logs/running.log"

func SetupLogger() {
	logFileLocation, _ := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(logFileLocation)
}

func MonitorFileSize(maxSize int64) {
	ticker := time.NewTicker(50 * time.Second) // 每隔50秒检查一次文件大小
	defer ticker.Stop()

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
