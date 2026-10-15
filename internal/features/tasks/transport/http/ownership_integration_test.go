package tasks_transport_http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	pool "github.com/Muxammednone/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	server "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/server"
	token "github.com/Muxammednone/golang-todoapp/internal/features/auth/token"
	statsrepo "github.com/Muxammednone/golang-todoapp/internal/features/statistics/repository/postgres"
	statsservice "github.com/Muxammednone/golang-todoapp/internal/features/statistics/service"
	statshttp "github.com/Muxammednone/golang-todoapp/internal/features/statistics/transport/http"
	tasksrepo "github.com/Muxammednone/golang-todoapp/internal/features/tasks/repository/postgres"
	tasksservice "github.com/Muxammednone/golang-todoapp/internal/features/tasks/service"
	taskshttp "github.com/Muxammednone/golang-todoapp/internal/features/tasks/transport/http"
	usersrepo "github.com/Muxammednone/golang-todoapp/internal/features/users/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Opt-in: requires migrated PostgreSQL and POSTGRES_* environment variables.
// Only profiles created by this test and their tasks are removed at cleanup.
func TestTaskOwnershipIntegration(t *testing.T) {
	if os.Getenv("TODOAPP_INTEGRATION") != "1" {
		t.Skip("set TODOAPP_INTEGRATION=1 to test against PostgreSQL")
	}
	ctx := context.Background()
	cfg, err := pool.NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	db, err := pool.NewPool(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(db.Close)
	profiles := usersrepo.NewUsersRepository(db)
	baseID := time.Now().UnixNano()
	owners := make([]domain.User, 2)
	secret := []byte("ownership-integration-test-key-only")
	signed := func(id int64, expiry time.Time) string {
		t.Helper()
		value, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject: strconv.FormatInt(id, 10), ExpiresAt: jwt.NewNumericDate(expiry),
		}).SignedString(secret)
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	tokens := make([]string, 2)
	for i := range owners {
		owner, err := profiles.CreateProfile(ctx, baseID+int64(i), domain.User{FullName: "Ownership integration test"})
		if err != nil {
			t.Fatal(err)
		}
		owners[i] = owner
		t.Cleanup(func() {
			if _, err := db.Exec(context.Background(), "DELETE FROM todoapp.tasks WHERE author_user_id=$1", owner.ID); err != nil {
				t.Error(err)
			}
			if _, err := db.Exec(context.Background(), "DELETE FROM todoapp.users WHERE id=$1 AND sso_user_id=$2", owner.ID, baseID+int64(i)); err != nil {
				t.Error(err)
			}
		})
		tokens[i] = signed(baseID+int64(i), time.Now().Add(time.Hour))
	}
	auth := middleware.Auth(token.NewVerifier(secret))
	profile := middleware.RequireProfile(profiles)
	router := server.NewAPIVersionRouter(server.ApiVersion1)
	router.RegisterRoutes(taskshttp.NewTasksHTTPHandler(tasksservice.NewTasksService(tasksrepo.NewTaskRepository(db)), auth, profile).Routes()...)
	router.RegisterRoutes(statshttp.NewStatisticsHTTPHandler(statsservice.NewStatisticsService(statsrepo.NewStatisticsRepository(db)), auth, profile).Routes()...)
	log := &logger.Logger{Logger: zap.NewNop()}
	request := func(method, path, body, bearer string, want int) []byte {
		t.Helper()
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r = r.WithContext(logger.ToContext(r.Context(), log))
		if bearer != "" {
			r.Header.Set("Authorization", "Bearer "+bearer)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if w.Code != want {
			t.Fatalf("%s %s: status %d, want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		return w.Body.Bytes()
	}
	for _, route := range []struct{ method, path string }{
		{"POST", "/tasks"}, {"GET", "/tasks"}, {"GET", "/tasks/1"}, {"PATCH", "/tasks/1"}, {"DELETE", "/tasks/1"}, {"GET", "/statistics"},
	} {
		request(route.method, route.path, `{}`, "", http.StatusUnauthorized)
	}
	request("GET", "/tasks", "", "broken-token", 401)
	request("GET", "/tasks", "", signed(baseID, time.Now().Add(-time.Hour)), 401)
	request("GET", "/tasks", "", signed(baseID+2, time.Now().Add(time.Hour)), 404)
	var created taskshttp.TaskDTOResponse
	body := fmt.Sprintf(`{"title":"Private task","author_user_id":%d}`, owners[1].ID)
	if err := json.Unmarshal(request("POST", "/tasks", body, tokens[0], 201), &created); err != nil {
		t.Fatal(err)
	}
	if created.AuthorUserID != owners[0].ID {
		t.Fatalf("client overrode owner: %d", created.AuthorUserID)
	}
	path := fmt.Sprintf("/tasks/%d", created.ID)
	request("GET", path, "", tokens[0], 200)
	request("GET", path, "", tokens[1], 404)
	request("PATCH", path, `{"title":"stolen"}`, tokens[1], 404)
	request("DELETE", path, "", tokens[1], 404)
	var list []taskshttp.TaskDTOResponse
	if err := json.Unmarshal(request("GET", fmt.Sprintf("/tasks?user_id=%d", owners[0].ID), "", tokens[1], 200), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatal("foreign tasks leaked")
	}
	if err := json.Unmarshal(request("GET", "/tasks?limit=1&offset=0", "", tokens[0], 200), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatal("own task missing")
	}
	for i, want := range []int{1, 0} {
		var stats statshttp.GetStatisticsResponse
		if err := json.Unmarshal(request("GET", fmt.Sprintf("/statistics?user_id=%d", owners[0].ID), "", tokens[i], 200), &stats); err != nil {
			t.Fatal(err)
		}
		if stats.TasksCreated != want {
			t.Fatalf("account %d: task count %d, want %d", i, stats.TasksCreated, want)
		}
	}
	var updated taskshttp.TaskDTOResponse
	if err := json.Unmarshal(request("PATCH", path, `{"completed":true}`, tokens[0], 200), &updated); err != nil {
		t.Fatal(err)
	}
	if !updated.Completed || updated.Title != "Private task" {
		t.Fatal("update failed or foreign patch changed task")
	}
	request("DELETE", path, "", tokens[0], 204)
	request("GET", path, "", tokens[0], 404)
}
