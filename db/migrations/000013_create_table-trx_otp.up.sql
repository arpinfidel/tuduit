-- postgres

create table trx_otp (
	id bigserial not null
	, created_at timestamptz not null default current_timestamp
	, created_by bigint not null
	, updated_at timestamptz not null default current_timestamp
	, updated_by bigint not null
	, deleted_at timestamptz
	, deleted_by bigint

	, whatsapp_number bigint not null
	, otp text not null
	, token text not null
	, invalidated_at timestamptz not null
) partition by range (invalidated_at);

alter table trx_otp
	add primary key (id, invalidated_at);
