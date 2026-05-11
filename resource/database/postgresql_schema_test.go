package database_test

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestPostgreSQLConvertedTinyintColumnsUseNumericTypes(t *testing.T) {
	sqlBytes, err := os.ReadFile("postgresql_converted.sql")
	if err != nil {
		t.Fatalf("read postgresql schema: %v", err)
	}

	expectedTypes := map[string]map[string]string{
		"demo_teacher": {
			"gender": "SMALLINT",
			"status": "SMALLINT",
		},
		"sys_department": {
			"status": "SMALLINT",
		},
		"sys_dict": {
			"status": "SMALLINT",
		},
		"sys_dict_item": {
			"status": "SMALLINT",
		},
		"sys_users": {
			"status": "SMALLINT",
		},
		"sys_user_tenant": {
			"is_default": "BOOLEAN",
		},
	}

	for table, columns := range expectedTypes {
		tableRe := regexp.MustCompile(`(?is)CREATE TABLE\s+` + regexp.QuoteMeta(table) + `\s*\((.*?)\);`)
		tableMatch := tableRe.FindSubmatch(sqlBytes)
		if len(tableMatch) != 2 {
			t.Fatalf("%s table definition not found", table)
		}

		for column, want := range columns {
			columnRe := regexp.MustCompile(`(?im)^\s*` + regexp.QuoteMeta(column) + `\s+([a-z0-9_]+)`)
			columnMatch := columnRe.FindSubmatch(tableMatch[1])
			if len(columnMatch) != 2 {
				t.Fatalf("%s.%s column definition not found", table, column)
			}

			if got := strings.ToUpper(string(columnMatch[1])); got != want {
				t.Fatalf("%s.%s type = %s, want %s", table, column, got, want)
			}
		}
	}
}

func TestPostgreSQLSchemaIncludesEduBenefitFoundationPatch(t *testing.T) {
	sqlBytes, err := os.ReadFile("patches/2026-05-11-edu-benefit-foundation-postgresql.sql")
	if err != nil {
		t.Fatalf("read edu benefit foundation patch: %v", err)
	}
	sqlText := string(sqlBytes)

	requiredSnippets := []string{
		"ALTER TABLE edu_course ADD COLUMN IF NOT EXISTS default_teaching_mode VARCHAR(32) DEFAULT 'offline'",
		"ALTER TABLE edu_course ADD COLUMN IF NOT EXISTS requires_room SMALLINT DEFAULT 1",
		"ALTER TABLE edu_class ADD COLUMN IF NOT EXISTS benefit_check_policy VARCHAR(32) DEFAULT 'required'",
		"CREATE TABLE IF NOT EXISTS edu_term",
		"CREATE TABLE IF NOT EXISTS edu_term_closed_day",
		"CREATE TABLE IF NOT EXISTS edu_benefit_product",
		"CREATE TABLE IF NOT EXISTS edu_benefit_product_course",
		"CREATE TABLE IF NOT EXISTS edu_student_benefit",
		"CREATE TABLE IF NOT EXISTS edu_benefit_ledger",
		"CREATE TABLE IF NOT EXISTS edu_benefit_event",
		"CREATE TABLE IF NOT EXISTS edu_benefit_external_sync",
		"CREATE TABLE IF NOT EXISTS edu_lesson_student_eligibility",
		"CONSTRAINT uk_edu_term_tenant_name UNIQUE (tenant_id, name)",
		"CONSTRAINT uk_edu_benefit_product_tenant_code UNIQUE (tenant_id, code)",
		"CREATE INDEX IF NOT EXISTS idx_edu_student_benefit_student_course ON edu_student_benefit (tenant_id, student_id, course_id)",
		"CREATE INDEX IF NOT EXISTS idx_edu_benefit_external_sync_provider_key ON edu_benefit_external_sync (tenant_id, provider_code, idempotency_key)",
		"INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES",
		"'/api/edu/terms/list', 'GET'",
		"'/api/edu/benefit-products/list', 'GET'",
		"'/api/edu/student-benefits/repair-schedule', 'POST'",
		"'/api/edu/benefit-external-sync/retry', 'POST'",
		"INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES",
		"'edu/term/term', '学期管理'",
		"'edu/benefit/product', '权益产品'",
		"'edu/benefit/student-benefit', '学生权益'",
		"INSERT INTO sys_menu_api (menu_id, api_id) VALUES",
		"INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)",
		"UPDATE sys_tenants",
		"SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 281), true)",
		"SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140404), true)",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(sqlText, snippet) {
			t.Fatalf("edu benefit foundation patch missing required snippet: %s", snippet)
		}
	}
}

