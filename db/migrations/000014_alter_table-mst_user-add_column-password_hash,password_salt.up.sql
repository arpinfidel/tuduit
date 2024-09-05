ALTER TABLE mst_user ADD COLUMN password_hash bytea not null;
ALTER TABLE mst_user ADD COLUMN password_salt bytea not null;
