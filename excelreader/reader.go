package excelreader

import (
	"context"
	"errors"

	"github.com/xuri/excelize/v2"
)

type ReadConfig struct {
	Minio     MinioConfig
	Bucket    string
	Object    string
	MaxColumn int // optional: 0 = no limit
}

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

	// 🔥 IMPORTANT: langsung pakai streaming reader
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, ErrInvalidExcel
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return nil, errors.New("no sheet found")
	}

	rowIter, err := f.Rows(sheet)
	if err != nil {
		return nil, err
	}
	defer rowIter.Close()

	var result [][]string

	for rowIter.Next() {

		// support context cancel
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		cols, err := rowIter.Columns()
		if err != nil {
			return nil, err
		}

		// Optional: limit column read (hemat memory)
		if cfg.MaxColumn > 0 && len(cols) > cfg.MaxColumn {
			cols = cols[:cfg.MaxColumn]
		}

		result = append(result, cols)
	}

	if err := rowIter.Error(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrEmptyExcelFile
	}

	return result, nil
}
