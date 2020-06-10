package controllers

import (
	"encoding/json"
	"sample_rmq/common/infrastructures"
	"sample_rmq/common/interfaces"
	"sample_rmq/common/utilities"
)

type InferenceController struct {
	validator    interfaces.IValidator
	pythonBroker interfaces.IRouter
}

func (controller InferenceController) Hello(req interfaces.IRouterRequest) {
	jsonBytes, _ := json.Marshal("hello")
	_, replyBytes := controller.pythonBroker.Invoke("sample.data", jsonBytes)
	req.ReplyBack(200, interfaces.RouterResponse{Data: string(replyBytes)})
}

// ProvideInferenceController returns a InferenceController
func ProvideInferenceController(validator utilities.GoValidator, pythonBroker *infrastructures.RMQRouter) InferenceController {
	return InferenceController{validator: validator, pythonBroker: pythonBroker}
}
