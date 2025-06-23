package dispatcher

import (
	"fmt"
	"io"
	"os"
	"time"
)

type FileLoggerHandler struct {
	file *os.File
}

func NewFileLogger(filename string) (*FileLoggerHandler, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return &FileLoggerHandler{file: f}, err
}

func (l *FileLoggerHandler) Handle(signal TradeSignal) {
	line := fmt.Sprintf("%s [%s] %.2f %.2f\n", signal.Timestamp.Format(time.RFC3339), signal.Type, signal.Price, signal.Volume)
	io.WriteString(l.file, line)
}
