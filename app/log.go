package app

import (
	"scaffold/logger"
	"scaffold/util"
)

var LogConf *LogConfig

// LogConfigBase 基础配置信息
type LogConfigBase struct {
	Level string `mapstructure:"level"`
	Mode  string `mapstructure:"mode"`
}

// LogConfigFileWriter 文件日志 配置结构体
type LogConfigFileWriter struct {
	On           bool   `mapstructure:"on"`
	Path         string `mapstructure:"path"`
	WfPath       string `mapstructure:"wf_path"`
	RotatePath   string `mapstructure:"rotate_path"`
	RotateWfPath string `mapstructure:"rotate_wf_path"`
}

// LogConfigConsoleWriter 控制台输出日志 配置结构体
type LogConfigConsoleWriter struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

// LogConfig 日志配置
type LogConfig struct {
	Base LogConfigBase          `mapstructure:"base"`
	FW   LogConfigFileWriter    `mapstructure:"file_writer"`
	CW   LogConfigConsoleWriter `mapstructure:"console_writer"`
}

// InitLogConfig 初始化 Log 配置
func InitLogConfig(path string) (err error) {
	LogConf = &LogConfig{}
	if err = util.ParseConfig(path, LogConf); err != nil {
		return
	}
	// 设置默认日志级别
	if LogConf.Base.Level == "" {
		LogConf.Base.Level = "trace"
	}

	//配置日志
	logConfig := logger.LogConfig{
		Level: LogConf.Base.Level,
		FW: logger.ConfigFileWriter{
			On:           LogConf.FW.On,
			Path:         LogConf.FW.Path,
			WfPath:       LogConf.FW.WfPath,
			RotatePath:   LogConf.FW.RotatePath,
			RotateWfPath: LogConf.FW.RotateWfPath,
		},
		CW: logger.ConfigConsoleWriter{
			On:    LogConf.CW.On,
			Color: LogConf.CW.Color,
		},
	}
	if err := logger.SetupDefaultLogWithConfig(logConfig); err != nil {
		panic(err)
	}
	logger.SetLayout("2006-01-02T15:04:05.000")
	return
}
