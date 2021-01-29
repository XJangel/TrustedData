package fabric

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//定义结构体
type FabricSetup struct {
	ConfigFile    string //sdk配置文件所在路径
	ChannelID     string //应用通道名称
	ChannelConfig string //应用通道交易配置文件所在路径
	OrgAdmin      string // 组织管理员名称
	OrgName       string //组织名称
	Initialized   bool   //是否初始化
	//Admin         *resmgmt.Client   			   //fabric环境中资源管理者
	ccClient channel.Client
	SDK      *fabsdk.FabricSDK //SDK实例
}

func (f *FabricSetup) Initialize() (*channel.Client, error) {
	//判断是否已经初始化
	if f.Initialized {
		return nil, fmt.Errorf("SDK已被实例化")
	}

	//创建SDK对象
	sdk, err := fabsdk.New(config.FromFile(f.ConfigFile))

	if err != nil {
		return nil, fmt.Errorf("SDK实例化失败:%v", err)
	}
	f.SDK = sdk

	ccp := sdk.ChannelContext(f.ChannelID, fabsdk.WithUser("User1"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}

	f.Initialized = true
	fmt.Println("SDK实例化成功")

	return cc, nil

}
