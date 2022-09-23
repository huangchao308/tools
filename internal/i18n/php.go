package i18n

import (
	"bufio"
	"errors"
	"fmt"
	"log"
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
	return g.BaseGenerator.Run(g.generateLine, g.getOutFiles, g.after, g.getOldKvFromFile())
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

func (g *PhpI18nGenerator) getOldKvFromFile() map[string]string {
	result := make(map[string]string)
	if g.params.OldFile == "" {
		return result
	}
	fs, err := os.Open(g.params.OldFile)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	for scanner.Scan() {
		key := ""
		value := ""
		line := scanner.Text()
		if !strings.Contains(line, "=>") {
			continue
		}
		str := strings.Split(line, "=>")
		if len(str) > 1 {
			key = str[0]
			key = strings.TrimSpace(key)
			key = strings.Trim(key, "\"")
			key = strings.Trim(key, "'")
			value = str[1]
			value = strings.TrimSpace(value)
			value = strings.Trim(value, "\"")
			value = strings.Trim(value, "'")
			result[key] = value
		}
	}

	return result
}
