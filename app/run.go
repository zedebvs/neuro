package main

import (
	"neuro/audio/micro"
	"neuro/audio/voice"
	"neuro/config"
	"neuro/core"
	"neuro/lim/stt"
	"neuro/lim/vad"
	"os"

	//"neuro/config"
	//embmodel "neuro/lim/EmbModel"
	"neuro/logger"
	"neuro/utils"

	ort "github.com/yalue/onnxruntime_go"
)

var con = config.Conf

func loadModel(name string) string {
	path, _ := config.Conf.DataPath(con.ModelsDir, con.Models[name])
	return path
}

func ChunkCalc(mc int) int {
	che := mc % 32
	if che != 0 {
		return mc/32 + 1
	}
	return mc / 32
}

func RunApp() func() {
	logger.Run()

	path, _ := utils.GetPath([]string{con.Data, con.LibsDir, "onnxruntime.dll"})
	ort.SetSharedLibraryPath(path)
	err := ort.InitializeEnvironment()
	if err != nil {
		utils.RedPanic("Не удалось запустить движок нейросетей: " + path + "\nОшибка: " + err.Error())
	}

	libsPath, _ := utils.GetPath([]string{con.Data, con.LibsDir})
	err = os.Setenv("PATH", libsPath+";"+os.Getenv("PATH"))

	if err != nil {
		utils.RedPanic("Не удалось запустить движок нейросети wisper: " + path + "\nОшибка: " + err.Error())
	}

	logger.Log.InfoLog("$2Движки нейросетей$ $6успешно подняты$", 3)

	voice, err := voice.New()
	if err != nil {
		utils.RedPanic("Ошибка при инициализации аудио\n")
	}

	micro, _ := micro.New()

	vad, err := vad.New(loadModel("vad"), con.Threshold)
	if err != nil {
		utils.RedPanic("Ошибка при инициализации модели обработки\n")
	}

	stt, err := stt.New(loadModel("stt_large"))
	if err != nil {
		utils.RedPanic("Ошибка при инициализации модели преобразования\n")
	}

	disp := core.New(micro, voice, vad, stt, ChunkCalc(con.StartTolk), ChunkCalc(con.EndTolk), con.PauseTime)

	return func() {
		disp.Shutdown()
		ort.DestroyEnvironment()
		voice.Shutdown()
		micro.Shutdown()
		stt.Shutdown()
		logger.Log.Shutdown()
	}
}
