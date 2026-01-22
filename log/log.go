package log

import stdlog "log"

// D prints debug message.
func D(v ...any) { stdlog.Println(v...) }

// I prints info message.
func I(v ...any) { stdlog.Println(v...) }

// E prints error message (string or any).
func E(v ...any) { stdlog.Println(v...) }

// Error prints error object(s).
func Error(v ...any) { stdlog.Println(v...) }

func Println(v ...any)               { stdlog.Println(v...) }
func Printf(format string, v ...any) { stdlog.Printf(format, v...) }

func Fatal(v ...any)                 { stdlog.Fatal(v...) }
func Fatalf(format string, v ...any) { stdlog.Fatalf(format, v...) }

// 如果你之前有自定义的开关/级别，也可以保留；先保证这两个存在即可
func Warnln(v ...any) {
	stdlog.Println(v...)
}

func Warnf(format string, v ...any) {
	stdlog.Printf(format, v...)
}
