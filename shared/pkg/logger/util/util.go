package util

import (
	"fmt"
	"log/syslog"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
)

func GetSysLogger(c config.LoggerOutputConfig, prefix string) (w *syslog.Writer, err error) {
	if c.Type != "syslog" || c.Facility > 7 || c.Facility < 0 {
		err = fmt.Errorf("invalid syslog config")

		return nil, err
	}

	w, err = syslog.New(GetPriority(c.Facility), prefix)

	if err != nil {
		w = nil
	}

	return w, err
}

func GetPriority(facility int) syslog.Priority {
	switch facility {
	case 0:
		return syslog.LOG_LOCAL0
	case 1:
		return syslog.LOG_LOCAL1
	case 2:
		return syslog.LOG_LOCAL2
	case 3:
		return syslog.LOG_LOCAL3
	case 4:
		return syslog.LOG_LOCAL4
	case 5:
		return syslog.LOG_LOCAL5
	case 6:
		return syslog.LOG_LOCAL6
	case 7:
		return syslog.LOG_LOCAL7
	}

	return syslog.LOG_LOCAL0
}
