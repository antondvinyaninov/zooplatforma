-- Apply all PostgreSQL migrations
-- This script creates all tables needed for the application

-- Run this on EasyPanel PostgreSQL database:
-- PGPASSWORD=lmLG7k2ed4vas19 psql -h zooplatforma-db -U zp -d zp-db -f apply_all_migrations.sql

\echo 'Starting PostgreSQL migrations...'

-- Include the full PostgreSQL migration
\i database/migrations/036_migrate_to_postgresql.sql

\echo 'Migrations completed!'
\echo 'Verifying tables...'

-- List all tables
\dt

\echo 'Done!'
