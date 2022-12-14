package cmd

import (
	"errors"
	"log"

	"github.com/huangchao308/tools/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	excelFile   string
	sheetName   string
	outDir      string
	outFile     string
	includeLang string
	oldFile     string
	overwrite   bool
)

var i18nCmd = &cobra.Command{
	Use:       "i18n",
	Short:     "生成多语言文件",
	Long:      "根据 Excel 文件、编程语言，生成对应的多语言配置文件",
	Example:   "./bin/tools i18n go -f excel.xlsx -s Sheet1 -O ./out -o data --include.language en",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"go", "php"},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("tartget lang:", args[0])
		log.Printf("excelFile: %s, sheetName: %s, outDir: %s, outFile: %s, includeLang: %s, oldFile: %s", excelFile, sheetName, outDir, outFile, includeLang, oldFile)
		params := &i18n.I18nGeneratorParams{
			ExcelFile:  excelFile,
			SheetName:  sheetName,
			OutDir:     outDir,
			OutFile:    outFile,
			TargetLang: args[0],
			TextLang:   includeLang,
			OldFile:    oldFile,
			Overwrite:  overwrite,
		}
		generator := i18n.NewI18nGenerator(params)
		if generator == nil {
			return errors.New("not implemented")
		}
		err := generator.Run()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	i18nCmd.Flags().StringVarP(&excelFile, "excel.file", "f", "", "Excel 文件路径")
	i18nCmd.MarkFlagRequired("excel.file")
	i18nCmd.Flags().StringVarP(&sheetName, "excel.sheet", "s", "Sheet1", "Excel 表格名称")
	i18nCmd.Flags().StringVarP(&outDir, "out.dir", "O", "./out", "输出目录")
	i18nCmd.Flags().StringVarP(&outFile, "out.file", "o", "data", "输出文件名, 仅生成 go 语言的配置文件时有效, PHP 的配置文件固定为对应语言(如 ar.php)")
	i18nCmd.Flags().StringVar(&includeLang, "include.language", "all", "想要生成的语言")
	i18nCmd.Flags().StringVar(&oldFile, "old.file", "", "老的多语言配置文件, 用于合并新旧的文件")
	i18nCmd.Flags().BoolVar(&overwrite, "overwrite", true, "是否覆盖旧的翻译")
}
