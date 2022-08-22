package i18n

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/huangchao308/tools/pkg"
)

type GolangI18nGenerator struct {
	BaseGenerator
}

func NewGolangI18nGenerator(params *I18nGeneratorParams) I18nGenerator {
	return &GolangI18nGenerator{NewBaseGenerator(params)}
}

func (g *GolangI18nGenerator) Run() error {
	return g.BaseGenerator.Run(g.generateLine, g.getOutFiles, g.after)
}

func (g *GolangI18nGenerator) generateLine(key, value string) string {
	return "\"" + key + "\" = \"" + value + "\""
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

func (g *GolangI18nGenerator) after(fs map[string]*os.File) error {
	var err error
	var errs []string
	for k, f := range fs {
		err = f.Close()
		if err != nil {
			errs = append(errs, fmt.Sprintf("close file %s error: %s", k, err.Error()))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
