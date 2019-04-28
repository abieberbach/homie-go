package homie_go

import (
	"strings"
)

type Node struct {
	id         string
	name       string
	nodeType   string
	homie      *Homie
	properties map[string]*Property
}

func (node *Node) AdvertiseProperty(id string, name string, dataType DataType, unit string, retained bool) *Property {
	return node.AdvertisePropertyWithFormat(id, name, dataType, unit, "", retained)
}
func (node *Node) AdvertisePropertyWithFormat(id string, name string, dataType DataType, unit string, format string, retained bool) *Property {
	return node.AdvertiseChangablePropertyWithFormat(id, name, dataType, unit, format, nil, retained)
}

func (node *Node) AdvertiseChangableProperty(id string, name string, dataType DataType, unit string, valueChangedHandler OnValueChanged, retained bool) *Property {
	return node.AdvertiseChangablePropertyWithFormat(id, name, dataType, unit, "", valueChangedHandler, retained)
}

func (node *Node) AdvertiseChangablePropertyWithFormat(id string, name string, dataType DataType, unit string, format string, valueChangedHandler OnValueChanged, retained bool) *Property {
	property := &Property{
		id:                   id,
		name:                 name,
		dataType:             dataType,
		unit:                 unit,
		retained:             retained,
		format:               format,
		valueChangedCallback: valueChangedHandler,
		node:                 node,
		homie:                node.homie,
	}
	node.properties[id] = property
	return property
}

func (node *Node) getNodeSubTopic(subTopic string) string {
	return strings.Join([]string{node.id, subTopic}, "/")
}

func (node *Node) publishProperties() {
	node.homie.sendValue(node.getNodeSubTopic("$name"), node.name)
	node.homie.sendValue(node.getNodeSubTopic("$type"), node.nodeType)
	propertyIds := make([]string, 0, len(node.properties))
	for propertyId, property := range node.properties {
		propertyIds = append(propertyIds, propertyId)
		property.publishAttributes()
	}
	node.homie.sendValue(node.getNodeSubTopic("$properties"), strings.Join(propertyIds, ","))
}
