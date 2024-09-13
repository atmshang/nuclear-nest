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

Nuclear Nest 提供了跨平台的数据路径管理功能，自动根据不同平台设置和管理数据目录路径。你可以通过 `SetAppName` 设置应用名称，从而在不同平台上生成相应的数据目录路径。

#### 使用方法

首先，通过 `SetAppName` 设置应用名称：

```go
import "github.com/atmshang/nuclear-nest/pkg/datautil"

func main() {
    datautil.SetAppName("MyApp")
    path := datautil.GetRelDataPath()
    fmt.Println("Data path:", path)
}
```

#### 平台构建方式

为了支持不同的平台，你需要在构建时指定目标平台。以下是三个主要平台的构建方式：

- **Android**：
  ```bash
  GOOS=android GOARCH=arm64 go build -o myapp-android
  ```

- **Linux**：
  ```bash
  GOOS=linux GOARCH=amd64 go build -o myapp-linux
  ```

- **Windows**：
  ```bash
  GOOS=windows GOARCH=amd64 go build -o myapp-windows.exe
  ```

#### 路径差异

根据不同的平台，数据目录路径会有所不同：

- **Android**：
  - 路径格式：`/data/.LincOS/data/<AppName>/`
  - 示例路径：`/data/.LincOS/data/MyApp/`

- **Linux**：
  - 路径格式：`/.LincOS/data/<AppName>/`
  - 示例路径：`/.LincOS/data/MyApp/`

- **Windows**：
  - 路径格式：`<ExecutableDir>/data/<AppName>/`
  - 示例路径：`C:\Program Files\MyApp\data\MyApp\`

在每个平台上，`GetRelDataPath` 函数会自动检查并创建必要的目录结构，以确保数据存储路径的可用性。

通过这种方式，Nuclear Nest 确保了在不同平台上的一致性和便利性，使得开发者可以专注于应用逻辑，而不必担心平台特定的路径管理问题。



当然，以下是关于版本信息管理功能的更详细文档，包括如何设置版本信息、获取版本信息，以及生成 MD5 校验文件和变更日志文件的详细说明：

### 版本信息管理

Nuclear Nest 提供了强大的版本信息管理功能，帮助开发者轻松管理应用的版本信息，并生成相关的校验和日志文件。

#### 设置版本信息

在应用启动时，你可以通过 `SetVersionList` 函数设置应用的版本信息：

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
}
```

- **VersionInfo**：这是一个结构体，用于存储应用的版本信息，包括应用名称、版本名称、版本号、作者、发布日期和描述。

#### 获取版本信息

你可以使用 `GetVersionInfo` 函数获取当前应用的版本信息：

```go
versionInfo := versionutil.GetVersionInfo()
fmt.Printf("Current version: %s\n", versionInfo.VersionName)
```

- **GetVersionInfo**：返回当前设置的版本信息，包括可执行文件的 MD5 校验值和从 `md5checksum` 文件中读取的校验值。

#### 打印版本信息

`PrintVersionInfo` 函数可以将版本信息以 JSON 格式打印到标准输出：

```go
versionutil.PrintVersionInfo()
```

- **PrintVersionInfo**：将版本信息转换为 JSON 格式并打印，便于调试和查看。

#### 生成 MD5 校验文件

`CreateMD5File` 函数用于生成当前可执行文件的 MD5 校验文件：

```go
versionutil.CreateMD5File()
```

- **CreateMD5File**：计算当前可执行文件的 MD5 校验值，并将其写入到 `md5checksum` 文件中。这有助于验证文件的完整性。

#### 生成变更日志文件

`CreateChangeLogFile` 函数用于生成变更日志文件，记录版本变更历史：

```go
versionutil.CreateChangeLogFile()
```

- **CreateChangeLogFile**：根据设置的版本信息生成 `CHANGELOG.md` 文件，记录每个版本的详细信息，包括版本号、作者、发布日期和描述。这有助于跟踪应用的演变和更新历史。

#### 工作流程

1. **设置版本信息**：在应用启动时，通过 `SetVersionList` 设置版本信息。
2. **获取和打印版本信息**：使用 `GetVersionInfo` 和 `PrintVersionInfo` 查看当前版本信息。
3. **生成校验和日志文件**：通过 `CreateMD5File` 和 `CreateChangeLogFile` 生成相关文件，确保版本管理的完整性和可追溯性。

