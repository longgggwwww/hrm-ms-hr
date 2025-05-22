package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/organization"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

func SeedDepartments(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "department.csv")
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
		orgCodeIdx, orgCodeExists := headerMap["org_code"]

		if !codeExists || !nameExists || !orgCodeExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		orgCode := record[orgCodeIdx]
		org, err := client.Organization.Query().
			Where(organization.Code(orgCode)).
			Only(ctx)
		if err != nil {
			log.Printf("Failed to find organization for department %s: %v", record[codeIdx], err)
			continue
		}

		code := record[codeIdx]
		name := record[nameIdx]
		log.Printf("Seeding Department: %s - %s", code, name)

		err = client.Department.Create().
			SetCode(code).
			SetName(name).
			SetOrgID(org.ID).
			OnConflict(sql.ConflictColumns("code", "org_id")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to upsert Department %s: %v", code, err)
		}
	}
	return nil
}
