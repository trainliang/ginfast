-- PostgreSQL incremental patch for education schedule rules and lessons.
-- Safe to execute more than once.

BEGIN;

CREATE TABLE IF NOT EXISTS edu_schedule_rule (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rule_type VARCHAR(32) NOT NULL,
    status SMALLINT DEFAULT 1,
    rule_content JSONB DEFAULT '{}'::jsonb,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_rule_deleted_at ON edu_schedule_rule (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_rule_tenant_rule_type ON edu_schedule_rule (tenant_id, rule_type);

CREATE TABLE IF NOT EXISTS edu_lesson (
    id SERIAL PRIMARY KEY,
    schedule_rule_id INTEGER DEFAULT 0,
    lesson_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    teacher_id INTEGER DEFAULT 0,
    room_id INTEGER DEFAULT 0,
    student_id INTEGER DEFAULT 0,
    class_id INTEGER DEFAULT 0,
    lesson_status VARCHAR(32) DEFAULT 'planned',
    source_type VARCHAR(32) DEFAULT '',
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_deleted_at ON edu_lesson (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_date ON edu_lesson (tenant_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_teacher_date ON edu_lesson (tenant_id, teacher_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_room_date ON edu_lesson (tenant_id, room_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_student_date ON edu_lesson (tenant_id, student_id, lesson_date);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_tenant_class_date ON edu_lesson (tenant_id, class_id, lesson_date);

CREATE TABLE IF NOT EXISTS edu_lesson_change_log (
    id SERIAL PRIMARY KEY,
    lesson_id INTEGER NOT NULL,
    change_type VARCHAR(32) NOT NULL,
    before_data JSONB DEFAULT '{}'::jsonb,
    after_data JSONB DEFAULT '{}'::jsonb,
    change_reason VARCHAR(500) DEFAULT '',
    changed_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_change_log_deleted_at ON edu_lesson_change_log (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_change_log_tenant_lesson ON edu_lesson_change_log (tenant_id, lesson_id);

CREATE TABLE IF NOT EXISTS edu_schedule_conflict_override (
    id SERIAL PRIMARY KEY,
    conflict_type VARCHAR(32) NOT NULL,
    target_type VARCHAR(32) NOT NULL,
    target_id INTEGER NOT NULL,
    lesson_date DATE,
    reason VARCHAR(500) DEFAULT '',
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_deleted_at ON edu_schedule_conflict_override (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_conflict ON edu_schedule_conflict_override (tenant_id, conflict_type);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_target ON edu_schedule_conflict_override (tenant_id, target_type, target_id);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(282, '排课规则列表', '/api/edu/schedule-rules/list', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(283, '新增排课规则', '/api/edu/schedule-rules/add', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(284, '编辑排课规则', '/api/edu/schedule-rules/edit', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(285, '预览排课变更', '/api/edu/schedule-rules/preview-change', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(286, '删除排课规则', '/api/edu/schedule-rules/delete', 'DELETE', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(287, '课次日历', '/api/edu/lessons/calendar', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(288, '课次列表', '/api/edu/lessons/list', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(289, '检查排课冲突', '/api/edu/schedules/check-conflicts', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140405, 140350, '/edu/schedule', 'EduSchedule', '', 'edu/schedule/schedule', '排课管理', 0, 0, 0, 1, 0, '', 0, '', 'IconCalendar', 8, 2, 0, 'edu:scheduleRule:list', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140406, 140405, '', '', '', '', '新增排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:scheduleRule:add', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140407, 140405, '', '', '', '', '编辑排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:scheduleRule:edit', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140408, 140405, '', '', '', '', '预览排课变更', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:scheduleRule:previewChange', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140409, 140405, '', '', '', '', '删除排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:scheduleRule:delete', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140410, 140405, '', '', '', '', '课次日历', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:lesson:calendar', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140411, 140405, '', '', '', '', '课次列表', 0, 0, 0, 1, 0, '', 0, '', '', 6, 3, 0, 'edu:lesson:list', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140412, 140405, '', '', '', '', '检查排课冲突', 0, 0, 0, 1, 0, '', 0, '', '', 7, 3, 0, 'edu:scheduleConflict:check', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140405),(140406),(140407),(140408),(140409),(140410),(140411),(140412)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140405, 282),
(140406, 283),
(140407, 284),
(140408, 285),
(140409, 286),
(140410, 287),
(140411, 288),
(140412, 289)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/schedule-rules/list', 'GET'),
    ('/api/edu/schedule-rules/add', 'POST'),
    ('/api/edu/schedule-rules/edit', 'PUT'),
    ('/api/edu/schedule-rules/preview-change', 'POST'),
    ('/api/edu/schedule-rules/delete', 'DELETE'),
    ('/api/edu/lessons/calendar', 'GET'),
    ('/api/edu/lessons/list', 'GET'),
    ('/api/edu/schedules/check-conflicts', 'POST')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140405','140406','140407','140408','140409','140410','140411','140412']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 289), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140412), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
