package homie_go

type DeviceState string

const (
	DeviceStateInit         DeviceState = "init"
	DeviceStateReady        DeviceState = "ready"
	DeviceStateDisconnected DeviceState = "disconnected"
	DeviceStateSleeping     DeviceState = "sleeping"
	DeviceStateLost         DeviceState = "lost"
	DeviceStateAlert        DeviceState = "alert"
)
