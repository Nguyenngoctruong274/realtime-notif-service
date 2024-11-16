package excel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"yes4all/ads-noti-api/pkg/logger"
)

const (
	OneKB int = 1024000
)

type RowSheet struct {
	ItemID      string   `json:"itemId"`
	Role        string   `json:"role"`
	ProfileID   string   `json:"ProfileID"`
	Name        string   `json:"name"`
	PortfolioID string   `json:"portfolioExternalId"`
	Emails      []string `json:"Username"`
}

func ReadByFileHeader(
	ctx context.Context,
	file multipart.File,
) (xlsxFile *excelize.File, xlsxClose func() error, err error) {
	funcName := "ReadByFileHeader"
	entry := logger.NewLogger().WithKeyword(ctx, funcName)

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		entry.WithError(err).Error()
		return
	}
	r := bytes.NewReader(bs)
	if r.Len() > OneKB*10 {
		err = fmt.Errorf("file không vượt quá 10MB")
		entry.WithError(err).Error()
		return
	}
	xlsxFile, err = excelize.OpenReader(r)
	if err != nil {
		entry.WithError(err).Error()
		return
	}

	entry.Info()
	return xlsxFile, xlsxFile.Close, nil
}

func GetDataSheets(ctx context.Context, f *excelize.File) (res map[string][]RowSheet, err error) {
	funcName := "GetDataSheets"
	entry := logger.NewLogger().WithKeyword(ctx, funcName)

	res = make(map[string][]RowSheet)
	sheetMap := f.GetSheetMap()
	for _, sheet := range sheetMap {
		rows, sheetErr := GetDataRowSheet(ctx, f, sheet)
		if sheetErr != nil {
			err = sheetErr
			entry.WithError(err).Error()
			return
		}

		res[sheet] = rows
	}

	entry.WithOutputField("res", res).Info()
	return res, nil
}

func GetDataRowSheet(ctx context.Context, f *excelize.File, sheet string) (res []RowSheet, err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "")

	rows, err := f.GetRows(sheet)
	if err != nil {
		entry.WithError(err).Error()
		return
	}

	res = make([]RowSheet, 0)
	headers := []string{}
	for indexRow, row := range rows {
		if indexRow == 0 {
			headers = row
			continue
		}

		rowData := make(map[string]interface{})
		for colIndex, colCell := range row {
			colCell = strings.TrimSpace(colCell)

			if headers[colIndex] == "Username" {
				emails := strings.Split(colCell, "\n")
				rowData[headers[colIndex]] = emails
			} else {
				rowData[headers[colIndex]] = colCell
			}
		}

		temp := RowSheet{}
		jsStr, _ := json.Marshal(rowData)
		json.Unmarshal(jsStr, &temp)
		res = append(res, temp)
	}

	entry.Info()
	return
}
