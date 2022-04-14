package logger

import (
	"errors"
)

// ConfigFileWriter 文件日志配置
type ConfigFileWriter struct {
	On           bool   `toml:"On"`         // 是否开启文件日志
	Path         string `toml:"Path"`       // 日志路径
	WfPath       string `toml:"WfPath"`     // 日志写入路径
	RotatePath   string `toml:"RotatePath"` // 反转日志路径
	RotateWfPath string `toml:"RotateWfPath"`
}

// ConfigConsoleWriter 控制台输出
type ConfigConsoleWriter struct {
	On    bool `toml:"On"`
	Color bool `toml:"Color"`
}

type LogConfig struct {
	Level string              `toml:"Level"`
	FW    ConfigFileWriter    `toml:"FileWriter"`
	CW    ConfigConsoleWriter `toml:"ConsoleWriter"`
}

// SetupLogInstanceWithConfig 日志模块 初始化配置
func SetupLogInstanceWithConfig(config LogConfig, logger *Logger) (err error) {
	// 判断是否开启 文件日志
	if config.FW.On {
		if len(config.FW.Path) > 0 { // 判断是否配置文件日志路径
			writer := NewFileWriter()
			writer.SetFileName(config.FW.Path)
			err := writer.SetPathPattern(config.FW.RotatePath)
			if err != nil {
				return err
			}
			writer.SetLogLevelFloor(TRACE)
			if len(config.FW.WfPath) > 0 {
				writer.SetLogLevelCeil(INFO)
			} else {
				writer.SetLogLevelCeil(ERROR)
			}
			logger.Register(writer)
		}

		if len(config.FW.WfPath) > 0 { // 判断是否配置转存地址
			rotateWriter := NewFileWriter()
			rotateWriter.SetFileName(config.FW.WfPath)
			err := rotateWriter.SetPathPattern(config.FW.RotateWfPath)
			if err != nil {
				return err
			}
			rotateWriter.SetLogLevelFloor(WARNING)
			rotateWriter.SetLogLevelCeil(ERROR)
			logger.Register(rotateWriter)
		}
	}
	// 判断是否开启 控制器输出
	if config.CW.On {
		writer := NewConsoleWriter()
		writer.SetColor(config.CW.Color)
		logger.Register(writer)
	}
	switch config.Level {
	case "trace":
		logger.SetLevel(TRACE)

	case "debug":
		logger.SetLevel(DEBUG)

	case "info":
		logger.SetLevel(INFO)

	case "warning":
		logger.SetLevel(WARNING)

	case "error":
		logger.SetLevel(ERROR)

	case "fatal":
		logger.SetLevel(FATAL)

	default:
		err = errors.New("InvalidLogLevel")
	}
	return
}

// SetupDefaultLogWithConfig 加载默认日志配置
func SetupDefaultLogWithConfig(logConfig LogConfig) (err error) {
	defaultLoggerInit()
	return SetupLogInstanceWithConfig(logConfig, loggerDefault)
}
