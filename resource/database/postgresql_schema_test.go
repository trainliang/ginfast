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
