package fileimport

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
	filesmodel "kpo-hw-2/internal/files/model"
)

type CSVImporter struct{}

func NewCSVImporter() *CSVImporter {
	return &CSVImporter{}
}

func (i *CSVImporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "csv",
		Title:       "CSV",
		Description: "Импорт данных из CSV-файла.",
		Extension:   "csv",
	}
}

func (i *CSVImporter) Parse(data []byte) (filesmodel.Payload, error) {
	if len(data) == 0 {
		return filesmodel.Payload{}, nil
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1

	var payload filesmodel.Payload
	var line int
	headerSkipped := false

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return filesmodel.Payload{}, err
		}

		line++

		if len(record) == 0 {
			continue
		}

		for i := range record {
			record[i] = strings.TrimSpace(strings.TrimPrefix(record[i], "\ufeff"))
		}

		if !headerSkipped && len(record) > 0 && strings.EqualFold(record[0], "entity") {
			headerSkipped = true
			continue
		}

		entity := strings.ToLower(recordValue(record, 0))
		switch entity {
		case "account":
			account, err := parseAccountRecord(record)
			if err != nil {
				return filesmodel.Payload{}, fmt.Errorf("csv: parse account line %d: %w", line, err)
			}
			payload.Accounts = append(payload.Accounts, account)
		case "category":
			category, err := parseCategoryRecord(record)
			if err != nil {
				return filesmodel.Payload{}, fmt.Errorf("csv: parse category line %d: %w", line, err)
			}
			payload.Categories = append(payload.Categories, category)
		case "operation":
			operation, err := parseOperationRecord(record)
			if err != nil {
				return filesmodel.Payload{}, fmt.Errorf("csv: parse operation line %d: %w", line, err)
			}
			payload.Operations = append(payload.Operations, operation)
		case "":
			continue
		default:
			return filesmodel.Payload{}, fmt.Errorf("csv: unknown entity %q on line %d", record[0], line)
		}
	}

	return payload, nil
}

func recordValue(record []string, idx int) string {
	if idx >= len(record) {
		return ""
	}
	return record[idx]
}

func parseAccountRecord(record []string) (filesmodel.Account, error) {
	id := recordValue(record, 1)
	name := recordValue(record, 2)
	balanceStr := recordValue(record, 4)

	var balance int64
	if balanceStr != "" {
		parsed, err := strconv.ParseInt(balanceStr, 10, 64)
		if err != nil {
			return filesmodel.Account{}, fmt.Errorf("balance: %w", err)
		}
		balance = parsed
	}

	return filesmodel.Account{
		ID:      id,
		Name:    name,
		Balance: balance,
	}, nil
}

func parseCategoryRecord(record []string) (filesmodel.Category, error) {
	return filesmodel.Category{
		ID:   recordValue(record, 1),
		Name: recordValue(record, 2),
		Type: recordValue(record, 3),
	}, nil
}

func parseOperationRecord(record []string) (filesmodel.Operation, error) {
	amountStr := recordValue(record, 7)
	var amount int64
	if amountStr != "" {
		parsed, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return filesmodel.Operation{}, fmt.Errorf("amount: %w", err)
		}
		amount = parsed
	}

	dateStr := recordValue(record, 8)
	var date time.Time
	if dateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return filesmodel.Operation{}, fmt.Errorf("date: %w", err)
		}
		date = parsed
	}

	return filesmodel.Operation{
		ID:            recordValue(record, 1),
		Type:          recordValue(record, 3),
		BankAccountID: recordValue(record, 5),
		CategoryID:    recordValue(record, 6),
		Amount:        amount,
		Date:          date,
		Description:   recordValue(record, 9),
	}, nil
}

var _ fileimport.Importer = (*CSVImporter)(nil)
