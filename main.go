package main

import (
	"flag"
	"fmt"
	"github.com/gertd/cf-config-broker/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"net/http"
	"os"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

type bindingCredentials struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var configFile = flag.String("config", "", "Location of the Service Broker config json file")
var brokerConfig *config.Config

var logger = lager.NewLogger("cf-config-broker")

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

type myServiceBroker struct{}

// catalog request
func (*myServiceBroker) Services() []brokerapi.Service {

	logger.Info("services-called")

	return brokerConfig.ServiceCatalog
}

// provision new service instance
func (*myServiceBroker) Provision(instanceID string, serviceDetails brokerapi.ServiceDetails) error {

	logger.Info("provision-called", lager.Data{"instanceId": instanceID, "serviceDetails": serviceDetails})

	if brokerConfig.ServiceInstances == nil {
		brokerConfig.ServiceInstances = make(map[string]brokerapi.ServiceDetails)
	}

	brokerConfig.ServiceInstances[instanceID] = serviceDetails

	return nil
}

// deprovision service instance
func (*myServiceBroker) Deprovision(instanceID string) error {

	logger.Info("deprovision-called", lager.Data{"instanceId": instanceID})

	delete(brokerConfig.ServiceInstances, instanceID)

	return nil
}

// bind service, return the "credentials payload" as the response
func (*myServiceBroker) Bind(instanceID, bindingID string) (interface{}, error) {

	logger.Info("bind-called", lager.Data{"instanceId": instanceID, "bindingId": bindingID})

	serviceDetails := brokerConfig.ServiceInstances[instanceID]

	//logger.Info("bind-called", lager.Data{"brokerConfig.ServiceInstances": brokerConfig.ServiceInstances})
	//logger.Info("bind-called", lager.Data{"serviceDetails": serviceDetails})

	result := brokerConfig.BindingCredentials[serviceDetails.ID]

	//logger.Info("bind-called", lager.Data{"brokerConfig.BindingCredentials": brokerConfig.BindingCredentials})
	//logger.Info("bind-called", lager.Data{"result": result})

	logger.Info("bind-result", lager.Data{"result": result})

	return result, nil
}

// unbind the service from the application
func (*myServiceBroker) Unbind(instanceID, bindingID string) error {

	logger.Info("unbind-called", lager.Data{"instanceId": instanceID, "bindingId": bindingID})

	return nil
}

func main() {

	var err error

	brokerConfig, err = config.LoadFromFile(*configFile)
	if err != nil {
		panic(fmt.Errorf("configuration load error from file %s. Err: %s", *configFile, err))
	}

	serviceBroker := &myServiceBroker{}

	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	logger.Debug("config-load-success", lager.Data{"file-source": *configFile, "config": brokerConfig})

	credentials := brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}

	addr := getListeningAddr(brokerConfig)

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)

	http.Handle("/", brokerAPI)

	http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("error-listenting", err)
	}
}
