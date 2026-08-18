package rdap

import (
	"encoding/json"
	"strings"
)

type Entity struct {
	ObjectClassName string          `json:"objectClassName"`
	Handle          string          `json:"handle,omitempty"`
	Roles           []string        `json:"roles,omitempty"`
	Entities        []Entity        `json:"entities,omitempty"`
	VCardArray      json.RawMessage `json:"vcardArray,omitempty"`
}

func (e Entity) HasRole(role string) bool {
	for _, candidate := range e.Roles {
		if strings.EqualFold(candidate, role) {
			return true
		}
	}

	return false
}

func (e Entity) Name() string {
	properties, ok := e.vcardProperties()
	if !ok {
		return ""
	}

	for _, property := range properties {
		if len(property) < 4 {
			continue
		}

		name, ok := property[0].(string)
		if !ok || !strings.EqualFold(name, "fn") {
			continue
		}

		valueType, ok := property[2].(string)
		if !ok || !strings.EqualFold(valueType, "text") {
			continue
		}

		value, ok := property[3].(string)
		if ok {
			return value
		}
	}

	return ""
}

func (e Entity) Organization() string {
	properties, ok := e.vcardProperties()
	if !ok {
		return ""
	}

	for _, property := range properties {
		if len(property) < 4 {
			continue
		}

		name, ok := property[0].(string)
		if !ok || !strings.EqualFold(name, "org") {
			continue
		}

		valueType, ok := property[2].(string)
		if !ok || !strings.EqualFold(valueType, "text") {
			continue
		}

		switch value := property[3].(type) {
		case string:
			return value

		case []any:
			var parts []string

			for _, item := range value {
				if s, ok := item.(string); ok && s != "" {
					parts = append(parts, s)
				}
			}

			return strings.Join(parts, " ")
		}
	}

	return ""
}

func (e Entity) vcardProperties() ([][]any, bool) {
	if len(e.VCardArray) == 0 {
		return nil, false
	}

	var card []any

	if err := json.Unmarshal(e.VCardArray, &card); err != nil {
		return nil, false
	}

	if len(card) != 2 {
		return nil, false
	}

	kind, ok := card[0].(string)
	if !ok || !strings.EqualFold(kind, "vcard") {
		return nil, false
	}

	rawProperties, ok := card[1].([]any)
	if !ok {
		return nil, false
	}

	properties := make([][]any, 0, len(rawProperties))

	for _, raw := range rawProperties {
		property, ok := raw.([]any)
		if !ok {
			continue
		}

		properties = append(properties, property)
	}

	return properties, true
}
