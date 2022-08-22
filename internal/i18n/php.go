package i18n

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/huangchao308/tools/pkg"
)

type PhpI18nGenerator struct {
	BaseGenerator
}

func NewPhpI18nGenerator(params *I18nGeneratorParams) I18nGenerator {
	return &PhpI18nGenerator{NewBaseGenerator(params)}
}

func (g *PhpI18nGenerator) Run() error {
	return g.BaseGenerator.Run(g.generateLine, g.getOutFiles, g.after)
}
func (g *PhpI18nGenerator) generateLine(key, value string) string {
	return "    \"" + key + "\" => \"" + value + "\","
}

func (g *PhpI18nGenerator) getOutFiles() (map[string]*os.File, error) {
	outFiles := make(map[string]*os.File)
	for _, lang := range g.langs {
		dir := path.Join(g.params.OutDir)
		err := pkg.MkDirIfNotExists(dir)
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(path.Join(g.params.OutDir, lang+".php"), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, err
		}
		f.WriteString("<?php\n$__messages = [\n")
		outFiles[lang] = f
	}
	return outFiles, nil
}

func (g *PhpI18nGenerator) after(fs map[string]*os.File) error {
	var err error
	var errs []string
	for k, f := range fs {
		f.WriteString("];\nreturn $__messages;\n")
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
