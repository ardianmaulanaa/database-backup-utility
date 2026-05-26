package logger

import (
	"fmt"
	"os"
	"time"
)

func WriteLog(message string) {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		fmt.Println("Failed to create logs folder:", err)
		return
	}

	file, err := os.OpenFile("logs/backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = file.WriteString(logMessage)
	if err != nil {
		fmt.Println("Failed to write log", err)
		return
	}
}
