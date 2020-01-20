package ocrmconfigs

import "github.com/sirupsen/logrus"

func ConfigureLogger(level string, dateFormat string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	if dateFormat != "" {
		customFormatter := &logrus.TextFormatter{}
		customFormatter.TimestampFormat = dateFormat //"02.01.2006 15:04:05"//"dd.mm.yyyy HH24:MI:SS"
		customFormatter.FullTimestamp = true
		logger.SetFormatter(customFormatter)
	}
	return logger, nil
}