通过这些功能，Nuclear Nest 的版本信息管理模块帮助开发者有效地管理和发布应用版本，确保应用的完整性和可追溯性。



### 日志记录

Nuclear Nest 提供了基于 `zap` 和 `lumberjack` 的高效日志记录系统，支持日志文件的自动分割和管理。该系统确保了日志的可靠性和可维护性。

#### 使用方法

要使用日志记录功能，你可以调用 `logutil` 包中的日志函数：

```go
import "github.com/atmshang/nuclear-nest/pkg/logutil"

func main() {
    logutil.Printf("Application started")
    logutil.Println("This is a log message")
    logutil.Errorf("This is an error message")
}
```

#### 日志生成位置

日志文件会根据平台生成在以下位置：

- **Android**、**Linux**：
  - 日志文件存放在应用的数据目录下的 `logs` 文件夹中。
  - 示例路径：`/data/.LincOS/data/MyApp/logs/` 或 `/.LincOS/data/MyApp/logs/`

- **Windows**：
  - 日志文件存放在可执行文件所在目录的 `logs` 文件夹中。
  - 示例路径：`C:\Program Files\MyApp\data\MyApp\logs\`

#### 存放规则

日志系统使用 `lumberjack` 实现日志文件的自动分割和管理，具体规则如下：

- **日志文件大小**：每个日志文件的最大大小为 16MB。
- **日志文件备份**：最多保留 5 个备份文件。
- **日志文件保留时间**：日志文件最多保留 7 天。
- **压缩**：日志文件不启用压缩。

这些规则确保了日志文件不会无限制地增长，从而占用过多的磁盘空间，同时也提供了足够的历史日志以供调试和审计使用。

通过这种方式，Nuclear Nest 的日志记录系统提供了灵活且高效的日志管理方案，使得开发者可以轻松地跟踪和调试应用程序的运行状态。



### 命令行参数解析

Nuclear Nest 提供了一个通用的命令行参数解析功能，帮助开发者轻松处理常见的命令行操作。通过 `flagutil` 包，你可以快速解析和响应命令行参数。

#### 使用方法

要使用命令行参数解析功能，你可以调用 `flagutil` 包中的 `ParseFlags` 函数：

```go
import "github.com/atmshang/nuclear-nest/pkg/flagutil"

func main() {
    flagutil.ParseFlags()
    // 其他初始化代码
}
```

#### 实现的参数

`flagutil` 实现了以下命令行参数：

- **`-v`**：打印版本信息。
  - 使用示例：`./myapp -v`
  - 功能：输出当前应用的版本信息，包括版本号、版本名称、作者等。

- **`-m`**：生成并写入 MD5 校验文件。
  - 使用示例：`./myapp -m`
  - 功能：计算当前可执行文件的 MD5 校验值，并将其写入到 `md5checksum` 文件中。

- **`-c`**：生成并写入变更日志文件。
  - 使用示例：`./myapp -c`
  - 功能：根据版本信息生成 `CHANGELOG.md` 文件，记录版本变更历史。

这些参数为开发者提供了便捷的工具来管理和发布应用版本信息，确保应用的完整性和可追溯性。通过这种方式，Nuclear Nest 的命令行参数解析功能帮助开发者更好地控制应用的运行行为和版本管理。



### API 处理

Nuclear Nest 提供了标准化的 API 响应结构和错误处理机制，帮助开发者更轻松地构建和维护 RESTful API。通过 `apiutil` 包，你可以实现一致的 API 返回格式和高效的错误处理。

#### 标准返回结构体

`apiutil` 提供了一个标准的 API 响应结构体 `Response`，用于统一 API 的返回格式：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

- **Code**：状态码，用于表示请求的处理结果。
- **Message**：消息文本，提供关于请求处理的简要说明。
- **Data**：返回的数据，可以是任意类型。

#### 空结构体

`apiutil` 还提供了一个 `EmptyResponse` 结构体，用于表示空数据的返回：

```go
type EmptyResponse struct{}
```

在需要返回空数据时，可以使用 `EmptyResponse` 作为 `Data` 字段的值。

#### 错误处理

`apiutil` 提供了一个全局错误处理中间件 `ErrorHandler`，用于捕获和处理未被捕获的错误：

```go
import (
    "github.com/atmshang/nuclear-nest/pkg/apiutil"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    apiutil.UseErrorHandler(r)
    // 其他路由和处理器
    r.Run()
}
```

#### 带锁的 API 超时处理

在某些情况下，你可能需要对某些 API 请求进行锁定，以防止并发修改。`apiutil` 提供了 `TryLock` 函数，用于在指定超时时间内尝试获取锁：

```go
import (
    "github.com/atmshang/nuclear-nest/pkg/apiutil"
    "github.com/gin-gonic/gin"
    "sync"
    "time"
)

