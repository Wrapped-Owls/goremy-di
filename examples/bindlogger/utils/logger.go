package utils

import "log"

type (
	Logger interface {
		Error(...any)
		Errorf(format string, args ...any)
		Info(...any)
		Infof(format string, args ...any)
		Warning(...any)
		Warningf(format string, args ...any)
	}
	LogLevel uint8

	DefaultLogger struct {
		Level LogLevel
	}
)

const (
	LevelError LogLevel = iota << 1
	LevelWarning
	LevelInfo
	LevelDefault = LevelWarning | LevelError
)

func (l DefaultLogger) hasInfo() bool {
	return (l.Level & LevelInfo) == LevelInfo
}

func (l DefaultLogger) hasWarning() bool {
	return (l.Level & LevelWarning) == LevelWarning
}

func (l DefaultLogger) hasError() bool {
	return (l.Level & LevelError) == LevelError
}

func (l DefaultLogger) Error(args ...any) {
	if l.hasError() {
		log.Println("[Error]", args)
	}
}

func (l DefaultLogger) Errorf(format string, args ...any) {
	if l.hasError() {
		log.Printf("[Errorf] "+format, args...)
	}
}

func (l DefaultLogger) Warning(args ...any) {
	if l.hasWarning() {
		log.Println("[Warning]", args)
	}
}

func (l DefaultLogger) Warningf(format string, args ...any) {
	if l.hasWarning() {
		log.Printf("[Warningf] "+format, args...)
	}
}

func (l DefaultLogger) Info(args ...any) {
	if l.hasInfo() {
		log.Println("[Info]", args)
	}
}

func (l DefaultLogger) Infof(format string, args ...any) {
	if l.hasInfo() {
		log.Printf("[Infof] "+format, args...)
	}
}
