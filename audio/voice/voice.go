package voice

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"neuro/logger"

	oto "github.com/ebitengine/oto/v3"
)

type VoiceAction struct {
	ch         chan []byte
	wg         sync.WaitGroup
	pipeWriter *io.PipeWriter
	player     *oto.Player
	skip       atomic.Bool
}

func (r *VoiceAction) worker() {
	defer r.wg.Done()

	for data := range r.ch {
		if r.skip.Load() {
			continue
		}
		_, err := r.pipeWriter.Write(data)
		if err != nil {
			logger.Log.ErrorLog("$4Во время воспроизведения$ $5аудиопотока$ $4произошла ошибка$", 4)
			//close(r.ch)
			continue //ToDo: Потом ченить придумать с риидером
		}
	}
}

func (r *VoiceAction) Shutdown() {
	close(r.ch)
	r.wg.Wait()
	r.pipeWriter.Close()

	for r.player.BufferedSize() > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	r.player.Close()
}

func (r *VoiceAction) ClearBufer() {
	r.skip.Store(true)

	for len(r.ch) > 0 {
		<-r.ch
	}
	logger.Log.InfoLog("$2Буффер аудиопотока$ $6успешно очищен$", 3) //2 6 3
}

func (r *VoiceAction) Resume() {
	r.skip.Store(false)
}

func New() (*VoiceAction, error) {

	op := &oto.NewContextOptions{
		SampleRate:   24000,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	ctx, ready, err := oto.NewContext(op)

	if err != nil {
		logger.Log.ErrorLog("$4Во время инициализации$ $5аудиопотока$ $4произошла ошибка$", 4) //4 5 4
		return nil, err
	}
	<-ready

	reader, writer := io.Pipe()
	player := ctx.NewPlayer(reader)

	go player.Play()
	logger.Log.InfoLog("$6Аудиопоток$ $2успешно инициализирован и готов принимать звук$", 3) //2 6 3

	Voice := &VoiceAction{ch: make(chan []byte, 2048), pipeWriter: writer, player: player}

	Voice.wg.Add(1)
	go Voice.worker()

	return Voice, nil
}
