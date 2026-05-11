-- PostgreSQL incremental patch for tenant-scoped education teacher role config.
-- Safe to execute more than once.

BEGIN;

CREATE TABLE IF NOT EXISTS edu_teacher_role_config (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER NOT NULL,
    CONSTRAINT uk_edu_teacher_role_config_tenant_role UNIQUE (tenant_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_edu_teacher_role_config_deleted_at ON edu_teacher_role_config (deleted_at);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(262, '教师角色配置', '/api/edu/classes/teacher-role-config', 'GET', '教培管理', '2026-05-11 10:00:00', '2026-05-11 10:00:00', NULL, 1),
(263, '保存教师角色配置', '/api/edu/classes/teacher-role-config', 'PUT', '教培管理', '2026-05-11 10:00:00', '2026-05-11 10:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140386, 140350, '/edu/settings/teacher-role', 'EduTeacherRoleConfig', '', 'edu/settings/teacher-role', '教师角色配置', 0, 0, 0, 1, 0, '', 0, '', 'IconSettings', 5, 2, 0, 'edu:teacherRole:view', '2026-05-11 10:00:00', '2026-05-11 10:00:00', NULL, 1),
(140387, 140386, '', '', '', '', '保存教师角色配置', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:teacherRole:save', '2026-05-11 10:00:00', '2026-05-11 10:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES (140386),(140387)) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140386, 262),
(140387, 263)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/classes/teacher-role-config', 'GET'),
    ('/api/edu/classes/teacher-role-config', 'PUT')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140386','140387']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 263), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140387), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
