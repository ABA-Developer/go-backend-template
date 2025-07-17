CREATE TABLE IF NOT EXISTS "sessions" (
    "guid" varchar UNIQUE PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    "user_id" integer NOT NULL,
    "access_token" text NOT NULL,
    "access_token_expired_at" timestamp NOT NULL,
    "refresh_token" text NOT NULL,
    "refresh_token_expired_at" timestamp NOT NULL,
    "ip_address" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "created_at" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp WITH TIME ZONE
);