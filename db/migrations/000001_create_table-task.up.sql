-- sqlite3

create table task (
	id bigserial primary key
	, created_at timestamp not null default current_timestamp
	, updated_at timestamp not null default current_timestamp

	, user_id bigint
	, task_schedule_id bigint

	, name text not null
	, description text
	, status text not null default 'pending'
	, started boolean not null default false
	, started_at timestamp
	, completed boolean not null default false
	, completed_at timestamp
	, archived boolean not null default false
	, archived_at timestamp
);
