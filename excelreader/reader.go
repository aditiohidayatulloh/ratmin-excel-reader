package excelreader

import (
	"bytes"
	"context"
	"io"

	"github.com/xuri/excelize/v2"
)

type ReadConfig struct {
	Minio  MinioConfig
	Bucket string
	Object string
}

// ReadExcelFromMinio
// return: [][]string (raw rows, header included)
func ReadExcelFromMinio(
	ctx context.Context,
	cfg ReadConfig,
) ([][]string, error) {

	if cfg.Bucket == "" || cfg.Object == "" {
		return nil, ErrInvalidPath
	}

	reader, err := getObject(ctx, cfg.Minio, cfg.Bucket, cfg.Object)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrEmptyExcelFile
	}

	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, ErrInvalidExcel
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
