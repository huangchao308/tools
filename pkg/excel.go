package pkg

import (
	"sync"

	"github.com/xuri/excelize/v2"
)

var excelHelper *ExcelHelper

type ExcelHelper struct {
	ExcelFile  string
	SheetName  string
	opened     bool
	mutex      sync.Mutex
	openedFile *excelize.File
}

func NewExcelHelper(excelFile string, sheetName string) *ExcelHelper {
	if excelHelper == nil {
		excelHelper = &ExcelHelper{
			ExcelFile: excelFile,
			SheetName: sheetName,
			opened:    false,
			mutex:     sync.Mutex{},
		}
	}
	return excelHelper
}

func (h *ExcelHelper) Open() error {
	if h.opened {
		return nil
	}
	h.mutex.Lock()
	defer h.mutex.Unlock()
	f, err := excelize.OpenFile(h.ExcelFile)
	if err != nil {
		return err
	}
	h.openedFile = f
	h.opened = true
	return nil
}

func (h *ExcelHelper) GetRows(keys []string, skipFirstRow bool) ([]map[string]string, error) {
	rows, err := h.openedFile.Rows(h.SheetName)
	if err != nil {
		return nil, err
	}
	i := 0
	results := make([]map[string]string, 0)
	for rows.Next() {
		i++
		if skipFirstRow && i == 1 {
			continue
		}
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		item := make(map[string]string)
		for j, c := range cols {
			item[keys[j]] = c
		}
		results = append(results, item)
	}

	return results, nil
}

func (h *ExcelHelper) GetKeys() ([]string, error) {
	rows, err := h.openedFile.Rows(h.SheetName)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	if rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		keys = append(keys, cols...)
	}

	return keys, nil
}
