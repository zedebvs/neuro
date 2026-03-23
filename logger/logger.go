package logger

import (
	"fmt"
	"neuro/config"
	"neuro/utils"
	"os"
	"sync"
	"time"
)

type LogMessage struct {
	consoleText string
	fileText    string
}

type logFile struct {
	file *os.File
	ch   chan LogMessage
	wg   sync.WaitGroup
}

type Logger struct {
	Error *logFile
	Info  *logFile
	Chat  *logFile
}

type LoggerWrite interface {
	SaveLog(Message string, LogLevel uint8) error
}

func (desc *logFile) worker() {
	defer desc.wg.Done()

	for msg := range desc.ch {
		fmt.Printf("%s\n", msg.consoleText)
		_, err := desc.file.WriteString(msg.fileText)
		if err != nil {
			colorText, _ := ParseLogString("$4Во время записи лога произошла ошибка$")
			fmt.Printf("%s\n", colorText)
		}
	}
}

func (desc *logFile) saveLog(Message string, LogLevel string, time time.Time) {

	timeConsole := fmt.Sprintf("$0 %s$", time.Format("15:04:05"))

	colorText, text := ParseLogString(fmt.Sprintf("%s %s %s", timeConsole, LogLevel, Message))

	timeStr := time.Format("2006-01-02 15:04:05")

	fileText := fmt.Sprintf("%s %s\n", timeStr, text)
	desc.ch <- LogMessage{consoleText: colorText, fileText: fileText}
}

func (desc *logFile) shutdown() {
	close(desc.ch)
	desc.wg.Wait()
	desc.file.Close()
}

func (desc *Logger) Shutdown() {
	desc.Error.shutdown()
	desc.Chat.shutdown()
	desc.Info.shutdown()
}

func (log *Logger) ErrorLog(message string, level uint) { //4
	tag := fmt.Sprintf("$%d[ERROR]$", level)
	log.Error.saveLog(message, tag, time.Now())
}

func (log *Logger) InfoLog(message string, level uint) { //3
	tag := fmt.Sprintf("$%d[INFO]$", level)
	log.Info.saveLog(message, tag, time.Now())
}

func (log *Logger) ChatLog(message string, level uint) { //6
	tag := fmt.Sprintf("$%d[CHAT]$", level)
	log.Chat.saveLog(message, tag, time.Now())
}

var Log *Logger

func openLogFile(FilePath string) *logFile {

	File, _ := config.Conf.DataPath(config.Conf.LogsDir, FilePath)
	LogDescriptor, err := os.OpenFile(File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		utils.RedPanic("Произошла ошибка во время получения дескриптора лог файла " + FilePath)
	}

	l := &logFile{file: LogDescriptor, ch: make(chan LogMessage, 100)}

	l.wg.Add(1)
	go l.worker()

	return l
}

func Run() {
	Log = &Logger{}

	PathDir, err := config.Conf.DataPath(config.Conf.LogsDir, config.Conf.LogFiles["Error"])
	if err != nil {
		utils.RedPanic("Произошла ошибка во время проверки и создания папки с логами")
	}

	err = utils.CreateDir(PathDir)

	if err != nil {
		utils.RedPanic("Произошла ошибка во время проверки и создания папки с логами")
	}
	Log.Error = openLogFile(config.Conf.LogFiles["Error"])
	Log.Info = openLogFile(config.Conf.LogFiles["Info"])
	Log.Chat = openLogFile(config.Conf.LogFiles["Chat"])

}
