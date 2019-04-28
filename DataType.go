package homie_go

type DataType string

const (
	DataTypeInteger DataType = "integer"
	DataTypeFloat   DataType = "float"
	DataTypeBoolean DataType = "boolean"
	DataTypeString  DataType = "string"
	DataTypeEnum    DataType = "enum"
	DataTypeColor   DataType = "color"
)
