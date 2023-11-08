package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	indexName = "objectType~key"
	RecordKey = "recordKey"
	//UserKey   = "userKey"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

type outputEvent struct {
	EventName string
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("init...")
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result [][]byte
	var err error
	switch fn {
	case "set":
		err = set(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "get":
		result, err = get(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "del":
		err = del(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "his":
		result, err = his(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}
	default:
		return shim.Error(fmt.Sprintf("no dunction can be found"))
	}

	// Return the result as success payload
	// 特殊字符进行分隔
	return shim.Success(bytes.Join(result, []byte("U+10FFFF")))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	// 创建复合主键
	val, err := stub.CreateCompositeKey(indexName, []string{RecordKey, args[0]})
	if err != nil {
		return fmt.Errorf("Failed to create compositekey: %s", args[0])
	}

	err = stub.PutState(val, []byte(args[1]))
	if err != nil {
		return fmt.Errorf("Failed to set asset: %s", args[0])
	}
	event := outputEvent{
		EventName: "set",
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = stub.SetEvent("chaincode-event", payload)
	return nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) ([][]byte, error) {
	var records [][]byte
	valIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{RecordKey})
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer valIterator.Close()

	for valIterator.HasNext() {
		queryResponse, err := valIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("faild to get next query result: %v", err)
		}

		records = append(records, queryResponse.Value)
	}

	event := outputEvent{
		EventName: "get",
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	err = stub.SetEvent("chaincode-event", payload)
	return records, nil
}

func del(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	// 创建复合主键
	val, err := stub.CreateCompositeKey(indexName, []string{RecordKey, args[0]})
	if err != nil {
		return fmt.Errorf("Failed to create compositekey: %s", args[0])
	}

	//err := stub.DelState(args[0])
	err = stub.DelState(val)
	if err != nil {
		return fmt.Errorf("Failed to del asset: %s with error: %s", args[0], err)
	}

	event := outputEvent{
		EventName: "del",
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = stub.SetEvent("chaincode-event", payload)
	return nil
}

func his(stub shim.ChaincodeStubInterface, args []string) ([][]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	val, err := stub.CreateCompositeKey(indexName, []string{RecordKey, args[0]})
	if err != nil {
		return nil, fmt.Errorf("Failed to create compositekey: %s", args[0])
	}
	historyIter, err := stub.GetHistoryForKey(val)
	if err != nil {
		return nil, err
	}

	var history [][]byte
	for historyIter.HasNext() {
		modification, err := historyIter.Next()
		if err != nil {
			return nil, err
		}

		modificationBytes, err := json.Marshal(modification)
		if err != nil {
			return nil, err
		}
		history = append(history, modificationBytes)
	}
	return history, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
