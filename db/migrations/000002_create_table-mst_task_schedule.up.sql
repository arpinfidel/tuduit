-- postgres

create table if not exists mst_task_schedule (
	id bigserial primary key not null 
	, created_at timestamp not null default current_timestamp
	, created_by bigint not null
	, updated_at timestamp not null default current_timestamp
	, updated_by bigint not null
	, deleted_at timestamp
	, deleted_by bigint
	
	, user_id bigint
	
	, start_date timestamp
	, schedule text -- cron expression
	, duration bigint -- duration before deadline in seconds
);
