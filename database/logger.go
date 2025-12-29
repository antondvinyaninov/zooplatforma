package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger - структура для логирования
type Logger struct {
	serviceName string
	logFile     *os.File
	logger      *log.Logger
}

// LogLevel - уровень логирования
type LogLevel string

const (
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
	DEBUG   LogLevel = "DEBUG"
)

// NewLogger - создание нового логгера
func NewLogger(serviceName string) (*Logger, error) {
	// Создаём папку для логов если не существует
	logDir := filepath.Join("logs", serviceName)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Имя файла с текущей датой
	logFileName := fmt.Sprintf("%s_%s.log", serviceName, time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	// Открываем файл для записи (создаём если не существует, добавляем в конец)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Создаём logger с выводом и в файл, и в консоль
	multiWriter := log.New(logFile, "", 0)

	return &Logger{
		serviceName: serviceName,
		logFile:     logFile,
		logger:      multiWriter,
	}, nil
}

// log - внутренний метод для записи лога
func (l *Logger) log(level LogLevel, message string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf(message, args...)
	logEntry := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, l.serviceName, level, formattedMessage)

	// Пишем в файл
	l.logger.Println(logEntry)

	// Также выводим в консоль
	fmt.Println(logEntry)
}

// Info - информационное сообщение
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(INFO, message, args...)
}

// Warning - предупреждение
func (l *Logger) Warning(message string, args ...interface{}) {
	l.log(WARNING, message, args...)
}

// Error - ошибка
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(ERROR, message, args...)
}

// Debug - отладочное сообщение
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(DEBUG, message, args...)
}

// Close - закрытие логгера
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// LogRequest - логирование HTTP запроса
func (l *Logger) LogRequest(method, path, ip string, statusCode int, duration time.Duration) {
	l.Info("HTTP %s %s from %s - Status: %d - Duration: %v", method, path, ip, statusCode, duration)
}

// LogError - логирование ошибки с контекстом
func (l *Logger) LogError(context string, err error) {
	l.Error("%s: %v", context, err)
}

// LogDBQuery - логирование SQL запроса
func (l *Logger) LogDBQuery(query string, duration time.Duration) {
	l.Debug("DB Query: %s - Duration: %v", query, duration)
}

// Глобальный логгер (опционально)
var GlobalLogger *Logger

// InitGlobalLogger - инициализация глобального логгера
func InitGlobalLogger(serviceName string) error {
	logger, err := NewLogger(serviceName)
	if err != nil {
		return err
	}
	GlobalLogger = logger
	return nil
}
