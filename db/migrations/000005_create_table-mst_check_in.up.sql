-- postgres

create table mst_check_in (
	id bigserial primary key not null 
	, created_at timestamp not null default current_timestamp
	, created_by bigint not null
	, updated_at timestamp not null default current_timestamp
	, updated_by bigint not null
	, deleted_at timestamp
	, deleted_by bigint

	, user_id bigint not null
	
	, check_in_time time not null
	, last_sent timestamp not null
);
