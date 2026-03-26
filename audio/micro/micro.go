package micro

import (
	"fmt"
	"neuro/logger"
	"sync"
	"sync/atomic"

	port "github.com/gordonklaus/portaudio"
)

type Microphone struct {
	Ch     chan *[]float32
	stream *port.Stream
	wg     sync.WaitGroup
	skip   atomic.Bool
	buffer []int16
	stop   atomic.Bool
	pool   *sync.Pool
}

func (mic *Microphone) worker() {
	defer mic.wg.Done()

	for {
		err := mic.stream.Read()
		if err != nil {
			logger.Log.ErrorLog("$4Во время$ $5записи аудио$ $4произошла ошибка$", 4)
			fmt.Printf("%v\n", err)
			if mic.stop.Load() {
				return
			}
			continue
		}
		if mic.skip.Load() {
			continue
		}

		bufferPoint := mic.pool.Get().(*[]float32)
		buffer := *bufferPoint

		for i, v := range mic.buffer {
			buffer[i] = float32(v) / 32768.0
		}
		mic.Ch <- bufferPoint
	}
}

func (mic *Microphone) Shutdown() {
	mic.wg.Wait() //TODO: воТ эту хрень потом переделать

	mic.stop.Store(true)
	mic.stream.Stop()
	mic.stream.Close()

	close(mic.Ch)
	port.Terminate()
}

func (mic *Microphone) Pause() {
	mic.skip.Store(true)
	logger.Log.InfoLog("$2Микрофон$ $6остановлен$", 3)
}

func (mic *Microphone) Resume() {
	mic.skip.Store(false)
	logger.Log.InfoLog("$2Микрофон$ $6успешно возобноввил работу$", 3)
}

func (mic *Microphone) Return(point *[]float32) {
	mic.pool.Put(point)
}

func New() (*Microphone, error) {
	port.Initialize()

	buffer := make([]int16, 512)
	stream, err := port.OpenDefaultStream(1, 0, 16000, len(buffer), buffer)
	if err != nil {
		logger.Log.ErrorLog("$4При инициализации$ $5потока ввода$ $4произошла ошибка$", 4) //4 5 4
		return nil, err
	}

	err = stream.Start()
	if err != nil {
		logger.Log.ErrorLog("$4При запуске$ $5потока ввода$ $4произошла ошибка$", 4) //4 5 4
		return nil, err
	}

	logger.Log.InfoLog("$6Поток ввода$ $2успешно инициализирован и готов принимать звук$", 3) //2 6 3

	pool := &sync.Pool{
		New: func() any {
			buff := make([]float32, 512)
			return &buff
		},
	}

	micro := &Microphone{
		Ch:     make(chan *[]float32, 2048),
		stream: stream,
		buffer: buffer,
		pool:   pool,
	}

	micro.wg.Add(1)
	go micro.worker()

	return micro, nil
}
