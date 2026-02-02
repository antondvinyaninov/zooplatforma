#!/bin/bash

# Script to apply SQL fix to EasyPanel PostgreSQL database
# Run this on EasyPanel server

echo "🔧 Applying organizations table fix to PostgreSQL..."

# Database credentials from EasyPanel
DB_HOST="zooplatforma-db"
DB_NAME="zp-db"
DB_USER="zp"
DB_PASSWORD="lmLG7k2ed4vas19"

# Apply SQL fix
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME << 'EOF'
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

-- Verify changes
\d organizations
\d organization_members
EOF

echo "✅ Fix applied successfully!"