var myLock sync.Mutex

func main() {
    r := gin.Default()

    r.GET("/locked-resource", func(c *gin.Context) {
        if !apiutil.TryLock(c, &myLock, 5*time.Second) {
            return // 如果获取锁失败，响应已经在 TryLock 中处理
        }
        defer myLock.Unlock()

        // 处理请求
        c.JSON(200, apiutil.Response{
            Code:    2000,
            Message: "Resource accessed successfully",
            Data:    "Your data here",
        })
    })

    r.Run()
}
```

- **TryLock**：尝试在指定的超时时间内获取锁。如果获取失败，将返回一个标准的错误响应。
- **锁定机制**：确保在处理关键资源时，避免并发访问导致的数据不一致或冲突。

通过这些功能，Nuclear Nest 的 API 处理模块帮助开发者实现一致的 API 设计，并提供了高效的并发控制机制。



### 认证工具

Nuclear Nest 提供了一个灵活的认证工具，用于模块间的内部认证。该工具基于 RSA 和 AES 加密，确保请求的安全性和完整性。

#### 设置公钥和私钥

在使用认证工具之前，你需要设置 RSA 公钥和私钥。这些密钥用于加密和解密认证信息。

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
}
```

##### 安全注意事项

- **密钥管理**：确保公钥和私钥的安全存储。私钥应严格保密，不应在代码库中硬编码。
- **环境变量**：考虑使用环境变量或安全的密钥管理服务来存储和加载密钥。
- **定期更换**：定期更换密钥对，以提高安全性。

#### 为请求增加认证信息

在发送请求时，你可以使用 `GenerateAuthHeaderValue` 函数为请求增加认证信息：

```go
import (
    "github.com/atmshang/nuclear-nest/pkg/authutil"
    "net/http"
)

func sendAuthenticatedRequest(serviceName string, url string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    authHeader := authutil.GenerateAuthHeaderValue(serviceName)
    req.Header.Set(authutil.HeaderInternalServiceAuth, authHeader)

    return client.Do(req)
}
```

- **GenerateAuthHeaderValue**：生成一个认证头部值，该值包含服务名称和过期时间，并使用 RSA 加密。

#### 接收请求的认证处理

在接收请求时，你可以使用 `InternalServiceAuth` 中间件来验证请求的认证信息：

```go
import (
    "github.com/atmshang/nuclear-nest/pkg/authutil"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.Use(authutil.InternalServiceAuth())

    r.GET("/secure-endpoint", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Authenticated successfully",
        })
    })

    r.Run()
}
```

- **InternalServiceAuth**：这是一个 Gin 中间件，用于验证请求的认证信息。如果认证失败，将返回 `401 Unauthorized`。

#### 工作流程

1. **发送方**：使用 `GenerateAuthHeaderValue` 生成认证信息，并将其添加到请求头中。
2. **接收方**：使用 `InternalServiceAuth` 中间件验证请求头中的认证信息。
3. **认证机制**：认证信息使用 RSA 加密，确保只有持有正确私钥的接收方能够解密和验证。

通过这些功能，Nuclear Nest 的认证工具为模块间通信提供了安全可靠的认证机制，确保数据的安全性和完整性。

## 贡献

不欢迎贡献代码！但可以报告问题。

## 许可证

该项目使用 [GNU Affero General Public License (AGPL)](https://www.gnu.org/licenses/agpl-3.0.html)。有关详细信息，请参阅 LICENSE 文件。shu