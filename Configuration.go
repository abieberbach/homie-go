package homie_go

import "net/url"

type Configuration struct {
	name               string
	deviceID           string
	brokerUrl          *url.URL
	baseTopic          string
	statisticsInterval int
	disconnectRetry    int
}

func DefaultConfiguration(name string, deviceID string, brokerUrl *url.URL) *Configuration {
	return DefaultConfigurationWithCustomTopic(name, deviceID, brokerUrl, "homie/")
}

func DefaultConfigurationWithCustomTopic(name string, deviceId string, brokerUrl *url.URL, baseTopic string) *Configuration {
	return NewConfiguration(name, deviceId, brokerUrl, baseTopic, 2000, 10000)
}

func NewConfiguration(name string, deviceId string, brokerUrl *url.URL, baseTopic string, statisticsInterval, disconnectRetry int) *Configuration {
	return &Configuration{
		name:               name,
		deviceID:           deviceId,
		brokerUrl:          brokerUrl,
		baseTopic:          baseTopic,
		statisticsInterval: statisticsInterval,
		disconnectRetry:    disconnectRetry,
	}
}
