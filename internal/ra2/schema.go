package ra2

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/pkg/errors"
)

type Schema struct {
	Country        []Property `json:"country"`
	CrateRules     []Property `json:"crate_rules"`
	CombatDamage   []Property `json:"combat_damage"`
	Radiation      []Property `json:"radiation"`
	ElevationModel []Property `json:"elevation_model"`
	WallModel      []Property `json:"wall_model"`

	Unit          []Property `json:"unit"`
	MovingUnit    []Property `json:"moving_unit"`
	TurretChanger []Property `json:"turret_changer"`
	Infantry      []Property `json:"infantry"`
	Vehicle       []Property `json:"vehicle"`
	Aircraft      []Property `json:"aircraft"`
	Building      []Property `json:"building"`
}

func LoadSchema(r io.Reader) (*Schema, error) {
	var schema Schema
	if err := json.NewDecoder(r).Decode(&schema); err != nil {
		return nil, errors.WithStack(err)
	}
	return &schema, nil
}

func (s *Schema) ListAvailableProperties(key string) []Property {
	switch key {
	case "country":
		return s.Country
	case "crate_rules":
		return s.CrateRules
	case "combat_damage":
		return s.CombatDamage
	case "radiation":
		return s.Radiation
	case "elevation_model":
		return s.ElevationModel
	case "wall_model":
		return s.WallModel
	case "infantry", "vehicle", "aircraft", "building":
		return s.ListAvailableUnitProperties(UnitType(key))
	default:
		return nil
	}
}

func (s *Schema) ListAvailableUnitProperties(unitType UnitType) []Property {
	switch unitType {
	case UnitTypeInfantry:
		return slices.Concat(s.Unit, s.MovingUnit, s.Infantry)
	case UnitTypeVehicle:
		return slices.Concat(s.Unit, s.MovingUnit, s.Vehicle)
	case UnitTypeAircraft:
		return slices.Concat(s.Unit, s.MovingUnit, s.Aircraft)
	case UnitTypeBuilding:
		return slices.Concat(s.Unit, s.Building)
	default:
		return nil
	}
}
