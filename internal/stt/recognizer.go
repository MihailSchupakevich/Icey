package stt

// Используем готовую обертку над vosk-api, например github.com/swormperger/vosk
import (
	vosk "github.com/swormperger/vosk"
	"os"
)

type Recognizer struct {
	model *vosk.Model
	rec   *vosk.Recognizer
}

func New(modelPath string, sampleRate int) (*Recognizer, error) {
	m, err := vosk.NewModel(modelPath)
	if err != nil {
		return nil, err
	}
	r := m.NewRecognizer(sampleRate)
	return &Recognizer{model: m, rec: r}, nil
}

func (r *Recognizer) AcceptWaveform(data []byte) (string, bool) {
	// data приходит как []int16, преобразуем в байты для API
	if r.rec.AcceptWaveform(data) {
		return r.rec.Result(), true // Фраза завершена (пауза)
	}
	return "", false // Говорят дальше
}

func (r *Recognizer) Close() {
	r.model.Close()
}
