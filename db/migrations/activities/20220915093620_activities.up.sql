BEGIN;

CREATE TABLE IF NOT EXISTS main.activities
(
    id          UUID         NOT NULL,
    user_id     UUID         NOT NULL,
    title       VARCHAR(128)         NOT NULL,
    description TEXT         NOT NULL,
    activities_type        VARCHAR(50)  NOT NULL,
    created_by  VARCHAR(128) NOT NULL,
    updated_by  VARCHAR(128) NOT NULL,
    deleted_by  VARCHAR(128),
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL,
    deleted_at  TIMESTAMPTZ,
    PRIMARY KEY (id)
    );

COMMIT;