package utilities

import (
	"sample_rmq/common/interfaces"
	"sample_rmq/common/packages/govalidator"
	// "sample_rmq/gateway/records"
)

type GoValidator struct {
	interfaces.IValidator
	recordsMapper map[string]govalidator.IValidationDBChecker
}

func (validator GoValidator) Validate(rules map[string][]string, body map[string]interface{}, messages map[string][]string) map[string][]string {
	opts := govalidator.Options{
		Body:            body,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
		Records:         validator.recordsMapper,
	}
	goVal := govalidator.New(opts)
	errs := goVal.Validate()

	if errs != nil && len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}

func (validator GoValidator) ValidateAndReply(rules map[string][]string, req interfaces.IRouterRequest) bool {
	body := req.GetBody()

	opts := govalidator.Options{
		Body:            body,
		Rules:           rules,
		Messages:        nil,
		RequiredDefault: false,
		Records:         validator.recordsMapper,
	}
	goVal := govalidator.New(opts)
	errs := goVal.Validate()

	if errs != nil && len(errs) > 0 {
		req.ReplyBack(400, interfaces.RouterResponse{Errors: errs})
		return false
	} else {
		return true
	}
}

//ProvideGoValidator returns a GoValidator
func ProvideGoValidator() GoValidator {
	recordsMap := make(map[string]govalidator.IValidationDBChecker)
	// recordsMap["user"] = userRecord
	// recordsMap["camera"] = cameraRecord
	return GoValidator{recordsMapper: recordsMap}
}
