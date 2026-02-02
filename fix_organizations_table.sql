-- Fix organizations table structure for PostgreSQL
-- Add missing columns from SQLite schema

-- Add missing columns to organizations table
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS short_name TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS cover_photo TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_city TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_region TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE;

-- Add missing column to organization_members table
ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS can_post BOOLEAN DEFAULT FALSE;

-- Update existing records: set can_post=TRUE for owners and admins
UPDATE organization_members SET can_post = TRUE WHERE role IN ('owner', 'admin');
