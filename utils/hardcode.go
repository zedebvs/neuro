package utils

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m" //голубой
	White  = "\033[37m"
)

var Colors = map[string]string{
	"0":   "\033[90m",
	"1":   "\033[37m",
	"2":   "\033[32m",
	"3":   "\033[33m",
	"4":   "\033[31m",
	"5":   "\033[35m",
	"6":   "\033[36m",
	"255": "\033[0m",
}

func RedPanic(message string) {
	panic(Red + message + Reset)
}
