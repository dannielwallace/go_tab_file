// tab_file project tab_file.go
package tab_file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var MAX_LINE_SIZE = 4096

type TabFile struct {
	filepath string
	colNames []string
	rowNames []string
	rows     [][]string
}

type TabRow struct {
	colNames []string
	cells    []string
}

func (this *TabRow) GetCellByColName(colName string) (val string, err error) {
	colNum := -1
	for k, v := range this.colNames {
		if colName == v {
			colNum = k
			break
		}
	}

	if colNum < 0 {
		err = errors.New(fmt.Sprintf("can't find specific colname:%s", colName))
		return
	}

	return this.GetCellByColNum(colNum)
}

func (this *TabRow) GetCellByColNum(colNum int) (val string, err error) {
	if colNum >= len(this.cells) {
		err = errors.New(fmt.Sprintf("colNum out of row cells:%d", len(this.cells)))
		return
	}

	val = this.cells[colNum]
	return
}

func (this *TabFile) GetRowNames() []string {
	return this.rowNames
}

func (this *TabFile) GetRowByName(rowName string) (row *TabRow, err error) {
	rowNum := -1
	for k, v := range this.rowNames {
		if rowName == v {
			rowNum = k
			break
		}
	}

	if rowNum < 0 {
		err = errors.New(fmt.Sprintf("can't find specific rowname:%s", rowName))
		return
	}

	return this.GetRowByNum(rowNum)
}

func (this *TabFile) GetRowByNum(rowNum int) (row *TabRow, err error) {
	if rowNum >= len(this.rows) {
		err = errors.New(fmt.Sprintf("rowNum out of file height:%d", len(this.rows)))
		return
	}

	row = &TabRow{
		colNames: this.colNames,
		cells:    this.rows[rowNum],
	}

	return
}

func OpenFile(filepath string) (tabFile *TabFile, err error) {
	var file *os.File
	file, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	tabFile = &TabFile{
		filepath: filepath,
	}

	reader := bufio.NewReaderSize(file, MAX_LINE_SIZE)

	var line = []byte{0}
	for i := 1; ; i++ {
		line, _, err = reader.ReadLine()
		if len(line) <= 0 || err != nil {
			break
		}

		strLine := string(line)
		if i == 1 {
			tabFile.colNames = strings.Split(strLine, "\t")
		}

		cells := strings.Split(strLine, "\t")
		for len(cells) < len(tabFile.colNames) {
			cells = append(cells, "")
		}
		cells = cells[:len(tabFile.colNames)]

		tabFile.rowNames = append(tabFile.rowNames, cells[0])
		tabFile.rows = append(tabFile.rows, cells)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
