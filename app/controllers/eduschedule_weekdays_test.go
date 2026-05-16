package controllers

import (
	"testing"

	"gin-fast/app/models"
)

func TestScheduleControllerBuildRuleMapsWeekdays(t *testing.T) {
	t.Run("add request maps weekdays", func(t *testing.T) {
		req := &models.EduScheduleRuleAddRequest{
			Name:         "规则A",
			RuleType:     "class",
			RepeatType:   "weekly",
			TermID:       "1",
			ClassID:      "2",
			StudentID:    "3",
			CourseID:     "4",
			TeacherID:    "5",
			RequiresRoom: 1,
			RoomID:       "6",
			Weekdays:     []int8{2, 4},
			StartTime:    "09:00",
			EndTime:      "10:00",
		}

		rule := buildEduScheduleRuleFromAddRequest(req, 11, 22)
		if got := rule.Weekdays; len(got) != 2 || got[0] != 2 || got[1] != 4 {
			t.Fatalf("unexpected weekdays: %+v", got)
		}
		if rule.TermID != 1 || rule.ClassID != 2 || rule.StudentID != 3 || rule.CourseID != 4 || rule.TeacherID != 5 || rule.RoomID != 6 {
			t.Fatalf("unexpected id conversion: %+v", rule)
		}
	})

	t.Run("update request maps weekdays", func(t *testing.T) {
		req := &models.EduScheduleRuleUpdateRequest{
			EduScheduleRuleAddRequest: models.EduScheduleRuleAddRequest{
				Name:         "规则B",
				RuleType:     "class",
				RepeatType:   "weekly",
				CourseID:     "8",
				TeacherID:    "9",
				RequiresRoom: 1,
				RoomID:       "10",
				Weekdays:     []int8{1, 3, 5},
				StartTime:    "11:00",
				EndTime:      "12:00",
			},
			ID: "99",
		}

		rule := buildEduScheduleRuleFromUpdateRequest(req, 33, 44)
		if got := rule.Weekdays; len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 5 {
			t.Fatalf("unexpected weekdays: %+v", got)
		}
		if rule.ID != 99 {
			t.Fatalf("unexpected id: %+v", rule.ID)
		}
	})
}