func TestPostgreSQLSchemaIncludesEduScheduleRulesPatch(t *testing.T) {
	sqlBytes, err := os.ReadFile("patches/2026-05-11-edu-schedule-rules-postgresql.sql")
	if err != nil {
		t.Fatalf("read edu schedule rules patch: %v", err)
	}
	sqlText := string(sqlBytes)

	requiredSnippets := []string{
		"CREATE TABLE IF NOT EXISTS edu_schedule_rule",
		"CREATE TABLE IF NOT EXISTS edu_lesson",
		"CREATE TABLE IF NOT EXISTS edu_lesson_change_log",
		"CREATE TABLE IF NOT EXISTS edu_schedule_conflict_override",
		"repeat_type VARCHAR(32)",
		"term_id INTEGER DEFAULT 0",
		"course_id INTEGER NOT NULL",
		"rule_id INTEGER DEFAULT 0",
		"rule_version INTEGER DEFAULT 1",
		"lesson_type VARCHAR(32) NOT NULL",
		"status VARCHAR(32) DEFAULT 'scheduled'",
		"is_manual_adjusted SMALLINT DEFAULT 0",
		"source_lesson_id INTEGER DEFAULT 0",
		"action_type VARCHAR(32) NOT NULL",
		"operator_id INTEGER DEFAULT 0",
		"occurred_at TIMESTAMP",
		"conflict_key VARCHAR(128) NOT NULL",
		"(tenant_id, rule_type)",
		"(tenant_id, lesson_date)",
		"(tenant_id, teacher_id, lesson_date)",
		"(tenant_id, room_id, lesson_date)",
		"(tenant_id, student_id, lesson_date)",
		"(tenant_id, class_id, lesson_date)",
		"INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES",
		"'/api/edu/schedule-rules/list', 'GET'",
		"'/api/edu/schedule-rules/add', 'POST'",
		"'/api/edu/schedule-rules/edit', 'PUT'",
		"'/api/edu/schedule-rules/preview-change', 'POST'",
		"'/api/edu/schedule-rules/delete', 'DELETE'",
		"'/api/edu/lessons/calendar', 'GET'",
		"'/api/edu/lessons/list', 'GET'",
		"'/api/edu/lessons/reschedule', 'PUT'",
		"'/api/edu/lessons/stop', 'PUT'",
		"'/api/edu/lessons/cancel', 'PUT'",
		"'/api/edu/lessons/restore', 'PUT'",
		"'/api/edu/lessons/makeup', 'POST'",
		"'/api/edu/lessons/:id/change-logs', 'GET'",
		"'/api/edu/schedules/check-conflicts', 'POST'",
		"INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES",
		"'/edu/schedule', 'EduSchedule', '', 'edu/schedule/schedule', '排课管理'",
		"INSERT INTO sys_menu_api (menu_id, api_id) VALUES",
		"INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)",
		"UPDATE sys_tenants",
		"SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 295), true)",
		"SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140418), true)",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(sqlText, snippet) {
			t.Fatalf("edu schedule rules patch missing required snippet: %s", snippet)
		}
	}
}

func TestPostgreSQLSchemaIncludesEduStatisticsPermissionsPatch(t *testing.T) {
	sqlBytes, err := os.ReadFile("patches/2026-05-11-edu-statistics-permissions-postgresql.sql")
	if err != nil {
		t.Fatalf("read edu statistics permissions patch: %v", err)
	}
	sqlText := string(sqlBytes)

	requiredSnippets := []string{
		"INSERT INTO sys_api (id, title, path, method, api_group, created_at, updated_at, deleted_at, created_by) VALUES",
		"'/api/edu/statistics/schedule', 'GET'",
		"'/api/edu/statistics/benefits', 'GET'",
		"'/api/edu/statistics/external-sync', 'GET'",
		"INSERT INTO sys_menu (id, parent_id, path, name, redirect, component, title, is_full, hide, disable, keep_alive, affix, link, iframe, svg_icon, icon, sort, type, is_link, permission, created_at, updated_at, deleted_at, created_by) VALUES",
		"'/edu/statistics', 'EduStatistics', '', 'edu/statistics/statistics', '运营统计'",
		"'edu:statistics:view'",
		"INSERT INTO sys_menu_api (menu_id, api_id) VALUES",
		"INSERT INTO sys_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)",
		"/api/edu/statistics/schedule",
		"/api/edu/statistics/benefits",
		"/api/edu/statistics/external-sync",
		"UPDATE sys_tenants",
		"(296, '排课统计', '/api/edu/statistics/schedule', 'GET'",
		"(297, '权益统计', '/api/edu/statistics/benefits', 'GET'",
		"(298, '外部同步统计', '/api/edu/statistics/external-sync', 'GET'",
		"(140419, 140350, '/edu/statistics'",
		"SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 298), true)",
		"SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140419), true)",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(sqlText, snippet) {
			t.Fatalf("edu statistics permissions patch missing required snippet: %s", snippet)
		}
	}
}

