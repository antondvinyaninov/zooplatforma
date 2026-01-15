-- Migration: Create treatments table
-- Description: Table for storing pet treatments (flea, tick, worm treatments) entered by owners

CREATE TABLE IF NOT EXISTS treatments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    date DATETIME NOT NULL,
    medication TEXT NOT NULL,
    dosage TEXT NOT NULL,
    next_date DATETIME,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_treatments_pet_id ON treatments(pet_id);
CREATE INDEX IF NOT EXISTS idx_treatments_user_id ON treatments(user_id);
CREATE INDEX IF NOT EXISTS idx_treatments_next_date ON treatments(next_date);
CREATE INDEX IF NOT EXISTS idx_treatments_date ON treatments(date);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_treatments_timestamp 
AFTER UPDATE ON treatments
FOR EACH ROW
BEGIN
    UPDATE treatments SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
