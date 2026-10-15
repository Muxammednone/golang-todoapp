ALTER TABLE todoapp.users
    ADD COLUMN sso_user_id BIGINT,
    ADD CONSTRAINT users_sso_user_id_key UNIQUE (sso_user_id),
    ADD CONSTRAINT users_sso_user_id_positive CHECK (sso_user_id > 0);
