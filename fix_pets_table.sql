-- Fix pets table - add missing columns

-- Add photo column if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='pets' AND column_name='photo') THEN
        ALTER TABLE pets ADD COLUMN photo TEXT;
    END IF;
END $$;

-- Add curator_id column if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='pets' AND column_name='curator_id') THEN
        ALTER TABLE pets ADD COLUMN curator_id INTEGER REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Add location column if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='pets' AND column_name='location') THEN
        ALTER TABLE pets ADD COLUMN location TEXT;
    END IF;
END $$;

-- Create index on curator_id for better performance
CREATE INDEX IF NOT EXISTS idx_pets_curator_id ON pets(curator_id);
