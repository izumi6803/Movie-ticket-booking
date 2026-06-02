-- Migration: Rename hall_id to screen_id in showtimes table
-- This fixes the schema mismatch from the old hall-based structure

DO $$
BEGIN
    -- Check if hall_id column exists
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'showtimes' 
        AND column_name = 'hall_id'
    ) THEN
        -- Drop the old screen_id if it exists (from failed migration)
        IF EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'showtimes' 
            AND column_name = 'screen_id'
        ) THEN
            ALTER TABLE showtimes DROP COLUMN screen_id;
        END IF;
        
        -- Rename hall_id to screen_id
        ALTER TABLE showtimes RENAME COLUMN hall_id TO screen_id;
        
        -- Update foreign key constraint
        ALTER TABLE showtimes DROP CONSTRAINT IF EXISTS showtimes_hall_id_fkey;
        ALTER TABLE showtimes ADD CONSTRAINT showtimes_screen_id_fkey 
            FOREIGN KEY (screen_id) REFERENCES screens(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Also ensure available_seats column exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'showtimes' 
        AND column_name = 'available_seats'
    ) THEN
        ALTER TABLE showtimes ADD COLUMN available_seats INTEGER NOT NULL DEFAULT 0;
    END IF;
END $$;
