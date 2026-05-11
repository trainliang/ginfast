package service

import (
	"context"
	"strings"
	"time"

	"gin-fast/app/global/app"
	"gin-fast/app/models"
)

type EduStatisticsService struct{}

func NewEduStatisticsService() *EduStatisticsService { return &EduStatisticsService{} }

func (s *EduStatisticsService) ScheduleStatistics(ctx context.Context, req *models.EduStatisticsRangeRequest) (*models.EduScheduleStatisticsResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return &models.EduScheduleStatisticsResponse{}, nil
	}

	query := app.DB().WithContext(ctx).Model(&models.EduLesson{}).Where("tenant_id = ?", tenantID)
	if req != nil {
		if req.StartDate != nil && !req.StartDate.Time.IsZero() {
			query = query.Where("lesson_date >= ?", req.StartDate.Time)
		}
		if req.EndDate != nil && !req.EndDate.Time.IsZero() {
			query = query.Where("lesson_date <= ?", req.EndDate.Time)
		}
		if req.TeacherID != nil && *req.TeacherID != 0 {
			query = query.Where("teacher_id = ?", *req.TeacherID)
		}
		if req.RoomID != nil && *req.RoomID != 0 {
			query = query.Where("room_id = ?", *req.RoomID)
		}
		if req.ClassID != nil && *req.ClassID != 0 {
			query = query.Where("class_id = ?", *req.ClassID)
		}
		if req.StudentID != nil && *req.StudentID != 0 {
			query = query.Where("student_id = ?", *req.StudentID)
		}
		if req.CourseID != nil && *req.CourseID != 0 {
			query = query.Where("course_id = ?", *req.CourseID)
		}
	}

	var lessons []models.EduLesson
	if err := query.Find(&lessons).Error; err != nil {
		return nil, err
	}

	resp := &models.EduScheduleStatisticsResponse{}
	teacherStats := map[uint]*models.EduTeacherLessonSummary{}
	roomStats := map[uint]*models.EduRoomUsageSummary{}
	classStats := map[uint]*models.EduClassLessonSummary{}
	oneToOneStats := map[uint]*models.EduOneToOneLessonSummary{}

	for i := range lessons {
		lesson := &lessons[i]
		if !isCountedLessonStatus(lesson.Status) {
			continue
		}
		minutes := lessonDurationMinutes(lesson.StartTime, lesson.EndTime)
		if minutes <= 0 {
			continue
		}

		if lesson.TeacherID != 0 {
			summary := teacherStats[lesson.TeacherID]
			if summary == nil {
				summary = &models.EduTeacherLessonSummary{TeacherID: lesson.TeacherID}
				teacherStats[lesson.TeacherID] = summary
			}
			summary.LessonCount++
			summary.TotalMinutes += minutes
		}

		if lesson.RequiresRoom == 1 && lesson.RoomID != 0 {
			summary := roomStats[lesson.RoomID]
			if summary == nil {
				summary = &models.EduRoomUsageSummary{RoomID: lesson.RoomID}
				roomStats[lesson.RoomID] = summary
			}
			summary.LessonCount++
			summary.TotalMinutes += minutes
		}

		if lesson.LessonType == "class" || lesson.ClassID != 0 {
			if lesson.ClassID != 0 {
				summary := classStats[lesson.ClassID]
				if summary == nil {
					summary = &models.EduClassLessonSummary{ClassID: lesson.ClassID}
					classStats[lesson.ClassID] = summary
				}
				summary.LessonCount++
				summary.TotalMinutes += minutes
			}
		}
		if lesson.LessonType == "one_to_one" || lesson.StudentID != 0 {
			if lesson.StudentID != 0 {
				summary := oneToOneStats[lesson.StudentID]
				if summary == nil {
					summary = &models.EduOneToOneLessonSummary{StudentID: lesson.StudentID}
					oneToOneStats[lesson.StudentID] = summary
				}
				summary.LessonCount++
				summary.TotalMinutes += minutes
			}
		}
	}

	resp.TeacherLessons = collectTeacherStats(teacherStats)
	resp.RoomUsages = collectRoomStats(roomStats)
	resp.ClassLessons = collectClassStats(classStats)
	resp.OneToOneLessons = collectOneToOneStats(oneToOneStats)

	var eligibilities []models.EduLessonStudentEligibility
	if err := app.DB().WithContext(ctx).Model(&models.EduLessonStudentEligibility{}).
		Where("tenant_id = ? AND eligibility_status IN ?", tenantID, []string{"ineligible", "warn"}).
		Order("id asc").
		Find(&eligibilities).Error; err != nil {
		return nil, err
	}
	for i := range eligibilities {
		item := &eligibilities[i]
		resp.EligibilityDetails = append(resp.EligibilityDetails, models.EduLessonEligibilityDetail{
			LessonID:         item.LessonID,
			StudentID:        item.StudentID,
			CourseID:         item.CourseID,
			ClassID:          item.ClassID,
			StudentBenefitID: item.StudentBenefitID,
			EligibilityStatus: item.EligibilityStatus,
			ReasonCode:       item.ReasonCode,
			CheckedAt:        item.CheckedAt,
		})
	}

	var overrides []struct {
		ConflictType string
		Count        int64
	}
	if err := app.DB().WithContext(ctx).Model(&models.EduScheduleConflictOverride{}).
		Select("conflict_type, COUNT(1) AS count").
		Where("tenant_id = ?", tenantID).
		Group("conflict_type").
		Order("conflict_type asc").
		Scan(&overrides).Error; err != nil {
		return nil, err
	}
	for _, item := range overrides {
		resp.ConflictOverrides = append(resp.ConflictOverrides, models.EduScheduleConflictTypeSummary{
			ConflictType: item.ConflictType,
			Count:        item.Count,
		})
	}

	var conflictRows []models.EduScheduleConflictOverride
	if err := app.DB().WithContext(ctx).Model(&models.EduScheduleConflictOverride{}).
		Where("tenant_id = ?", tenantID).
		Order("id asc").
		Find(&conflictRows).Error; err != nil {
		return nil, err
	}
	for i := range conflictRows {
		item := &conflictRows[i]
		resp.ConflictDetails = append(resp.ConflictDetails, models.EduScheduleConflictOverrideDetail{
			RuleID:       item.RuleID,
			LessonID:     item.LessonID,
			ConflictType: item.ConflictType,
			ConflictKey:  item.ConflictKey,
			Reason:       item.Reason,
			OperatorID:   item.OperatorID,
			OccurredAt:   item.OccurredAt,
		})
	}
	return resp, nil
}

