CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(100),
    "middle_name" VARCHAR(100),
    "last_name" VARCHAR(100),
    "email" VARCHAR(100) NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "role" VARCHAR(100) DEFAULT 'user',
    "is_active" BOOLEAN DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "created_by" INTEGER,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "updated_by" INTEGER
);

INSERT INTO 
    users (first_name, middle_name, last_name, email, password, role)
VALUES 
    (
        'Nusapala', 'Berkah', 'Autonomous', 'user@gmail.com', '$2a$12$LQi1CpKB/dUNMKko2sHd/.umM9hdOYSoMRF7b8JbgiV3ZvSWIEqQC', 'user'
    )
ON CONFLICT DO NOTHING;