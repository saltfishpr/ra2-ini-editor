package ra2

import (
	"gopkg.in/ini.v1"
)

type I18NString map[string]string

type Property struct {
	Key     string `json:"key"`     // 属性键
	Value   string `json:"value"`   // 属性值
	Comment string `json:"comment"` // 属性注释

	Name string     `json:"name"` // 属性名称
	Desc I18NString `json:"desc"` // 属性描述
}

func parseProperties(sec *ini.Section) []Property {
	var properties []Property
	for _, key := range sec.Keys() {
		properties = append(properties, Property{
			Key:     key.Name(),
			Value:   key.Value(),
			Comment: key.Comment,
		})
	}
	return properties
}
