package common

import log "github.com/sirupsen/logrus"

func convertVerbositytoLogLevel(verbosity int) log.Level {
	switch verbosity {
	case 0:
		return log.WarnLevel
	case 1:
		return log.InfoLevel
	}

	return log.DebugLevel
}

func SetupLogging(verbosity int) {
	log.SetLevel(convertVerbositytoLogLevel(verbosity))
}
