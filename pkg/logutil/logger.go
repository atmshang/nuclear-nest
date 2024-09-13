package logutil

import (
	"fmt"
	"github.com/atmshang/nuclear-nest/pkg/datautil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	once        sync.Once
	lastLogTime time.Time
)

func init() {
	once.Do(func() {
		dataPath := datautil.GetRelDataPath()

		logsDir := filepath.Join(dataPath, "logs")
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating logs directory: %v", err)
			os.Exit(1)
		}

		// 在 logs 文件夹中创建日志文件
		logFile := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".log")
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(config)
		consoleEncoder := zapcore.NewConsoleEncoder(config)

		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    16,    // 限制日志文件大小为 16MB
			MaxBackups: 5,     // 保留最多 5 个备份文件
			MaxAge:     7,     // 日志文件最多保留 7 天
			Compress:   false, // 不启用压缩
		})

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)

		logger = zap.New(core)

		go func() {
			// 定期检查并sync
			ticker := time.NewTicker(5 * time.Second) // 可以根据需要调整间隔时间
			defer ticker.Stop()

			for range ticker.C {
				if time.Since(lastLogTime) > 0 {
					Sync()
				}
			}
		}()

		// 注册一个函数，在程序退出时自动调用 Sync()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			Sync()
			os.Exit(0)
		}()
	})
}

func Printf(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
	lastLogTime = time.Now()
}

func Errorf(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
	lastLogTime = time.Now()
}

func Print(args ...interface{}) {
	logger.Info(fmt.Sprint(args...))
	lastLogTime = time.Now()
}

func Println(args ...interface{}) {
	logger.Info(fmt.Sprintln(args...))
	lastLogTime = time.Now()
}

func Fatal(args ...interface{}) {
	logger.Fatal(fmt.Sprint(args...))
	lastLogTime = time.Now()
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, args...))
	lastLogTime = time.Now()
}

func Fatalln(args ...interface{}) {
	logger.Fatal(fmt.Sprintln(args...))
	lastLogTime = time.Now()
}

func Sync() {
	_ = logger.Sync()
}
