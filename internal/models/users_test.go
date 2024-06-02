package models

import (
	"testing"

	"snippetbox.micypac.io/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test for UserModel")
	}

	// Set up a suite of table-driven test and expected results.
	tests := []struct{
		name string
		userID int
		want bool
	}{
		{
			name: "Valid ID",
			userID: 1,
			want: true,
		},
		{
			name: "Zero ID",
			userID: 0,
			want: false,
		},
		{
			name: "Non-existent ID",
			userID: 2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() to get a connection pool to test our db.
			// Calling it here, inside t.Run(), means that fresh db tables and data will be set up 
			// and torn down for each test item.
			db := newTestDB(t)

			// Create a new instance of UserModel.
			m := UserModel{db}

			// Call the Exists() method to check return value and error match expected value
			// for the sub-test.
			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
