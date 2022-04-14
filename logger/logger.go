package logger

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LevelFlags = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnelSizeDefault = 1024 // 默认隧道大小

// Record 日志记录结构体
type Record struct {
	time  string // 时间
	code  string // 状态码
	info  string // 日志信息
	level int    // 日志等级
}

func (record *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LevelFlags[record.level], record.time, record.code, record.info)
}

// Writer 记录器抽象类
type Writer interface {
	Init() error
	Write(*Record) error
}

// Rotator 转存器抽象类
type Rotator interface {
	Rotate() error
	SetPathPattern(string) error
}

// Flusher 刷新 抽象类
type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record // 通道类的记录
	level       int          // 日志等级
	lastTime    int64        // 日志 时间戳
	lastTimeStr string       // 日志 时间
	c           chan bool    // 通道状态
	layout      string       // 输出时间
	recordPool  *sync.Pool
}

// NewLogger 创建 Logger
func NewLogger() *Logger {
	if loggerDefault != nil && !takeUp {
		takeUp = true // 默认启动标志
		return loggerDefault
	}
	logger := new(Logger)
	logger.writers = []Writer{}
	logger.tunnel = make(chan *Record, tunnelSizeDefault) // 创建 隧道
	logger.c = make(chan bool, 2)
	logger.level = DEBUG                  // 默认日志级别
	logger.layout = "2006/01/02 15:04:05" // 时间输出格式
	logger.recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	go boostrapLogWriter(logger)

	return logger
}

// Register 注册
func (l *Logger) Register(writer Writer) {
	if err := writer.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, writer)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

// Trace 追踪日志
func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

// Debug Bug 日志
func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Warn 警告日志
func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

// Info 输出信息
func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

// 错误日志
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

// Fatal 致命错误
func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

// Close 关闭日志
func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, writer := range l.writers {
		if flush, ok := writer.(Flusher); ok {
			if err := flush.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

// 发送记录到 写入器
func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	// info 日志信息 code 日志代码
	var info, code string

	if level < l.level {
		return
	} // 如果日志等级 大于当前等级，则不发送

	if format != "" {
		info = fmt.Sprintf(format, args...)
	} else {
		info = fmt.Sprint(args...)
	}

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	record := l.recordPool.Get().(*Record)
	record.info = info
	record.code = code
	record.time = l.lastTimeStr
	record.level = level

	l.tunnel <- record
}

// 推送日志
func boostrapLogWriter(l *Logger) {
	if l == nil {
		panic("logger is nil") // 抛出异常
	}
	var (
		record *Record
		ok     bool
	)

	if record, ok = <-l.tunnel; !ok {
		l.c <- true
		return
	}

	for _, writer := range l.writers {
		if err := writer.Write(record); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case record, ok = <-l.tunnel:
			if !ok {
				l.c <- true
				return
			}
			for _, writer := range l.writers {
				if err := writer.Write(record); err != nil {
					log.Println(err)
				}
			}

			l.recordPool.Put(record)

		case <-flushTimer.C:
			for _, writer := range l.writers {
				if flush, ok := writer.(Flusher); ok {
					if err := flush.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			for _, writer := range l.writers {
				if rotator, ok := writer.(Rotator); ok {
					if err := rotator.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

// 默认 Logger
var (
	loggerDefault *Logger
	takeUp        = false // 是否启动
)

// SetLevel 设置日志级别
func SetLevel(level int) {
	defaultLoggerInit()
	loggerDefault.level = level
}

// SetLayout 设置时间输出格式
func SetLayout(layout string) {
	defaultLoggerInit()
	loggerDefault.layout = layout
}

func Trace(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(FATAL, fmt, args...)
}

// Register 注册
func Register(writer Writer) {
	defaultLoggerInit()
	loggerDefault.Register(writer)
}

// Close 关闭
func Close() {
	defaultLoggerInit()
	loggerDefault.Close()
	loggerDefault = nil
	takeUp = false
}

// 默认日志初始化
func defaultLoggerInit() {
	if !takeUp {
		loggerDefault = NewLogger()
	}
}
