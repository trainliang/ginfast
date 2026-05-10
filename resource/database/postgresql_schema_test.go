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
