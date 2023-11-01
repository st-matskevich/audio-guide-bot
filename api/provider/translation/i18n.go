package translation

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18NTranslationProvider struct {
	bundle *i18n.Bundle
}

func (provider *I18NTranslationProvider) TranslateMessage(messageID string, locale string, templateData TemplateData) (string, error) {
	localizer := i18n.NewLocalizer(provider.bundle, locale)
	translation, tag, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil && tag == language.Und {
		return "", err
	}

	return translation, nil
}

func CreateI18NTranslationProvider() (TranslationProvider, error) {
	bundle := i18n.NewBundle(DEFAULT_LANGUAGE)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	locales := []string{
		"locales/en.json",
		"locales/be.json",
		"locales/ru.json",
	}

	for _, locale := range locales {
		_, err := bundle.LoadMessageFile(locale)
		if err != nil {
			return nil, err
		}
	}

	provider := I18NTranslationProvider{
		bundle: bundle,
	}

	return &provider, nil
}
