-- 新增 options API 权限记录
INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(282, '学生选项', '/api/edu/students/options', 'GET', '教培管理', '2026-05-12 00:00:00', '2026-05-12 00:00:00', NULL, 1),
(283, '权益产品选项', '/api/edu/benefit-products/options', 'GET', '教培管理', '2026-05-12 00:00:00', '2026-05-12 00:00:00', NULL, 1),
(284, '班级选项', '/api/edu/classes/options', 'GET', '教培管理', '2026-05-12 00:00:00', '2026-05-12 00:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

-- 将新API分配给超级管理员角色(role_id=1)和租户管理员(role_id=2)
INSERT INTO sys_role_api (role_id, api_id)
SELECT role_id, api_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES (282),(283),(284)) AS apis(api_id)
ON CONFLICT DO NOTHING;
