package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EduBenefitController 学生权益中心控制器
type EduBenefitController struct {
	Common
	EduBenefitService *service.EduBenefitService
}

func NewEduBenefitController() *EduBenefitController {
	return &EduBenefitController{Common: Common{}, EduBenefitService: service.NewEduBenefitService()}
}

func (ctl *EduBenefitController) ProductList(c *gin.Context) {
	var req models.EduBenefitProductListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduBenefitProductList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计权益产品数量失败", err)
	}
	list := models.NewEduBenefitProductList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取权益产品列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduBenefitController) AddProduct(c *gin.Context) {
	var req models.EduBenefitProductAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	product := buildBenefitProductFromAddRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduBenefitService.CreateProduct(c, product); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "权益产品创建成功", product)
}

func (ctl *EduBenefitController) UpdateProduct(c *gin.Context) {
	var req models.EduBenefitProductUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	product := buildBenefitProductFromUpdateRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduBenefitService.UpdateProduct(c, product); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "权益产品更新成功", product)
}

func (ctl *EduBenefitController) DeleteProduct(c *gin.Context) {
	var req models.EduBenefitProductDeleteRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduBenefitService.DeleteProduct(c, ctl.RequireTenant(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "权益产品删除成功", nil)
}

func (ctl *EduBenefitController) StudentBenefitList(c *gin.Context) {
	var req models.EduStudentBenefitListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduStudentBenefitList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计学生权益数量失败", err)
	}
	list := models.NewEduStudentBenefitList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取学生权益列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduBenefitController) AddStudentBenefit(c *gin.Context) {
	var req models.EduStudentBenefitAddRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	benefit := buildStudentBenefitFromAddRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduBenefitService.CreateStudentBenefit(c, benefit); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生权益创建成功", benefit)
}

func (ctl *EduBenefitController) UpdateStudentBenefit(c *gin.Context) {
	var req models.EduStudentBenefitUpdateRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	benefit := buildStudentBenefitFromUpdateRequest(&req, ctl.RequireTenant(c), ctl.GetCurrentUserID(c))
	if err := ctl.EduBenefitService.UpdateStudentBenefit(c, benefit); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "学生权益更新成功", benefit)
}

func (ctl *EduBenefitController) Check(c *gin.Context) {
	var req models.EduBenefitCheckRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.RequireTenant(c)
	req.TenantID = &tenantID
	result, err := ctl.EduBenefitService.CheckEligibility(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, result)
}

func (ctl *EduBenefitController) RepairSchedule(c *gin.Context) {
	var req models.EduBenefitRepairRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.RequireTenant(c)
	req.TenantID = &tenantID
	result, err := ctl.EduBenefitService.RepairScheduleEligibility(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, result)
}

func (ctl *EduBenefitController) LedgerList(c *gin.Context) {
	var req models.EduBenefitLedgerListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduBenefitLedgerList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计权益流水数量失败", err)
	}
	list := models.NewEduBenefitLedgerList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取权益流水列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduBenefitController) ExternalSyncList(c *gin.Context) {
	var req models.EduBenefitExternalSyncListRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	tenantID := ctl.GetCurrentTenantID(c)
	scope := func(db *gorm.DB) *gorm.DB { return db.Where("tenant_id = ?", tenantID) }
	total, err := models.NewEduBenefitExternalSyncList().GetTotal(c, req.Handler(), scope)
	if err != nil {
		ctl.FailAndAbort(c, "统计外部同步任务数量失败", err)
	}
	list := models.NewEduBenefitExternalSyncList()
	if err := list.Find(c, req.Paginate(), req.Handler(), scope); err != nil {
		ctl.FailAndAbort(c, "获取外部同步任务列表失败", err)
	}
	ctl.Success(c, gin.H{"list": list, "total": total})
}

func (ctl *EduBenefitController) RetryExternalSync(c *gin.Context) {
	var req models.EduBenefitExternalSyncRetryRequest
	if err := req.Validate(c); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	if err := ctl.EduBenefitService.RetryExternalSync(c, ctl.RequireTenant(c), req.ID); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.SuccessWithMessage(c, "外部同步任务已加入重试", nil)
}

func buildBenefitProductFromAddRequest(req *models.EduBenefitProductAddRequest, tenantID, userID uint) *models.EduBenefitProduct {
	product := &models.EduBenefitProduct{
		Name:            req.Name,
		Code:            req.Code,
		BenefitType:     req.BenefitType,
		CalculationMode: req.CalculationMode,
		TotalCount:      req.TotalCount,
		ValidDays:       req.ValidDays,
		Status:          1,
		Remark:          req.Remark,
		CreatedBy:       userID,
		TenantID:        tenantID,
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	return product
}

func buildBenefitProductFromUpdateRequest(req *models.EduBenefitProductUpdateRequest, tenantID, userID uint) *models.EduBenefitProduct {
	product := &models.EduBenefitProduct{
		BaseModel:       models.BaseModel{ID: req.ID},
		Name:            req.Name,
		Code:            req.Code,
		BenefitType:     req.BenefitType,
		CalculationMode: req.CalculationMode,
		TotalCount:      req.TotalCount,
		ValidDays:       req.ValidDays,
		Status:          1,
		Remark:          req.Remark,
		CreatedBy:       userID,
		TenantID:        tenantID,
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	return product
}

func buildStudentBenefitFromAddRequest(req *models.EduStudentBenefitAddRequest, tenantID, userID uint) *models.EduStudentBenefit {
	benefit := &models.EduStudentBenefit{
		StudentID:       req.StudentID,
		ProductID:       req.ProductID,
		BenefitType:     req.BenefitType,
		CalculationMode: req.CalculationMode,
		CourseID:        req.CourseID,
		ClassID:         req.ClassID,
		TeacherID:       req.TeacherID,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
		TotalCount:      req.TotalCount,
		UsedCount:       req.UsedCount,
		RemainingCount:  req.RemainingCount,
		Status:          1,
		SourceType:      req.SourceType,
		TenantID:        tenantID,
		CreatedBy:       userID,
	}
	if req.Status != nil {
		benefit.Status = *req.Status
	}
	return benefit
}

func buildStudentBenefitFromUpdateRequest(req *models.EduStudentBenefitUpdateRequest, tenantID, userID uint) *models.EduStudentBenefit {
	benefit := &models.EduStudentBenefit{
		BaseModel:       models.BaseModel{ID: req.ID},
		StudentID:       req.StudentID,
		ProductID:       req.ProductID,
		BenefitType:     req.BenefitType,
		CalculationMode: req.CalculationMode,
		CourseID:        req.CourseID,
		ClassID:         req.ClassID,
		TeacherID:       req.TeacherID,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
		TotalCount:      req.TotalCount,
		UsedCount:       req.UsedCount,
		RemainingCount:  req.RemainingCount,
		Status:          1,
		SourceType:      req.SourceType,
		TenantID:        tenantID,
		CreatedBy:       userID,
	}
	if req.Status != nil {
		benefit.Status = *req.Status
	}
	return benefit
}
