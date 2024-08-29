ALTER TABLE mst_task ADD COLUMN user_ids bigint[];
UPDATE mst_task SET user_ids = ARRAY[user_id];
-- ALTER TABLE mst_task DROP COLUMN user_id;
ALTER TABLE mst_task ALTER COLUMN user_ids SET NOT NULL;
CREATE INDEX idx_mst_task_user_ids ON mst_task USING gin (user_ids);
