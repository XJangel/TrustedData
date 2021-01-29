package data

import (
	"github.com/XJangel/TrustedData/pkg/common"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
	"k8s.io/klog"
	"os"
)

//type MqttMode int

// 从edgecore配置文件中获取这些信息
type mqtt struct {
	MqttServerInternal string
	MqttServerExternal string
	// 0: internal mqtt broker enable only.
	// 1: internal and external mqtt broker enable.
	// 2: external mqtt broker enable only
	MqttMode common.MqttMode `json:"mqttMode"`
}

func (m *mqtt) Start() {
	if m.MqttMode >= common.MqttModeBoth {
		hub := &mqttBus.Client{
			MQTTUrl: eventconfig.Config.MqttServerExternal,
		}
		mqttBus.MQTTHub = hub
		hub.InitSubClient()
		hub.InitPubClient()
		klog.Infof("Init Sub And Pub Client for externel mqtt broker %v successfully", eventconfig.Config.MqttServerExternal)
	}

	if eventconfig.Config.MqttMode <= v1alpha1.MqttModeBoth {
		// launch an internal mqtt server only
		mqttServer = mqttBus.NewMqttServer(
			int(eventconfig.Config.MqttSessionQueueSize),
			eventconfig.Config.MqttServerInternal,
			eventconfig.Config.MqttRetain,
			int(eventconfig.Config.MqttQOS))
		mqttServer.InitInternalTopics()
		err := mqttServer.Run()
		if err != nil {
			klog.Errorf("Launch internel mqtt broker failed, %s", err.Error())
			os.Exit(1)
		}
		klog.Infof("Launch internel mqtt broker %v successfully", eventconfig.Config.MqttServerInternal)
	}

	eb.pubCloudMsgToEdge()
}
