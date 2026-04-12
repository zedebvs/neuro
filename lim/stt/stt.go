package stt

import (
	"fmt"
	"neuro/logger"
	"sync"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

type STT struct {
	model whisper.Model

	Input  chan []float32
	Output chan string

	wg sync.WaitGroup
}

func (s *STT) worker() {
	defer s.wg.Done()

	ctx, err := s.model.NewContext()
	if err != nil {
		logger.Log.ErrorLog("$4При инициализации$ $5контекста модели преобразования звука$ $4произошла ошибка$ "+err.Error(), 4)
		return
	}
	ctx.SetLanguage("ru")
	ctx.SetTemperature(1.0) //

	for data := range s.Input {
		err := ctx.Process(data, nil, nil, nil)
		if err != nil {
			logger.Log.ErrorLog("$4При обработке$ $5аудиофрагмента$ $4произошла ошибка$ "+err.Error(), 4)
			continue
		}
		var fullText string

		for {
			segment, err := ctx.NextSegment()
			if err != nil {
				break
			}
			fullText += fmt.Sprintf(" %s", segment.Text)
			//s.Output <- segment.Text
		}
		s.Output <- fullText
	}
}

func (s *STT) Shutdown() {
	s.wg.Wait()
	s.model.Close()
}

func New(modeleName string) (*STT, error) {
	model, err := whisper.New(modeleName)
	if err != nil {
		logger.Log.ErrorLog("$4При инициализации$ $5модели преобразования звука$ $4произошла ошибка$ "+err.Error(), 4)
		return nil, err
	}

	stt := &STT{
		model:  model,
		Input:  make(chan []float32, 5),
		Output: make(chan string, 15),
	}

	stt.wg.Add(1)
	go stt.worker()

	return stt, nil
}
