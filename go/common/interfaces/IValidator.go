package interfaces

type IValidator interface{
	Validate(rules map[string][]string, body map[string]interface{}, messages map[string][]string) (map[string][]string)
	ValidateAndReply(rules map[string][]string, req IRouterRequest) bool
}
