-- sqlite3

create table if not exists user (
	id integer primary key autoincrement not null
	, created_at timestamp not null default current_timestamp
	, created_by integer not null
	, updated_at timestamp not null default current_timestamp
	, updated_by integer not null
	, deleted_at timestamp
	, deleted_by integer

	, name text not null
	, username text not null
	, whatsapp_number text not null
);

