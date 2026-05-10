-- PostgreSQL incremental patch for the education foundation module.
-- Safe to execute more than once.

BEGIN;

CREATE TABLE IF NOT EXISTS edu_course (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(64) NOT NULL,
    type VARCHAR(64) DEFAULT '',
    grade_range VARCHAR(128) DEFAULT '',
    status SMALLINT DEFAULT 1,
    sort INTEGER DEFAULT 0,
    description VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0,
    CONSTRAINT uk_edu_course_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX IF NOT EXISTS idx_edu_course_deleted_at ON edu_course (deleted_at);

CREATE TABLE IF NOT EXISTS edu_student (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(32) DEFAULT '',
    birthday TIMESTAMP,
    phone VARCHAR(64) DEFAULT '',
    status SMALLINT DEFAULT 1,
    avatar VARCHAR(255) DEFAULT '',
    school VARCHAR(128) DEFAULT '',
    grade VARCHAR(64) DEFAULT '',
    school_class VARCHAR(64) DEFAULT '',
    source_channel VARCHAR(64) DEFAULT '',
    enroll_date TIMESTAMP,
    health_note VARCHAR(500) DEFAULT '',
    allergy_note VARCHAR(500) DEFAULT '',
    emergency_contact VARCHAR(128) DEFAULT '',
    pickup_note VARCHAR(500) DEFAULT '',
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_student_deleted_at ON edu_student (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_student_tenant_name ON edu_student (tenant_id, name);
CREATE INDEX IF NOT EXISTS idx_edu_student_tenant_phone ON edu_student (tenant_id, phone);

CREATE TABLE IF NOT EXISTS edu_student_contact (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    relation VARCHAR(64) DEFAULT '',
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(64) NOT NULL,
    is_primary SMALLINT DEFAULT 0,
    can_pickup SMALLINT DEFAULT 0,
    remark VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_student_contact_deleted_at ON edu_student_contact (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_student_contact_student ON edu_student_contact (tenant_id, student_id);
CREATE INDEX IF NOT EXISTS idx_edu_student_contact_phone ON edu_student_contact (tenant_id, phone);

CREATE TABLE IF NOT EXISTS edu_class (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(64) NOT NULL,
    class_type VARCHAR(32) NOT NULL,
    course_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    room_id INTEGER DEFAULT 0,
    capacity INTEGER NOT NULL,
    status SMALLINT DEFAULT 1,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0,
    CONSTRAINT uk_edu_class_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX IF NOT EXISTS idx_edu_class_deleted_at ON edu_class (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_class_tenant_course ON edu_class (tenant_id, course_id);
CREATE INDEX IF NOT EXISTS idx_edu_class_tenant_teacher ON edu_class (tenant_id, teacher_id);
CREATE INDEX IF NOT EXISTS idx_edu_class_tenant_room ON edu_class (tenant_id, room_id);

CREATE TABLE IF NOT EXISTS edu_class_member (
    id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    join_date TIMESTAMP,
    leave_date TIMESTAMP,
    status VARCHAR(32) NOT NULL,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_class_member_deleted_at ON edu_class_member (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_class_member_class ON edu_class_member (tenant_id, class_id);
CREATE INDEX IF NOT EXISTS idx_edu_class_member_student ON edu_class_member (tenant_id, student_id);

CREATE TABLE IF NOT EXISTS edu_room (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(64) NOT NULL,
    type VARCHAR(64) DEFAULT '',
    capacity INTEGER NOT NULL,
    location VARCHAR(255) DEFAULT '',
    status SMALLINT DEFAULT 1,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0,
    CONSTRAINT uk_edu_room_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX IF NOT EXISTS idx_edu_room_deleted_at ON edu_room (deleted_at);

CREATE TABLE IF NOT EXISTS edu_room_weekly_rule (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL,
    weekday SMALLINT NOT NULL,
    start_time VARCHAR(5) NOT NULL,
    end_time VARCHAR(5) NOT NULL,
    available SMALLINT DEFAULT 1,
    remark VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_room_weekly_rule_deleted_at ON edu_room_weekly_rule (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_room_weekly_room ON edu_room_weekly_rule (tenant_id, room_id, weekday);

CREATE TABLE IF NOT EXISTS edu_room_exception (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL,
    exception_date TIMESTAMP NOT NULL,
    type VARCHAR(32) NOT NULL,
    start_time VARCHAR(5) DEFAULT '',
    end_time VARCHAR(5) DEFAULT '',
    reason VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_room_exception_deleted_at ON edu_room_exception (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_room_exception_room_date ON edu_room_exception (tenant_id, room_id, exception_date);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(217, '学生列表', '/api/edu/students/list', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(218, '学生详情', '/api/edu/students/:id', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(219, '新增学生', '/api/edu/students/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(220, '编辑学生', '/api/edu/students/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(221, '删除学生', '/api/edu/students/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(222, '导入学生', '/api/edu/students/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(223, '导出学生', '/api/edu/students/export', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(224, '课程/项目列表', '/api/edu/courses/list', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(225, '课程/项目选项', '/api/edu/courses/options', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(226, '课程/项目详情', '/api/edu/courses/:id', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(227, '新增课程/项目', '/api/edu/courses/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(228, '编辑课程/项目', '/api/edu/courses/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(229, '删除课程/项目', '/api/edu/courses/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(230, '导入课程/项目', '/api/edu/courses/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(231, '导出课程/项目', '/api/edu/courses/export', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(232, '班级列表', '/api/edu/classes/list', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(233, '教师选项', '/api/edu/classes/teacher-options', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(234, '班级详情', '/api/edu/classes/:id', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(235, '新增班级', '/api/edu/classes/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(236, '编辑班级', '/api/edu/classes/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(237, '删除班级', '/api/edu/classes/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(238, '班级成员列表', '/api/edu/classes/:id/members', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(239, '新增班级成员', '/api/edu/classes/:id/members/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(240, '编辑班级成员', '/api/edu/classes/:id/members/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(241, '删除班级成员', '/api/edu/classes/:id/members/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(242, '导入班级', '/api/edu/classes/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(243, '导出班级', '/api/edu/classes/export', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(244, '导入班级成员', '/api/edu/classes/members/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(245, '导出班级成员', '/api/edu/classes/members/export', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(246, '教室/场地列表', '/api/edu/rooms/list', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(247, '教室/场地选项', '/api/edu/rooms/options', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(248, '教室/场地详情', '/api/edu/rooms/:id', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(249, '新增教室/场地', '/api/edu/rooms/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(250, '编辑教室/场地', '/api/edu/rooms/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(251, '删除教室/场地', '/api/edu/rooms/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(252, '场地周规则列表', '/api/edu/rooms/:id/weekly-rules', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(253, '保存场地周规则', '/api/edu/rooms/:id/weekly-rules/save', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(254, '场地例外列表', '/api/edu/rooms/:id/exceptions', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(255, '新增场地例外', '/api/edu/rooms/:id/exceptions/add', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(256, '编辑场地例外', '/api/edu/rooms/:id/exceptions/edit', 'PUT', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(257, '删除场地例外', '/api/edu/rooms/:id/exceptions/delete', 'DELETE', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(258, '导入教室/场地', '/api/edu/rooms/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(259, '导出教室/场地', '/api/edu/rooms/export', 'GET', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(260, '导入场地周规则', '/api/edu/rooms/weekly-rules/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(261, '导入场地例外', '/api/edu/rooms/exceptions/import', 'POST', '教培管理', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140350, 0, '/edu', 'Edu', '', '', '教培管理', 0, 0, 0, 1, 0, '', 0, 'classify', '', 3, 1, 0, '', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140351, 140350, '/edu/student', 'EduStudent', '', 'edu/student/student', '学生管理', 0, 0, 0, 1, 0, '', 0, '', 'IconUser', 1, 2, 0, 'edu:student:list', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140352, 140350, '/edu/course', 'EduCourse', '', 'edu/course/course', '课程/项目管理', 0, 0, 0, 1, 0, '', 0, '', 'IconBook', 2, 2, 0, 'edu:course:list', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140353, 140350, '/edu/class', 'EduClass', '', 'edu/class/class', '班级管理', 0, 0, 0, 1, 0, '', 0, '', 'IconUserGroup', 3, 2, 0, 'edu:class:list', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140354, 140350, '/edu/room', 'EduRoom', '', 'edu/room/room', '教室/场地管理', 0, 0, 0, 1, 0, '', 0, '', 'IconHome', 4, 2, 0, 'edu:room:list', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140360, 140351, '', '', '', '', '新增学生', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:student:add', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140361, 140351, '', '', '', '', '编辑学生', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:student:edit', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140362, 140351, '', '', '', '', '删除学生', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:student:delete', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140363, 140351, '', '', '', '', '导入学生', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:student:import', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140364, 140351, '', '', '', '', '导出学生', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:student:export', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140365, 140352, '', '', '', '', '新增课程/项目', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:course:add', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140366, 140352, '', '', '', '', '编辑课程/项目', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:course:edit', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140367, 140352, '', '', '', '', '删除课程/项目', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:course:delete', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140368, 140352, '', '', '', '', '导入课程/项目', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:course:import', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140369, 140352, '', '', '', '', '导出课程/项目', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:course:export', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140370, 140353, '', '', '', '', '新增班级', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:class:add', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140371, 140353, '', '', '', '', '编辑班级', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:class:edit', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140372, 140353, '', '', '', '', '删除班级', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:class:delete', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140373, 140353, '', '', '', '', '导入班级', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:class:import', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140374, 140353, '', '', '', '', '导出班级', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:class:export', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140375, 140353, '', '', '', '', '班级成员', 0, 0, 0, 1, 0, '', 0, '', '', 6, 3, 0, 'edu:class:member', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140380, 140354, '', '', '', '', '新增教室/场地', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:room:add', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140381, 140354, '', '', '', '', '编辑教室/场地', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:room:edit', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140382, 140354, '', '', '', '', '删除教室/场地', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:room:delete', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140383, 140354, '', '', '', '', '导入教室/场地', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:room:import', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140384, 140354, '', '', '', '', '导出教室/场地', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:room:export', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(140385, 140354, '', '', '', '', '场地开放时间', 0, 0, 0, 1, 0, '', 0, '', '', 6, 3, 0, 'edu:room:time', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140350),(140351),(140352),(140353),(140354),(140360),(140361),(140362),(140363),(140364),
    (140365),(140366),(140367),(140368),(140369),(140370),(140371),(140372),(140373),(140374),
    (140375),(140380),(140381),(140382),(140383),(140384),(140385)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140351, 217), (140351, 218), (140360, 219), (140361, 220), (140362, 221), (140363, 222), (140364, 223),
(140352, 224), (140352, 225), (140352, 226), (140365, 227), (140366, 228), (140367, 229), (140368, 230), (140369, 231),
(140353, 232), (140353, 233), (140353, 234), (140353, 238), (140370, 235), (140371, 236), (140372, 237), (140373, 242), (140374, 243),
(140375, 239), (140375, 240), (140375, 241), (140375, 244), (140375, 245),
(140354, 246), (140354, 247), (140354, 248), (140354, 252), (140354, 254), (140380, 249), (140381, 250), (140382, 251), (140383, 258), (140384, 259),
(140385, 253), (140385, 255), (140385, 256), (140385, 257), (140385, 260), (140385, 261)
ON CONFLICT DO NOTHING;

INSERT INTO sys_dict (id, name, code, status, description, created_at, updated_at, deleted_at, created_by) VALUES
(101, '学生状态', 'edu_student_status', 1, '教培学生在读状态', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(102, '教培性别', 'edu_gender', 1, '教培学生性别', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(103, '年级', 'edu_grade', 1, '学生年级', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(104, '来源渠道', 'edu_source_channel', 1, '学生来源渠道', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(105, '课程/项目类型', 'edu_course_type', 1, '课程或托管项目类型', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(106, '班级类型', 'edu_class_type', 1, '托管班或班课，不包含一对一预约', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(107, '班级成员状态', 'edu_class_member_status', 1, '班级学员状态', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(108, '场地类型', 'edu_room_type', 1, '教室或场地类型', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1),
(109, '场地例外类型', 'edu_room_exception_type', 1, '场地开放或停用例外', '2026-05-10 10:00:00', '2026-05-10 10:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_dict_item (id, name, value, status, dict_id) VALUES
(10101, '在读', '1', 1, 101), (10102, '暂停', '0', 1, 101), (10103, '退学', '2', 1, 101),
(10201, '未知', 'unknown', 1, 102), (10202, '男', 'male', 1, 102), (10203, '女', 'female', 1, 102),
(10301, '幼儿园小班', 'k1', 1, 103), (10302, '幼儿园中班', 'k2', 1, 103), (10303, '幼儿园大班', 'k3', 1, 103),
(10304, '一年级', 'g1', 1, 103), (10305, '二年级', 'g2', 1, 103), (10306, '三年级', 'g3', 1, 103),
(10307, '四年级', 'g4', 1, 103), (10308, '五年级', 'g5', 1, 103), (10309, '六年级', 'g6', 1, 103),
(10401, '自然到访', 'walk_in', 1, 104), (10402, '转介绍', 'referral', 1, 104), (10403, '线上咨询', 'online', 1, 104),
(10404, '活动报名', 'campaign', 1, 104), (10405, '其他', 'other', 1, 104),
(10501, '托管', 'daycare', 1, 105), (10502, '艺术', 'art', 1, 105), (10503, '体育', 'sport', 1, 105), (10504, '其他', 'other', 1, 105),
(10601, '托管班', 'daycare', 1, 106), (10602, '班课', 'group', 1, 106),
(10701, '在读', 'studying', 1, 107), (10702, '暂停', 'paused', 1, 107), (10703, '退班', 'left', 1, 107),
(10801, '教室', 'classroom', 1, 108), (10802, '托管室', 'daycare_room', 1, 108), (10803, '球场', 'court', 1, 108),
(10804, '舞蹈室', 'dance_room', 1, 108), (10805, '其他', 'other', 1, 108),
(10901, '特殊开放', 'open', 1, 109), (10902, '停用', 'closed', 1, 109)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/students/list', 'GET'), ('/api/edu/students/:id', 'GET'), ('/api/edu/students/add', 'POST'), ('/api/edu/students/edit', 'PUT'),
    ('/api/edu/students/delete', 'DELETE'), ('/api/edu/students/import', 'POST'), ('/api/edu/students/export', 'GET'),
    ('/api/edu/courses/list', 'GET'), ('/api/edu/courses/options', 'GET'), ('/api/edu/courses/:id', 'GET'), ('/api/edu/courses/add', 'POST'),
    ('/api/edu/courses/edit', 'PUT'), ('/api/edu/courses/delete', 'DELETE'), ('/api/edu/courses/import', 'POST'), ('/api/edu/courses/export', 'GET'),
    ('/api/edu/classes/list', 'GET'), ('/api/edu/classes/teacher-options', 'GET'), ('/api/edu/classes/:id', 'GET'), ('/api/edu/classes/add', 'POST'),
    ('/api/edu/classes/edit', 'PUT'), ('/api/edu/classes/delete', 'DELETE'), ('/api/edu/classes/:id/members', 'GET'), ('/api/edu/classes/:id/members/add', 'POST'),
    ('/api/edu/classes/:id/members/edit', 'PUT'), ('/api/edu/classes/:id/members/delete', 'DELETE'), ('/api/edu/classes/import', 'POST'), ('/api/edu/classes/export', 'GET'),
    ('/api/edu/classes/members/import', 'POST'), ('/api/edu/classes/members/export', 'GET'),
    ('/api/edu/rooms/list', 'GET'), ('/api/edu/rooms/options', 'GET'), ('/api/edu/rooms/:id', 'GET'), ('/api/edu/rooms/add', 'POST'),
    ('/api/edu/rooms/edit', 'PUT'), ('/api/edu/rooms/delete', 'DELETE'), ('/api/edu/rooms/:id/weekly-rules', 'GET'), ('/api/edu/rooms/:id/weekly-rules/save', 'POST'),
    ('/api/edu/rooms/:id/exceptions', 'GET'), ('/api/edu/rooms/:id/exceptions/add', 'POST'), ('/api/edu/rooms/:id/exceptions/edit', 'PUT'), ('/api/edu/rooms/:id/exceptions/delete', 'DELETE'),
    ('/api/edu/rooms/import', 'POST'), ('/api/edu/rooms/export', 'GET'), ('/api/edu/rooms/weekly-rules/import', 'POST'), ('/api/edu/rooms/exceptions/import', 'POST')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140350','140351','140352','140353','140354']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 261), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140385), true);
SELECT setval('sys_dict_id_seq', GREATEST((SELECT last_value FROM sys_dict_id_seq), 109), true);
SELECT setval('sys_dict_item_id_seq', GREATEST((SELECT last_value FROM sys_dict_item_id_seq), 10902), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
