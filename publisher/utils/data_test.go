package utils

import (
	"os"
	"testing"
	"time"

	"github.com/irononet/publisher/structs"
	"github.com/stretchr/testify/assert"
)

func TestLoadUserRecords(t *testing.T) {
	csvContent := `1,John,Doe,john.doe@example.com,1637133600,1637220000,1637306400,2
2,Jane,Smith,jane.smith@example.com,1637133600,1637220000,1637306400,0`

	filePath := "testdata/mock_users.csv"
	err := createTempCSV(filePath, csvContent)
	if err != nil {
		t.Fatal("Error creating temporary CSV file:", err)
	}
	defer removeTempFile(filePath)

	// Load user records
	records, err := LoadUserRecords(filePath)

	// define expected user records
	expectedRecords := []*structs.User{
		{
			ID:           1,
			FirstName:    "John",
			LastName:     "Doe",
			EmailAddress: "john.doe@example.com",
			CreatedAt:    time.Unix(1637133600, 0),
			DeletedAt:    time.Unix(1637220000, 0),
			MergedAt:     time.Unix(1637306400, 0),
			ParentUserId: 2,
		}, {
			ID:           2,
			FirstName:    "Jane",
			LastName:     "Smith",
			EmailAddress: "jane.smith@example.com",
			CreatedAt:    time.Unix(1637133600, 0),
			DeletedAt:    time.Unix(1637220000, 0),
			MergedAt:     time.Unix(1637306400, 0),
			ParentUserId: 0,
		},
	}

	// Assertions
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedRecords, records)

}

func createTempCSV(filepath, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func removeTempFile(filepath string) {
	_ = os.Remove(filepath)
}
