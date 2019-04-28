package homie_go

import (
	"strconv"
	"strings"
	"time"
)
import mqtt "github.com/eclipse/paho.mqtt.golang"

const HomieSpec = "3.0.1"

type Homie struct {
	firmwareName    string
	firmwareVersion string
	currentState    DeviceState
	configuration   *Configuration
	startTime       time.Time
	statisticsTimer *time.Timer
	mqttClient      mqtt.Client
	eventController *eventController
	nodes           map[string]*Node
}

func NewHomie(firmwareName, firmwareVersion string, configuration *Configuration) *Homie {
	return &Homie{
		firmwareName:    firmwareName,
		firmwareVersion: firmwareVersion,
		currentState:    DeviceStateInit,
		configuration:   configuration,
		startTime:       time.Now(),
		eventController: NewEventController(),
		nodes:           make(map[string]*Node),
	}
}

func (homie *Homie) AddNode(nodeId, nodeName, nodeType string) *Node {
	node := &Node{
		id:         nodeId,
		name:       nodeName,
		nodeType:   nodeType,
		homie:      homie,
		properties: make(map[string]*Property),
	}
	homie.nodes[nodeId] = node
	return node
}

func (homie *Homie) Start() {
	homie.connectToMqtt()
}

func (homie *Homie) Stop() {
	homie.statisticsTimer.Stop()
	homie.setState(DeviceStateDisconnected)
	homie.mqttClient.Disconnect(10000)
}

func (homie *Homie) connectToMqtt() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(homie.configuration.brokerUrl.String())
	opts.SetClientID("homie_go_" + homie.configuration.deviceID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Duration(homie.configuration.disconnectRetry) * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		homie.setState(DeviceStateReady)
		homie.sendDeviceAttributes()
		homie.eventController.SendEvent(&MQTTEvent{eventType: EventTypeMQTTReady})
		homie.sendStatistics()
	})
	opts.SetConnectionLostHandler(func(c mqtt.Client, e error) {
		homie.setState(DeviceStateDisconnected)
		homie.eventController.SendEvent(&MQTTEvent{
			eventType:      EventTypeMQTTDisconnected,
			DisonnectError: e,
		})

	})
	opts.SetWill(homie.getTopicName("$state"), string(DeviceStateLost), 1, true)
	homie.mqttClient = mqtt.NewClient(opts)
	homie.mqttClient.Connect().Wait()
}

func (homie *Homie) sendDeviceAttributes() {
	homie.sendValue("$homie", HomieSpec)
	homie.sendValue("$name", homie.configuration.name)
	homie.sendValue("$fw/name", homie.firmwareName)
	homie.sendValue("$fw/version", homie.firmwareVersion)
	ip, mac := getIpAndMacAddr(homie.configuration.brokerUrl)
	homie.sendValue("$localip", ip)
	homie.sendValue("$mac", mac)
	homie.sendValue("$implementation", "go-homie")
	homie.sendIntValue("$stats/interval", homie.configuration.statisticsInterval)
	homie.sendValue("$state", string(homie.currentState))
	homie.publishNodes()
}

func (homie *Homie) publishNodes() {
	nodeIds := make([]string, 0, len(homie.nodes))
	for id, node := range homie.nodes {
		nodeIds = append(nodeIds, id)
		node.publishProperties()
	}
	homie.sendValue("$nodes", strings.Join(nodeIds, ","))
}

func (homie *Homie) sendStatistics() {

	statistics := NewStatistics(homie.startTime)
	homie.sendIntValue("$stats/uptime", int(statistics.uptime.Seconds()))

	homie.eventController.SendEvent(&StatisticsEvent{statistics})
	homie.statisticsTimer = time.AfterFunc(time.Duration(homie.configuration.statisticsInterval)*time.Millisecond, homie.sendStatistics)
}

func (homie *Homie) getTopicName(subTopic string) string {
	return strings.Join([]string{homie.configuration.baseTopic, homie.configuration.deviceID, subTopic}, "/")
}

func (homie *Homie) sendValue(subTopic string, value string) {
	homie.sendValueWithRetained(subTopic, value, true)
}

func (homie *Homie) sendValueWithRetained(subTopic string, value string, retained bool) {
	token := homie.mqttClient.Publish(homie.getTopicName(subTopic), 1, retained, value)
	token.Wait()
	homie.eventController.SendEvent(&MQTTEvent{eventType: EventTypeMQTTPacketAcknowledged})
}

func (homie *Homie) sendIntValue(subTopic string, value int) {
	homie.sendValue(subTopic, strconv.Itoa(value))
}

func (homie *Homie) setState(state DeviceState) {
	homie.currentState = DeviceStateReady
	if homie.mqttClient != nil && homie.mqttClient.IsConnected() {
		homie.sendValue("$state", string(state))
	}
}

func (homie *Homie) registerCallback(property *Property) {
	subscribe := homie.mqttClient.Subscribe(homie.getTopicName(property.getPropertySubTopic("")+"/set"), 1, func(client mqtt.Client, message mqtt.Message) {
		value := string(message.Payload())
		valueAccepted := property.valueChangedCallback(property.name, value)
		message.Ack()
		if valueAccepted {
			property.SendValue(value)
		}
	})
	subscribe.Wait()
}