func (s *EduStatisticsService) BenefitStatistics(ctx context.Context, req *models.EduStatisticsRangeRequest) (*models.EduBenefitStatisticsResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return &models.EduBenefitStatisticsResponse{}, nil
	}
	now := time.Now().UTC()
	upper := now.Add(30 * 24 * time.Hour)
	if req != nil && req.EndDate != nil && !req.EndDate.Time.IsZero() {
		upper = req.EndDate.Time
	}
	lower := now
	if req != nil && req.StartDate != nil && !req.StartDate.Time.IsZero() {
		lower = req.StartDate.Time
	}

	var benefits []models.EduStudentBenefit
	query := app.DB().WithContext(ctx).Model(&models.EduStudentBenefit{}).Where("tenant_id = ?", tenantID)
	if req != nil {
		if req.StudentID != nil && *req.StudentID != 0 {
			query = query.Where("student_id = ?", *req.StudentID)
		}
		if req.CourseID != nil && *req.CourseID != 0 {
			query = query.Where("course_id = ?", *req.CourseID)
		}
	}
	if err := query.Find(&benefits).Error; err != nil {
		return nil, err
	}
	resp := &models.EduBenefitStatisticsResponse{}
	courseCounts := map[uint]int64{}
	for i := range benefits {
		benefit := &benefits[i]
		if benefit.Status == 1 {
			resp.ValidCount++
		}
		if benefit.ValidTo != nil {
			switch {
			case benefit.ValidTo.Time.Before(now):
				resp.ExpiredCount++
			case !benefit.ValidTo.Time.Before(lower) && !benefit.ValidTo.Time.After(upper):
				resp.ExpiringSoonCount++
			}
		}
		if benefit.CalculationMode == "count_limited" && benefit.RemainingCount <= 0 {
			resp.ZeroRemainingCount++
		}
		if benefit.CourseID != 0 {
			courseCounts[benefit.CourseID]++
		}
	}
	for courseID, count := range courseCounts {
		resp.CourseBenefitCounts = append(resp.CourseBenefitCounts, models.EduBenefitCourseSummary{
			CourseID:     courseID,
			BenefitCount: count,
		})
	}
	return resp, nil
}

