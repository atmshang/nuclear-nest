package flagutil

import (
	"flag"
	"fmt"
	"github.com/atmshang/nuclear-nest/pkg/versionutil"
	"os"
)

// ParseFlags 解析命令行参数并执行相应的操作
func ParseFlags() {
	// 定义版本信息的命令行参数
	versionFlag := flag.Bool("v", false, "print version information")

	// 定义写 MD5 文件的命令行参数
	writeMD5Flag := flag.Bool("m", false, "write MD5 checksum to file")

	// 定义写变更日志的命令行参数
	writeChangeLogFlag := flag.Bool("c", false, "write change log to file")

	// 解析命令行参数
	flag.Parse()

	// 如果指定了 -v 参数，则打印版本信息并退出
	if *versionFlag {
		versionutil.PrintVersionInfo()
		os.Exit(0)
	}

	// 如果指定了 -m 参数，则写 MD5 文件并退出
	if *writeMD5Flag {
		versionutil.CreateMD5File()
		fmt.Println("MD5 checksum written to file.")
		os.Exit(0)
	}

	// 如果指定了 -c 参数，则写变更日志并退出
	if *writeChangeLogFlag {
		versionutil.CreateChangeLogFile()
		fmt.Println("Change log written to file.")
		os.Exit(0)
	}
}
