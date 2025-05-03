package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"runtime/debug"

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
	schema, err := ra2.LoadSchema("data/schema.json")
	if err != nil {
		panic(err)
	}
	origin, err := ra2.NewRules("data/rulesmd.ini")
	if err != nil {
		panic(err)
	}
	translation, err := ra2.LoadTranslation("data/ra2md.ini", "zh-TW")
	if err != nil {
		panic(err)
	}
	return &App{
		schema:      schema,
		origin:      origin,
		translation: translation,

		rules: lo.Must(ra2.NewRules("")),
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

	rules, err := ra2.NewRules(filename)
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

func (a *App) getRules() *ra2.Rules {
	r := a.origin
	if a.rules != nil {
		r = lo.Must(r.Merge(a.rules))
	}
	return r
}

type Property struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type Unit struct {
	Type       string     `json:"type"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}

func (a *App) getUnitName(unit *ra2.Unit) string {
	name := a.translation.Get(unit.UIName())
	if name == "" {
		name = unit.Name
	}
	return name
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
			Type: string(unit.Type),
			ID:   unit.ID,
			Name: a.getUnitName(unit),
		})
	}

	return units, nil
}

func (a *App) GetUnit(unitType string, id int) (*Unit, error) {
	r := a.getRules()

	unit := r.GetUnit(ra2.UnitType(unitType), id)
	if unit == nil {
		return nil, NewAppErrorf(404, "unit not found")
	}

	props := make([]Property, 0)
	for _, prop := range unit.Properties() {
		props = append(props, Property{
			Key:     prop.Key,
			Value:   prop.Value,
			Comment: prop.Comment,
		})
	}

	return &Unit{
		Type:       string(unit.Type),
		ID:         id,
		Name:       a.getUnitName(unit),
		Properties: props,
	}, nil
}
