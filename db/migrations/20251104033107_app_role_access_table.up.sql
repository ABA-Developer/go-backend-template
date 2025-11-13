CREATE TABLE IF NOT EXISTS app_role_access (
    role_id int4 NOT NULL,
    menu_permission_id int4 NOT NULL, 
    CONSTRAINT app_role_access_pkey PRIMARY KEY (role_id, menu_permission_id),
    CONSTRAINT app_role_access_role_id_fkey 
    FOREIGN KEY (role_id) REFERENCES app_role(id) ON DELETE CASCADE,    
    CONSTRAINT app_role_access_menu_perm_id_fkey 
    FOREIGN KEY (menu_permission_id) REFERENCES app_menu_permission(id) ON DELETE CASCADE
);