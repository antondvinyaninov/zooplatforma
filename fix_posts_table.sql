-- Fix posts table structure for PostgreSQL
-- Add missing columns from migration 036

-- Add missing columns to posts table
DO $$ 
BEGIN
    -- Add author_id (copy from user_id if not exists)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'author_id') THEN
        ALTER TABLE posts ADD COLUMN author_id INTEGER;
        -- Copy existing user_id to author_id
        UPDATE posts SET author_id = user_id WHERE author_id IS NULL;
        -- Make it NOT NULL after copying data
        ALTER TABLE posts ALTER COLUMN author_id SET NOT NULL;
        RAISE NOTICE 'Added author_id column and copied data from user_id';
    END IF;

    -- Add author_type
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'author_type') THEN
        ALTER TABLE posts ADD COLUMN author_type TEXT NOT NULL DEFAULT 'user';
        RAISE NOTICE 'Added author_type column';
    END IF;

    -- Add attached_pets (JSON array)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'attached_pets') THEN
        ALTER TABLE posts ADD COLUMN attached_pets TEXT DEFAULT '[]';
        RAISE NOTICE 'Added attached_pets column';
    END IF;

    -- Add attachments (JSON array)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'attachments') THEN
        ALTER TABLE posts ADD COLUMN attachments TEXT DEFAULT '[]';
        -- Copy media_urls to attachments if media_urls exists
        IF EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'media_urls') THEN
            UPDATE posts SET attachments = COALESCE(media_urls, '[]') WHERE attachments = '[]';
        END IF;
        RAISE NOTICE 'Added attachments column';
    END IF;

    -- Add tags (JSON array)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'tags') THEN
        ALTER TABLE posts ADD COLUMN tags TEXT DEFAULT '[]';
        RAISE NOTICE 'Added tags column';
    END IF;

    -- Add is_deleted (soft delete)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' 
                   AND table_name = 'posts' 
                   AND column_name = 'is_deleted') THEN
        ALTER TABLE posts ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
        RAISE NOTICE 'Added is_deleted column';
    END IF;

    -- Make content nullable (it has default '')
    ALTER TABLE posts ALTER COLUMN content DROP NOT NULL;
    ALTER TABLE posts ALTER COLUMN content SET DEFAULT '';
    
    RAISE NOTICE 'Posts table structure fixed successfully';
END $$;

-- Create index on author_id if not exists
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id, author_type);

-- Create index on is_deleted if not exists
CREATE INDEX IF NOT EXISTS idx_posts_is_deleted ON posts(is_deleted);
