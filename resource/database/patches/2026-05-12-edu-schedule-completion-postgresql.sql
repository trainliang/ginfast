-- PostgreSQL completion patch for education schedule operations and statistics.
-- Safe to execute after earlier education patches; fills columns and permissions added after the first schedule patch.

BEGIN;

ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS repeat_type VARCHAR(32) DEFAULT 'weekly';
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS term_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS start_date TIMESTAMP;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS end_date TIMESTAMP;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS class_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS student_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS course_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS teacher_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS teaching_mode VARCHAR(32) DEFAULT 'offline';
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS requires_room SMALLINT DEFAULT 1;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS room_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS weekday SMALLINT DEFAULT 0;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS start_time VARCHAR(16) DEFAULT '';
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS end_time VARCHAR(16) DEFAULT '';
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS effective_from TIMESTAMP;

ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS rule_id INTEGER DEFAULT 0;
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS rule_version INTEGER DEFAULT 1;
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS lesson_type VARCHAR(32) DEFAULT 'class';
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS course_id INTEGER DEFAULT 0;
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS teaching_mode VARCHAR(32) DEFAULT 'offline';
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS requires_room SMALLINT DEFAULT 1;
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS status VARCHAR(32) DEFAULT 'scheduled';
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS is_manual_adjusted SMALLINT DEFAULT 0;
ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS source_lesson_id INTEGER DEFAULT 0;
ALTER TABLE edu_lesson ALTER COLUMN start_time TYPE VARCHAR(16) USING start_time::text;
ALTER TABLE edu_lesson ALTER COLUMN end_time TYPE VARCHAR(16) USING end_time::text;

ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS rule_id INTEGER DEFAULT 0;
ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS action_type VARCHAR(32);
ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS reason VARCHAR(500) DEFAULT '';
ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS operator_id INTEGER DEFAULT 0;
ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS occurred_at TIMESTAMP;
ALTER TABLE edu_lesson_change_log ALTER COLUMN before_data TYPE TEXT USING before_data::text;
ALTER TABLE edu_lesson_change_log ALTER COLUMN after_data TYPE TEXT USING after_data::text;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'edu_lesson_change_log' AND column_name = 'change_type'
    ) THEN
        ALTER TABLE edu_lesson_change_log ALTER COLUMN change_type DROP NOT NULL;
    END IF;
