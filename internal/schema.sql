CREATE TABLE IF NOT EXISTS entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    intensity INTEGER NOT NULL CHECK (intensity >= 0 AND intensity <= 10),
    mood TEXT NOT NULL CHECK (mood IN (
        'happy', 'sad', 'angry', 'anxious', 'excited', 
        'calm', 'stressed', 'tired', 'neutral'
    )),
    message TEXT NOT NULL DEFAULT '',
    tags TEXT DEFAULT '[]', -- JSON array stored as text
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- indexing
CREATE INDEX IF NOT EXISTS idx_entries_mood ON entries(mood);
CREATE INDEX IF NOT EXISTS idx_entries_intensity ON entries(intensity);
CREATE INDEX IF NOT EXISTS idx_entries_created_at ON entries(created_at);

-- update updated_at on entry modification // not planning to support entry updation though
CREATE TRIGGER IF NOT EXISTS update_entries_updated_at 
    AFTER UPDATE ON entries
    FOR EACH ROW
BEGIN
    UPDATE entries SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
