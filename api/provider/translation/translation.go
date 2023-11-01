package translation

import "golang.org/x/text/language"

var DEFAULT_LANGUAGE = language.English

type TemplateData map[string]interface{}

type TranslationProvider interface {
	TranslateMessage(messageID string, locale string, templateData TemplateData) (string, error)
}
