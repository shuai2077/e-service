package main

import (
	"e-service/sdkInit"
	"e-service/server"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
)

const (
	cc_name      = "simplecc"
	cc_version   = "1.0.0"
	user_message = "./db"
)

var App sdkInit.Application

// 匹配http请求信息
type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// init orgs information

	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/e-service/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User2",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/e-service/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/e-service/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH") + "/src/e-service/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {

		fmt.Println("InitService successful")
		os.Exit(-1)
	}

	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	defer info.EvClient.Unregister(sdkInit.BlockListener(info.EvClient))
	defer info.EvClient.Unregister(sdkInit.ChainCodeEventListener(info.EvClient, info.ChaincodeID))

	opts := badger.DefaultOptions(user_message)
	opts.Logger = nil

	// set admin and user role
	db, err := badger.Open(opts)
	if err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("admin"), []byte("admin"))
		return err
	}); err != nil {
		fmt.Println("failed to store admin message")
	}

	if err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("user"), []byte("user"))
		return err
	}); err != nil {
		fmt.Println("failed to store user message")
	}

	db.Close()

	server.WebStart(App)
}
