package ocrmconfigs

import (
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"log"
	"os"
	"strconv"
)

func ParseIntEnvConfig(envName string) (int, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return 0, ocrmerrors.ErrParseConfig
	}
	intConf, err := strconv.Atoi(strConf)
	if err != nil {
		log.Printf("err when parse int env variable, err: %v, env parameter: %s", err, envName)
		return 0, ocrmerrors.ErrParseConfig
	}
	if intConf < 0 {
		log.Printf("err when parse env variable, int parameter less than 0, env parameter: %s", envName)
		return 0, ocrmerrors.ErrParseConfig
	}
	if intConf == 0 {
		log.Printf("err when parse env variable, int parameter equal 0, env parameter: %s", envName)
		return 0, ocrmerrors.ErrParseConfigEquals0
	}
	return intConf, nil
}

func ParseStrEnvConfig(envName string) (string, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return "", ocrmerrors.ErrParseConfig
	}
	return strConf, nil
}

func ParseBoolEnvConfig(envName string) (bool, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return false, ocrmerrors.ErrParseConfig
	}
	bolConf, err := strconv.ParseBool(strConf)
	if err != nil {
		log.Printf("err when parse bool env variable, err: %v, env parameter: %s", err, envName)
		return false, err
	}
	return bolConf, nil
}
