-- +goose Up
-- +goose StatementBegin
CREATE TABLE immediate_available_women (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    store_id   INT NOT NULL COMMENT '店舗ID',
    woman_id   INT NOT NULL COMMENT '女性ID',
    expires_at DATETIME NOT NULL COMMENT '有効期限',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    PRIMARY KEY (id),
    INDEX idx_store_expires (store_id, expires_at),
    INDEX idx_expires_at    (expires_at),
    CONSTRAINT fk_iaw_store  FOREIGN KEY (store_id)  REFERENCES stores (id),
    CONSTRAINT fk_iaw_woman  FOREIGN KEY (woman_id)  REFERENCES women  (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'すぐに予約が取れる女性';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE immediate_available_women;
-- +goose StatementEnd
