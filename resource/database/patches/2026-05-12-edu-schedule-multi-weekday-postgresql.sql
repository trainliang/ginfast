BEGIN;

CREATE TABLE IF NOT EXISTS edu_schedule_rule_weekday (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT NOT NULL,
    weekday SMALLINT NOT NULL,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at BIGINT DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_edu_schedule_rule_weekday_unique
ON edu_schedule_rule_weekday (tenant_id, rule_id, weekday)
WHERE deleted_at = 0;

INSERT INTO edu_schedule_rule_weekday (rule_id, weekday, tenant_id, created_at, updated_at, deleted_at)
SELECT id, weekday, tenant_id, created_at, updated_at, deleted_at
FROM edu_schedule_rule
WHERE weekday > 0
ON CONFLICT DO NOTHING;

ALTER TABLE edu_schedule_rule DROP COLUMN IF EXISTS weekday;

COMMIT;
