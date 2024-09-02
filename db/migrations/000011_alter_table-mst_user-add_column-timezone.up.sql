ALTER TABLE mst_user ADD COLUMN timezone text;
UPDATE mst_user SET timezone = 'Asia/Jakarta';
ALTER TABLE mst_user ALTER COLUMN timezone SET NOT NULL;
