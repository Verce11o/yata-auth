-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verification_codes (
    id SERIAL PRIMARY KEY ,
    type varchar(255) NOT NULL ,
    code UUID DEFAULT uuid_generate_v4(),
    user_id uuid,
    expire_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '1 day',
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
