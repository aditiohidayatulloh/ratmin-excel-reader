package excelreader

import "errors"

var (
	ErrInvalidPath    = errors.New("invalid minio object path")
	ErrEmptyExcelFile = errors.New("excel file is empty")
	ErrInvalidExcel   = errors.New("invalid excel content")
)
