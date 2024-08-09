-- sqlite3

create table task (
	id integer primary key autoincrement not null
	, created_at timestamp not null default current_timestamp
	, created_by integer not null
	, updated_at timestamp not null default current_timestamp
	, updated_by integer not null
	, deleted_at timestamp
	, deleted_by integer

	, user_id integer
	, task_schedule_id integer

	, name text not null
	, description text
	, status text not null default 'pending'
	, started_at timestamp
	, completed_at timestamp
	, archived_at timestamp
);
