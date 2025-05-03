package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"ra2-ini-editor/internal/ra2"
)

//go:embed data
var dataFs embed.FS

// App struct
type App struct {
	ctx context.Context

	schema      *ra2.Schema
	origin      *ra2.Rules
	translation *ra2.Translation

	rules *ra2.Rules
}

// NewApp creates a new App application struct
func NewApp() *App {
	schemaFile, err := dataFs.Open("data/schema.json")
	if err != nil {
		panic(err)
	}
	defer schemaFile.Close()
	schema, err := ra2.LoadSchema(schemaFile)
	if err != nil {
		panic(err)
	}

	originFile, err := dataFs.Open("data/rulesmd.ini")
	if err != nil {
		panic(err)
	}
	defer originFile.Close()
	origin, err := ra2.NewRules(originFile)
	if err != nil {
		panic(err)
	}

	translationFile, err := dataFs.Open("data/ra2md.ini")
	if err != nil {
		panic(err)
	}
	defer translationFile.Close()
	translation, err := ra2.LoadTranslation(translationFile, "zh-TW")
	if err != nil {
		panic(err)
	}
	return &App{
		schema:      schema,
		origin:      origin,
		translation: translation,

		rules: ra2.NewEmptyRules(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

func NewAppErrorf(code int, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func (a *App) NewULID() string {
	return ulid.Make().String()
}

func (a *App) Open() error {
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择一个文件",
		Filters: []runtime.FileFilter{
			{Pattern: "*.ini", DisplayName: "INI Files (*.ini)"},
		},
	})
	if err != nil {
		return NewAppErrorf(500, "open file dialog error: %v", err)
	}
	if filename == "" {
		return NewAppError(400, "no file selected")
	}

	rulesFile, err := os.Open(filename)
	if err != nil {
		return NewAppErrorf(500, "open file error: %v", err)
	}
	defer rulesFile.Close()
	rules, err := ra2.NewRules(rulesFile)
	if err != nil {
		return NewAppErrorf(500, "load rules error: %v", err)
	}
	a.rules = rules
	return nil
}

func (a *App) Save() error {
	if a.rules == nil {
		return NewAppError(400, "no file opened")
	}

	filename, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "保存文件",
		Filters: []runtime.FileFilter{
			{Pattern: "*.ini", DisplayName: "INI Files (*.ini)"},
		},
	})
	if err != nil {
		return NewAppErrorf(500, "save file dialog error: %v", err)
	}

	if filename == "" {
		return NewAppError(400, "no file selected")
	}
	rulesFile, err := os.Create(filename)
	if err != nil {
		return NewAppErrorf(500, "create file error: %v", err)
	}
	defer rulesFile.Close()
	if err := a.rules.Save(rulesFile); err != nil {
		return NewAppErrorf(500, "save rules error: %v", err)
	}
	return nil
}

func (a *App) UserRules() (string, error) {
	bts, err := a.rules.Content()
	if err != nil {
		return "", NewAppErrorf(500, "get rules content error: %v", err)
	}
	return string(bts), nil
}

func (a *App) getRules() *ra2.Rules {
	r := a.origin
	if a.rules != nil {
		r = lo.Must(r.Merge(a.rules))
	}
	return r
}

type Property struct {
	UKey    string `json:"ukey"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Comment string `json:"comment"`

	Desc *string `json:"desc"`
}

type Unit struct {
	Type       string     `json:"type"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	UIName     string     `json:"ui_name"`
	Properties []Property `json:"properties"`
}

func (a *App) ListAllUnits() ([]*Unit, error) {
	defer func() {
		if err := recover(); err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("panic: %v", err))
			runtime.LogError(a.ctx, string(debug.Stack()))
		}
	}()

	r := a.getRules()

	units := make([]*Unit, 0)
	for _, unit := range r.Units() {
		units = append(units, &Unit{
			Type:   string(unit.Type),
			ID:     unit.ID,
			Name:   unit.Name,
			UIName: a.translation.Get(unit.UIName()),
		})
	}

	return units, nil
}

