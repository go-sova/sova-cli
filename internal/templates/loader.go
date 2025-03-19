package templates

import (
	"text/template"

	"github.com/go-sova/sova-cli/pkg/utils"
)

type TemplateLoader struct {
	logger *utils.Logger
}

func NewTemplateLoader() *TemplateLoader {
	return &TemplateLoader{
		logger: utils.NewLoggerWithPrefix(utils.Info, "TemplateLoader"),
	}
}

func (l *TemplateLoader) SetLogger(logger *utils.Logger) {
	l.logger = logger
}

func (l *TemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	return template.New(name).ParseFiles(name)
}

func (l *TemplateLoader) LoadTemplateWithFuncs(name string, funcs template.FuncMap) (*template.Template, error) {
	return template.New(name).Funcs(funcs).ParseFiles(name)
}
