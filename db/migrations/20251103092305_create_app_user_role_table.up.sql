CREATE TABLE IF NOT EXISTS app_user_role (
    user_id varchar(50) NOT NULL,
    role_id int4 NOT NULL,
    
    CONSTRAINT app_user_role_pkey PRIMARY KEY (user_id, role_id),
    CONSTRAINT app_user_role_user_id_fkey FOREIGN KEY (user_id) REFERENCES app_user(id) ON DELETE CASCADE,
    CONSTRAINT app_user_role_role_id_fkey FOREIGN KEY (role_id) REFERENCES app_role(id) ON DELETE CASCADE
);