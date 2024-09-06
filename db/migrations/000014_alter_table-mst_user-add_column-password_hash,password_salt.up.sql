ALTER TABLE mst_user ADD COLUMN password_hash bytea;
ALTER TABLE mst_user ADD COLUMN password_salt bytea;

UPDATE mst_user SET
	  password_hash = decode('00', 'hex')
	, password_salt = decode('00', 'hex') WHERE password_hash IS NULL OR password_salt IS NULL;

ALTER TABLE mst_user ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE mst_user ALTER COLUMN password_salt SET NOT NULL;
