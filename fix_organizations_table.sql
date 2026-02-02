-- Fix organizations table structure for PostgreSQL
-- Add missing columns from SQLite schema
-- This script is safe to run multiple times

-- Add missing columns to organizations table (if table exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'organizations') THEN
        -- Add columns one by one with IF NOT EXISTS
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS short_name TEXT;
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS bio TEXT;
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS cover_photo TEXT;
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_city TEXT;
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_region TEXT;
        ALTER TABLE organizations ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE;
        RAISE NOTICE 'Organizations table columns added/verified';
    ELSE
        RAISE NOTICE 'Organizations table does not exist, skipping';
    END IF;
END $$;

-- Add missing column to organization_members table (if table exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'organization_members') THEN
        -- Add column with IF NOT EXISTS
        ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS can_post BOOLEAN DEFAULT FALSE;
        
        -- Update existing records: set can_post=TRUE for owners and admins
        UPDATE organization_members 
        SET can_post = TRUE 
        WHERE role IN ('owner', 'admin') AND (can_post IS NULL OR can_post = FALSE);
        
        RAISE NOTICE 'Organization_members table columns added/verified';
    ELSE
        RAISE NOTICE 'Organization_members table does not exist, skipping';
    END IF;
END $$;
