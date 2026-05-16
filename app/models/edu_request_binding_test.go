package models

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestEduRequestBindingSupportsStringNumbers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("schedule rule add", func(t *testing.T) {
		var req EduScheduleRuleAddRequest
		err := bindJSONForEduRequest(t, `{
			"name":"规则A",
			"ruleType":"class",
			"repeatType":"weekly",
			"termId":"1",
			"classId":"2",
			"studentId":"3",
			"courseId":"4",
			"teacherId":"5",
			"requiresRoom":1,
			"roomId":"6",
			"weekdays":[2,4],
			"startTime":"09:00",
			"endTime":"10:00",
			"status":"1",
			"version":"2"
		}`, &req)
		if err != nil {
			t.Fatalf("bind schedule rule add: %v", err)
		}
		if req.GetTermIDUint() != 1 || req.GetClassIDUint() != 2 || req.GetStudentIDUint() != 3 || req.GetCourseIDUint() != 4 || req.GetTeacherIDUint() != 5 || req.GetRoomIDUint() != 6 {
			t.Fatalf("unexpected converted ids: %+v", req)
		}
		if req.Status == nil || int8(*req.Status) != 1 {
			t.Fatalf("unexpected status: %+v", req.Status)
		}
		if req.Version == nil || int(*req.Version) != 2 {
			t.Fatalf("unexpected version: %+v", req.Version)
		}
		if len(req.Weekdays) != 2 || req.Weekdays[0] != 2 || req.Weekdays[1] != 4 {
			t.Fatalf("unexpected weekdays: %+v", req.Weekdays)
		}
	})

	t.Run("statistics range", func(t *testing.T) {
		var req EduStatisticsRangeRequest
		err := bindJSONForEduRequest(t, `{
			"teacherId":"7",
			"roomId":"8",
			"classId":"9",
			"studentId":"10",
			"courseId":"11"
		}`, &req)
		if err != nil {
			t.Fatalf("bind statistics range: %v", err)
		}
		if req.GetTeacherIDUint() != 7 || req.GetRoomIDUint() != 8 || req.GetClassIDUint() != 9 || req.GetStudentIDUint() != 10 || req.GetCourseIDUint() != 11 {
			t.Fatalf("unexpected converted range ids: %+v", req)
		}
	})

	t.Run("student benefit add", func(t *testing.T) {
		var req EduStudentBenefitAddRequest
		err := bindJSONForEduRequest(t, `{
			"studentId":"1",
			"productId":"2",
			"courseId":"3",
			"classId":"4",
			"teacherId":"5",
			"totalCount":"20",
			"usedCount":"6",
			"remainingCount":"14",
			"status":"1"
		}`, &req)
		if err != nil {
			t.Fatalf("bind student benefit add: %v", err)
		}
		if uint(req.StudentID) != 1 || uint(req.ProductID) != 2 || uint(req.CourseID) != 3 || uint(req.ClassID) != 4 || uint(req.TeacherID) != 5 {
			t.Fatalf("unexpected converted benefit ids: %+v", req)
		}
		if int(req.TotalCount) != 20 || int(req.UsedCount) != 6 || int(req.RemainingCount) != 14 {
			t.Fatalf("unexpected benefit counts: %+v", req)
		}
		if req.Status == nil || int8(*req.Status) != 1 {
			t.Fatalf("unexpected benefit status: %+v", req.Status)
		}
	})

	t.Run("student update contacts", func(t *testing.T) {
		var req EduStudentUpdateRequest
		err := bindJSONForEduRequest(t, `{
			"id":"12",
			"name":"学生A",
			"status":"1",
			"contacts":[
				{"id":"21","studentId":"12","name":"家长A","phone":"13800000000","isPrimary":"1","canPickup":"0"}
			]
		}`, &req)
		if err != nil {
			t.Fatalf("bind student update: %v", err)
		}
		if uint(req.ID) != 12 {
			t.Fatalf("unexpected student id: %+v", req.ID)
		}
		if req.Status == nil || int8(*req.Status) != 1 {
			t.Fatalf("unexpected student status: %+v", req.Status)
		}
		if len(req.Contacts) != 1 {
			t.Fatalf("unexpected contacts: %+v", req.Contacts)
		}
		contact := req.Contacts[0]
		if uint(contact.ID) != 21 || uint(contact.StudentID) != 12 {
			t.Fatalf("unexpected contact ids: %+v", contact)
		}
		if contact.IsPrimary == nil || int8(*contact.IsPrimary) != 1 || contact.CanPickup == nil || int8(*contact.CanPickup) != 0 {
			t.Fatalf("unexpected contact flags: %+v", contact)
		}
	})
}

func bindJSONForEduRequest(t *testing.T, body string, req interface{}) error {
	t.Helper()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	httpReq := httptest.NewRequest("POST", "/", strings.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	ctx.Request = httpReq

	validator, ok := req.(interface {
		Validate(*gin.Context) error
	})
	if !ok {
		t.Fatalf("request does not implement Validate: %T", req)
	}
	return validator.Validate(ctx)
}
