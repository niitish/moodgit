package internal

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

var db *sql.DB

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dbPath := path.Join(homeDir, ".moodgit", "moodgit.db")
	db, err = sql.Open("sqlite", dbPath)

	if err != nil {
		panic(err)
	}

	if err := migrate(); err != nil {
		panic(err)
	}
}

func migrate() error {
	schemaBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schemaBytes))
	return err
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
