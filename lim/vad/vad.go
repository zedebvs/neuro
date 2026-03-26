package vad

import (
	"neuro/logger"

	ort "github.com/yalue/onnxruntime_go"
)

type Detector struct {
	session   *ort.AdvancedSession
	inputData []float32
	hData     []float32
	cData     []float32
	srData    []int64
	probData  []float32
	hnData    []float32
	cnData    []float32
	threshold float32
}

func New(modelPath string, threshold float32) (*Detector, error) {
	vad := &Detector{
		threshold: threshold,
		inputData: make([]float32, 512),
		hData:     make([]float32, 128),
		cData:     make([]float32, 128),
		srData:    []int64{16000},
		probData:  make([]float32, 1),
		hnData:    make([]float32, 128),
		cnData:    make([]float32, 128),
	}

	inputTensor, _ := ort.NewTensor(ort.NewShape(1, 512), vad.inputData)
	srTensor, _ := ort.NewTensor(ort.NewShape(1), vad.srData)
	hTensor, _ := ort.NewTensor(ort.NewShape(2, 1, 64), vad.hData)
	cTensor, _ := ort.NewTensor(ort.NewShape(2, 1, 64), vad.cData)

	outProbTensor, _ := ort.NewTensor(ort.NewShape(1, 1), vad.probData)
	outHTensor, _ := ort.NewTensor(ort.NewShape(2, 1, 64), vad.hnData)
	outCTensor, _ := ort.NewTensor(ort.NewShape(2, 1, 64), vad.cnData)

	options, err := ort.NewSessionOptions()
	if err != nil {
		logger.Log.ErrorLog("$4При инициализации$ $5потоков для модели vad$ $4произошла ошибка$"+err.Error(), 4)
		return nil, err
	}

	options.SetIntraOpNumThreads(1)
	options.SetInterOpNumThreads(1)

	session, err := ort.NewAdvancedSession(
		modelPath,
		[]string{"input", "sr", "h", "c"},
		[]string{"output", "hn", "cn"},
		[]ort.ArbitraryTensor{inputTensor, srTensor, hTensor, cTensor},
		[]ort.ArbitraryTensor{outProbTensor, outHTensor, outCTensor},
		options,
	)
	if err != nil {
		logger.Log.ErrorLog("$4При инициализации$ $5модели обработки звука$ $4произошла ошибка$"+err.Error(), 4)
		return nil, err
	}

	vad.session = session
	return vad, nil
}

func (v *Detector) Process(chunk *[]float32) bool {
	copy(v.inputData, *chunk)
	v.session.Run()

	copy(v.hData, v.hnData)
	copy(v.cData, v.cnData)

	return v.probData[0] > v.threshold
}

func (v *Detector) Reset() {
	for i := range v.hData {
		v.hData[i] = 0
		v.cData[i] = 0
	}
}

func (v *Detector) Close() {
	v.session.Destroy()
}
