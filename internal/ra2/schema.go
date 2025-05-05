package ra2

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/pkg/errors"
)

type Schema struct {
	Version string    `json:"version"`
	Flags   []IniFlag `json:"flags"`
}

type IniFlag struct {
	Category     string `json:"category"`
	Filename     string `json:"filename"`
	Section      string `json:"section"`
	Key          string `json:"key"`
	ValueType    string `json:"value_type"`
	DefaultValue string `json:"default_value"`
	AddsToList   string `json:"adds_to_list"`
	Desc         string `json:"desc"`
}

func LoadSchema(r io.Reader) (*Schema, error) {
	var schema Schema
	if err := json.NewDecoder(r).Decode(&schema); err != nil {
		return nil, errors.WithStack(err)
	}
	return &schema, nil
}

func (s *Schema) getFlags(category string) []Property {
	var res []Property
	for _, flag := range s.Flags {
		if flag.Category == category {
			res = append(res, Property{
				Key: flag.Key,
				// Name:    flag.Key,
				Desc: I18NString{
					"zh": flag.Desc, // TODO
				},
			})
		}
	}
	return res
}

func (s *Schema) ListAvailableUnitProperties(unitType UnitType) []Property {
	switch unitType {
	case UnitTypeInfantry:
		return slices.Concat(s.getFlags("AbstractTypes"), s.getFlags("ObjectTypes"), s.getFlags("TechnoTypes"), s.getFlags("InfantryTypes"))
	case UnitTypeVehicle:
		return slices.Concat(s.getFlags("AbstractTypes"), s.getFlags("ObjectTypes"), s.getFlags("TechnoTypes"), s.getFlags("VehicleTypes"))
	case UnitTypeAircraft:
		return slices.Concat(s.getFlags("AbstractTypes"), s.getFlags("ObjectTypes"), s.getFlags("TechnoTypes"), s.getFlags("AircraftTypes"))
	case UnitTypeBuilding:
		return slices.Concat(s.getFlags("AbstractTypes"), s.getFlags("ObjectTypes"), s.getFlags("TechnoTypes"), s.getFlags("BuildingTypes"))
	default:
		return nil
	}
}
