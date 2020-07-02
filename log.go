package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func logInit() {
	moduleDirectory := []string{"experiment/"}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	consoleWriter := zerolog.ConsoleWriter{
		Out:         os.Stdout,
		NoColor:     true,
		TimeFormat:  time.RFC3339,
		FormatLevel: func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("[%v]", i)) },
		FormatCaller: func(i interface{}) string {
			filePath := i.(string)
			if len(filePath) == 0 {
				return filePath
			}

			for _, dir := range moduleDirectory {
				if strings.Contains(filePath, dir) {
					path := strings.Split(filePath, dir)
					filePath = path[len(path)-1] + " >>"
					break
				}
			}
			
			return filePath
		},
	}

	log.Logger = log.Output(consoleWriter).With().Caller().Logger()
}
