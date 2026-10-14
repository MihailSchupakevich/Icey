package main

import (
	"Icey/internal/audio"
	"Icey/internal/stt"
	"Icey/internal/wake"
	"context"
	"log"
	"os/exec"
	"strings"
)

const (
	voskModelPath = "./models/vosk-model-small-ru-0.22/model"
	sampleRate    = 16000
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := audio.New()
	if err := listener.Start(ctx); err != nil {
		log.Fatal(err)
	}

	wakeDet, err := wake.New(0.5)
	if err != nil {
		log.Fatal(err)
	}
	defer wakeDet.Delete()

	sttRec, err := stt.New(voskModelPath, sampleRate)
	if err != nil {
		log.Fatal(err)
	}
	defer sttRec.Close()

	log.Println("Ассистент запущен. Скажите 'Яндекс'...")

	// Буфер для накопления речи после активации
	var speechBuffer []int16
	isListening := false

	for chunk := range listener.Samples {
		if !isListening {
			// Режим ожидания: ищем ключевое слово
			if wakeDet.Process(chunk) {
				log.Println("[АКТИВАЦИЯ] Слушаю команду...")
				isListening = true
				speechBuffer = append(speechBuffer, chunk...)
			}
		} else {
			// Режим записи команды
			speechBuffer = append(speechBuffer, chunk...)

			// Пробуем распознать каждые N миллисекунд
			if len(speechBuffer) > sampleRate { // Накопили ~1 секунду
				// Преобразуем int16 в []byte для Vosk
				bytesData := int16ToBytes(speechBuffer)
				result, isFinal := sttRec.AcceptWaveform(bytesData)

				if isFinal && strings.TrimSpace(result) != "" {
					log.Printf("[РАСПОЗНАНО]: %s", result)
					executeCommand(result)
					isListening = false
					speechBuffer = nil
					log.Println("Снова в режиме ожидания...")
				} else if isListening && len(speechBuffer) > sampleRate*5 {
					// Таймаут: если говорят слишком долго без паузы
					log.Println("[ТАЙМАУТ]")
					isListening = false
					speechBuffer = nil
				}
			}
		}
	}
}

func executeCommand(text string) {
	switch {
	case strings.Contains(text, "открой блокнот"):
		exec.Command("notepad.exe").Start()
	case strings.Contains(text, "выключи компьютер"):
		exec.Command("shutdown", "/s", "/t", "5").Run()
	case strings.Contains(text, "какой сегодня день"):
		// Здесь можно подключить локальную LLM через Ollama API
		log.Println("Сегодня суббота, 22 августа 2026 года.")
	}
}

// Вспомогательная функция преобразования
func int16ToBytes(in []int16) []byte {
	out := make([]byte, len(in)*2)
	for i, v := range in {
		out[i*2] = byte(v)
		out[i*2+1] = byte(v >> 8)
	}
	return out
}
