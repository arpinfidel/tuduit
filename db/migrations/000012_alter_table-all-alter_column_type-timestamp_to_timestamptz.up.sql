set timezone = 'UTC';

alter table mst_task
	alter created_at type timestamptz
	, alter updated_at type timestamptz
	, alter deleted_at type timestamptz
	, alter started_at type timestamptz
	, alter completed_at type timestamptz
	, alter archived_at type timestamptz
	, alter start_date type timestamptz
	, alter end_date type timestamptz
	;

alter table mst_task_schedule
	alter created_at type timestamptz
	, alter updated_at type timestamptz
	, alter deleted_at type timestamptz
	, alter start_date type timestamptz
	, alter end_date type timestamptz
	, alter next_schedule type timestamptz
	;

alter table mst_user
	alter created_at type timestamptz
	, alter updated_at type timestamptz
	, alter deleted_at type timestamptz
	;

alter table mst_check_in
	alter created_at type timestamptz
	, alter updated_at type timestamptz
	, alter deleted_at type timestamptz
	, alter last_sent type timestamptz
	;
