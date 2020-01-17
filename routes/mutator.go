package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ExpediaDotCom/kubernetes-sidecar-injector/webhook"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
)

/*SideCars is an array of named SideCar instances*/
type SideCars struct {
	Sidecars []SideCar `yaml:"sidecars"`
}

/*SideCar is a named sidecar to be injected*/
type SideCar struct {
	Name    string          `yaml:"name"`
	Sidecar webhook.SideCar `yaml:"sidecar"`
}

func loadConfig(sideCarConfigFile string) (map[string]*webhook.SideCar, error) {
	data, err := ioutil.ReadFile(sideCarConfigFile)
	if err != nil {
		return nil, err
	}
	glog.Infof("New sideCar configuration: %s", data)

	var cfg SideCars
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	mapOfSideCar := make(map[string]*webhook.SideCar)
	for _, configuration := range cfg.Sidecars {
		mapOfSideCar[configuration.Name] = &configuration.Sidecar
	}

	return mapOfSideCar, nil
}

/*MutatorController is an interface that implements mutation method*/
type MutatorController interface {
	Mutate(http.ResponseWriter, *http.Request)
}

/*NewMutatorController is a factory method to create an instance of MutatorController*/
func NewMutatorController(sideCarConfigFile string) (MutatorController, error) {
	mapOfSideCars, err := loadConfig(sideCarConfigFile)
	if mapOfSideCars != nil {
		return mutatorController{mutator: webhook.Mutator{SideCars: mapOfSideCars}}, nil
	}
	return nil, err
}

type mutatorController struct {
	mutator webhook.Mutator
}

func (controller mutatorController) Mutate(writer http.ResponseWriter, request *http.Request) {
	body, err := readRequestBody(request)
	if err != nil {
		writeError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := controller.mutator.Mutate(body)
	if err != nil {
		writeError(writer, fmt.Sprintf("Failed to process request: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err := writer.Write(resp); err != nil {
		writeError(writer, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}
