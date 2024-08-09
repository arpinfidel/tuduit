-- sqlite3

create table check_in (
	id integer primary key autoincrement not null
	, created_at timestamp not null default current_timestamp
	, created_by integer not null
	, updated_at timestamp not null default current_timestamp
	, updated_by integer not null
	, deleted_at timestamp
	, deleted_by integer

	, user_id integer not null
	, check_in_time time not null
	, last_sent timestamp not null
);
