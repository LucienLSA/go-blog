package logging

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var lg *zap.Logger
// 不建议在此使用全局的logger变量

// InitLogger 初始化Logger
func InitLogger(cfg *settings.LogConfig, mode string) (err error) {
	writeSyncer := getLogWriter(
		// viper.GetString("log.filepath")+viper.GetString("log.filename"),
		// viper.GetInt("log.max_size"),
		// viper.GetInt("log.max_backups"),
		// viper.GetInt("log.max_age"),
		cfg.FilePath+cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups, cfg.MaxAge,
	)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	// err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	if mode == settings.Conf.Mode {
		// 开发模式，日志输出到终端和日志文件中
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// EnhancedGinLogger 增强的Gin日志中间件，支持TraceID和用户信息关联
func EnhancedGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 生成或获取TraceID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}
		c.Header("X-Trace-ID", traceID)

		// 获取用户信息（如果已认证）
		var userID int64
		var username string
		if userInfo, exists := c.Get("user_info"); exists {
			if u, ok := userInfo.(*ctl.UserInfo); ok {
				userID = u.UserId
				username = u.UserName
			}
		}

		c.Next()

		cost := time.Since(start)

		// 构建结构化日志字段
		logFields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		}

		// 如果有用户信息，添加到日志中
		if userID > 0 {
			logFields = append(logFields,
				zap.Int64("user_id", userID),
				zap.String("username", username),
			)
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			logFields = append(logFields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		zap.L().Info("HTTP Request", logFields...)
	}
}

// generateTraceID 生成唯一的TraceID
func generateTraceID() string {
	return fmt.Sprintf("trace_%d_%d", time.Now().UnixNano(), rand.Int63())
}

// GetLoggerWithContext 获取带有上下文信息的logger
func GetLoggerWithContext(ctx context.Context) *zap.Logger {
	logger := zap.L()

	// 尝试从上下文获取TraceID
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		logger = logger.With(zap.String("trace_id", traceID))
	}

	// 尝试从上下文获取用户信息
	if userInfo, ok := ctx.Value("user_info").(*ctl.UserInfo); ok {
		logger = logger.With(
			zap.Int64("user_id", userInfo.UserId),
			zap.String("username", userInfo.UserName),
		)
	}

	return logger
}

// LogUserAction 记录用户操作日志
func LogUserAction(ctx context.Context, action string, details map[string]interface{}) {
	logger := GetLoggerWithContext(ctx)

	fields := []zap.Field{
		zap.String("action", action),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}

	logger.Info("User Action", fields...)
}
