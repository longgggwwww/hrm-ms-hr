package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/organization"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

func SeedLabels(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "label.csv")
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
		nameIdx, nameExists := headerMap["name"]
		descriptionIdx := headerMap["description"]
		colorIdx, colorExists := headerMap["color"]
		orgCodeIdx, orgCodeExists := headerMap["org_code"]

		if !nameExists || !colorExists || !orgCodeExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		name := record[nameIdx]
		var description *string
		if descriptionIdx < len(record) && record[descriptionIdx] != "" {
			description = &record[descriptionIdx]
		}
		color := record[colorIdx]
		orgCode := record[orgCodeIdx]

		// Find organization by code
		org, err := client.Organization.Query().Where(organization.Code(orgCode)).Only(ctx)
		if err != nil {
			log.Printf("Failed to find organization with code %s: %v", orgCode, err)
			continue
		}

		log.Printf("Seeding Label: %s - %s for org %s", name, color, orgCode)
		_, err = client.Label.Create().
			SetName(name).
			SetNillableDescription(description).
			SetColor(color).
			SetOrgID(org.ID).
			Save(ctx)
		if err != nil {
			log.Printf("Failed to create Label %s: %v", name, err)
		}
	}
	return nil
}
