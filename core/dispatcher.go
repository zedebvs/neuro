package core

import (
	"fmt"
	"neuro/audio/micro"
	"neuro/audio/voice"
	"neuro/lim/stt"
	"neuro/lim/vad"
	"neuro/logger"
	"sync"
	"time"
)

type Dispatcher struct {
	mic   *micro.Microphone
	voice *voice.VoiceAction
	vad   *vad.Detector
	stt   *stt.STT

	wg sync.WaitGroup

	isTolk     bool
	startTalk  int
	endTalk    int
	trueCount  int
	falseCount int
	timing     int

	storeFrames []*[]float32
	fullFrames  []float32
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
		d.SendSTT(frame)
		d.mic.Return(frame)
		return
	}

	d.appendBuffer(frame)
	d.trueCount++

	if d.trueCount >= d.startTalk {
		d.isTolk = true
		logger.Log.InfoLog("$3Пользователь$ $4начал говорить$", 3)
		d.voice.ClearBufer()

		for i, v := range d.storeFrames {
			d.SendSTT(v)
			d.mic.Return(v)
			d.storeFrames[i] = nil
		}

		d.storeFrames = d.storeFrames[:0]
	}
}

func (d *Dispatcher) falseSpeesh(frame *[]float32) {
	if d.trueCount != 0 {
		d.trueCount = 0
	}

	if !d.isTolk {
		d.appendBuffer(frame)
		return
	}

	d.falseCount++
	d.SendSTT(frame)
	d.mic.Return(frame)

	if d.falseCount >= d.endTalk {
		d.isTolk = false
		logger.Log.InfoLog("$3Пользователь$ $4закончил говорить$", 3)
		d.mic.Pause()
		d.wg.Add(1)
		go d.StopTalk(d.timing)
		d.voice.Resume()
		d.vad.Reset()
	}

}

func (d *Dispatcher) appendBuffer(frame *[]float32) {

	if len(d.storeFrames) < d.startTalk+10 {
		d.storeFrames = append(d.storeFrames, frame)
		return
	}

	ret := d.storeFrames[0]
	copy(d.storeFrames, d.storeFrames[1:])
	d.storeFrames[d.startTalk+10-1] = frame
	d.mic.Return(ret)
}

func (d *Dispatcher) SendSTT(frame *[]float32) {
	d.fullFrames = append(d.fullFrames, *frame...)
}

func (d *Dispatcher) WorkerGetSTT() {
	defer d.wg.Done()
	for data := range d.stt.Output {
		fmt.Printf("%s", data)
	}
}

func (d *Dispatcher) StopTalk(timing int) {
	defer d.wg.Done()

	frames := make([]float32, len(d.fullFrames)-d.endTalk)
	copy(frames, d.fullFrames[:len(d.fullFrames)-d.endTalk])

	d.stt.Input <- frames

	d.fullFrames = d.fullFrames[:0]

	time.Sleep(time.Duration(timing) * time.Millisecond)
	fmt.Println("\n")
	d.mic.Resume()
}

func (d *Dispatcher) Shutdown() {
	d.wg.Wait()
}

func New(
	mic *micro.Microphone,
	voice *voice.VoiceAction,
	vad *vad.Detector,
	stt *stt.STT,
	start int,
	end int,
	mc int,
) *Dispatcher {
	disp := &Dispatcher{
		mic:         mic,
		voice:       voice,
		vad:         vad,
		stt:         stt,
		startTalk:   start,
		endTalk:     end,
		timing:      mc,
		storeFrames: make([]*[]float32, 0, 10+start),
		fullFrames:  make([]float32, 0, 16000*10),
	}

	disp.wg.Add(2)
	go disp.Worker()
	go disp.WorkerGetSTT()

	return disp
}
