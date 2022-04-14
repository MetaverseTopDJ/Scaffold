package logger

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"scaffold/util"
	"time"
)

var pathVariableTable map[byte]func(*time.Time) int // 日志路径变量，避免当个日志过大

type FileWriter struct {
	logLevelFloor    int           // 日志级别
	logLevelCeil     int           // 日志单元
	filename         string        // 文件名
	pathFmt          string        // 文件路径
	file             *os.File      // 打开文件描述符
	fileBufferWriter *bufio.Writer // 为io.Writer对象实现缓冲
	actions          []func(*time.Time) int
	variables        []interface{}
}

// NewFileWriter 创建一个 Writer
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (writer *FileWriter) Init() error {
	return writer.CreateLogFile()
}

// SetFileName 设置文件名称
func (writer *FileWriter) SetFileName(filename string) {
	writer.filename = filename
}

func (writer *FileWriter) SetLogLevelFloor(floor int) {
	writer.logLevelFloor = floor
}

func (writer *FileWriter) SetLogLevelCeil(ceil int) {
	writer.logLevelCeil = ceil
}

func (writer *FileWriter) SetPathPattern(pattern string) error {
	n := 0
	for _, c := range pattern {
		if c == '%' {
			n++
		}
	}

	if n == 0 {
		writer.pathFmt = pattern
		return nil
	}
	writer.actions = make([]func(*time.Time) int, 0, n)
	writer.variables = make([]interface{}, n)
	tmp := []byte(pattern)
	variable := 0
	for _, c := range tmp {
		if variable == 1 {
			act, ok := pathVariableTable[c]
			if !ok {
				return errors.New("Invalid rotate pattern (" + pattern + ")") // 无效的……
			}
			writer.actions = append(writer.actions, act)
			variable = 0
			continue
		}
		if c == '%' {
			variable = 1
		}
	}
	for i, act := range writer.actions {
		now := time.Now()
		writer.variables[i] = act(&now)
	}

	writer.pathFmt = convertPatternToFmt(tmp)

	return nil
}

func (writer *FileWriter) Write(record *Record) error {
	if record.level < writer.logLevelFloor || record.level > writer.logLevelCeil {
		return nil
	}
	if writer.fileBufferWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := writer.fileBufferWriter.WriteString(record.String()); err != nil {
		return err
	}
	return nil
}

// CreateLogFile 创建日志文件
func (writer *FileWriter) CreateLogFile() (err error) {
	var f *os.File
	if err = os.MkdirAll(path.Dir(writer.filename), 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	if f, err = os.OpenFile(writer.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	}
	writer.file = f
	if writer.fileBufferWriter = bufio.NewWriterSize(writer.file, 8192); writer.fileBufferWriter == nil {
		return errors.New("new fileBufferWriter failed. ") // 新建缓存失败
	}
	return nil
}

// Rotate 实现具体的转换器
func (writer *FileWriter) Rotate() error {
	now := time.Now()
	v := 0
	rotator := false
	oldVariables := make([]interface{}, len(writer.variables))
	copy(oldVariables, writer.variables)
	for i, act := range writer.actions {
		v = act(&now)
		if v != writer.variables[i] {
			writer.variables[i] = v
			rotator = true
		}
	}
	if !rotator {
		return nil
	}
	if writer.fileBufferWriter != nil {
		if err := writer.fileBufferWriter.Flush(); err != nil {
			return err
		}
	}
	if writer.file != nil {
		// 将文件以 pattern 形式改名并关闭
		filePath := fmt.Sprintf(writer.pathFmt, oldVariables...)
		if err := os.Rename(writer.filename, filePath); err != nil {
			return err
		}
		if err := writer.file.Close(); err != nil {
			return err
		}
	}
	return writer.CreateLogFile()
}

// Flush 清空缓存
func (writer *FileWriter) Flush() error {
	if writer.fileBufferWriter != nil {
		return writer.fileBufferWriter.Flush()
	}
	return nil
}

func convertPatternToFmt(pattern []byte) string {
	pattern = bytes.Replace(pattern, []byte("%Y"), []byte("%d"), -1)
	pattern = bytes.Replace(pattern, []byte("%M"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%D"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%H"), []byte("%02d"), -1)
	pattern = bytes.Replace(pattern, []byte("%m"), []byte("%02d"), -1)
	return string(pattern)
}

func init() {
	pathVariableTable = make(map[byte]func(*time.Time) int, 5)
	pathVariableTable['Y'] = util.GetYear
	pathVariableTable['M'] = util.GetMonth
	pathVariableTable['D'] = util.GetDay
	pathVariableTable['H'] = util.GetHour
	pathVariableTable['m'] = util.GetMin
}
