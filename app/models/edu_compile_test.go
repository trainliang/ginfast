package models

import "testing"

func TestEduModelsExposeExpectedTableNames(t *testing.T) {
	cases := map[string]string{
		"course":         (EduCourse{}).TableName(),
		"student":        (EduStudent{}).TableName(),
		"studentContact": (EduStudentContact{}).TableName(),
		"class":          (EduClass{}).TableName(),
		"classMember":    (EduClassMember{}).TableName(),
		"room":           (EduRoom{}).TableName(),
		"weeklyRule":     (EduRoomWeeklyRule{}).TableName(),
		"exception":      (EduRoomException{}).TableName(),
	}

	expected := map[string]string{
		"course":         "edu_course",
		"student":        "edu_student",
		"studentContact": "edu_student_contact",
		"class":          "edu_class",
		"classMember":    "edu_class_member",
		"room":           "edu_room",
		"weeklyRule":     "edu_room_weekly_rule",
		"exception":      "edu_room_exception",
	}

	for name, table := range cases {
		if table != expected[name] {
			t.Fatalf("%s table = %s, want %s", name, table, expected[name])
		}
	}
}
