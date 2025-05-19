package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/department"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

func SeedPositions(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "position.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return utils.WrapError("opening file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read the header row to map column names to indices
	header, err := reader.Read()
	if err != nil {
		return utils.WrapError("reading header row", err)
	}

	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}

	records, err := reader.ReadAll()
	if err != nil {
		return utils.WrapError("reading CSV records", err)
	}

	for _, record := range records {
		codeIdx, codeExists := headerMap["code"]
		nameIdx, nameExists := headerMap["name"]
		deptCodeIdx, deptCodeExists := headerMap["department_code"]

		if !codeExists || !nameExists || !deptCodeExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		deptCode := record[deptCodeIdx]
		dept, err := client.Department.Query().Where(department.Code(deptCode)).Only(ctx)
		if err != nil {
			log.Printf("Failed to find department for position %s: %v", record[codeIdx], err)
			continue
		}

		code := record[codeIdx]
		name := record[nameIdx]

		log.Printf("Seeding Position: %s - %s (dept: %s)", code, name, deptCode)

		err = client.Position.Create().
			SetCode(code).
			SetName(name).
			SetDepartmentID(dept.ID).
			OnConflict(sql.ConflictColumns("code", "department_id")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to upsert Position %s: %v", code, err)
		}
	}
	return nil
}
