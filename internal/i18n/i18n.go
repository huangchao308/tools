package i18n

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/huangchao308/tools/pkg"
)

type I18nGenerator interface {
	Run() error
}

type I18nGeneratorParams struct {
	ExcelFile  string
	SheetName  string
	OutDir     string
	OutFile    string
	TargetLang string
	TextLang   string
	OldFile    string
	Overwrite  bool
}

type BaseGenerator struct {
	params      *I18nGeneratorParams
	excelHelper *pkg.ExcelHelper
	langs       []string
}

func NewBaseGenerator(params *I18nGeneratorParams) BaseGenerator {
	langs := []string{"en", "zh_cn", "zh_tw", "ar", "id", "ko", "ms", "th", "tr", "vi", "ja", "bn", "hi", "ur"}
	if params.TextLang != "all" {
		langs = []string{params.TextLang}
	}
	return BaseGenerator{
		params:      params,
		excelHelper: pkg.NewExcelHelper(params.ExcelFile, params.SheetName),
		langs:       langs,
	}
}

func (g *BaseGenerator) formatString(key string) string {
	key = strings.TrimPrefix(key, "'")
	key = strings.TrimPrefix(key, "\"")
	key = strings.TrimSuffix(key, ",")
	key = strings.TrimSuffix(key, "'")
	key = strings.TrimSuffix(key, "\"")
	key = strings.ReplaceAll(key, "\"", "\\\"")
	key = strings.ReplaceAll(key, "\n", "\\n")
	return key
}

func (g *BaseGenerator) Run(generateLine func(k, v string) string, getOutFiles func() (map[string]*os.File, error), after func(fs map[string]*os.File) error, old map[string]string) error {
	err := g.excelHelper.Open()
	if err != nil {
		return err
	}
	keys, err := g.excelHelper.GetKeys()
	log.Printf("keys: %v\n", keys)
	if err != nil {
		return err
	}
	rows, err := g.excelHelper.GetRows(keys, true)
	if err != nil {
		return err
	}
	outFiles, err := getOutFiles()
	if err != nil {
		return err
	}
	defer func() {
		after(outFiles)
	}()
	hadnledKeys := make(map[string]bool)
	for _, row := range rows {
		key := row["key"]
		if key == "" {
			continue
		}
		key = g.formatString(key)
		key = g.params.OutFile + "." + key
		if _, ok := hadnledKeys[key]; ok {
			log.Println("This key has been handled:", key)
			continue
		}
		if !g.params.Overwrite && len(old) > 0 {
			if _, ok := old[key]; !ok {
				log.Printf("This key %s is ignored", key)
				continue
			}
		}
		for k, v := range row {
			if v == "" {
				continue
			}
			v = g.formatString(v)
			lang := strings.ToLower(k)
			if f, ok := outFiles[lang]; ok {
				line := generateLine(key, v)
				_, err := fmt.Fprintln(f, line)
				if err != nil {
					return err
				}
			}
		}
		hadnledKeys[key] = true
	}

	// 合并旧的 key
	if g.params.TextLang != "" && g.params.TextLang != "all" {
		f, ok := outFiles[g.params.TextLang]
		log.Println(ok)
		if !ok {
			return nil
		}
		for k, v := range old {
			if _, ok := hadnledKeys[k]; !ok {
				// log.Printf("merge old key: %s, value: %s", k, v)
				line := generateLine(k, v)
				_, err := fmt.Fprintln(f, line)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
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
