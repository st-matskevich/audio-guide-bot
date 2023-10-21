package translation

type TemplateData map[string]interface{}

type TranslationProvider interface {
	TranslateMessage(messageID string, locale string, templateData TemplateData) (string, error)
}
