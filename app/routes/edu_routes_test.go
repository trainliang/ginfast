package routes

import "testing"

func TestEduControllersAreRegistered(t *testing.T) {
	if eduStudentControllers == nil {
		t.Fatal("eduStudentControllers is nil")
	}
	if eduCourseControllers == nil {
		t.Fatal("eduCourseControllers is nil")
	}
	if eduClassControllers == nil {
		t.Fatal("eduClassControllers is nil")
	}
	if eduRoomControllers == nil {
		t.Fatal("eduRoomControllers is nil")
	}
}
