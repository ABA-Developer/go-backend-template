CREATE TABLE IF NOT EXISTS "sessions" (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	user_id varchar(50) NOT NULL,
	access_token text NOT NULL,
	access_token_expired_at timestamptz NOT NULL,
	refresh_token text NOT NULL,
	refresh_token_expired_at timestamptz NOT NULL,
	ip_address varchar(45) NOT NULL,
	user_agent varchar(255) NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz NULL,
	CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES app_user(id) ON DELETE CASCADE
);