func TestPostgreSQLSchemaIncludesEduScheduleCompletionPatch(t *testing.T) {
	sqlBytes, err := os.ReadFile("patches/2026-05-12-edu-schedule-completion-postgresql.sql")
	if err != nil {
		t.Fatalf("read edu schedule completion patch: %v", err)
	}
	sqlText := string(sqlBytes)

	requiredSnippets := []string{
		"ALTER TABLE edu_schedule_rule ADD COLUMN IF NOT EXISTS repeat_type VARCHAR(32)",
		"ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS rule_id INTEGER DEFAULT 0",
		"ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS lesson_type VARCHAR(32)",
		"ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS status VARCHAR(32) DEFAULT 'scheduled'",
		"ALTER TABLE edu_lesson ADD COLUMN IF NOT EXISTS is_manual_adjusted SMALLINT DEFAULT 0",
		"ALTER TABLE edu_lesson_change_log ADD COLUMN IF NOT EXISTS action_type VARCHAR(32)",
		"ALTER TABLE edu_schedule_conflict_override ADD COLUMN IF NOT EXISTS conflict_key VARCHAR(128)",
		"ALTER TABLE edu_lesson ALTER COLUMN start_time TYPE VARCHAR(16) USING start_time::text",
		"ALTER TABLE edu_lesson_change_log ALTER COLUMN before_data TYPE TEXT USING before_data::text",
		"ALTER TABLE edu_lesson_change_log ALTER COLUMN change_type DROP NOT NULL",
		"ALTER TABLE edu_schedule_conflict_override ALTER COLUMN target_type DROP NOT NULL",
		"ALTER TABLE edu_schedule_conflict_override ALTER COLUMN target_id SET DEFAULT 0",
		"ON CONFLICT (id) DO UPDATE SET",
		"DELETE FROM sys_menu_api WHERE menu_id = 140413 AND api_id IN (291, 292)",
		"'/api/edu/lessons/reschedule', 'PUT'",
		"'/api/edu/lessons/:id/change-logs', 'GET'",
		"(140419, 140350, '/edu/statistics'",
		"SELECT setval('sys_api_id_seq', GREATEST((SELECT last_value FROM sys_api_id_seq), 298), true)",
		"SELECT setval('sys_menu_id_seq', GREATEST((SELECT last_value FROM sys_menu_id_seq), 140419), true)",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(sqlText, snippet) {
			t.Fatalf("edu schedule completion patch missing required snippet: %s", snippet)
		}
	}
}

func TestPostgreSQLConvertedIncludesEducationFoundationSeed(t *testing.T) {
	sqlBytes, err := os.ReadFile("postgresql_converted.sql")
	if err != nil {
		t.Fatalf("read postgresql schema: %v", err)
	}
	sqlText := string(sqlBytes)

	requiredSnippets := []string{
		"CREATE TABLE edu_course",
		"CREATE TABLE edu_student",
		"CREATE TABLE edu_student_contact",
		"CREATE TABLE edu_class",
		"CREATE TABLE edu_class_member",
		"CREATE TABLE edu_room",
		"CREATE TABLE edu_room_weekly_rule",
		"CREATE TABLE edu_room_exception",
		"INSERT INTO sys_api VALUES (217, '学生列表', '/api/edu/students/list'",
		"INSERT INTO sys_dict VALUES (101, '学生状态', 'edu_student_status'",
		"INSERT INTO sys_dict_item VALUES (10101, '在读', '1'",
		"INSERT INTO sys_menu VALUES (140350, 0, '/edu', 'Edu'",
		"INSERT INTO sys_role_menu VALUES (1, 140351)",
		"INSERT INTO sys_role_menu VALUES (2, 140351)",
		"INSERT INTO sys_menu_api VALUES (140351, 217)",
		"140350,140351,140352,140353,140354",
		"SELECT setval('sys_api_id_seq', 261, true)",
		"SELECT setval('sys_casbin_rule_id_seq', 7650, true)",
		"SELECT setval('sys_menu_id_seq', 140385, true)",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(sqlText, snippet) {
			t.Fatalf("postgresql_converted.sql missing required education seed snippet: %s", snippet)
		}
	}
}
