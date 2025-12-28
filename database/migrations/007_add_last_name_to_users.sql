-- Add last_name column to users table
ALTER TABLE users ADD COLUMN last_name TEXT DEFAULT '';

-- Update existing users: split name into first_name and last_name
-- This is a simple split - you may want to adjust based on your data
UPDATE users 
SET last_name = CASE 
    WHEN instr(name, ' ') > 0 THEN substr(name, instr(name, ' ') + 1)
    ELSE ''
END,
name = CASE 
    WHEN instr(name, ' ') > 0 THEN substr(name, 1, instr(name, ' ') - 1)
    ELSE name
END
WHERE name LIKE '% %';
