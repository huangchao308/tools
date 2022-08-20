package i18n

type I18nGenerator interface {
	Run() error
}

type I18nGeneratorParams struct {
	ExcelFile  string
	SheetName  string
	OutDir     string
	OutFile    string
	TargetLang string
}

func NewI18nGenerator(params *I18nGeneratorParams) I18nGenerator {
	lang := params.TargetLang
	switch lang {
	case "go":
		return NewGolangI18nGenerator(params)
	case "php":
		return NewPhpI18nGenerator(params)
	}
	return nil
}
