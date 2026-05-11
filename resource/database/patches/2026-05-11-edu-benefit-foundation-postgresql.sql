-- PostgreSQL incremental patch for education terms and benefit foundation.
-- Safe to execute more than once.

BEGIN;

ALTER TABLE edu_course ADD COLUMN IF NOT EXISTS default_teaching_mode VARCHAR(32) DEFAULT 'offline';
ALTER TABLE edu_course ADD COLUMN IF NOT EXISTS requires_room SMALLINT DEFAULT 1;
ALTER TABLE edu_class ADD COLUMN IF NOT EXISTS benefit_check_policy VARCHAR(32) DEFAULT 'required';

CREATE TABLE IF NOT EXISTS edu_term (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0,
    CONSTRAINT uk_edu_term_tenant_name UNIQUE (tenant_id, name)
);
CREATE INDEX IF NOT EXISTS idx_edu_term_deleted_at ON edu_term (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_term_tenant_date ON edu_term (tenant_id, start_date, end_date);

CREATE TABLE IF NOT EXISTS edu_term_closed_day (
    id SERIAL PRIMARY KEY,
    term_id INTEGER NOT NULL,
    closed_date TIMESTAMP NOT NULL,
    reason VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_term_closed_day_deleted_at ON edu_term_closed_day (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_term_closed_day_term_date ON edu_term_closed_day (tenant_id, term_id, closed_date);

CREATE TABLE IF NOT EXISTS edu_benefit_product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(64) NOT NULL,
    benefit_type VARCHAR(32) NOT NULL,
    calculation_mode VARCHAR(32) NOT NULL,
    total_count INTEGER DEFAULT 0,
    valid_days INTEGER DEFAULT 0,
    status SMALLINT DEFAULT 1,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INTEGER DEFAULT 0,
    tenant_id INTEGER DEFAULT 0,
    CONSTRAINT uk_edu_benefit_product_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_product_deleted_at ON edu_benefit_product (deleted_at);

CREATE TABLE IF NOT EXISTS edu_benefit_product_course (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_product_course_deleted_at ON edu_benefit_product_course (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_product_course_product ON edu_benefit_product_course (tenant_id, product_id);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_product_course_course ON edu_benefit_product_course (tenant_id, course_id);

CREATE TABLE IF NOT EXISTS edu_student_benefit (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    benefit_type VARCHAR(32) NOT NULL,
    calculation_mode VARCHAR(32) NOT NULL,
    course_id INTEGER DEFAULT 0,
    class_id INTEGER DEFAULT 0,
    teacher_id INTEGER DEFAULT 0,
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    total_count INTEGER DEFAULT 0,
    used_count INTEGER DEFAULT 0,
    remaining_count INTEGER DEFAULT 0,
    status SMALLINT DEFAULT 1,
    source_type VARCHAR(32) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0,
    created_by INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_student_benefit_deleted_at ON edu_student_benefit (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_student_benefit_student_course ON edu_student_benefit (tenant_id, student_id, course_id);
CREATE INDEX IF NOT EXISTS idx_edu_student_benefit_product ON edu_student_benefit (tenant_id, product_id);

CREATE TABLE IF NOT EXISTS edu_benefit_ledger (
    id SERIAL PRIMARY KEY,
    student_benefit_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    action_type VARCHAR(32) NOT NULL,
    biz_type VARCHAR(32) NOT NULL,
    biz_id INTEGER DEFAULT 0,
    change_count INTEGER DEFAULT 0,
    before_count INTEGER DEFAULT 0,
    after_count INTEGER DEFAULT 0,
    occurred_at TIMESTAMP,
    operator_id INTEGER DEFAULT 0,
    remark VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_ledger_deleted_at ON edu_benefit_ledger (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_ledger_student ON edu_benefit_ledger (tenant_id, student_id);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_ledger_benefit ON edu_benefit_ledger (tenant_id, student_benefit_id);

CREATE TABLE IF NOT EXISTS edu_benefit_event (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    student_benefit_id INTEGER DEFAULT 0,
    event_type VARCHAR(32) NOT NULL,
    course_id INTEGER DEFAULT 0,
    class_id INTEGER DEFAULT 0,
    effective_from TIMESTAMP,
    effective_to TIMESTAMP,
    payload TEXT,
    processed_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_event_deleted_at ON edu_benefit_event (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_event_student_course ON edu_benefit_event (tenant_id, student_id, course_id);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_event_benefit ON edu_benefit_event (tenant_id, student_benefit_id);

CREATE TABLE IF NOT EXISTS edu_benefit_external_sync (
    id SERIAL PRIMARY KEY,
    student_benefit_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    provider_code VARCHAR(64) NOT NULL,
    idempotency_key VARCHAR(128) NOT NULL,
    payload TEXT,
    status VARCHAR(32) NOT NULL,
    external_order_no VARCHAR(128) DEFAULT '',
    last_error VARCHAR(500) DEFAULT '',
    retry_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_external_sync_deleted_at ON edu_benefit_external_sync (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_external_sync_provider_key ON edu_benefit_external_sync (tenant_id, provider_code, idempotency_key);
CREATE INDEX IF NOT EXISTS idx_edu_benefit_external_sync_student ON edu_benefit_external_sync (tenant_id, student_id);

CREATE TABLE IF NOT EXISTS edu_lesson_student_eligibility (
    id SERIAL PRIMARY KEY,
    lesson_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    class_id INTEGER DEFAULT 0,
    student_benefit_id INTEGER DEFAULT 0,
    eligibility_status VARCHAR(32) NOT NULL,
    reason_code VARCHAR(64) DEFAULT '',
    checked_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    tenant_id INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_eligibility_deleted_at ON edu_lesson_student_eligibility (deleted_at);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_eligibility_student_course ON edu_lesson_student_eligibility (tenant_id, student_id, course_id);
CREATE INDEX IF NOT EXISTS idx_edu_lesson_student_eligibility_lesson ON edu_lesson_student_eligibility (tenant_id, lesson_id);

INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES
(264, '学期列表', '/api/edu/terms/list', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(265, '新增学期', '/api/edu/terms/add', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(266, '编辑学期', '/api/edu/terms/edit', 'PUT', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(267, '删除学期', '/api/edu/terms/delete', 'DELETE', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(268, '学期停课日列表', '/api/edu/terms/:id/closed-days', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(269, '保存学期停课日', '/api/edu/terms/:id/closed-days/save', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(270, '权益产品列表', '/api/edu/benefit-products/list', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(271, '新增权益产品', '/api/edu/benefit-products/add', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(272, '编辑权益产品', '/api/edu/benefit-products/edit', 'PUT', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(273, '删除权益产品', '/api/edu/benefit-products/delete', 'DELETE', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(274, '学生权益列表', '/api/edu/student-benefits/list', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(275, '新增学生权益', '/api/edu/student-benefits/add', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(276, '编辑学生权益', '/api/edu/student-benefits/edit', 'PUT', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(277, '学生权益校验', '/api/edu/student-benefits/check', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(278, '续费修复课次资格', '/api/edu/student-benefits/repair-schedule', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(279, '权益流水列表', '/api/edu/benefit-ledgers/list', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(280, '外部同步列表', '/api/edu/benefit-external-sync/list', 'GET', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(281, '重试外部同步', '/api/edu/benefit-external-sync/retry', 'POST', '教培管理', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES
(140388, 140350, '/edu/term', 'EduTerm', '', 'edu/term/term', '学期管理', 0, 0, 0, 1, 0, '', 0, '', 'IconCalendar', 5, 2, 0, 'edu:term:list', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140389, 140350, '/edu/benefit/product', 'EduBenefitProduct', '', 'edu/benefit/product', '权益产品', 0, 0, 0, 1, 0, '', 0, '', 'IconGift', 6, 2, 0, 'edu:benefitProduct:list', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140390, 140350, '/edu/benefit/student-benefit', 'EduStudentBenefit', '', 'edu/benefit/student-benefit', '学生权益', 0, 0, 0, 1, 0, '', 0, '', 'IconSafe', 7, 2, 0, 'edu:studentBenefit:list', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140391, 140388, '', '', '', '', '新增学期', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:term:add', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140392, 140388, '', '', '', '', '编辑学期', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:term:edit', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140393, 140388, '', '', '', '', '删除学期', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:term:delete', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140394, 140388, '', '', '', '', '学期停课日', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:term:time', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140395, 140389, '', '', '', '', '新增权益产品', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:benefitProduct:add', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140396, 140389, '', '', '', '', '编辑权益产品', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:benefitProduct:edit', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140397, 140389, '', '', '', '', '删除权益产品', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:benefitProduct:delete', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140398, 140390, '', '', '', '', '新增学生权益', 0, 0, 0, 1, 0, '', 0, '', '', 1, 3, 0, 'edu:studentBenefit:add', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140399, 140390, '', '', '', '', '编辑学生权益', 0, 0, 0, 1, 0, '', 0, '', '', 2, 3, 0, 'edu:studentBenefit:edit', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140400, 140390, '', '', '', '', '学生权益校验', 0, 0, 0, 1, 0, '', 0, '', '', 3, 3, 0, 'edu:studentBenefit:check', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140401, 140390, '', '', '', '', '续费修复', 0, 0, 0, 1, 0, '', 0, '', '', 4, 3, 0, 'edu:studentBenefit:repair', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140402, 140390, '', '', '', '', '权益流水', 0, 0, 0, 1, 0, '', 0, '', '', 5, 3, 0, 'edu:benefitLedger:list', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140403, 140390, '', '', '', '', '外部同步', 0, 0, 0, 1, 0, '', 0, '', '', 6, 3, 0, 'edu:benefitExternalSync:list', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1),
(140404, 140390, '', '', '', '', '重试外部同步', 0, 0, 0, 1, 0, '', 0, '', '', 7, 3, 0, 'edu:benefitExternalSync:retry', '2026-05-11 12:00:00', '2026-05-11 12:00:00', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id)
SELECT role_id, menu_id
FROM (VALUES (1),(2)) AS roles(role_id)
CROSS JOIN (VALUES
    (140388),(140389),(140390),(140391),(140392),(140393),(140394),(140395),(140396),(140397),
    (140398),(140399),(140400),(140401),(140402),(140403),(140404)
) AS menus(menu_id)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menu_api (menu_id, api_id) VALUES
(140388, 264), (140391, 265), (140392, 266), (140393, 267), (140394, 268), (140394, 269),
(140389, 270), (140395, 271), (140396, 272), (140397, 273),
(140390, 274), (140398, 275), (140399, 276), (140400, 277), (140401, 278), (140402, 279), (140403, 280), (140404, 281)
ON CONFLICT DO NOTHING;

INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', role_name, path, method, '*', '', ''
FROM (VALUES ('role_1'), ('role_2')) AS roles(role_name)
CROSS JOIN (VALUES
    ('/api/edu/terms/list', 'GET'),
    ('/api/edu/terms/add', 'POST'),
    ('/api/edu/terms/edit', 'PUT'),
    ('/api/edu/terms/delete', 'DELETE'),
    ('/api/edu/terms/:id/closed-days', 'GET'),
    ('/api/edu/terms/:id/closed-days/save', 'POST'),
    ('/api/edu/benefit-products/list', 'GET'),
    ('/api/edu/benefit-products/add', 'POST'),
    ('/api/edu/benefit-products/edit', 'PUT'),
    ('/api/edu/benefit-products/delete', 'DELETE'),
    ('/api/edu/student-benefits/list', 'GET'),
    ('/api/edu/student-benefits/add', 'POST'),
    ('/api/edu/student-benefits/edit', 'PUT'),
    ('/api/edu/student-benefits/check', 'POST'),
    ('/api/edu/student-benefits/repair-schedule', 'POST'),
    ('/api/edu/benefit-ledgers/list', 'GET'),
    ('/api/edu/benefit-external-sync/list', 'GET'),
    ('/api/edu/benefit-external-sync/retry', 'POST')
) AS apis(path, method)
ON CONFLICT DO NOTHING;

UPDATE sys_tenants
SET menu_permission = array_to_string(
    ARRAY(
        SELECT DISTINCT id
        FROM unnest(
            string_to_array(COALESCE(NULLIF(menu_permission, ''), ''), ',') ||
            ARRAY['140388','140389','140390','140391','140392','140393','140394','140395','140396','140397','140398','140399','140400','140401','140402','140403','140404']
        ) AS ids(id)
        WHERE id <> ''
    ),
    ','
)
WHERE id = 1;

SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 281), true);
SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140404), true);
SELECT setval('sys_casbin_rule_id_seq', GREATEST((SELECT last_value FROM sys_casbin_rule_id_seq), COALESCE((SELECT MAX(id) FROM sys_casbin_rule), 1)), true);

COMMIT;
