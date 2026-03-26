package core

import (
	"neuro/audio/micro"
	"neuro/audio/voice"
	"neuro/lim/vad"
	"neuro/logger"
	"sync"
	"time"
)

type Dispatcher struct {
	mic   *micro.Microphone
	voice *voice.VoiceAction
	vad   *vad.Detector
	//stt *
	wg   sync.WaitGroup
	pool *sync.Pool

	isTolk     bool
	startTalk  int
	endTalk    int
	trueCount  int
	falseCount int
	timing     int

	storeFrames []*[]float32
}

func (d *Dispatcher) Worker() {

	defer d.wg.Done()

	for data := range d.mic.Ch {
		if d.vad.Process(data) {
			d.trueSpeesh(data)
		} else {
			d.falseSpeesh(data)
		}
	}

}

func (d *Dispatcher) trueSpeesh(frame *[]float32) {
	if d.falseCount != 0 {
		d.falseCount = 0
	}

	if d.isTolk {
		//	d.SendSTT(frame)
		d.mic.Return(frame)
		return
	}

	d.storeFrames = append(d.storeFrames, frame)
	d.trueCount++

	if d.trueCount >= d.startTalk {
		d.isTolk = true
		logger.Log.InfoLog("$3Пользователь$ $6начал говорить$", 3)
		d.voice.ClearBufer()

		for _, v := range d.storeFrames {
			//d.SendSTT(v)
			d.mic.Return(v)
		}

		d.storeFrames = d.storeFrames[:0]
	}
}

func (d *Dispatcher) falseSpeesh(frame *[]float32) {
	if d.trueCount != 0 {
		d.trueCount = 0
	}

	if !d.isTolk {
		for _, c := range d.storeFrames {
			d.mic.Return(c)
		}
		d.storeFrames = d.storeFrames[:0]
		d.mic.Return(frame)
		return
	}

	d.falseCount++
	//d.SendSTT(frame)
	d.mic.Return(frame)

	if d.falseCount >= d.endTalk {
		d.isTolk = false
		logger.Log.InfoLog("$3Пользователь$ $4закончил говорить$", 3)
		d.voice.Resume()
		d.vad.Reset()
		//d.sst.state.store(true)
		d.wg.Add(1)
		go d.StopTalk(d.timing)
	}

}

func (d *Dispatcher) SendSTT(frame *[]float32) {
	//bufferPoint := d.pool.Get().(*[]float32)
	//copy(*bufferPoint, *frame)

	//sst.ch <- bufferPoint
}

func (d *Dispatcher) WorkerGetSTT() {
	//defer d.sst.wg.Done()
	//for data := range d.sst.ch {
	//} чета типа такого будет
}

func (d *Dispatcher) StopTalk(timing int) {
	defer d.wg.Done()

	d.mic.Pause()
	time.Sleep(time.Duration(timing) * time.Millisecond)
	d.mic.Resume()
}

func (d *Dispatcher) Shutdown() {
	d.wg.Wait()
}

func New(
	mic *micro.Microphone,
	voice *voice.VoiceAction,
	vad *vad.Detector,
	start int,
	end int,
	mc int,
) *Dispatcher {
	disp := &Dispatcher{
		mic:         mic,
		voice:       voice,
		vad:         vad,
		startTalk:   start,
		endTalk:     end,
		timing:      mc,
		storeFrames: make([]*[]float32, 0, start),
		pool: &sync.Pool{
			New: func() any {
				buff := make([]float32, 512)
				return &buff
			},
		},
	}

	disp.wg.Add(1)
	go disp.Worker()

	return disp
}
