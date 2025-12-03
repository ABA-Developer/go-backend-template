CREATE TABLE IF NOT EXISTS "app_user" (
	id varchar(50) DEFAULT gen_random_uuid() NOT NULL,
	name varchar(50) NOT NULL,
	full_name varchar(255) NOT NULL,
	email varchar(50) NOT NULL,
	password varchar(255) NOT NULL,
	phone varchar(15) NULL,
	active bool DEFAULT true NOT NULL,
	img_path varchar(255) NULL,
	img_name varchar(255) NULL,
	created_by varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_by varchar(255) NULL,
	updated_at timestamp NULL,
	CONSTRAINT app_user_pkey PRIMARY KEY (id),
	CONSTRAINT app_user_user_id_key UNIQUE (id)
);