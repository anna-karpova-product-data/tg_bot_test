package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp   string `json:"timestamp"`
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	MessageType string `json:"message_type"` // "voice" или "text"
	Content     string `json:"content"`
	BotResponse string `json:"bot_response"`
	Mood        string `json:"mood,omitempty"`
}

type Logger struct {
	logFile *os.File
}

func New(logPath string) (*Logger, error) {
	// Создаем директорию для логов, если она не существует
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Открываем файл для записи логов
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{
		logFile: file,
	}, nil
}

func (l *Logger) Log(userID int64, username, messageType, content, botResponse, mood string) error {
	entry := LogEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		UserID:      userID,
		Username:    username,
		MessageType: messageType,
		Content:     content,
		BotResponse: botResponse,
		Mood:        mood,
	}

	// Преобразуем запись в JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	// Добавляем перенос строки
	jsonData = append(jsonData, '\n')

	// Записываем в файл
	if _, err := l.logFile.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
