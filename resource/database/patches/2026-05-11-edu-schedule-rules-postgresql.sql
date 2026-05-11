-- PostgreSQL incremental patch for education schedule rules and lessons.
-- Safe to execute more than once.

BEGIN;

CREATE TABLE IF NOT EXISTS edu_schedule_rule (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rule_type VARCHAR(32) NOT NULL,
    repeat_type VARCHAR(32) NOT NULL,
    term_id INTEGER DEFAULT 0,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    class_id INTEGER DEFAULT 0,
    student_id INTEGER DEFAULT 0,
    course_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    teaching_mode VARCHAR(32) DEFAULT 'offline',
    requires_room SMALLINT DEFAULT 1,
    room_id INTEGER DEFAULT 0,
    weekday SMALLINT DEFAULT 0,
    start_time VARCHAR(16) NOT NULL,
    end_time VARCHAR(16) NOT NULL,
    status SMALLINT DEFAULT 1,
    version INTEGER DEFAULT 1,
    effective_from TIMESTAMP,
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
    rule_id INTEGER DEFAULT 0,
    rule_version INTEGER DEFAULT 1,
    lesson_type VARCHAR(32) NOT NULL,
    lesson_date DATE NOT NULL,
    start_time VARCHAR(16) NOT NULL,
    end_time VARCHAR(16) NOT NULL,
    class_id INTEGER DEFAULT 0,
    student_id INTEGER DEFAULT 0,
    course_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    teaching_mode VARCHAR(32) DEFAULT 'offline',
    requires_room SMALLINT DEFAULT 1,
    room_id INTEGER DEFAULT 0,
    status VARCHAR(32) DEFAULT 'scheduled',
    is_manual_adjusted SMALLINT DEFAULT 0,
    source_lesson_id INTEGER DEFAULT 0,
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
    rule_id INTEGER DEFAULT 0,
    action_type VARCHAR(32) NOT NULL,
    before_data TEXT,
    after_data TEXT,
    reason VARCHAR(500) DEFAULT '',
    operator_id INTEGER DEFAULT 0,
    occurred_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_change_log_deleted_at ON edu_lesson_change_log (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_change_log_tenant_lesson ON edu_lesson_change_log (tenant_id, lesson_id);

CREATE TABLE IF NOT EXISTS edu_schedule_conflict_override (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER DEFAULT 0,
    lesson_id INTEGER DEFAULT 0,
    conflict_type VARCHAR(32) NOT NULL,
    conflict_key VARCHAR(128) NOT NULL,
    reason VARCHAR(500) NOT NULL,
    operator_id INTEGER DEFAULT 0,
    occurred_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_deleted_at ON edu_schedule_conflict_override (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_conflict ON edu_schedule_conflict_override (tenant_id, conflict_type);
CREATE INDEX IF NOT EXISTS idx_edu_schedule_conflict_override_tenant_lesson ON edu_schedule_conflict_override (tenant_id, lesson_id);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(282, '排课规则列表', '/api/edu/schedule-rules/list', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(283, '新增排课规则', '/api/edu/schedule-rules/add', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(284, '编辑排课规则', '/api/edu/schedule-rules/edit', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(285, '预览排课变更', '/api/edu/schedule-rules/preview-change', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(286, '删除排课规则', '/api/edu/schedule-rules/delete', 'DELETE', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(287, '课次日历', '/api/edu/lessons/calendar', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(288, '课次列表', '/api/edu/lessons/list', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(289, '检查排课冲突', '/api/edu/schedules/check-conflicts', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(290, '课次调课', '/api/edu/lessons/reschedule', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(291, '课次停课', '/api/edu/lessons/stop', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(292, '课次取消', '/api/edu/lessons/cancel', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(293, '课次恢复', '/api/edu/lessons/restore', 'PUT', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(294, '课次补课', '/api/edu/lessons/makeup', 'POST', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(295, '课次变更历史', '/api/edu/lessons/:id/change-logs', 'GET', '教培管理', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140405, 140350, '/edu/schedule', 'EduSchedule', '', 'edu/schedule/schedule', '排课管理', 0, 0, 0, 1, 0, '', 0, '', 'IconCalendar', 8, 2, 0, 'edu:scheduleRule:list', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140406, 140405, '', '', '', '', '新增排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:scheduleRule:add', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140407, 140405, '', '', '', '', '编辑排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:scheduleRule:edit', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140408, 140405, '', '', '', '', '预览排课变更', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:scheduleRule:previewChange', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140409, 140405, '', '', '', '', '删除排课规则', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:scheduleRule:delete', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140410, 140405, '', '', '', '', '课次日历', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:lesson:calendar', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140411, 140405, '', '', '', '', '课次列表', 0, 0, 0, 1, 0, '', 0, '', '', 6, 3, 0, 'edu:lesson:list', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140412, 140405, '', '', '', '', '检查排课冲突', 0, 0, 0, 1, 0, '', 0, '', '', 7, 3, 0, 'edu:scheduleConflict:check', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140413, 140405, '', '', '', '', '课次调课', 0, 0, 0, 1, 0, '', 0, '', '', 8, 3, 0, 'edu:lesson:reschedule', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140414, 140405, '', '', '', '', '课次停课', 0, 0, 0, 1, 0, '', 0, '', '', 9, 3, 0, 'edu:lesson:stop', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140415, 140405, '', '', '', '', '课次取消', 0, 0, 0, 1, 0, '', 0, '', '', 10, 3, 0, 'edu:lesson:cancel', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140416, 140405, '', '', '', '', '课次恢复', 0, 0, 0, 1, 0, '', 0, '', '', 11, 3, 0, 'edu:lesson:restore', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140417, 140405, '', '', '', '', '课次补课', 0, 0, 0, 1, 0, '', 0, '', '', 12, 3, 0, 'edu:lesson:makeup', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1),
(140418, 140405, '', '', '', '', '课次变更历史', 0, 0, 0, 1, 0, '', 0, '', '', 13, 3, 0, 'edu:lesson:changeLogs', '2026-05-11 13:00:00', '2026-05-11 13:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140405),(140406),(140407),(140408),(140409),(140410),(140411),(140412),
    (140413),(140414),(140415),(140416),(140417),(140418)
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
(140412, 289),
(140413, 290),
(140414, 291),
(140415, 292),
(140416, 293),
(140417, 294),
(140418, 295)
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
    ('/api/edu/schedules/check-conflicts', 'POST'),
    ('/api/edu/lessons/reschedule', 'PUT'),
    ('/api/edu/lessons/stop', 'PUT'),
    ('/api/edu/lessons/cancel', 'PUT'),
    ('/api/edu/lessons/restore', 'PUT'),
    ('/api/edu/lessons/makeup', 'POST'),
    ('/api/edu/lessons/:id/change-logs', 'GET')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140405','140406','140407','140408','140409','140410','140411','140412','140413','140414','140415','140416','140417','140418']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 295), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140418), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
