package logger

import (
	"fmt"
	"neuro/config"
	"neuro/utils"
	"os"
	"sync"
	"time"
)

type logFile struct {
	file *os.File
	mu   sync.Mutex
}

type Logger struct {
	Error *logFile
	Info  *logFile
	Chat  *logFile
}

type LoggerWrite interface {
	SaveLog(Message string, LogLevel uint8) error
}

func (desc *logFile) SaveLog(Message string, LogLelvel string, time time.Time) error {
	desc.mu.Lock()
	defer desc.mu.Unlock()

	colorText, text := ParseLogString(LogLelvel + Message)
	timeStr := time.Format("2006-01-02 15:04:05")

	fileText := fmt.Sprintf("%s %s\n", timeStr, text)

	fmt.Printf("%s\n", colorText)
	_, err := desc.file.WriteString(fileText)
	return err
}

func (log *Logger) ErrorLog(message string, level uint) { //4
	tag := fmt.Sprintf("$%d[ERROR]$ ", level)
	log.Error.SaveLog(message, tag, time.Now())
}

func (log *Logger) InfoLog(message string, level uint) { //3
	tag := fmt.Sprintf("$%d[INFO]$ ", level)
	log.Info.SaveLog(message, tag, time.Now())
}

func (log *Logger) ChatLog(message string, level uint) { //6
	tag := fmt.Sprintf("$%d[CHAT]$ ", level)
	log.Chat.SaveLog(message, tag, time.Now())
}

var Log *Logger

func OpenLogFile(FilePath string) *logFile {

	File, _ := config.Conf.DataPath(config.Conf.LogsDir, FilePath)
	LogDescriptor, err := os.OpenFile(File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		utils.RedPanic("Произошла ошибка во время получения дескриптора лог файла " + FilePath)
	}

	return &logFile{file: LogDescriptor}
}

func init() {
	Log = &Logger{}

	PathDir, _ := config.Conf.DataPath(config.Conf.LogsDir, config.Conf.LogFiles["Error"])
	err := utils.CreateDir(PathDir)

	if err != nil {
		utils.RedPanic("Произошла ошибка во время проверки и создания папки с логами")
	}
	Log.Error = OpenLogFile(config.Conf.LogFiles["Error"])
	Log.Info = OpenLogFile(config.Conf.LogFiles["Info"])
	Log.Chat = OpenLogFile(config.Conf.LogFiles["Chat"])

}
