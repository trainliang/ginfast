package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduClassController 教培班级控制器
type EduClassController struct {
	Common
	EduClassService *service.EduClassService
}

func NewEduClassController() *EduClassController {
	return &EduClassController{Common: Common{}, EduClassService: service.NewEduClassService()}
}

func (ctl *EduClassController) List(c *gin.Context) {
	var req models.EduClassListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduClassList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计班级数量失败", err)
	}
	list := models.NewEduClassList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取班级列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduClassController) GetByID(c *gin.Context) {
	var req models.EduClassGetRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	class := models.NewEduClass()
	if err := class.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? AND tenant_id = ?", req.ID, ctl.GetCurrentTenantID(c))
	}); err != nil {
		ctl.FailAndAbort(c, "查询班级失败", err)
	}
	if class.IsEmpty() {
		ctl.FailAndAbort(c, "班级不存在", nil)
	}
	ctl.Success(c, class)
}

func (ctl *EduClassController) Add(c *gin.Context) {
	var req models.EduClassAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	class := buildClassFromAddRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduClassService.Create(c, class); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级创建成功", class)
}

func (ctl *EduClassController) Update(c *gin.Context) {
	var req models.EduClassUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	class := buildClassFromUpdateRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduClassService.Update(c, class); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级更新成功", class)
}

func (ctl *EduClassController) Delete(c *gin.Context) {
	var req models.EduClassDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduClassService.Delete(c, ctl.GetCurrentTenantID(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级删除成功", nil)
}

func (ctl *EduClassController) Members(c *gin.Context) {
	var req models.EduClassMemberListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if classID, err := parseOptionalPathUint(c, "id"); err != nil {
		ctl.FailAndAbort(c, "班级ID格式错误", err)
	} else if classID > 0 {
		req.ClassID = &classID
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduClassMemberList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计班级成员数量失败", err)
	}
	list := models.NewEduClassMemberList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取班级成员失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduClassController) AddMember(c *gin.Context) {
	var req models.EduClassMemberAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	member := buildClassMemberFromAddRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduClassService.AddMember(c, member); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级成员添加成功", member)
}

func (ctl *EduClassController) UpdateMember(c *gin.Context) {
	var req models.EduClassMemberUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	member := buildClassMemberFromUpdateRequest(&req, ctl.GetCurrentTenantID(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduClassService.UpdateMember(c, member); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级成员更新成功", member)
}

func (ctl *EduClassController) DeleteMember(c *gin.Context) {
	var req models.EduClassMemberDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduClassService.DeleteMember(c, ctl.GetCurrentTenantID(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "班级成员删除成功", nil)
}

func (ctl *EduClassController) TeacherOptions(c *gin.Context) {
	list, err := ctl.EduClassService.TeacherOptions(c, ctl.GetCurrentTenantID(c))
	if err != nil {
		ctl.FailAndAbort(c, "获取教师选项失败", err)
	}
	ctl.Success(c, gin.H{"list": list})
}

func (ctl *EduClassController) Import(c *gin.Context) {
	var req models.EduClassImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduClassService.ImportRows(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入班级失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduClassController) ImportMembers(c *gin.Context) {
	var req models.EduClassMemberImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduClassService.ImportMemberRows(c, ctl.GetCurrentTenantID(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入班级成员失败", err)
	}
	ctl.Success(c, result)
}

func (ctl *EduClassController) Export(c *gin.Context) {
	var req models.EduClassExportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduClassService.ExportRows(c, ctl.GetCurrentTenantID(c), req.IDs)
	if err != nil {
		ctl.FailAndAbort(c, "导出班级失败", err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func (ctl *EduClassController) ExportMembers(c *gin.Context) {
	var req models.EduClassExportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	rows, err := ctl.EduClassService.ExportMemberRows(c, ctl.GetCurrentTenantID(c), req.IDs)
	if err != nil {
		ctl.FailAndAbort(c, "导出班级成员失败", err)
	}
	ctl.Success(c, gin.H{"list": rows})
}

func buildClassFromAddRequest(req *models.EduClassAddRequest, tenantID, userID uint) *models.EduClass {
	class := &models.EduClass{Name: req.Name, Code: req.Code, ClassType: req.ClassType, CourseID: req.CourseID, TeacherID: req.TeacherID, Capacity: req.Capacity, Status: 1, StartDate: req.StartDate, EndDate: req.EndDate, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
	if req.RoomID != nil {
		class.RoomID = *req.RoomID
	}
	if req.Status != nil {
		class.Status = *req.Status
	}
	return class
}

func buildClassFromUpdateRequest(req *models.EduClassUpdateRequest, tenantID, userID uint) *models.EduClass {
	class := &models.EduClass{BaseModel: models.BaseModel{ID: req.ID}, Name: req.Name, Code: req.Code, ClassType: req.ClassType, CourseID: req.CourseID, TeacherID: req.TeacherID, Capacity: req.Capacity, Status: 1, StartDate: req.StartDate, EndDate: req.EndDate, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
	if req.RoomID != nil {
		class.RoomID = *req.RoomID
	}
	if req.Status != nil {
		class.Status = *req.Status
	}
	return class
}

func buildClassMemberFromAddRequest(req *models.EduClassMemberAddRequest, tenantID, userID uint) *models.EduClassMember {
	return &models.EduClassMember{ClassID: req.ClassID, StudentID: req.StudentID, JoinDate: req.JoinDate, LeaveDate: req.LeaveDate, Status: req.Status, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
}

func buildClassMemberFromUpdateRequest(req *models.EduClassMemberUpdateRequest, tenantID, userID uint) *models.EduClassMember {
	return &models.EduClassMember{BaseModel: models.BaseModel{ID: req.ID}, ClassID: req.ClassID, StudentID: req.StudentID, JoinDate: req.JoinDate, LeaveDate: req.LeaveDate, Status: req.Status, Remark: req.Remark, CreatedBy: userID, TenantID: tenantID}
}

func parseOptionalPathUint(c *gin.Context, key string) (uint, error) {
	value := c.Param(key)
	if value == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
