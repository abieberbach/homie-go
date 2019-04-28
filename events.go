package homie_go

import "fmt"

type eventController struct {
	handler map[string]EventHandler
}

func NewEventController() *eventController {
	return &eventController{
		handler: make(map[string]EventHandler),
	}
}

func (controller *eventController) AddEventHandler(eventHandler EventHandler) {
	controller.handler[fmt.Sprint(eventHandler)] = eventHandler
}

func (controller *eventController) RemoveEventHandler(eventHandler EventHandler) {
	delete(controller.handler, fmt.Sprint(eventHandler))
}

func (controller eventController) SendEvent(event Event) {
	for _, handler := range controller.handler {
		handler(event);
	}
}
