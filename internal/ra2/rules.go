package ra2

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

type SectionName string

const (
	SectionNameCountry  SectionName = "Countries"
	SectionNameInfantry SectionName = "InfantryTypes"
	SectionNameVehicle  SectionName = "VehicleTypes"
	SectionNameAircraft SectionName = "AircraftTypes"
	SectionNameBuilding SectionName = "BuildingTypes"
)

type BaseSetting struct {
	sec *ini.Section
}

func (s *BaseSetting) UIName() string {
	return s.sec.Key("UIName").String()
}

func (s *BaseSetting) Properties() []Property {
	return parseProperties(s.sec)
}

func (s *BaseSetting) Get(key string) string {
	if s.sec.HasKey(key) {
		return s.sec.Key(key).String()
	}
	return ""
}

func (s *BaseSetting) Set(key, value string, comment ...string) error {
	k, err := s.sec.NewKey(key, value)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(comment) > 0 {
		k.Comment = comment[0]
	}
	return nil
}

func (s *BaseSetting) Del(key string) {
	s.sec.DeleteKey(key)
}

type Country struct {
	BaseSetting

	ID   int
	Name string
}

type Unit struct {
	BaseSetting

	Type UnitType
	ID   int    // e.g. 0
	Name string // e.g. "E1"
}

type UnitType string

const (
	UnitTypeUnknown  UnitType = ""
	UnitTypeInfantry UnitType = "infantry"
	UnitTypeVehicle  UnitType = "vehicle"
	UnitTypeAircraft UnitType = "aircraft"
	UnitTypeBuilding UnitType = "building"
)

var UnitTypes = []UnitType{
	UnitTypeInfantry,
	UnitTypeVehicle,
	UnitTypeAircraft,
	UnitTypeBuilding,
}

func NewUnitType(name string) UnitType {
	switch name {
	case "infantry":
		return UnitTypeInfantry
	case "vehicle":
		return UnitTypeVehicle
	case "aircraft":
		return UnitTypeAircraft
	case "building":
		return UnitTypeBuilding
	default:
		return UnitTypeUnknown
	}
}

func (ut UnitType) Section() SectionName {
	switch ut {
	case UnitTypeInfantry:
		return SectionNameInfantry
	case UnitTypeVehicle:
		return SectionNameVehicle
	case UnitTypeAircraft:
		return SectionNameAircraft
	case UnitTypeBuilding:
		return SectionNameBuilding
	default:
		panic("unknown unit type")
	}
}

type Rules struct {
	f *ini.File
}

func NewRules(r io.ReadCloser) (*Rules, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		KeyValueDelimiters: "=",
	}, r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Rules{
		f: f,
	}, nil
}

func NewEmptyRules() *Rules {
	return &Rules{
		f: ini.Empty(),
	}
}

func (r *Rules) Save(w io.Writer) error {
	if _, err := r.f.WriteTo(w); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *Rules) Content() ([]byte, error) {
	var buf bytes.Buffer
	if err := r.Save(&buf); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func (r *Rules) Merge(others ...*Rules) (*Rules, error) {
	f := ini.Empty()
	if err := mergeIni(f, r.f); err != nil {
		return nil, errors.WithStack(err)
	}
	for _, other := range others {
		if err := mergeIni(f, other.f); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &Rules{
		f: f,
	}, nil
}

func (r *Rules) Units() []*Unit {
	var units []*Unit
	units = append(units, r.UnitsByType(UnitTypeInfantry)...)
	units = append(units, r.UnitsByType(UnitTypeVehicle)...)
	units = append(units, r.UnitsByType(UnitTypeAircraft)...)
	units = append(units, r.UnitsByType(UnitTypeBuilding)...)
	return units
}

func (r *Rules) UnitsByType(unitType UnitType) []*Unit {
	var units []*Unit
	for _, key := range r.f.Section(string(unitType.Section())).Keys() {
		id := cast.ToInt(key.Name())
		name := key.Value()
		units = append(units, &Unit{
			BaseSetting: BaseSetting{sec: r.f.Section(name)},

			Type: unitType,
			ID:   id,
			Name: name,
		})
	}
	return units
}

func (r *Rules) GetUnit(unitType UnitType, unitID int) *Unit {
	for _, unit := range r.Units() {
		if unit.Type == unitType && unit.ID == unitID {
			return unit
		}
	}
	return nil
}

func (r *Rules) AddUnit(unitType UnitType, unitID int, unitName string, properties []Property) (*Unit, error) {
	if r.f.HasSection(unitName) {
		return nil, errors.New("unit name already used")
	}
	if r.GetUnit(unitType, unitID) != nil {
		return nil, errors.New("unit ID already exists")
	}

	defSec := r.f.Section(string(unitType.Section()))
	_, err := defSec.NewKey(cast.ToString(unitID), unitName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sec, err := r.f.NewSection(unitName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	unit := &Unit{
		BaseSetting: BaseSetting{sec: sec},

		Type: unitType,
		ID:   unitID,
		Name: unitName,
	}
	for _, prop := range properties {
		if err := unit.Set(prop.Key, prop.Value, prop.Comment); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return unit, nil
}

func (r *Rules) DelUnit(unitType UnitType, unitID int) error {
	unit := r.GetUnit(unitType, unitID)
	if unit == nil {
		return errors.New("unit not found")
	}
	r.f.DeleteSection(unit.Name)
	defSec := r.f.Section(string(unitType.Section()))
	defSec.DeleteKey(cast.ToString(unitID))
	return nil
}

func mergeIni(dst, src *ini.File) error {
	for _, sec := range src.Sections() {
		newSec, err := dst.NewSection(sec.Name())
		if err != nil {
			return errors.WithStack(err)
		}
		for _, key := range sec.Keys() {
			newKey, err := newSec.NewKey(key.Name(), key.Value())
			if err != nil {
				return errors.WithStack(err)
			}
			newKey.Comment = key.Comment
		}
	}
	return nil
}

func compareIni(f1, f2 *ini.File) bool {
	// 比较 section 数量是否一致
	sections1 := f1.Sections()
	sections2 := f2.Sections()

	if len(sections1) != len(sections2) {
		return false
	}

	// 比较每一个 section
	for _, s1 := range sections1 {
		sectionName := s1.Name()
		s2, err := f2.GetSection(sectionName)
		if err != nil {
			return false
		}

		// 比较 key 数量是否一致
		keys1 := s1.Keys()
		keys2 := s2.Keys()
		if len(keys1) != len(keys2) {
			return false
		}

		// 比较每一个 key 和其对应的值
		for _, k1 := range keys1 {
			k2, err := s2.GetKey(k1.Name())
			if err != nil {
				return false
			}
			if k1.Value() != k2.Value() {
				return false
			}
		}
	}

	return true
}
