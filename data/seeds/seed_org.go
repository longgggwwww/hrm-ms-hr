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

func SeedOrganizations(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "org.csv")
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
		logoIdx := headerMap["logo_url"]
		addressIdx := headerMap["address"]
		phoneIdx := headerMap["phone"]
		emailIdx := headerMap["email"]
		websiteIdx := headerMap["website"]
		parentCodeIdx := headerMap["parent_code"]

		if !codeExists || !nameExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		code := record[codeIdx]
		name := record[nameIdx]
		var logo *string
		if logoIdx < len(record) && record[logoIdx] != "" {
			logo = &record[logoIdx]
		}
		var address *string
		if addressIdx < len(record) && record[addressIdx] != "" {
			address = &record[addressIdx]
		}
		var phone *string
		if phoneIdx < len(record) && record[phoneIdx] != "" {
			phone = &record[phoneIdx]
		}
		var email *string
		if emailIdx < len(record) && record[emailIdx] != "" {
			email = &record[emailIdx]
		}
		var website *string
		if websiteIdx < len(record) && record[websiteIdx] != "" {
			website = &record[websiteIdx]
		}
		var parentID *int = nil
		if parentCodeIdx < len(record) {
			parentCode := record[parentCodeIdx]
			if parentCode != "" {
				parentOrg, err := client.Organization.Query().Where(organization.Code(parentCode)).Only(ctx)
				if err == nil {
					pid := parentOrg.ID
					parentID = &pid
				}
			}
		}

		log.Printf("Seeding Organization: %s - %s", code, name)
		create := client.Organization.Create().
			SetCode(code).
			SetName(name).
			SetNillableLogoURL(logo).
			SetNillableAddress(address).
			SetNillablePhone(phone).
			SetNillableEmail(email).
			SetNillableWebsite(website).
			SetNillableParentID(parentID)
		err := create.OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to upsert Organization %s: %v", code, err)
		}
	}
	return nil
}