func (a *App) GetUnit(unitType string, id int) (*Unit, error) {
	r := a.getRules()

	unit := r.GetUnit(ra2.NewUnitType(unitType), id)
	if unit == nil {
		return nil, NewAppErrorf(404, "unit not found")
	}

	availableProps := a.schema.ListAvailableUnitProperties(unit.Type)
	props := make([]Property, 0)
	for _, prop := range unit.Properties() {
		prop := Property{
			UKey:    ulid.Make().String(),
			Key:     prop.Key,
			Value:   prop.Value,
			Comment: prop.Comment,
		}
		schemaProp, ok := lo.Find(availableProps, func(p ra2.Property) bool {
			return p.Key == prop.Key
		})
		if ok {
			prop.Desc = lo.ToPtr(schemaProp.Desc.Get("zh"))
		}
		props = append(props, prop)
	}

	return &Unit{
		Type:       string(unit.Type),
		ID:         unit.ID,
		Name:       unit.Name,
		UIName:     a.translation.Get(unit.UIName()),
		Properties: props,
	}, nil
}

func (a *App) ListAvailableProperties(unitType string) ([]Property, error) {
	availableProps := a.schema.ListAvailableUnitProperties(ra2.NewUnitType(unitType))
	props := make([]Property, 0)
	for _, prop := range availableProps {
		props = append(props, Property{
			UKey:    ulid.Make().String(),
			Key:     prop.Key,
			Value:   prop.Value,
			Comment: prop.Comment,
			Desc:    lo.ToPtr(prop.Desc.Get("zh")),
		})
	}
	return props, nil
}

func (a *App) NextUnitID(unitType string) (int, error) {
	r := a.getRules()

	maxID := 0
	for _, unit := range r.UnitsByType(ra2.NewUnitType(unitType)) {
		if unit.Type == ra2.NewUnitType(unitType) && unit.ID > maxID {
			maxID = unit.ID
		}
	}

	return maxID + 1, nil
}

func (a *App) SaveUnit(mod *Unit) error {
	modProps := make([]ra2.Property, 0)
	for _, prop := range mod.Properties {
		modProps = append(modProps, ra2.Property{
			Key:     prop.Key,
			Value:   prop.Value,
			Comment: prop.Comment,
		})
	}

	originUnit := a.origin.GetUnit(ra2.NewUnitType(mod.Type), mod.ID)
	if originUnit == nil {
		_, err := a.rules.AddUnit(ra2.NewUnitType(mod.Type), mod.ID, mod.Name, modProps)
		if err != nil {
			return NewAppErrorf(500, "add unit error: %v", err)
		}
		return nil
	}

	userUnit := a.rules.GetUnit(ra2.NewUnitType(mod.Type), mod.ID)
	if userUnit == nil {
		// 新建用户级 unit
		unit, err := a.rules.AddUnit(ra2.NewUnitType(mod.Type), mod.ID, mod.Name, nil)
		if err != nil {
			return NewAppErrorf(500, "add unit error: %v", err)
		}
		userUnit = unit
	}

	originProps := originUnit.Properties()
	userProps := userUnit.Properties()
	for _, modProp := range modProps {
		prop, ok := lo.Find(originProps, func(p ra2.Property) bool {
			return p.Key == modProp.Key
		})
		if !ok {
			userUnit.Set(modProp.Key, modProp.Value, modProp.Comment)
		} else {
			if prop.Value != modProp.Value || prop.Comment != modProp.Comment {
				userUnit.Set(modProp.Key, modProp.Value, modProp.Comment)
			}
		}
	}
	for _, userProp := range userProps {
		_, ok := lo.Find(modProps, func(p ra2.Property) bool {
			return p.Key == userProp.Key
		})
		if !ok {
			userUnit.Del(userProp.Key)
		}
	}
	for _, originProp := range originProps {
		_, ok := lo.Find(modProps, func(p ra2.Property) bool {
			return p.Key == originProp.Key
		})
		if !ok {
			userUnit.Set(originProp.Key, "")
		}
	}

	return nil
}

func (a *App) DeleteUnit(unitType string, id int) error {
	r := a.getRules()

	unit := r.GetUnit(ra2.NewUnitType(unitType), id)
	if unit == nil {
		return NewAppErrorf(404, "unit not found")
	}

	userUnit := a.rules.GetUnit(ra2.NewUnitType(unitType), id)
	if userUnit == nil {
		return NewAppErrorf(400, "cannot delete unit from origin rules")
	}

	if err := a.rules.DelUnit(userUnit.Type, userUnit.ID); err != nil {
		return NewAppErrorf(500, "delete unit error: %v", err)
	}

	return nil
}
