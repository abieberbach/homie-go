package homie_go

type EventType string

const (
	EventTypeOTAStarted             EventType = "OTAStarted"
	EventTypeOTAProgress            EventType = "OTAProgress"
	EventTypeOTAFailed              EventType = "OTAFailed"
	EventTypeMQTTReady              EventType = "MQTTReady"
	EventTypeMQTTDisconnected       EventType = "MQTTDisconnected"
	EventTypeMQTTPacketAcknowledged EventType = "MQTTPacketAcknowledged"
	EventTypeSendingStatistics      EventType = "SendingStatistics"
)

type EventHandler func(Event)

type Event interface {
	GetEventType() EventType
}

type OTAEvent struct {
	eventType EventType
	SizeDone  int
	SizeTotal int
}

func (event *OTAEvent) GetEventType() EventType {
	return event.eventType
}

type MQTTEvent struct {
	eventType      EventType
	DisonnectError error
}

func (event *MQTTEvent) GetEventType() EventType {
	return event.eventType
}

type StatisticsEvent struct {
	Statistics *Statistics
}

func (event *StatisticsEvent) GetEventType() EventType {
	return EventTypeSendingStatistics
}
