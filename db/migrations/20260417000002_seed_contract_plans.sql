-- +goose Up
-- +goose StatementBegin
INSERT INTO contract_plans (code) VALUES
    ('standard'),
    ('premium');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM contract_plans;
-- +goose StatementEnd
