package homie_go

import (
	"strconv"
	"strings"
)

type OnValueChanged func(name, value string) bool

type Property struct {
	id                   string
	name                 string
	valueChangedCallback OnValueChanged
	retained             bool
	unit                 string
	dataType             DataType
	format               string
	node                 *Node
	homie                *Homie
}

func (property *Property) getPropertySubTopic(subTopic string) string {
	if len(subTopic) == 0 {
		return property.node.getNodeSubTopic(property.id)
	}
	return property.node.getNodeSubTopic(strings.Join([]string{property.id, subTopic}, "/"))
}

func (property *Property) publishAttributes() {
	property.homie.sendValue(property.getPropertySubTopic("$name"), property.name)
	property.homie.sendValue(property.getPropertySubTopic("$settable"), strconv.FormatBool(property.valueChangedCallback != nil))
	property.homie.sendValue(property.getPropertySubTopic("$unit"), property.unit)
	property.homie.sendValue(property.getPropertySubTopic("$datatype"), string(property.dataType))
	property.homie.sendValue(property.getPropertySubTopic("$retained"), strconv.FormatBool(property.retained))
	if len(property.format) > 0 {
		property.homie.sendValue(property.getPropertySubTopic("$format"), property.format)
	}
	if property.valueChangedCallback != nil {
		property.homie.registerCallback(property)
	}
}

func (property *Property) SendValue(value string) {
	property.homie.sendValueWithRetained(property.getPropertySubTopic(""), value, property.retained)
}
