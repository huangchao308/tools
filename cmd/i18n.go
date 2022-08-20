package cmd

import (
	"errors"
	"log"

	"github.com/huangchao308/tools/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	excelFile string
	sheetName string
	outDir    string
	outFile   string
)

var i18nCmd = &cobra.Command{
	Use:       "i18n",
	Short:     "生成多语言文件",
	Long:      "根据 Excel 文件、编程语言，生成对应的多语言配置文件",
	Example:   "./bin/tools i18n go -f excel.xlsx -s Sheet1 -O ./out -o data",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"go", "php"},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("tartget lang:", args[0])
		params := &i18n.I18nGeneratorParams{
			ExcelFile:  excelFile,
			SheetName:  sheetName,
			OutDir:     outDir,
			OutFile:    outFile,
			TargetLang: args[0],
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
}
