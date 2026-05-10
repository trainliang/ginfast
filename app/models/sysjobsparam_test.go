package models

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSysJobsListRequestHandleUsesGroupNameColumn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}

	group := "default"
	req := SysJobsListRequest{Group: &group}

	stmt := db.Model(&SysJobs{}).Scopes(req.Handle()).Find(&SysJobsList{}).Statement
	sql := stmt.SQL.String()

	if !strings.Contains(sql, "group_name = ?") {
		t.Fatalf("expected SQL to filter by group_name, got %q", sql)
	}
	if strings.Contains(sql, "group = ?") {
		t.Fatalf("SQL still filters by reserved column group: %q", sql)
	}
}

func TestSysJobsModelMapsGroupToGroupNameColumn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}

	job := SysJobs{
		Id:    "job-1",
		Group: "default",
	}

	stmt := db.Create(&job).Statement
	sql := stmt.SQL.String()

	if !strings.Contains(sql, "`group_name`") {
		t.Fatalf("expected SQL to write group_name column, got %q", sql)
	}
	if strings.Contains(sql, "`group`") {
		t.Fatalf("SQL still writes reserved column group: %q", sql)
	}
}
