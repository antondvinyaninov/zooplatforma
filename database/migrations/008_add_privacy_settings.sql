-- Add privacy settings to users table
ALTER TABLE users ADD COLUMN profile_visibility TEXT DEFAULT 'public'; -- public, friends, private
ALTER TABLE users ADD COLUMN show_phone TEXT DEFAULT 'nobody'; -- everyone, friends, nobody
ALTER TABLE users ADD COLUMN show_email TEXT DEFAULT 'nobody'; -- everyone, friends, nobody
ALTER TABLE users ADD COLUMN allow_messages TEXT DEFAULT 'everyone'; -- everyone, friends, nobody
ALTER TABLE users ADD COLUMN show_online TEXT DEFAULT 'yes'; -- yes, no
