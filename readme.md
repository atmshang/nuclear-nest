# Nuclear Nest

Nuclear Nest 是一个用 Go 编写的工具库，提供了一系列用于管理数据路径、版本信息、日志记录、命令行参数解析和 API 处理的功能。该库旨在为跨平台应用程序提供一致的开发体验。

![kbn](./kbn.png)

## 功能

- **数据路径管理**：根据不同平台（Android、Linux、Windows）自动设置和管理数据目录路径。
- **版本信息管理**：支持版本信息的设置和获取，并提供生成 MD5 校验文件和变更日志文件的功能。
- **日志记录**：基于 `zap` 和 `lumberjack` 实现的高效日志记录系统，支持日志文件的自动分割和管理。
- **命令行参数解析**：提供通用的命令行参数解析功能，支持打印版本信息、生成 MD5 校验文件和变更日志文件。
- **API 处理**：提供标准化的 API 响应结构和错误处理机制。
- **认证工具**：支持模块间的内部认证，基于 RSA 和 AES 加密。

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

### API 处理

使用标准化的 API 响应结构和错误处理：

```go
import (
    "github.com/atmshang/nuclear-nest/pkg/apiutil"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    apiutil.UseErrorHandler(r)
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, apiutil.Response{
            Code:    2000,
            Message: "Success",
            Data:    "Example data",
        })
    })
    r.Run()
}
```

### 认证工具

设置公钥和私钥进行模块间认证：

```go
import "github.com/atmshang/nuclear-nest/pkg/authutil"

func main() {
    err := authutil.SetPublicKey("your-public-key-pem")
    if err != nil {
        log.Fatal(err)
    }

    err = authutil.SetPrivateKey("your-private-key-pem")
    if err != nil {
        log.Fatal(err)
    }

    // 使用认证中间件
    r := gin.Default()
    r.Use(authutil.InternalServiceAuth())
    r.Run()
}
```

## 贡献

不欢迎贡献代码！但可以报告问题。

## 许可证

该项目使用 [GNU Affero General Public License (AGPL)](https://www.gnu.org/licenses/agpl-3.0.html)。有关详细信息，请参阅 LICENSE 文件。