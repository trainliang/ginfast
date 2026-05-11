package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduStudentController 教培学生控制器
type EduStudentController struct {
	Common
	EduStudentService *service.EduStudentService
}

func NewEduStudentController() *EduStudentController {
	return &EduStudentController{Common: Common{}, EduStudentService: service.NewEduStudentService()}
}

func (ctl *EduStudentController) List(c *gin.Context) {
	var req models.EduStudentListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduStudentList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计学生数量失败", err)
	}
	list := models.NewEduStudentList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Contacts")
	}); err != nil {
		ctl.FailAndAbort(c, "获取学生列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduStudentController) GetByID(c *gin.Context) {
	var req models.EduStudentGetRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	student := models.NewEduStudent()
	if err := student.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? AND tenant_id = ?", req.ID, ctl.GetCurrentTenantID(c)).Preload("Contacts")
	}); err != nil {
		ctl.FailAndAbort(c, "查询学生失败", err)
	}
	if student.IsEmpty() {
		ctl.FailAndAbort(c, "学生不存在", nil)
	}
	ctl.Success(c, student)
}

func (ctl *EduStudentController) Add(c *gin.Context) {
	var req models.EduStudentAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	student := buildStudentFromAddRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduStudentService.Create(c, student); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生创建成功", student)
}

func (ctl *EduStudentController) Update(c *gin.Context) {
	var req models.EduStudentUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	student := buildStudentFromUpdateRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduStudentService.Update(c, student); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生更新成功", student)
}

func (ctl *EduStudentController) Delete(c *gin.Context) {
	var req models.EduStudentDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduStudentService.Delete(c, ctl.RequireTenant(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生删除成功", nil)
}

func (ctl *EduStudentController) Import(c *gin.Context) {
	var req models.EduStudentImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduStudentService.ImportRows(c, ctl.RequireTenant(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入学生失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduStudentController) Export(c *gin.Context) {
	var req models.EduStudentExportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduStudentService.ExportRows(c, ctl.RequireTenant(c), req.IDs)
	if err != nil {
		ctl.FailAndAbort(c, "导出学生失败", err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func buildStudentFromAddRequest(req *models.EduStudentAddRequest, tenantID, userID uint) *models.EduStudent {
	student := &models.EduStudent{
		Name:             req.Name,
		Gender:           req.Gender,
		Birthday:         req.Birthday,
		Phone:            req.Phone,
		Status:           1,
		Avatar:           req.Avatar,
		School:           req.School,
		Grade:            req.Grade,
		SchoolClass:      req.SchoolClass,
		SourceChannel:    req.SourceChannel,
		EnrollDate:       req.EnrollDate,
		HealthNote:       req.HealthNote,
		AllergyNote:      req.AllergyNote,
		EmergencyContact: req.EmergencyContact,
		PickupNote:       req.PickupNote,
		Remark:           req.Remark,
		CreatedBy:        userID,
		TenantID:         tenantID,
	}
	if req.Status != nil {
		student.Status = *req.Status
	}
	for _, item := range req.Contacts {
		student.Contacts = append(student.Contacts, buildStudentContact(item.Relation, item.Name, item.Phone, item.Remark, item.IsPrimary, item.CanPickup, tenantID, userID))
	}
	return student
}

func buildStudentFromUpdateRequest(req *models.EduStudentUpdateRequest, tenantID, userID uint) *models.EduStudent {
	student := &models.EduStudent{
		BaseModel:        models.BaseModel{ID: req.ID},
		Name:             req.Name,
		Gender:           req.Gender,
		Birthday:         req.Birthday,
		Phone:            req.Phone,
		Status:           1,
		Avatar:           req.Avatar,
		School:           req.School,
		Grade:            req.Grade,
		SchoolClass:      req.SchoolClass,
		SourceChannel:    req.SourceChannel,
		EnrollDate:       req.EnrollDate,
		HealthNote:       req.HealthNote,
		AllergyNote:      req.AllergyNote,
		EmergencyContact: req.EmergencyContact,
		PickupNote:       req.PickupNote,
		Remark:           req.Remark,
		CreatedBy:        userID,
		TenantID:         tenantID,
	}
	if req.Status != nil {
		student.Status = *req.Status
	}
	for _, item := range req.Contacts {
		student.Contacts = append(student.Contacts, buildStudentContact(item.Relation, item.Name, item.Phone, item.Remark, item.IsPrimary, item.CanPickup, tenantID, userID))
	}
	return student
}

func buildStudentContact(relation, name, phone, remark string, isPrimary, canPickup *int8, tenantID, userID uint) *models.EduStudentContact {
	contact := &models.EduStudentContact{
		Relation:  relation,
		Name:      name,
		Phone:     phone,
		Remark:    remark,
		CreatedBy: userID,
		TenantID:  tenantID,
	}
	if isPrimary != nil {
		contact.IsPrimary = *isPrimary
	}
	if canPickup != nil {
		contact.CanPickup = *canPickup
	}
	return contact
}
