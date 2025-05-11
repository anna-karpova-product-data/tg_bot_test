package speech

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ConvertOggToWav конвертирует OGG файл в WAV формат используя ffmpeg
func ConvertOggToWav(inputPath string) (string, error) {
	// Создаем временный файл для WAV
	outputPath := filepath.Join(filepath.Dir(inputPath), filepath.Base(inputPath)+".wav")

	// Проверяем наличие ffmpeg
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found: %v", err)
	}

	// Выполняем конвертацию с логированием вывода
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error converting audio: %v, ffmpeg output: %s", err, string(output))
	}

	// Проверяем, что файл создан и не пустой
	info, err := os.Stat(outputPath)
	if err != nil {
		return "", fmt.Errorf("wav file not created: %v, ffmpeg output: %s", err, string(output))
	}
	if info.Size() < 1000 {
		return "", fmt.Errorf("wav file too small (%d bytes), likely invalid. ffmpeg output: %s", info.Size(), string(output))
	}

	return outputPath, nil
}

// CleanupAudioFiles удаляет временные аудио файлы
func CleanupAudioFiles(paths ...string) {
	for _, path := range paths {
		if path != "" {
			os.Remove(path)
		}
	}
} 