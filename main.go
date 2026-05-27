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

	if len(os.Args) < 4 {
		fmt.Println("Please provide backup folder path.")
		fmt.Println("Example:")
		fmt.Println("  go run main.go mongodb restore backups/mongodb_2026-05-27_14-30-10")
		logger.WriteLog("MongoDB restore failed: backup path not provided")
		return
	}

	backupPath := os.Args[3]

	if !checkMongoRestoreVersion() {
		fmt.Println("MongoDB restore stopped because mongorestore is not available.")
		logger.WriteLog("MongoDB restore stopped because mongorestore is not available")
		return
	}

	_, err := os.Stat(backupPath)
	if os.IsNotExist(err) {
		fmt.Println("Backup folder not found:", backupPath)
		logger.WriteLog("MongoDB restore failed: backup folder not found: " + backupPath)
		return
	}

	cmd := exec.Command("mongorestore", backupPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("MongoDB restore failed:", err)
		fmt.Println(string(output))
		logger.WriteLog("MongoDB restore failed: " + err.Error())
		return
	}

	fmt.Println("MongoDB restore completed successfully")
	fmt.Println(string(output))
	logger.WriteLog("MongoDB restore completed successfully from: " + backupPath)
}

func backupRedis() {
	fmt.Println("Starting Redis backup...")
	logger.WriteLog("Starting Redis backup")

	if !checkRedisCLI() {
		fmt.Println("Redis backup stopped because redis-cli not available")
		logger.WriteLog("Redis backup stopped because redis-cli is not available")
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := "backup/redis_" + timestamp

	err := os.MkdirAll(backupPath, 0755)
	if err != nil {
		fmt.Println("Failed to create Redis backup folder", err)
		logger.WriteLog("Redis backup failed: " + err.Error())
		return
	}

	cmdSave := exec.Command("redis-cli", "SAVE")
	output, err := cmdSave.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to save Redis data:", err)
		fmt.Println(string(output))
		logger.WriteLog("Redis SAVE failed: " + err.Error())
		return
	}

	sourceFile := "/opt/homebrew/var/db/redis/dump.rdb"
	destinationFile := backupPath + "/dump.rdb"

	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println("Failed to read Redis dump file:", err)
		logger.WriteLog("Redis backup failed while reading dump.rdb: " + err.Error())
		return
	}

	err = os.WriteFile(destinationFile, input, 0644)

	if err != nil {
		fmt.Println("Failed to write Redis backup file:", err)
		logger.WriteLog("Redis backup failed while writing dump.rdb: " + err.Error())
		return
	}

	fmt.Println("Redis backup completed successfully")
	fmt.Println("Backup saved to:", destinationFile)
	logger.WriteLog("Redis backup completed successfully: " + destinationFile)
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

func checkRedisCLI() bool {
	fmt.Println("Checking redis-cli version")

	cmd := exec.Command("redis-cli", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("redis-cli is not available:", err)
		logger.WriteLog("redis-cli check failed: " + err.Error())
		return false
	}
	fmt.Println(string(output))
	logger.WriteLog("redis-cli is available")
	return true
}

func checkMongoRestoreVersion() bool {
	fmt.Println("Checking mongorestore version")

	cmd := exec.Command("mongorestore", "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("mongorestore is not available:", err)
		logger.WriteLog("mongorestore check failed: " + err.Error())
		return false
	}

	fmt.Println(string(output))
	logger.WriteLog("mongorestore is available")
	return true
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
