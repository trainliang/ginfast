-- PostgreSQL patch for education lesson completion frontend menu and permissions.

BEGIN;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140425, 140350, '/edu/lesson-completion', 'EduLessonCompletion', '/edu/lesson-completion/rules', '', '消课管理', 0, 0, 0, 1, 0, '', 0, '', 'IconCheckCircle', 13, 1, 0, '', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140426, 140425, '/edu/lesson-completion/rules', 'EduLessonCompletionRules', '', 'edu/lesson-completion/rules', '消课规则', 0, 0, 0, 1, 0, '', 0, '', 'IconSettings', 1, 2, 0, 'edu:lessonCompletion:rule', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140427, 140425, '/edu/lesson-completion/lessons', 'EduLessonCompletionLessons', '', 'edu/lesson-completion/lessons', '课次消课', 0, 0, 0, 1, 0, '', 0, '', 'IconCalendarClock', 2, 2, 0, 'edu:lessonCompletion:view', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140428, 140426, '', '', '', '', '保存消课规则', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:lessonCompletion:rule:save', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140429, 140427, '', '', '', '', '提交学生消课', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:lessonCompletion:submit', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140430, 140427, '', '', '', '', '撤销学生消课', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:lessonCompletion:revoke', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1),
(140431, 140427, '', '', '', '', '整节课完成', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:lessonCompletion:complete', '2026-05-12 10:00:00', '2026-05-12 10:00:00', NULL, 1)
ON CONFLICT (id) DO UPDATE SET
    parent_id = EXCLUDED.parent_id,
    path = EXCLUDED.path,
    name = EXCLUDED.name,
    redirect = EXCLUDED.redirect,
    component = EXCLUDED.component,
    title = EXCLUDED.title,
    icon = EXCLUDED.icon,
    sort = EXCLUDED.sort,
    type = EXCLUDED.type,
    permission = EXCLUDED.permission,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140425),(140426),(140427),(140428),(140429),(140430),(140431)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140426, 299),
(140428, 300),
(140427, 301),
(140429, 302),
(140430, 303),
(140431, 304)
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
            ARRAY['140425','140426','140427','140428','140429','140430','140431']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140431), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
