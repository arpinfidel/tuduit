-- postgres

create table if not exists "mst_user" (
	id bigserial primary key not null 
	, created_at timestamp not null default current_timestamp
	, created_by bigint not null
	, updated_at timestamp not null default current_timestamp
	, updated_by bigint not null
	, deleted_at timestamp
	, deleted_by bigint

	, name text not null
	, username text not null
	, whatsapp_number text not null
);
