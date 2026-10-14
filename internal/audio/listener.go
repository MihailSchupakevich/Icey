package audio

import (
	"context"
	"github.com/gordonklaus/portaudio"
)

// Listener захватывает аудио и шлет чанки в канал.
type Listener struct {
	Samples chan []int16 // Канал для передачи сырого звука
}

func New() *Listener {
	return &Listener{Samples: make(chan []int16, 10)}
}

func (l *Listener) Start(ctx context.Context) error {
	portaudio.Initialize()
	defer portaudio.Terminate()

	go func() {
		stream, err := portaudio.OpenDefaultStream(1, 0, 16000, 1600, l.Samples) // 1 канал (Mono), 16кГц, буфер 100мс
		if err != nil {
			close(l.Samples)
			return
		}
		defer stream.Close()

		err = stream.Start()
		if err != nil {
			close(l.Samples)
			return
		}
		<-ctx.Done() // Ждем сигнала остановки
		stream.Stop()
		close(l.Samples)
	}()
	return nil
}
