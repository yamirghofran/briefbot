-- +goose Up
-- Add processing status columns to items table
ALTER TABLE items ADD COLUMN processing_status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE items ADD COLUMN processing_error TEXT;

-- Create index for efficient job queue queries
CREATE INDEX idx_items_processing_status ON items(processing_status);
CREATE INDEX idx_items_pending_processing ON items(processing_status, created_at) WHERE processing_status = 'pending';

-- Add check constraint for valid processing statuses
ALTER TABLE items ADD CONSTRAINT check_processing_status 
    CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed'));

-- +goose Down
-- Remove processing status columns and indexes
DROP INDEX IF EXISTS idx_items_pending_processing;
DROP INDEX IF EXISTS idx_items_processing_status;
ALTER TABLE items DROP CONSTRAINT IF EXISTS check_processing_status;
ALTER TABLE items DROP COLUMN IF EXISTS processing_status;
ALTER TABLE items DROP COLUMN IF EXISTS processing_error;