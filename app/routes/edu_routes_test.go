package routes

import (
	"testing"
	"time"

	"gin-fast/app/global/app"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type fakeRouteEduConfig struct{}

func (fakeRouteEduConfig) ConfigFileChangeListen(...func()) {}
func (fakeRouteEduConfig) Get(string) interface{}           { return nil }
func (fakeRouteEduConfig) GetString(string) string          { return "mysql" }
func (fakeRouteEduConfig) GetBool(string) bool              { return false }
func (fakeRouteEduConfig) GetInt(string) int                { return 0 }
func (fakeRouteEduConfig) GetInt32(string) int32            { return 0 }
func (fakeRouteEduConfig) GetInt64(string) int64            { return 0 }
func (fakeRouteEduConfig) GetFloat64(string) float64        { return 0 }
func (fakeRouteEduConfig) GetDuration(string) time.Duration { return 0 }
func (fakeRouteEduConfig) GetStringSlice(string) []string   { return nil }
func (fakeRouteEduConfig) GetUintSlice(string) []uint       { return nil }
func (fakeRouteEduConfig) Set(string, interface{})          {}
func (fakeRouteEduConfig) SaveConfig() error                { return nil }

type fakeRouteCasbin struct{}

func (fakeRouteCasbin) InitCasbin(*gorm.DB, string) error                         { return nil }
func (fakeRouteCasbin) Enforce(string, string, string, ...string) (bool, error)   { return true, nil }
func (fakeRouteCasbin) GetEnforcer() *casbin.Enforcer                             { return nil }
func (fakeRouteCasbin) CasbinMiddleware() gin.HandlerFunc                         { return func(c *gin.Context) { c.Next() } }
func (fakeRouteCasbin) AddRolesForUserByID(uint, []uint, ...string) error         { return nil }
func (fakeRouteCasbin) DeleteRolesForUserByID(uint, []uint, ...string) error      { return nil }
func (fakeRouteCasbin) AddRoleInheritance(uint, uint, ...string) error            { return nil }
func (fakeRouteCasbin) DeleteRoleInheritance(uint, uint, ...string) error         { return nil }
func (fakeRouteCasbin) AddPolicyForRole(uint, string, string, ...string) error    { return nil }
func (fakeRouteCasbin) RemovePolicyForRole(uint, string, string, ...string) error { return nil }
func (fakeRouteCasbin) AddPoliciesForRole(uint, [][]string, ...string) error      { return nil }
func (fakeRouteCasbin) RemoveAllPoliciesForRole(uint, ...string) error            { return nil }
func (fakeRouteCasbin) GetRolesForUserByID(uint, ...string) ([]uint, error)       { return nil, nil }
func (fakeRouteCasbin) GetUsersForRole(uint, ...string) ([]uint, error)           { return nil, nil }
func (fakeRouteCasbin) HasRoleForUser(uint, uint, ...string) (bool, error)        { return true, nil }
func (fakeRouteCasbin) GetPermissionsForUser(uint, ...string) ([][]string, error) { return nil, nil }
func (fakeRouteCasbin) StopAutoLoadPolicy()                                       {}
func (fakeRouteCasbin) PrefixDomain(uint) string                                  { return "domain_1" }

func TestEduControllersAreRegistered(t *testing.T) {
	if eduStudentControllers == nil || eduCourseControllers == nil || eduClassControllers == nil || eduRoomControllers == nil || eduTermControllers == nil || eduBenefitControllers == nil || eduScheduleControllers == nil || eduStatisticsControllers == nil || eduLessonCompletionControllers == nil {
		t.Fatal("edu controllers have nil registrations")
	}
}

func TestEduRoutes(t *testing.T) {
	app.ConfigYml = fakeRouteEduConfig{}
	app.CasbinV2 = fakeRouteCasbin{}
	engine := gin.New()
	InitRoutes(engine)

	got := map[string]bool{}
	for _, rt := range engine.Routes() {
		got[rt.Method+" "+rt.Path] = true
	}

	want := []string{
		"GET /api/edu/terms/list",
		"POST /api/edu/terms/add",
		"PUT /api/edu/terms/edit",
		"DELETE /api/edu/terms/delete",
		"GET /api/edu/terms/:id/closed-days",
		"POST /api/edu/terms/:id/closed-days/save",
		"GET /api/edu/benefit-products/list",
		"POST /api/edu/benefit-products/add",
		"PUT /api/edu/benefit-products/edit",
		"DELETE /api/edu/benefit-products/delete",
		"GET /api/edu/student-benefits/list",
		"POST /api/edu/student-benefits/add",
		"PUT /api/edu/student-benefits/edit",
		"POST /api/edu/student-benefits/check",
		"POST /api/edu/student-benefits/repair-schedule",
		"GET /api/edu/benefit-ledgers/list",
		"GET /api/edu/benefit-external-sync/list",
		"POST /api/edu/benefit-external-sync/retry",
		"GET /api/edu/schedule-rules/list",
		"POST /api/edu/schedule-rules/add",
		"PUT /api/edu/schedule-rules/edit",
		"POST /api/edu/schedule-rules/preview-change",
		"DELETE /api/edu/schedule-rules/delete",
		"GET /api/edu/lessons/calendar",
		"GET /api/edu/lessons/list",
		"PUT /api/edu/lessons/reschedule",
		"PUT /api/edu/lessons/stop",
		"PUT /api/edu/lessons/cancel",
		"PUT /api/edu/lessons/restore",
		"POST /api/edu/lessons/makeup",
		"GET /api/edu/lessons/:id/change-logs",
		"POST /api/edu/schedules/check-conflicts",
		"GET /api/edu/statistics/schedule",
		"GET /api/edu/statistics/benefits",
		"GET /api/edu/statistics/external-sync",
		"GET /api/edu/lesson-completion/rules",
		"PUT /api/edu/lesson-completion/rules",
		"GET /api/edu/lessons/:id/completions",
		"POST /api/edu/lessons/:id/completions",
		"POST /api/edu/lessons/:id/completions/:studentId/revoke",
		"POST /api/edu/lessons/:id/complete",
	}
	for _, route := range want {
		if !got[route] {
			t.Fatalf("route %s not registered", route)
		}
	}
}
