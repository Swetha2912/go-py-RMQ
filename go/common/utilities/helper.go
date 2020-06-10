package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
)

func RecoverPanic() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r, string(debug.Stack()))
		//PushToSlack(fmt.Sprint(r), "panic", string(debug.Stack()))
	}
}

// ConvertMap -- converts any type to map of interface
func ConvertMap(input interface{}, output *Payload) error {
	var err error
	switch input.(type) {
	case []byte:
		byteArray := input.([]byte)
		err = json.Unmarshal(byteArray, output)
	default:
		bytesArray, err2 := json.Marshal(input)
		err2 = json.Unmarshal(bytesArray, output)
		err = err2
	}
	return err
}

// ExitOnErr exits program when error is thrown in critical modules
func ExitOnErr(err error, msg string){
	if err != nil{
		fmt.Println("............. Major Error - Exiting ...........")
		fmt.Println(msg)
		fmt.Println(err.Error())
		fmt.Println("...............................................")
		os.Exit(3)
	}
}