func (s *EduStatisticsService) ExternalSyncStatistics(ctx context.Context, req *models.EduStatisticsRangeRequest) (*models.EduExternalSyncStatisticsResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if tenantID == 0 {
		return &models.EduExternalSyncStatisticsResponse{}, nil
	}
	query := app.DB().WithContext(ctx).Model(&models.EduBenefitExternalSync{}).Where("tenant_id = ?", tenantID)
	if req != nil {
		if req.StudentID != nil && *req.StudentID != 0 {
			query = query.Where("student_id = ?", *req.StudentID)
		}
	}
	rows := make([]struct {
		Status string
		Count  int64
	}, 0)
	if err := query.Select("status, COUNT(1) AS count").Group("status").Order("status asc").Scan(&rows).Error; err != nil {
		return nil, err
	}
	resp := &models.EduExternalSyncStatisticsResponse{}
	for _, row := range rows {
		resp.StatusCounts = append(resp.StatusCounts, models.EduExternalSyncStatusSummary{
			Status: row.Status,
			Count:  row.Count,
		})
		switch row.Status {
		case "failed":
			resp.FailedCount = row.Count
		case "retrying":
			resp.RetryingCount = row.Count
		}
	}
	if err := app.DB().WithContext(ctx).Model(&models.EduBenefitExternalSync{}).
		Where("tenant_id = ? AND status = ? AND next_retry_at IS NOT NULL AND next_retry_at <= ?", tenantID, "retrying", time.Now().UTC()).
		Count(&resp.DueRetryCount).Error; err != nil {
		return nil, err
	}
	return resp, nil
}

func isCountedLessonStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "scheduled", "completed":
		return true
	default:
		return false
	}
}

func lessonDurationMinutes(startTime, endTime string) int64 {
	start, err := time.Parse("15:04", strings.TrimSpace(startTime))
	if err != nil {
		return 0
	}
	end, err := time.Parse("15:04", strings.TrimSpace(endTime))
	if err != nil {
		return 0
	}
	return int64(end.Hour()*60 + end.Minute() - (start.Hour()*60 + start.Minute()))
}

func collectTeacherStats(src map[uint]*models.EduTeacherLessonSummary) []models.EduTeacherLessonSummary {
	out := make([]models.EduTeacherLessonSummary, 0, len(src))
	for _, item := range src {
		out = append(out, *item)
	}
	return out
}

func collectRoomStats(src map[uint]*models.EduRoomUsageSummary) []models.EduRoomUsageSummary {
	out := make([]models.EduRoomUsageSummary, 0, len(src))
	for _, item := range src {
		out = append(out, *item)
	}
	return out
}

func collectClassStats(src map[uint]*models.EduClassLessonSummary) []models.EduClassLessonSummary {
	out := make([]models.EduClassLessonSummary, 0, len(src))
	for _, item := range src {
		out = append(out, *item)
	}
	return out
}

func collectOneToOneStats(src map[uint]*models.EduOneToOneLessonSummary) []models.EduOneToOneLessonSummary {
	out := make([]models.EduOneToOneLessonSummary, 0, len(src))
	for _, item := range src {
		out = append(out, *item)
	}
	return out
}
