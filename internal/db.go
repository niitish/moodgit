package internal

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

var db *sql.DB

func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dbPath := path.Join(homeDir, ".moodgit", "moodgit.db")
	db, err = sql.Open("sqlite", dbPath)

	if err != nil {
		return fmt.Errorf("failed to open database.\ndid you run moodgit init?\n%w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("failed to migrate database.\ndid you run moodgit init?\n%w", err)
	}

	return nil
}

func migrate() error {
	schemaBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schemaBytes))
	return err
}

func getFilteredHistory(pageSize int, offset int, filter string, search string) ([]Entry, int, error) {
	var countQuery strings.Builder
	var dataQuery strings.Builder
	var countArgs []interface{}
	var dataArgs []interface{}

	countQuery.WriteString(`
		SELECT COUNT(*)
		FROM entries 
		WHERE 1=1`)

	dataQuery.WriteString(`
		SELECT id, intensity, mood, message, tags, created_at, updated_at
		FROM entries 
		WHERE 1=1`)

	if filter != "all" && filter != "" {
		countQuery.WriteString(" AND mood = ?")
		dataQuery.WriteString(" AND mood = ?")
		countArgs = append(countArgs, filter)
		dataArgs = append(dataArgs, filter)
	}

	if search != "" {
		countQuery.WriteString(" AND (message LIKE ? OR tags LIKE ?)")
		dataQuery.WriteString(" AND (message LIKE ? OR tags LIKE ?)")
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
		dataArgs = append(dataArgs, searchPattern, searchPattern)
	}

	var totalCount int
	err := db.QueryRow(countQuery.String(), countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	dataQuery.WriteString(" ORDER BY created_at DESC LIMIT ? OFFSET ?")
	dataArgs = append(dataArgs, pageSize, offset)

	rows, err := db.Query(dataQuery.String(), dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		var tagsJSON string
		if err := rows.Scan(&entry.ID, &entry.Intensity, &entry.Mood, &entry.Message, &tagsJSON, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return nil, 0, err
		}

		if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, totalCount, nil
}

func AddEntry(entry Entry) error {
	tagsJSON, err := json.Marshal(entry.Tags)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO entries (intensity, mood, message, tags) 
		VALUES (?, ?, ?, ?)`,
		entry.Intensity, entry.Mood, entry.Message, string(tagsJSON))
	return err
}

func AmendLastEntry(entry Entry) error {
	tagsJSON, err := json.Marshal(entry.Tags)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE entries 
		SET intensity = ?, mood = ?, message = ?, tags = ? 
		WHERE id = (SELECT id FROM entries ORDER BY created_at DESC LIMIT 1)`,
		entry.Intensity, entry.Mood, entry.Message, string(tagsJSON))
	return err
}

func GetHistory(limit uint16) error {
	rows, err := db.Query(
		`SELECT id, intensity, mood, message, tags, created_at, updated_at
		FROM entries ORDER BY created_at DESC LIMIT ?`, limit)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		var tagsJSON string
		if err := rows.Scan(&entry.ID, &entry.Intensity, &entry.Mood, &entry.Message, &tagsJSON, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(tagsJSON), &entry.Tags); err != nil {
			return err
		}

		fmt.Println(entry.String())
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
