package logger

import (
	"neuro/utils"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\$(\d+)(.*?)\$`)

func ParseLogString(message string) (string, string) {

	match := re.FindAllStringSubmatch(message, -1)
	consoleText, fileText := message, message

	for _, mat := range match {
		notFiltred := mat[0]
		colorNum := mat[1]
		filtred := mat[2]

		color := utils.Colors[colorNum]

		colordText := color + filtred + utils.Colors["255"]

		consoleText = strings.Replace(consoleText, notFiltred, colordText, 1)
		fileText = strings.Replace(fileText, notFiltred, filtred, 1)

	}
	return consoleText, fileText
}
