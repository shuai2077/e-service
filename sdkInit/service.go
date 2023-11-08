package sdkInit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *Application) Get(args []string) ([]byte, error) {
	response, err := t.SdkEnvInfo.ChClient.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return response.Payload, nil
}

func (t *Application) Set(args []string) (string, error) {
	var tempArgs [][]byte
	for i := 1; i < len(args); i++ {
		tempArgs = append(tempArgs, []byte(args[i]))
	}

	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// 资产转移失败
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}

func (t *Application) Del(args []string) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// 资产转移失败
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}

func (t *Application) His(args []string) ([]byte, error) {
	response, err := t.SdkEnvInfo.ChClient.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return response.Payload, nil
}
