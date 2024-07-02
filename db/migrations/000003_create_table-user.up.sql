-- sqlite3

create table if not exists user (
	id bigserial primary key
	, created_at timestamp not null default current_timestamp
	, updated_at timestamp not null default current_timestamp

	, name text not null
	, username text not null
	, whatsapp_number text not null
);

