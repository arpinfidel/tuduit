-- postgres

create table mst_task (
	id bigserial primary key not null 
	, created_at timestamp not null default current_timestamp
	, created_by bigint not null
	, updated_at timestamp not null default current_timestamp
	, updated_by bigint not null
	, deleted_at timestamp
	, deleted_by bigint

	, user_id bigint
	, task_schedule_id bigint

	, name text not null
	, description text
	, status text not null default 'pending'
	, started_at timestamp
	, completed_at timestamp
	, archived_at timestamp
);
