-- +goose Up
-- +goose StatementBegin
INSERT INTO business_types (code) VALUES
    ('ソープ'),
    ('デリヘル'),
    ('ファッションヘルス'),
    ('イメクラ'),
    ('ピンサロ');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM business_types;
-- +goose StatementEnd
