-- sqlite3

create table if not exists schedule (
	id bigserial primary key
	, created_at timestamp not null default current_timestamp
	, updated_at timestamp not null default current_timestamp
	
	, user_id bigint
	
	, start_date timestamp
	, schedule text -- cron expression
	, duration bigint -- duration before deadline in seconds
);
