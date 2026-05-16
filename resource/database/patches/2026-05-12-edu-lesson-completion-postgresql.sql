-- PostgreSQL patch for education lesson completion backend.

BEGIN;

CREATE TABLE IF NOT EXISTS edu_lesson_completion_rule (
    id SERIAL PRIMARY KEY,
    result_type VARCHAR(64) NOT NULL,
    deduct_enabled SMALLINT DEFAULT 0,
    deduct_mode VARCHAR(32) DEFAULT 'fixed_count',
    fixed_count INTEGER DEFAULT 1,
    duration_unit_minutes INTEGER DEFAULT 60,
    insufficient_policy VARCHAR(32) DEFAULT 'block',
    status SMALLINT DEFAULT 1,
    remark VARCHAR(500) DEFAULT '',
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS edu_lesson_student_completion (
    id SERIAL PRIMARY KEY,
    lesson_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    class_id INTEGER DEFAULT 0,
    course_id INTEGER DEFAULT 0,
    teacher_id INTEGER DEFAULT 0,
    result_type VARCHAR(64) NOT NULL,
    deduct_enabled SMALLINT DEFAULT 0,
    deduct_count INTEGER DEFAULT 0,
    student_benefit_id INTEGER DEFAULT 0,
    status VARCHAR(32) NOT NULL,
    ledger_id INTEGER DEFAULT 0,
    revoke_ledger_id INTEGER DEFAULT 0,
    reason VARCHAR(500) DEFAULT '',
    operator_id INTEGER DEFAULT 0,
    completed_at TIMESTAMP,
    revoked_at TIMESTAMP,
    revoked_by INTEGER DEFAULT 0,
    revoke_reason VARCHAR(500) DEFAULT '',
    version INTEGER DEFAULT 1,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_edu_lesson_completion_rule_unique
ON edu_lesson_completion_rule (tenant_id, result_type)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_edu_lesson_completion_rule_deleted_at
ON edu_lesson_completion_rule (deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_edu_lesson_student_completion_active_unique
ON edu_lesson_student_completion (tenant_id, lesson_id, student_id)
WHERE deleted_at IS NULL AND status IN ('completed', 'arrears');

CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_completion_deleted_at
ON edu_lesson_student_completion (deleted_at);

CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_completion_lesson
ON edu_lesson_student_completion (tenant_id, lesson_id);

CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_completion_student
ON edu_lesson_student_completion (tenant_id, student_id);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(299, '消课规则列表', '/api/edu/lesson-completion/rules', 'GET', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(300, '保存消课规则', '/api/edu/lesson-completion/rules', 'PUT', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(301, '课次消课状态', '/api/edu/lessons/:id/completions', 'GET', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(302, '提交学生消课', '/api/edu/lessons/:id/completions', 'POST', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(303, '撤销学生消课', '/api/edu/lessons/:id/completions/:studentId/revoke', 'POST', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(304, '整节课完成', '/api/edu/lessons/:id/complete', 'POST', '教培管理', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    path = EXCLUDED.path,
    method = EXCLUDED.method,
    api_group = EXCLUDED.api_group,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140420, 140405, '', '', '', '', '消课规则配置', 0, 0, 0, 1, 0, '', 0, '', '', 14, 3, 0, 'edu:lessonCompletion:rule', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140421, 140405, '', '', '', '', '课次消课状态', 0, 0, 0, 1, 0, '', 0, '', '', 15, 3, 0, 'edu:lessonCompletion:list', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140422, 140405, '', '', '', '', '提交学生消课', 0, 0, 0, 1, 0, '', 0, '', '', 16, 3, 0, 'edu:lessonCompletion:submit', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140423, 140405, '', '', '', '', '撤销学生消课', 0, 0, 0, 1, 0, '', 0, '', '', 17, 3, 0, 'edu:lessonCompletion:revoke', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140424, 140405, '', '', '', '', '整节课完成', 0, 0, 0, 1, 0, '', 0, '', '', 18, 3, 0, 'edu:lessonCompletion:complete', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1)
ON CONFLICT (id) DO UPDATE SET
    parent_id = EXCLUDED.parent_id,
    title = EXCLUDED.title,
    permission = EXCLUDED.permission,
    sort = EXCLUDED.sort,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140420, 299),
(140420, 300),
(140421, 301),
(140422, 302),
(140423, 303),
(140424, 304)
ON CONFLICT DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140420),(140421),(140422),(140423),(140424)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/lesson-completion/rules', 'GET'),
    ('/api/edu/lesson-completion/rules', 'PUT'),
    ('/api/edu/lessons/:id/completions', 'GET'),
    ('/api/edu/lessons/:id/completions', 'POST'),
    ('/api/edu/lessons/:id/completions/:studentId/revoke', 'POST'),
    ('/api/edu/lessons/:id/complete', 'POST')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140420','140421','140422','140423','140424']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 304), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140424), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
