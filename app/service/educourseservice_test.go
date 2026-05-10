package service

import (
	"testing"

	"gin-fast/app/models"
)

func TestEduCourseServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduCourseService()
}

func TestEduCourseServiceRejectsDuplicateCodeWithinTenant(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduCourseService()
	ctx := contextWithTenant(1)
	if err := db.Create(&models.EduCourse{Name: "A", Code: "C001", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed course: %v", err)
	}
	err := svc.Create(ctx, &models.EduCourse{Name: "B", Code: "C001", TenantID: 1})
	if err == nil {
		t.Fatalf("expected duplicate code error")
	}
}
