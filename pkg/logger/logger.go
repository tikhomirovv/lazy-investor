package logger

import "log"

type Logger struct {
}

func (l *Logger) Infof(template string, args ...any) {
	log.Printf("[INFO] "+template, args...)
}
func (l *Logger) Errorf(template string, args ...any) {
	log.Printf("[ERROR] "+template, args...)

}
func (l *Logger) Fatalf(template string, args ...any) {
	log.Fatalf(template, args...)
}
