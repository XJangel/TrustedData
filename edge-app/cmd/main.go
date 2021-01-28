package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	ccID           = "mycc"
)

const (
	fcn_query  = "query"
	fcn_invoke = "invoke"
)

func main() {

	//读取配置文件，创建SDK
	configProvider := config.FromFile("./config.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}
	defer sdk.Close()

	//读取配置文件(config.yaml)中的组织(member1.example.com)的用户(Admin)
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Fatalf("create msp client fail: %s\n", err.Error())
	}

	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		log.Fatalf("get admin identify fail: %s\n", err.Error())
	} else {
		fmt.Println("AdminIdentify is found:")
		fmt.Println(adminIdentity)
	}

	//初始化channel信息
	channelProvider := sdk.ChannelContext(channelID, fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Fatalf("create channel client fail: %s\n", err.Error())
	}

	//操作chaincode
	querycc(channelClient, "a")
	executeCC(channelClient, "a.b.c")
}

func querycc(channelClient *channel.Client, key string) []byte {
	//var args [][]byte
	//args = append(args, []byte("key1"))

	request := channel.Request{
		ChaincodeID: ccID,
		Fcn:         fcn_query,
		Args:        makeArgs(key),
	}
	// the proposal responses from peer(s) 所以这个返回值是什么？？
	// 查询的worldState??
	response, err := channelClient.Query(request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
	} else {
		fmt.Printf("the tx ID is: %s\n", response.TransactionID)
		fmt.Printf("response is %s\n", response.Payload)
	}
	return response.Payload
}

func executeCC(client *channel.Client, who string) error {
	request := channel.Request{
		ChaincodeID: ccID,
		Fcn:         fcn_invoke,
		Args:        makeArgs(who),
	}
	res, err := client.Execute(request, channel.WithRetry(retry.DefaultChannelOpts))
	fmt.Println("exe tx:", res.TransactionID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return err
}

func makeArgs(args ...string) [][]byte {
	var ccargs [][]byte
	for _, arg := range args {
		ccargs = append(ccargs, []byte(arg))
	}
	return ccargs
}
