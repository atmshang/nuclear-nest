//go:build lincos_w
// +build lincos_w

package datautil

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetRelDataPath 函数用于获取相对于可执行文件的 data 目录路径
func GetRelDataPath() string {
	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v", err)
		os.Exit(1)
	}

	// 获取可执行文件所在目录的路径
	execDir := filepath.Dir(execPath)

	// 构建 data 目录的路径
	// 在可执行文件所在目录下的 data 目录
	dataDir := filepath.Join(execDir, "data", appName)

	// 创建 data 目录（如果不存在）
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating data directory: %v", err)
		os.Exit(1)
	}

	// 返回 data 目录的路径
	return dataDir
}
