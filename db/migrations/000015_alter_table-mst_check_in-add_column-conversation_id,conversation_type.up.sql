ALTER TABLE mst_check_in ADD COLUMN conversation_id text;
ALTER TABLE mst_check_in ADD COLUMN conversation_type text;
UPDATE mst_check_in SET conversation_id = '';
UPDATE mst_check_in SET conversation_type = '';
ALTER TABLE mst_check_in ALTER COLUMN conversation_id SET NOT NULL;
ALTER TABLE mst_check_in ALTER COLUMN conversation_type SET NOT NULL;
