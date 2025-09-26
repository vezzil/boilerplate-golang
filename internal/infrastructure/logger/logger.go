package logger

import (
	"log"
	"os"
)

// Simple wrapper around the stdlib logger. Swappable later (zap/logrus).
var (
	std = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func Info(msg string, args ...any)  { std.Printf("INFO: "+msg+"\n", args...) }
func Warn(msg string, args ...any)  { std.Printf("WARN: "+msg+"\n", args...) }
func Error(msg string, args ...any) { std.Printf("ERROR: "+msg+"\n", args...) }
func Fatal(msg string, args ...any) { std.Fatalf("FATAL: "+msg+"\n", args...) }
