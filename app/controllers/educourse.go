package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduCourseController 教培课程/项目控制器
type EduCourseController struct {
	Common
	EduCourseService *service.EduCourseService
}

// NewEduCourseController 创建教培课程/项目控制器
func NewEduCourseController() *EduCourseController {
	return &EduCourseController{
		Common:           Common{},
		EduCourseService: service.NewEduCourseService(),
	}
}

// List 课程/项目分页列表
func (ctl *EduCourseController) List(c *gin.Context) {
	var req models.EduCourseListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
	total, err := models.NewEduCourseList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计课程/项目数量失败", err)
	}
	list := models.NewEduCourseList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取课程/项目列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

// Options 启用课程/项目选项
func (ctl *EduCourseController) Options(c *gin.Context) {
	tenantID := ctl.GetCurrentTenantID(c)
	list := models.NewEduCourseList()
	if err := list.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ? AND status = ?", tenantID, 1).Order("sort asc, id asc")
	}); err != nil {
		ctl.FailAndAbort(c, "获取课程/项目选项失败", err)
	}
	ctl.Success(c, gin.H{"list": list})
}

// GetByID 根据ID获取课程/项目
func (ctl *EduCourseController) GetByID(c *gin.Context) {
	var req models.EduCourseGetRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	course := models.NewEduCourse()
	if err := course.Find(c, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? AND tenant_id = ?", req.ID, ctl.GetCurrentTenantID(c))
	}); err != nil {
		ctl.FailAndAbort(c, "查询课程/项目失败", err)
	}
	if course.IsEmpty() {
		ctl.FailAndAbort(c, "课程/项目不存在", nil)
	}
	ctl.Success(c, course)
}

// Add 新增课程/项目
func (ctl *EduCourseController) Add(c *gin.Context) {
	var req models.EduCourseAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.RequireTenant(c)
	course := &models.EduCourse{
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		GradeRange:   req.GradeRange,
		RequiresRoom: -1,
		Status:       1,
		Description:  req.Description,
		CreatedBy:    ctl.GetCurrentUserID(c),
		TenantID:     tenantID,
	}
	if req.DefaultTeachingMode != "" {
		course.DefaultTeachingMode = req.DefaultTeachingMode
	}
	if req.RequiresRoom != nil {
		course.RequiresRoom = int8(*req.RequiresRoom)
	}
	if req.Status != nil {
		course.Status = int8(*req.Status)
	}
	if req.Sort != nil {
		course.Sort = int(*req.Sort)
	}
	if err := ctl.EduCourseService.Create(c, course); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "课程/项目创建成功", course)
}

// Update 更新课程/项目
func (ctl *EduCourseController) Update(c *gin.Context) {
	var req models.EduCourseUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.RequireTenant(c)
	course := &models.EduCourse{
		BaseModel:    models.BaseModel{ID: uint(req.ID)},
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		GradeRange:   req.GradeRange,
		RequiresRoom: -1,
		Status:       1,
		Description:  req.Description,
		CreatedBy:    ctl.GetCurrentUserID(c),
		TenantID:     tenantID,
	}
	if req.DefaultTeachingMode != "" {
		course.DefaultTeachingMode = req.DefaultTeachingMode
	}
	if req.RequiresRoom != nil {
		course.RequiresRoom = int8(*req.RequiresRoom)
	}
	if req.Status != nil {
		course.Status = int8(*req.Status)
	}
	if req.Sort != nil {
		course.Sort = int(*req.Sort)
	}
	if err := ctl.EduCourseService.Update(c, course); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "课程/项目更新成功", course)
}

// Delete 删除课程/项目
func (ctl *EduCourseController) Delete(c *gin.Context) {
	var req models.EduCourseDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduCourseService.Delete(c, ctl.RequireTenant(c), uint(req.ID)); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "课程/项目删除成功", nil)
}

// Import 导入课程/项目
func (ctl *EduCourseController) Import(c *gin.Context) {
	var req models.EduCourseImportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	result, err := ctl.EduCourseService.ImportRows(c, ctl.RequireTenant(c), req.Rows)
	if err != nil && len(result.Errors) == 0 {
		ctl.FailAndAbort(c, "导入课程/项目失败", err)
	}
	ctl.Success(c, result)
}

// Export 导出课程/项目数据
func (ctl *EduCourseController) Export(c *gin.Context) {
	var req models.EduCourseExportRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ids := make([]uint, 0, len(req.IDs))
	for _, id := range req.IDs {
		ids = append(ids, uint(id))
	}
	rows, err := ctl.EduCourseService.ExportRows(c, ctl.RequireTenant(c), ids)
	if err != nil {
		ctl.FailAndAbort(c, "导出课程/项目失败", err)
	}
	ctl.Success(c, gin.H{"list": rows})
}
