package excel

import (
	"os"
	"encoding/csv"
)

func CreateNewExcelOrCsv(fileName string, title []string, vals [][]string) (err error) {

	f, err := os.Create(fileName)

	if err != nil {
		return
	}

	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	var w = csv.NewWriter(f)

	w.Write(title)
	w.WriteAll(vals)
	w.Flush()
	return
}
