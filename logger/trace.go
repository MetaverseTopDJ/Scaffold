package logger

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/MetaverseTopDJ/Scaffold/util"
)

const (
	DLTagUndefined     = "_undef"
	DLTagMySQLFailed   = "_com_mysql_failure"
	DLTagRedisFailed   = "_com_redis_failure"
	DLTagMySQLSuccess  = "_com_mysql_success"
	DLTagRedisSuccess  = "_com_redis_success"
	DLTagThriftFailed  = "_com_thrift_failure"
	DLTagThriftSuccess = "_com_thrift_success"
	DLTagHTTPSuccess   = "_com_http_success"
	DLTagHTTPFailed    = "_com_http_failure"
	DLTagTCPFailed     = "_com_tcp_failure"
	DLTagRequestIn     = "_com_request_in"
	DLTagRequestOut    = "_com_request_out"
)

const (
	_dlTag          = "dl_tag"
	_traceID        = "trace_id"
	_spanID         = "span_id"
	_childSpanID    = "child_span_id"
	_dlTagBizPrefix = "_com_"
	// _dlTagBizUndef  = "_com_undef"
)

var Log *Logger

// TraceLogic 追踪结构体
type TraceBody struct {
	TraceID     string // 链路ID
	SpanID      string // Span ID
	Caller      string //
	SrcMethod   string //
	HintCode    int64  //
	HintContent string //
}

// TContext Trace 上下文
type TContext struct {
	TraceBody
	CSpanID string
}

// TagInfo 消息追踪
func (l *Logger) TagInfo(trace *TContext, DLTag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(DLTag)
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	Info(ParseParams(m))
}

// TagWarn 警告追踪
func (l *Logger) TagWarn(trace *TContext, DLTag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(DLTag)
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	Warn(ParseParams(m))
}

// TagError 错误追踪
func (l *Logger) TagError(trace *TContext, DLTag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(DLTag)
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	Error(ParseParams(m))
}

func (l *Logger) TagTrace(trace *TContext, DLTag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(DLTag)
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	Trace(ParseParams(m))
}

// TagDebug Bug 追踪
func (l *Logger) TagDebug(trace *TContext, DLTag string, m map[string]interface{}) {
	m[_dlTag] = checkDLTag(DLTag)
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	Debug(ParseParams(m))
}

func NewTrace() *TContext {
	trace := &TContext{}
	trace.TraceID = GetTraceID()
	trace.SpanID = NewSpanID()
	return trace
}

func NewSpanID() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(util.LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

// 计算 TraceID
func calcTraceID(ip string) (traceID string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}

// GetTraceId 根据IP地址 计算链路ID
func GetTraceID() (traceID string) {
	return calcTraceID(util.LocalIP.String())
}

// 校验 DLTag 合法性
func checkDLTag(tag string) string {
	if strings.HasPrefix(tag, _dlTagBizPrefix) {
		return tag
	}
	if strings.HasPrefix(tag, "_com_") {
		return tag
	}
	if tag == DLTagUndefined {
		return tag
	}
	return tag
}

// ParseParams map格式化为string
func ParseParams(m map[string]interface{}) string {
	var tag = "_undef"
	if _tag, _have := m["dl_tag"]; _have {
		if _Val, _ok := _tag.(string); _ok {
			tag = _Val
		}
	}
	for _key, _val := range m {
		if _key == "dl_tag" {
			continue
		}
		tag = tag + "||" + fmt.Sprintf("%v=%+v", _key, _val)
	}
	tag = strings.Trim(fmt.Sprintf("%q", tag), "\"")
	return tag
}
