package models

// EduStatisticsRangeRequest 统计范围请求
type EduStatisticsRangeRequest struct {
	StartDate *JSONTime `json:"startDate" form:"startDate"`
	EndDate   *JSONTime `json:"endDate" form:"endDate"`
	TeacherID *uint     `json:"teacherId" form:"teacherId"`
	RoomID    *uint     `json:"roomId" form:"roomId"`
	ClassID   *uint     `json:"classId" form:"classId"`
	StudentID *uint     `json:"studentId" form:"studentId"`
	CourseID  *uint     `json:"courseId" form:"courseId"`
}

// EduTeacherLessonSummary 教师课时汇总
type EduTeacherLessonSummary struct {
	TeacherID    uint  `json:"teacherId"`
	LessonCount  int64 `json:"lessonCount"`
	TotalMinutes int64 `json:"totalMinutes"`
}

// EduRoomUsageSummary 场地使用汇总
type EduRoomUsageSummary struct {
	RoomID       uint  `json:"roomId"`
	LessonCount  int64 `json:"lessonCount"`
	TotalMinutes int64 `json:"totalMinutes"`
}

// EduClassLessonSummary 班级课次汇总
type EduClassLessonSummary struct {
	ClassID      uint  `json:"classId"`
	LessonCount  int64 `json:"lessonCount"`
	TotalMinutes int64 `json:"totalMinutes"`
}

// EduOneToOneLessonSummary 一对一课次汇总
type EduOneToOneLessonSummary struct {
	StudentID    uint  `json:"studentId"`
	LessonCount  int64 `json:"lessonCount"`
	TotalMinutes int64 `json:"totalMinutes"`
}

// EduScheduleConflictTypeSummary 冲突覆盖类型汇总
type EduScheduleConflictTypeSummary struct {
	ConflictType string `json:"conflictType"`
	Count        int64  `json:"count"`
}

// EduLessonEligibilityDetail 课次权益异常明细
type EduLessonEligibilityDetail struct {
	LessonID          uint      `json:"lessonId"`
	StudentID         uint      `json:"studentId"`
	CourseID          uint      `json:"courseId"`
	ClassID           uint      `json:"classId"`
	StudentBenefitID   uint      `json:"studentBenefitId"`
	EligibilityStatus  string    `json:"eligibilityStatus"`
	ReasonCode        string    `json:"reasonCode"`
	CheckedAt         *JSONTime `json:"checkedAt"`
}

// EduScheduleConflictOverrideDetail 冲突覆盖明细
type EduScheduleConflictOverrideDetail struct {
	RuleID       uint      `json:"ruleId"`
	LessonID     uint      `json:"lessonId"`
	ConflictType string    `json:"conflictType"`
	ConflictKey  string    `json:"conflictKey"`
	Reason       string    `json:"reason"`
	OperatorID   uint      `json:"operatorId"`
	OccurredAt   *JSONTime `json:"occurredAt"`
}

// EduScheduleStatisticsResponse 排课统计响应
type EduScheduleStatisticsResponse struct {
	TeacherLessons    []EduTeacherLessonSummary          `json:"teacherLessons"`
	RoomUsages        []EduRoomUsageSummary              `json:"roomUsages"`
	ClassLessons      []EduClassLessonSummary            `json:"classLessons"`
	OneToOneLessons   []EduOneToOneLessonSummary         `json:"oneToOneLessons"`
	EligibilityDetails []EduLessonEligibilityDetail      `json:"eligibilityDetails"`
	ConflictOverrides []EduScheduleConflictTypeSummary   `json:"conflictOverrides"`
	ConflictDetails   []EduScheduleConflictOverrideDetail `json:"conflictDetails"`
}

// EduBenefitCourseSummary 按课程汇总的权益数量
type EduBenefitCourseSummary struct {
	CourseID     uint  `json:"courseId"`
	BenefitCount int64 `json:"benefitCount"`
}

// EduBenefitStatisticsResponse 权益统计响应
type EduBenefitStatisticsResponse struct {
	ValidCount         int64                  `json:"validCount"`
	ExpiringSoonCount   int64                  `json:"expiringSoonCount"`
	ExpiredCount        int64                  `json:"expiredCount"`
	ZeroRemainingCount  int64                  `json:"zeroRemainingCount"`
	CourseBenefitCounts []EduBenefitCourseSummary `json:"courseBenefitCounts"`
}

// EduExternalSyncStatusSummary 外部同步状态汇总
type EduExternalSyncStatusSummary struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// EduExternalSyncStatisticsResponse 外部同步统计响应
type EduExternalSyncStatisticsResponse struct {
	StatusCounts  []EduExternalSyncStatusSummary `json:"statusCounts"`
	FailedCount   int64                         `json:"failedCount"`
	RetryingCount int64                         `json:"retryingCount"`
	DueRetryCount int64                         `json:"dueRetryCount"`
}
