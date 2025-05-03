package ra2

import (
	"github.com/charmbracelet/log"
	"gopkg.in/ini.v1"
)

type Translation struct {
	lang string
	f    *ini.File
	sec  *ini.Section
}

func LoadTranslation(filename string, lang string) (*Translation, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		KeyValueDelimiters: "=",
	}, filename)
	if err != nil {
		return nil, err
	}
	sec, err := cfg.GetSection(lang)
	if err != nil {
		return nil, err
	}
	return &Translation{
		lang: lang,
		f:    cfg,
		sec:  sec,
	}, nil
}

func (t *Translation) Get(key string) string {
	k, err := t.sec.GetKey(key)
	if err != nil {
		log.Errorf("failed to get key %s: %v", key, err)
		return key
	}
	return k.String()
}
