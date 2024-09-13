# Nuclear Nest

Nuclear Nest 是一个用 Go 编写的工具库，提供了一系列用于管理数据路径、版本信息、日志记录和命令行参数解析的功能。该库旨在为跨平台应用程序提供一致的开发体验。

![kbn](./kbn.png)

## 功能

- **数据路径管理**：根据不同平台（Android、Linux、Windows）自动设置和管理数据目录路径。
- **版本信息管理**：支持版本信息的设置和获取，并提供生成 MD5 校验文件和变更日志文件的功能。
- **日志记录**：基于 `zap` 和 `lumberjack` 实现的高效日志记录系统，支持日志文件的自动分割和管理。
- **命令行参数解析**：提供通用的命令行参数解析功能，支持打印版本信息、生成 MD5 校验文件和变更日志文件。

## 安装

确保你已经安装了 Go 语言环境，然后使用以下命令获取该库：

```bash
go get github.com/atmshang/nuclear-nest
```

## 使用

### 数据路径管理

根据不同的平台，自动设置数据目录路径。你可以通过 `SetAppName` 设置应用名称：

```go
import "github.com/atmshang/nuclear-nest/pkg/datautil"

func main() {
    datautil.SetAppName("MyApp")
    path := datautil.GetRelDataPath()
    fmt.Println("Data path:", path)
}
```

### 版本信息管理

设置版本信息并生成相关文件：

```go
import "github.com/atmshang/nuclear-nest/pkg/versionutil"

func main() {
    versions := []versionutil.VersionInfo{
        {
            ApplicationName: "MyApp",
            VersionName:     "1.0.0",
            VersionCode:     100,
            Author:          "Your Name",
            ReleaseDate:     "2024-09-13",
            Description:     "Initial release",
        },
    }

    versionutil.SetVersionList(versions)
    versionutil.PrintVersionInfo()
    versionutil.CreateMD5File()
    versionutil.CreateChangeLogFile()
}
```

### 日志记录

初始化日志记录系统并记录日志：

```go
import "github.com/atmshang/nuclear-nest/pkg/logutil"

func main() {
    logutil.Printf("Application started")
    logutil.Println("This is a log message")
}
```

### 命令行参数解析

使用 `flagutil` 解析命令行参数：

```go
import "github.com/atmshang/nuclear-nest/pkg/flagutil"

func main() {
    flagutil.ParseFlags()
    // 其他初始化代码
}
```

## 贡献

不欢迎贡献代码！但可以报告问题。

## 许可证

该项目使用 [GNU Affero General Public License (AGPL)](https://www.gnu.org/licenses/agpl-3.0.html)。有关详细信息，请参阅 LICENSE 文件。