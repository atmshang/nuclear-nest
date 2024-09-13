//go:build android
// +build android

package datautil

import (
	"fmt"
	"os"
)

// GetRelDataPath 函数用于获取 LincOS 下的 data 目录路径
func GetRelDataPath() string {
	// 定义 data 目录的路径
	dataPath := fmt.Sprintf("/data/.LincOS/data/%s/", appName)

	// 检查 data 目录是否存在
	_, err := os.Stat(dataPath)
	if err != nil {
		// 如果 data 目录不存在，尝试创建它
		err = os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating data directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 如果 data 目录已存在，检查它是否是一个目录
		fileInfo, err := os.Stat(dataPath)
		if err != nil {
			fmt.Printf("Error checking data directory: %v\n", err)
			os.Exit(1)
		}
		if !fileInfo.IsDir() {
			fmt.Printf("Error: %s is not a directory\n", dataPath)
			os.Exit(1)
		}
	}

	// 返回 data 目录的路径
	return dataPath
}
