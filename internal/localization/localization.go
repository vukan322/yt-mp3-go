package localization

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewBundle(localesDir string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	messageFiles, err := filepath.Glob(filepath.Join(localesDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to find locale files: %w", err)
	}

	for _, file := range messageFiles {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			return nil, fmt.Errorf("failed to load message file %s: %w", file, err)
		}
	}

	return bundle, nil
}

