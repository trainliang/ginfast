package service

import (
	"testing"

	"gin-fast/app/models"
)

func TestEduStudentServiceCompile(t *testing.T) {
	setupEduTestDB(t)
	_ = NewEduStudentService()
}

func TestEduStudentServiceUpdateReplacesContacts(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStudentService()
	ctx := contextWithTenant(1)

	student := &models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}
	student.Contacts = models.EduStudentContactList{
		&models.EduStudentContact{Relation: "母亲", Name: "老联系人", Phone: "13800000001", TenantID: 1},
	}
	if err := svc.Create(ctx, student); err != nil {
		t.Fatalf("create student: %v", err)
	}

	student.Name = "学生A-更新"
	student.Contacts = models.EduStudentContactList{
		&models.EduStudentContact{Relation: "父亲", Name: "新联系人", Phone: "13800000002", TenantID: 1},
	}
	if err := svc.Update(ctx, student); err != nil {
		t.Fatalf("update student: %v", err)
	}

	var contacts []models.EduStudentContact
	if err := db.Where("student_id = ?", student.ID).Order("id asc").Find(&contacts).Error; err != nil {
		t.Fatalf("load contacts: %v", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(contacts))
	}
	if contacts[0].Name != "新联系人" {
		t.Fatalf("expected replaced contact, got %+v", contacts[0])
	}
}

func TestEduStudentServiceRejectsZeroTenantCreate(t *testing.T) {
	setupEduTestDB(t)
	svc := NewEduStudentService()
	err := svc.Create(contextWithTenant(0), &models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 0})
	if err == nil {
		t.Fatalf("expected zero tenant create to be rejected")
	}
}

func TestEduStudentServiceRejectsCrossTenantUpdate(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStudentService()
	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 0}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}
	err := svc.Update(contextWithTenant(1), &models.EduStudent{BaseModel: models.BaseModel{ID: 1}, Name: "学生A", Phone: "13800000000", TenantID: 1})
	if err == nil {
		t.Fatalf("expected cross tenant update to be rejected")
	}
}

func TestEduStudentServiceDeleteRejectsActiveClassMember(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStudentService()
	ctx := contextWithTenant(1)

	if err := db.Create(&models.EduStudent{Name: "学生A", Phone: "13800000000", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed student: %v", err)
	}
	if err := db.Create(&models.EduClassMember{ClassID: 1, StudentID: 1, Status: "studying", TenantID: 1}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	if err := svc.Delete(ctx, 1, 1); err == nil {
		t.Fatalf("expected delete to be rejected")
	}

	var count int64
	if err := db.Model(&models.EduStudent{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count student: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected student to remain, got %d", count)
	}
}

func TestEduStudentServiceImportRollsBackOnError(t *testing.T) {
	db := setupEduTestDB(t)
	svc := NewEduStudentService()
	ctx := contextWithTenant(1)

	result, err := svc.ImportRows(ctx, 1, []models.EduStudentImportRow{
		{Name: "新增", Phone: "13800000002"},
		{Name: "错误行", Phone: ""},
	})
	if err == nil {
		t.Fatalf("expected import error")
	}
	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("expected row errors")
	}

	var count int64
	if err := db.Model(&models.EduStudent{}).Where("tenant_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count students: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback without partial insert, got %d", count)
	}
}
