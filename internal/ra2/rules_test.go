package ra2

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func TestRules(t *testing.T) {
	rules, err := NewRules("../../data/rulesmd.ini")
	if err != nil {
		t.Fatalf("failed to parse file: %v", err)
	}

	for _, unit := range rules.Units() {
		t.Logf("Unit: %s, ID: %d, Type: %s", unit.Name, unit.ID, unit.Type)
		properties := unit.Properties()
		for _, prop := range properties {
			t.Logf("  Property: %s, Value: %s", prop.Key, prop.Value)
		}
	}

	assert.NotNil(t, rules)
}

func TestRules_Merge(t *testing.T) {
	type args struct {
		others []*Rules
	}
	tests := []struct {
		name      string
		r         *Rules
		args      args
		want      *Rules
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "1",
			r: &Rules{
				f: lo.Must(ini.Load([]byte("[section]\nkey1=value1"))),
			},
			args: args{
				others: []*Rules{
					{f: lo.Must(ini.Load([]byte("[section]\nkey1=value\nkey2=value2")))},
				},
			},
			want: &Rules{
				f: lo.Must(ini.Load([]byte("[section]\nkey1=value\nkey2=value2"))),
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Merge(tt.args.others...)
			tt.assertion(t, err)
			assert.True(t, compareIni(got.f, tt.want.f))
		})
	}
}