END $$;

ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS rule_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS lesson_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS conflict_key VARCHAR(128) DEFAULT '';
ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS operator_id INTEGER DEFAULT 0;
ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS occurred_at TIMESTAMP;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'edu_schedule_conflict_override' AND column_name = 'target_type'
    ) THEN
        ALTER TABLE edu_schedule_conflict_override ALTER COLUMN target_type DROP NOT NULL;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'edu_schedule_conflict_override' AND column_name = 'target_id'
    ) THEN
        ALTER TABLE edu_schedule_conflict_override ALTER COLUMN target_id DROP NOT NULL;
        ALTER TABLE edu_schedule_conflict_override ALTER COLUMN target_id SET DEFAULT 0;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_edu_schedule_rule_tenant_rule_type ON edu_schedule_rule (tenant_id, rule_type);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_date ON edu_lesson (tenant_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_teacher_date ON edu_lesson (tenant_id, teacher_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_room_date ON edu_lesson (tenant_id, room_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_student_date ON edu_lesson (tenant_id, student_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_class_date ON edu_lesson (tenant_id, class_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_change_log_tenant_lesson ON edu_lesson_change_log (tenant_id, lesson_id);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_conflict ON edu_schedule_conflict_override (tenant_id, conflict_type);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_lesson ON edu_schedule_conflict_override (tenant_id, lesson_id);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(290, '课次调课', '/api/edu/lessons/reschedule', 'PUT', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(291, '课次停课', '/api/edu/lessons/stop', 'PUT', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(292, '课次取消', '/api/edu/lessons/cancel', 'PUT', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(293, '课次恢复', '/api/edu/lessons/restore', 'PUT', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(294, '课次补课', '/api/edu/lessons/makeup', 'POST', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(295, '课次变更历史', '/api/edu/lessons/:id/change-logs', 'GET', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(296, '排课统计', '/api/edu/statistics/schedule', 'GET', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(297, '权益统计', '/api/edu/statistics/benefits', 'GET', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(298, '外部同步统计', '/api/edu/statistics/external-sync', 'GET', '教培管理', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    path = EXCLUDED.path,
    method = EXCLUDED.method,
    api_group = EXCLUDED.api_group,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140413, 140405, '', '', '', '', '课次调课', 0, 0, 0, 1, 0, '', 0, '', '', 8, 3, 0, 'edu:lesson:reschedule', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140414, 140405, '', '', '', '', '课次停课', 0, 0, 0, 1, 0, '', 0, '', '', 9, 3, 0, 'edu:lesson:stop', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140415, 140405, '', '', '', '', '课次取消', 0, 0, 0, 1, 0, '', 0, '', '', 10, 3, 0, 'edu:lesson:cancel', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140416, 140405, '', '', '', '', '课次恢复', 0, 0, 0, 1, 0, '', 0, '', '', 11, 3, 0, 'edu:lesson:restore', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140417, 140405, '', '', '', '', '课次补课', 0, 0, 0, 1, 0, '', 0, '', '', 12, 3, 0, 'edu:lesson:makeup', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140418, 140405, '', '', '', '', '课次变更历史', 0, 0, 0, 1, 0, '', 0, '', '', 13, 3, 0, 'edu:lesson:changeLogs', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1),
(140419, 140350, '/edu/statistics', 'EduStatistics', '', 'edu/statistics/statistics', '运营统计', 0, 0, 0, 1, 0, '', 0, '', 'IconBarChart', 9, 2, 0, 'edu:statistics:view', '2026-05-12 07:00:00', '2026-05-12 07:00:00', NULL, 1)
ON CONFLICT (id) DO UPDATE SET
    parent_id = EXCLUDED.parent_id,
    path = EXCLUDED.path,
    name = EXCLUDED.name,
    redirect = EXCLUDED.redirect,
    component = EXCLUDED.component,
    title = EXCLUDED.title,
    is_full = EXCLUDED.is_full,
    hide = EXCLUDED.hide,
    disable = EXCLUDED.disable,
    keep_alive = EXCLUDED.keep_alive,
    affix = EXCLUDED.affix,
    link = EXCLUDED.link,
    iframe = EXCLUDED.iframe,
    svg_icon = EXCLUDED.svg_icon,
    icon = EXCLUDED.icon,
    sort = EXCLUDED.sort,
    type = EXCLUDED.type,
    is_link = EXCLUDED.is_link,
    permission = EXCLUDED.permission,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

DELETE FROM sys_menu_api WHERE menu_id = 140413 AND api_id IN (291, 292);

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140413, 290),
(140414, 291),
(140415, 292),
(140416, 293),
(140417, 294),
(140418, 295),
(140419, 296),
(140419, 297),
(140419, 298)
ON CONFLICT DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140413),(140414),(140415),(140416),(140417),(140418),(140419)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/lessons/reschedule', 'PUT'),
    ('/api/edu/lessons/stop', 'PUT'),
    ('/api/edu/lessons/cancel', 'PUT'),
    ('/api/edu/lessons/restore', 'PUT'),
    ('/api/edu/lessons/makeup', 'POST'),
    ('/api/edu/lessons/:id/change-logs', 'GET'),
    ('/api/edu/statistics/schedule', 'GET'),
    ('/api/edu/statistics/benefits', 'GET'),
    ('/api/edu/statistics/external-sync', 'GET')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140413','140414','140415','140416','140417','140418','140419']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 298), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140419), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
