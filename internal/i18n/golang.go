package i18n

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/huangchao308/tools/pkg"
)

type GolangI18nGenerator struct {
	params      *I18nGeneratorParams
	excelHelper *pkg.ExcelHelper
	langs       []string
}

func NewGolangI18nGenerator(params *I18nGeneratorParams) I18nGenerator {
	return &GolangI18nGenerator{
		params:      params,
		excelHelper: pkg.NewExcelHelper(params.ExcelFile, params.SheetName),
		langs:       []string{"en", "zh_cn", "zh_tw", "ar", "id", "ko", "ms", "th", "tr", "vi", "ja", "bn", "hi", "ur"},
	}
}

func (g *GolangI18nGenerator) Run() error {
	err := g.excelHelper.Open()
	if err != nil {
		return err
	}
	keys, err := g.excelHelper.GetKeys()
	if err != nil {
		return err
	}
	rows, err := g.excelHelper.GetRows(keys, true)
	if err != nil {
		return err
	}
	outFiles, err := g.getOutFiles()
	if err != nil {
		return err
	}
	defer func() {
		for _, f := range outFiles {
			if f != nil {
				f.Close()
			}
		}
	}()
	hadnledKeys := make(map[string]bool)
	for _, row := range rows {
		key := row["key"]
		if key == "" {
			continue
		}
		if _, ok := hadnledKeys[key]; ok {
			log.Println("This key has been handled:", key)
			continue
		}
		key = strings.TrimPrefix(key, "'")
		key = strings.TrimPrefix(key, "\"")
		key = strings.TrimSuffix(key, "'")
		key = strings.TrimSuffix(key, "\"")
		for k, v := range row {
			if v == "" {
				continue
			}
			v = strings.TrimPrefix(v, "'")
			v = strings.TrimPrefix(v, "\"")
			v = strings.TrimSuffix(v, "'")
			v = strings.TrimSuffix(v, "\"")
			v = strings.ReplaceAll(v, "\"", "\\\"")
			lang := strings.ToLower(k)
			if f, ok := outFiles[lang]; ok {
				line := "\"" + key + "\" = \"" + v + "\""
				_, err := fmt.Fprintln(f, line)
				if err != nil {
					return err
				}
			}
		}
		hadnledKeys[key] = true
	}

	return nil
}

func (g *GolangI18nGenerator) getOutFiles() (map[string]*os.File, error) {
	outFiles := make(map[string]*os.File)
	for _, lang := range g.langs {
		dir := path.Join(g.params.OutDir, lang)
		err := pkg.MkDirIfNotExists(dir)
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(path.Join(dir, g.params.OutFile+".toml"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}
		outFiles[lang] = f
	}
	return outFiles, nil
}
