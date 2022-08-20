package i18n

import "github.com/huangchao308/tools/pkg"

type PhpI18nGenerator struct {
	params      *I18nGeneratorParams
	excelHelper *pkg.ExcelHelper
	langs       []string
}

func NewPhpI18nGenerator(params *I18nGeneratorParams) I18nGenerator {
	return &PhpI18nGenerator{
		params:      params,
		excelHelper: pkg.NewExcelHelper(params.ExcelFile, params.SheetName),
		langs:       []string{"en", "zh_cn", "zh_tw", "ar", "id", "ko", "ms", "th", "tr", "vi", "ja", "bn", "hi", "ur"},
	}
}

func (g *PhpI18nGenerator) Run() error {
	return nil
}
