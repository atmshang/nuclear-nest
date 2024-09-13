package versionutil

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type VersionInfo struct {
	ApplicationName string `json:"applicationName"`
	VersionName     string `json:"versionName"`
	VersionCode     int    `json:"versionCode"`
	ExecutableMD5   string `json:"executableMD5"`
	MD5Checksum     string `json:"md5Checksum"`
	Author          string `json:"author"`
	ReleaseDate     string `json:"releaseDate"`
	Description     string `json:"description"`
}

var versionList []VersionInfo

// SetVersionList 设置版本信息列表
func SetVersionList(versions []VersionInfo) {
	versionList = versions
}

func GetVersionInfo() VersionInfo {
	if len(versionList) == 0 {
		log.Fatal("Version list is empty. Please set version information using SetVersionList.")
	}

	// 获取当前可执行文件的路径
	execFile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// 计算可执行文件的 MD5 摘要值
	file, err := os.Open(execFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	execMD5 := fmt.Sprintf("%x", hash.Sum(nil))

	// 获取可执行文件所在目录
	execDir := filepath.Dir(execFile)

	// 构建 md5checksum 文件的完整路径
	md5ChecksumFile := filepath.Join(execDir, "md5checksum")

	// 从文件中读取 MD5 摘要值
	md5Checksum, err := os.ReadFile(md5ChecksumFile)
	if err != nil {
		md5Checksum = []byte("unknown")
	}

	// 去除 MD5 摘要值中的换行符
	md5Checksum = []byte(strings.TrimSpace(string(md5Checksum)))

	// 获取最新的版本信息对象
	latestVersion := versionList[len(versionList)-1]

	// 更新 ExecutableMD5 和 MD5Checksum 字段
	latestVersion.ExecutableMD5 = execMD5
	latestVersion.MD5Checksum = string(md5Checksum)

	return latestVersion
}

func PrintVersionInfo() {
	// 调用 GetVersionInfo 函数获取版本信息
	versionInfo := GetVersionInfo()

	// 将 VersionInfo 结构体转换为 JSON 并打印
	jsonData, err := json.Marshal(versionInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func CreateMD5File() {
	// 获取当前可执行文件的路径
	execFile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// 计算可执行文件的 MD5 摘要值
	file, err := os.Open(execFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	execMD5 := fmt.Sprintf("%x", hash.Sum(nil))

	// 获取可执行文件所在目录
	execDir := filepath.Dir(execFile)

	// 构建 md5checksum 文件的完整路径
	md5ChecksumFile := filepath.Join(execDir, "md5checksum")

	// 将 MD5 摘要值写入 md5checksum 文件
	err = os.WriteFile(md5ChecksumFile, []byte(execMD5), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateChangeLogFile() {
	if len(versionList) == 0 {
		log.Fatal("Version list is empty. Please set version information using SetVersionList.")
	}

	// 获取当前可执行文件所在目录
	execFile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(execFile)

	// 构建 CHANGELOG.md 文件的完整路径
	changeLogFile := filepath.Join(execDir, "CHANGELOG.md")

	// 创建 CHANGELOG.md 文件
	file, err := os.Create(changeLogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 生成变更日志内容
	var changeLog string
	for _, version := range versionList {
		changeLog += fmt.Sprintf("## %s\n", version.VersionName)
		changeLog += fmt.Sprintf("- 应用：%s\n\n", version.ApplicationName)
		changeLog += fmt.Sprintf("- 版本名：%s\n\n", version.VersionName)
		changeLog += fmt.Sprintf("- 版本号：%d\n\n", version.VersionCode)
		changeLog += fmt.Sprintf("- 作者：%s\n\n", version.Author)
		changeLog += fmt.Sprintf("- 发布日期：%s\n\n", version.ReleaseDate)
		changeLog += fmt.Sprintf("- 描述：%s\n\n", version.Description)
	}

	// 将变更日志内容写入文件
	_, err = file.WriteString(changeLog)
	if err != nil {
		log.Fatal(err)
	}
}
