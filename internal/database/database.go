package database

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/paolojulian/wedding-be/internal/models"
)

func ReadInvitations() ([]models.Invitation, error) {
	file, err := readDatabase()
	if err != nil {
		return nil, err
	}

	var database models.Database
	decoder := json.NewDecoder(file)
	defer file.Close()

	if err := decoder.Decode(&database); err != nil {
		return nil, err
	}

	return database.Invitations, nil
}

func ReadUsers() ([]models.User, error) {
	file, err := readDatabase()
	if err != nil {
		return nil, err
	}

	var database models.Database
	decoder := json.NewDecoder(file)
	defer file.Close()

	if err := decoder.Decode(&database); err != nil {
		return nil, err
	}

	return database.Users, nil
}

func getDatabasePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dbPath := filepath.Join(currentDir, "internal", "database", "database.json")

	return dbPath, nil
}

func readDatabase() (*os.File, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// func writeDatabase(database *models.Database) error {
// 	dbPath, err := getDatabsePath()
// 	if err != nil {
// 		return err
// 	}

// 	file, err := os.Create(dbPath)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	if err := encoder.Encode(database); err != nil {
// 		return err
// 	}

// 	return nil
// }
