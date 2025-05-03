package ra2

import (
	"os"
	"testing"
)

func TestTranslation(t *testing.T) {
	translationFile, err := os.Open("../../data/ra2md.ini")
	if err != nil {
		t.Fatalf("failed to open translation file: %v", err)
	}
	defer translationFile.Close()
	translation, err := LoadTranslation(translationFile, "zh-TW")
	if err != nil {
		t.Fatalf("failed to load translation: %v", err)
	}

	for _, key := range translation.sec.Keys() {
		t.Logf("Key: %s, Value: %s", key.Name(), key.Value())
	}
}
