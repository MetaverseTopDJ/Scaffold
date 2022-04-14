package logger

import (
	"fmt"
	"os"
)

type colorRecord Record

func (record *colorRecord) String() string {
	switch record.level {
	case TRACE:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	case INFO:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[32m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[33m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[31m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[35m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			record.time, LevelFlags[record.level], record.code, record.info)
	}
	return ""
}

type ConsoleWriter struct {
	color bool
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// 执行日志输出
func (writer *ConsoleWriter) Write(record *Record) error {
	if writer.color {
		_, err := fmt.Fprint(os.Stdout, ((*colorRecord)(record)).String())
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprint(os.Stdout, record.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *ConsoleWriter) Init() error {
	return nil
}

func (writer *ConsoleWriter) SetColor(c bool) {
	writer.color = c
}
