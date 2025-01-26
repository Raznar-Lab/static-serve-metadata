package logger

import "fmt"

type Logger struct {
	prefix string
}

const debugPrefix = "DEBUG"
const (
	LOG_PREFIX_SYSTEM   = "SYSTEM"
)

const (
	LEVEL_INFO  = "INFO"
	LEVEL_WARN  = "WARN"
	LEVEL_ERROR = "ERROR"
)

var isDebug = false

func InitDebugValue() {
	isDebug = true
}

func (l *Logger) LogInfo(v ...interface{}) {
	fmt.Printf("[%s] [%s] %v\n", l.prefix, LEVEL_INFO, fmt.Sprint(v...))
}

func (l *Logger) LogWarn(v ...interface{}) {
	fmt.Printf("[%s] [%s] %v\n", l.prefix, LEVEL_WARN, fmt.Sprint(v...))
}

func (l *Logger) LogError(v ...interface{}) {
	fmt.Printf("[%s] [%s] %v\n", l.prefix, LEVEL_ERROR, fmt.Sprint(v...))
}

func (l *Logger) DebugInfo(v ...interface{}) {
	if isDebug {
		fmt.Printf("[%s] [%s] %v\n", debugPrefix + "@" + l.prefix, LEVEL_INFO, fmt.Sprint(v...))
	}
}

func (l *Logger) DebugWarn(v ...interface{}) {
	if isDebug {
		fmt.Printf("[%s] [%s] %v\n", debugPrefix + "@" + l.prefix, LEVEL_WARN, fmt.Sprint(v...))
	}
}

func (l *Logger) DebugError(v ...interface{}) {
	if isDebug {
		fmt.Printf("[%s]  [%s] %v\n", debugPrefix + "@" + l.prefix, LEVEL_ERROR, fmt.Sprint(v...))
	}
}


var System = &Logger{LOG_PREFIX_SYSTEM}
