package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"database-backup-utility/internal/logger"
)

func main() {
	if len(os.Args) < 3 {
		showHelp()
		return
	}

	dbType := os.Args[1]
	command := os.Args[2]

	switch dbType {
	case "mongodb":
		handleMongoDB(command)

	case "redis":
		handleRedis(command)

	default:
		fmt.Println("Unsupported database type:", dbType)
		logger.WriteLog("Unsupported database type: " + dbType)
		showHelp()
	}
}

func handleMongoDB(command string) {
	switch command {
	case "backup":
		backupMongoDB()

	case "restore":
		restoreMongoDB()

	default:
		fmt.Println("Unknown MongoDB command:", command)
	}
}

func handleRedis(command string) {
	switch command {
	case "backup":
		backupRedis()

	case "restore":
		restoreRedis()

	default:
		fmt.Println("Unknown Redis command:", command)
	}
}

func backupMongoDB() {
	fmt.Println("Starting MongoDB backup...")
	logger.WriteLog("Starting MongoDB backup")

	checkMongoDumpVersion()

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := "backups/mongodb_" + timestamp

	err := os.MkdirAll("backups", 0755)
	if err != nil {
		fmt.Println("Failed to create backups folder:", err)
		logger.WriteLog("MongoDB backup failed:" + err.Error())
		return
	}

	cmd := exec.Command("mongodump", "--out", backupPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("MongoDB backup failed:", err)
		fmt.Println(string(output))
		logger.WriteLog("MongoDB backup failed: " + err.Error())
		return
	}

	fmt.Println("MongoDB backup completed successfully")
	fmt.Println("Backup saved to:", backupPath)
	logger.WriteLog("MongoDB backup completed successfully: " + backupPath)

}

func restoreMongoDB() {
	fmt.Println("Starting MongoDB restore...")
	logger.WriteLog("Starting MongoDB restore")

	checkMongoDumpVersion()

	fmt.Println("MongoDB restore process will be implemented here.")
	logger.WriteLog("MongoDB restore process executed")
}

func backupRedis() {
	fmt.Println("Starting Redis backup...")
	logger.WriteLog("Starting Redis backup")

	fmt.Println("Redis backup process will be implemented here.")
	logger.WriteLog("Redis backup process executed")
}

func restoreRedis() {
	fmt.Println("Starting Redis restore...")
	logger.WriteLog("Starting Redis restore")

	fmt.Println("Redis restore process will be implemented here.")
	logger.WriteLog("Redis restore process executed")
}

func checkMongoDumpVersion() {
	fmt.Println("Checking mongodump version")

	cmd := exec.Command("mongodump", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("mongodump is not available:", err)
		logger.WriteLog("mongodump check failed: " + err.Error())
		return
	}

	fmt.Println(string(output))
	logger.WriteLog("mongodump is available")
}

func showHelp() {
	fmt.Println("Database Backup Utility")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go mongodb backup")
	fmt.Println("  go run main.go mongodb restore")
	fmt.Println("  go run main.go redis backup")
	fmt.Println("  go run main.go redis restore")
}
