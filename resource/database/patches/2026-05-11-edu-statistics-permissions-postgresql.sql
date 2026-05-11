-- PostgreSQL incremental patch for education statistics permissions.
-- Safe to execute more than once.

BEGIN;

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(290, '排课统计', '/api/edu/statistics/schedule', 'GET', '教培管理', '2026-05-11 14:00:00', '2026-05-11 14:00:00', NULL, 1),
(291, '权益统计', '/api/edu/statistics/benefits', 'GET', '教培管理', '2026-05-11 14:00:00', '2026-05-11 14:00:00', NULL, 1),
(292, '外部同步统计', '/api/edu/statistics/external-sync', 'GET', '教培管理', '2026-05-11 14:00:00', '2026-05-11 14:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140413, 140350, '/edu/statistics', 'EduStatistics', '', 'edu/statistics/statistics', '运营统计', 0, 0, 0, 1, 0, '', 0, '', 'IconBarChart', 9, 2, 0, 'edu:statistics:view', '2026-05-11 14:00:00', '2026-05-11 14:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, 140413
FROM (VALUES (1),(2)) AS roles(role_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140413, 290),
(140413, 291),
(140413, 292)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
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
            ARRAY['140413']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 292), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140413), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
