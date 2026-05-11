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
	if eduStudentControllers == nil || eduCourseControllers == nil || eduClassControllers == nil || eduRoomControllers == nil || eduTermControllers == nil {
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
	}
	for _, route := range want {
		if !got[route] {
			t.Fatalf("route %s not registered", route)
		}
	}
}
