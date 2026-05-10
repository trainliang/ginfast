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
