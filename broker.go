package main

import (
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type serviceBroker struct{}

// service catalog request
func (*serviceBroker) Services() []brokerapi.Service {

	logger.Info("services-called")

	return brokerConfig.ServiceCatalog
}

// provision new service instance
func (*serviceBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	logger.Info("provision-called", lager.Data{"instanceId": instanceID, "details": details, "asyncAllowed": asyncAllowed})

	if brokerConfig.ServiceInstances == nil {
		brokerConfig.ServiceInstances = make(map[string]brokerapi.ProvisionDetails)
	}

	brokerConfig.ServiceInstances[instanceID] = details

	spec := brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: ""}

	return spec, nil
}

// deprovision service instance
func (*serviceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {

	logger.Info("deprovision-called", lager.Data{"instanceId": instanceID, "details": details, "asyncAllowed": asyncAllowed})

	delete(brokerConfig.ServiceInstances, instanceID)

	return false, nil
}

// bind service, return the "credentials payload" as the response
func (*serviceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	logger.Info("bind-called", lager.Data{"instanceId": instanceID, "bindingId": bindingID, "details": details})

	provisionDetails := brokerConfig.ServiceInstances[instanceID]

	logger.Info("bind-called", lager.Data{"brokerConfig.ServiceInstances": brokerConfig.ServiceInstances})
	logger.Info("bind-called", lager.Data{"provisionDetails": provisionDetails})
	logger.Info("bind-called", lager.Data{"brokerConfig.BindingCredentials": brokerConfig.BindingCredentials})

	binding := brokerapi.Binding{
		Credentials:     brokerConfig.BindingCredentials[provisionDetails.ServiceID],
		SyslogDrainURL:  "",
		RouteServiceURL: "",
	}

	logger.Info("bind-return", lager.Data{"binding": binding})

	return binding, nil
}

// unbind the service from the application
func (*serviceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {

	logger.Info("unbind-called", lager.Data{"instanceId": instanceID, "bindingId": bindingID, "details": details})

	return nil
}

// update service instance
func (*serviceBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {

	logger.Info("update-called", lager.Data{"instanceId": instanceID, "details": details, "asyncAllowed": asyncAllowed})

	return false, nil
}

// lastoperation
func (*serviceBroker) LastOperation(instanceID string) (brokerapi.LastOperation, error) {

	logger.Info("lastOperation-called", lager.Data{"instanceId": instanceID})

	lastOperation := brokerapi.LastOperation{
		State:       brokerapi.Succeeded,
		Description: "done",
	}

	logger.Info("lastOperation-return", lager.Data{"lastOperation": lastOperation})

	return lastOperation, nil
}
