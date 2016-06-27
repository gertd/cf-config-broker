package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gertd/cf-config-broker/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

var configFile = flag.String("config", "", "Location of the Service Broker config json file")
var brokerConfig *config.Config

var logger = lager.NewLogger("cf-config-broker")

func main() {

	var err error

	brokerConfig, err = config.LoadFromFile(*configFile)
	if err != nil {
		panic(fmt.Errorf("configuration load error from file %s. Err: %s", *configFile, err))
	}

	serviceBroker := &serviceBroker{}

	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	logger.Debug("config-load-success", lager.Data{"file-source": *configFile, "config": brokerConfig})

	addr := getListeningAddr(brokerConfig)

	brokerAPI := brokerapi.New(serviceBroker, logger, brokerConfig.Credentials)

	http.Handle("/", brokerAPI)

	http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("error-listenting", err)
	}
}

func getListeningAddr(config *config.Config) string {

	envPort := os.Getenv("PORT")
	if envPort == "" {
		if len(config.ListeningAddr) == 0 {
			return ":3000"
		}
		return config.ListeningAddr
	}
	return ":" + envPort
}

func getLogLevel(config *config.Config) lager.LogLevel {
	var minLogLevel lager.LogLevel
	switch config.LogLevel {
	case DEBUG:
		minLogLevel = lager.DEBUG
	case INFO:
		minLogLevel = lager.INFO
	case ERROR:
		minLogLevel = lager.ERROR
	case FATAL:
		minLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("invalid log level: %s", config.LogLevel))
	}

	return minLogLevel
}
