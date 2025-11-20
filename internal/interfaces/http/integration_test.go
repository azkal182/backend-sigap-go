package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	appService "github.com/your-org/go-backend-starter/internal/application/service"
	"github.com/your-org/go-backend-starter/internal/application/usecase"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainRepo "github.com/your-org/go-backend-starter/internal/domain/repository"
	"github.com/your-org/go-backend-starter/internal/domain/service"
	"github.com/your-org/go-backend-starter/internal/infrastructure/database"
	infraRepo "github.com/your-org/go-backend-starter/internal/infrastructure/repository"
	infraService "github.com/your-org/go-backend-starter/internal/infrastructure/service"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/handler"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/middleware"
	"github.com/your-org/go-backend-starter/internal/interfaces/http/router"
	"github.com/your-org/go-backend-starter/internal/testutil"
	"gorm.io/gorm"
)

// testUserRepository wraps userRepository to use a specific DB
type testUserRepository struct {
	db *gorm.DB
}

func seedClass(t *testing.T, db *gorm.DB, fanID uuid.UUID) entity.Class {
	classEntity := entity.Class{
		ID:        uuid.New(),
		FanID:     fanID,
		Name:      fmt.Sprintf("Class-%d", time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&classEntity).Error; err != nil {
		t.Fatalf("failed to seed class: %v", err)
	}
	return classEntity
}

func seedTeacher(t *testing.T, db *gorm.DB) entity.Teacher {
	teacher := entity.Teacher{
		ID:          uuid.New(),
		TeacherCode: fmt.Sprintf("TCH-%d", time.Now().UnixNano()),
		FullName:    "Integration Teacher",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := db.Create(&teacher).Error; err != nil {
		t.Fatalf("failed to seed teacher: %v", err)
	}
	return teacher
}

func seedSubject(t *testing.T, db *gorm.DB, name string) entity.Subject {
	subject := entity.Subject{
		ID:        uuid.New(),
		Name:      fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&subject).Error; err != nil {
		t.Fatalf("failed to seed subject: %v", err)
	}
	return subject
}

func seedScheduleSlot(t *testing.T, db *gorm.DB, dormID uuid.UUID, slotNumber int) entity.ScheduleSlot {
	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(45 * time.Minute)
	slot := entity.ScheduleSlot{
		ID:          uuid.New(),
		DormitoryID: dormID,
		SlotNumber:  slotNumber,
		Name:        fmt.Sprintf("Slot-%d", slotNumber),
		StartTime:   start,
		EndTime:     end,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := db.Create(&slot).Error; err != nil {
		t.Fatalf("failed to seed schedule slot: %v", err)
	}
	return slot
}

func TestScheduleSlotEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Slot Dorm")
	user, token := createTestUser(t, db, "slot-admin", tokenService, "schedule_slots:read", "schedule_slots:create", "schedule_slots:update", "schedule_slots:delete")
	assignPermissionsToUser(t, db, user.ID, []string{"schedule_slots:read", "schedule_slots:create", "schedule_slots:update", "schedule_slots:delete"})

	start := time.Now().UTC().Add(time.Hour)
	end := start.Add(time.Hour)
	createPayload := map[string]interface{}{
		"dormitory_id": dorm.ID.String(),
		"slot_number":  1,
		"name":         "Morning Slot",
		"start_time":   start.Format(time.RFC3339),
		"end_time":     end.Format(time.RFC3339),
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/schedule-slots", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.ScheduleSlotResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	slotID := createResp.Data.ID
	assert.NotEmpty(t, slotID)

	listReq := httptest.NewRequest(http.MethodGet, "/api/schedule-slots?page=1&page_size=5&dormitory_id="+dorm.ID.String(), nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/api/schedule-slots/"+slotID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	updatePayload := map[string]interface{}{
		"slot_number": 2,
		"name":        "Updated Slot",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/schedule-slots/"+slotID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/schedule-slots/"+slotID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func TestClassScheduleEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Schedule Dorm")
	fan := seedFan(t, db)
	classEntity := seedClass(t, db, fan.ID)
	teacher := seedTeacher(t, db)
	subject := seedSubject(t, db, "Integration Subject")
	slot := seedScheduleSlot(t, db, dorm.ID, 1)

	user, token := createTestUser(
		t,
		db,
		"schedule-admin",
		tokenService,
		"class_schedules:read",
		"class_schedules:create",
		"class_schedules:update",
		"class_schedules:delete",
	)
	assignPermissionsToUser(t, db, user.ID, []string{"class_schedules:read", "class_schedules:create", "class_schedules:update", "class_schedules:delete"})

	createPayload := map[string]interface{}{
		"class_id":     classEntity.ID.String(),
		"dormitory_id": dorm.ID.String(),
		"subject_id":   subject.ID.String(),
		"teacher_id":   teacher.ID.String(),
		"slot_id":      slot.ID.String(),
		"day_of_week":  "mon",
		"location":     "Room 101",
		"notes":        "Initial schedule",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/class-schedules", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.ClassScheduleResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	scheduleID := createResp.Data.ID
	require.NotEmpty(t, scheduleID)

	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/class-schedules?class_id=%s", classEntity.ID.String()), nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/api/class-schedules/"+scheduleID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	updatePayload := map[string]interface{}{
		"notes": "Updated schedule",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/class-schedules/"+scheduleID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/class-schedules/"+scheduleID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func TestSKSDefinitionEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	fan := seedFan(t, db)
	subject := seedSubject(t, db, "Definition Subject")
	user, token := createTestUser(
		t,
		db,
		"sks-def-admin",
		tokenService,
		"sks_definitions:read",
		"sks_definitions:create",
		"sks_definitions:update",
		"sks_definitions:delete",
	)
	assignPermissionsToUser(t, db, user.ID, []string{"sks_definitions:read", "sks_definitions:create", "sks_definitions:update", "sks_definitions:delete"})

	createPayload := map[string]interface{}{
		"fan_id":      fan.ID.String(),
		"subject_id":  subject.ID.String(),
		"code":        fmt.Sprintf("SKS-%d", time.Now().UnixNano()),
		"name":        "Integration SKS",
		"kkm":         80,
		"description": "Initial definition",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/sks", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.SKSDefinitionResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	sksID := createResp.Data.ID

	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/sks?fan_id=%s", fan.ID.String()), nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/api/sks/"+sksID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	updatePayload := map[string]interface{}{
		"name": "Updated SKS",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/sks/"+sksID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/sks/"+sksID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func TestSKSExamEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	fan := seedFan(t, db)
	subject := seedSubject(t, db, "Exam Subject")
	sksDefinition := entity.SKSDefinition{
		ID:          uuid.New(),
		FanID:       fan.ID,
		Code:        fmt.Sprintf("SKS-%d", time.Now().UnixNano()),
		Name:        "Exam Definition",
		KKM:         75,
		Description: "Exam definition",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		SubjectID:   &subject.ID,
	}
	require.NoError(t, db.Create(&sksDefinition).Error)
	teacher := seedTeacher(t, db)

	user, token := createTestUser(
		t,
		db,
		"sks-exam-admin",
		tokenService,
		"sks_exams:read",
		"sks_exams:create",
		"sks_exams:update",
		"sks_exams:delete",
	)
	assignPermissionsToUser(t, db, user.ID, []string{"sks_exams:read", "sks_exams:create", "sks_exams:update", "sks_exams:delete"})

	createPayload := map[string]interface{}{
		"sks_id":      sksDefinition.ID.String(),
		"examiner_id": teacher.ID.String(),
		"exam_date":   "2025-01-10",
		"exam_time":   "09:00",
		"location":    "Hall A",
		"notes":       "Initial exam",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/sks-exams", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.SKSExamScheduleResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	examID := createResp.Data.ID

	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/sks-exams?sks_id=%s", sksDefinition.ID.String()), nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/api/sks-exams/"+examID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	updatePayload := map[string]interface{}{
		"notes": "Updated exam",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/sks-exams/"+examID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/sks-exams/"+examID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func TestTeacherEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	user, token := createTestUser(t, db, "teacher-admin", tokenService, "teachers:read", "teachers:create", "teachers:update", "teachers:delete")
	assignPermissionsToUser(t, db, user.ID, []string{"teachers:read", "teachers:create", "teachers:update", "teachers:delete"})

	createPayload := map[string]interface{}{
		"teacher_code": "TCHINT001",
		"full_name":    "Integration Teacher",
		"gender":       "male",
		"phone":        "+620000001",
		"email":        "integration.teacher@example.com",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/teachers", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.TeacherResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	teacherID := createResp.Data.ID
	assert.NotEmpty(t, teacherID)

	listReq := httptest.NewRequest(http.MethodGet, "/api/teachers?page=1&page_size=5", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/api/teachers/"+teacherID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	updatePayload := map[string]interface{}{
		"phone":     "+620000009",
		"is_active": true,
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/teachers/"+teacherID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/teachers/"+teacherID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func seedFan(t *testing.T, db *gorm.DB) entity.Fan {
	fan := entity.Fan{ID: uuid.New(), Name: "FAN Seed", Level: "junior", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(&fan).Error; err != nil {
		t.Fatalf("failed to seed fan: %v", err)
	}
	return fan
}

func TestFanEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	user, token := createTestUser(t, db, "fan-admin", tokenService, "fans:read", "fans:create", "fans:update", "fans:delete")
	assignPermissionsToUser(t, db, user.ID, []string{"fans:read", "fans:create", "fans:update", "fans:delete"})

	// Create
	createPayload := map[string]string{
		"name":  "Integration FAN",
		"level": "senior",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/fans", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.FanResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	fanID := createResp.Data.ID

	// List
	listReq := httptest.NewRequest(http.MethodGet, "/api/fans?page=1&page_size=10", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	// Get
	getReq := httptest.NewRequest(http.MethodGet, "/api/fans/"+fanID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	// Update
	updatePayload := map[string]string{"name": "Updated FAN"}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/fans/"+fanID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes := httptest.NewRecorder()
	router.ServeHTTP(updateRes, updateReq)
	assert.Equal(t, http.StatusOK, updateRes.Code)

	// Delete
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/fans/"+fanID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func TestClassEndpoints(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	user, token := createTestUser(t, db, "class-admin", tokenService, "classes:read", "classes:create", "classes:update", "classes:delete")
	assignPermissionsToUser(t, db, user.ID, []string{"classes:read", "classes:create", "classes:update", "classes:delete"})
	fan := seedFan(t, db)

	// Create class
	createPayload := map[string]interface{}{
		"fan_id": fan.ID.String(),
		"name":   "Integration Class",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/classes", bytes.NewBuffer(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)
	assert.Equal(t, http.StatusCreated, createRes.Code)

	var createResp struct {
		Data dto.ClassResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
	classID := createResp.Data.ID

	// List by fan
	listReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/classes?fan_id=%s", fan.ID), nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRes := httptest.NewRecorder()
	router.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)

	// Enroll student
	student := entity.Student{ID: uuid.New(), StudentNumber: "ST-CL-1", FullName: "Class Student", Gender: "male", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, db.Create(&student).Error)
	enrollPayload := map[string]string{
		"student_id": student.ID.String(),
		"start_date": time.Now().Format(time.RFC3339),
	}
	enrollBody, _ := json.Marshal(enrollPayload)
	enrollReq := httptest.NewRequest(http.MethodPost, "/api/classes/"+classID+"/students", bytes.NewBuffer(enrollBody))
	enrollReq.Header.Set("Authorization", "Bearer "+token)
	enrollReq.Header.Set("Content-Type", "application/json")
	enrollRes := httptest.NewRecorder()
	router.ServeHTTP(enrollRes, enrollReq)
	assert.Equal(t, http.StatusNoContent, enrollRes.Code)

	// Assign staff
	staffUser := entity.User{ID: uuid.New(), Username: "staff-class", Password: "hash", Name: "Staff", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, db.Create(&staffUser).Error)
	assignPayload := map[string]string{
		"user_id": staffUser.ID.String(),
		"role":    "class_manager",
	}
	assignBody, _ := json.Marshal(assignPayload)
	assignReq := httptest.NewRequest(http.MethodPost, "/api/classes/"+classID+"/staff", bytes.NewBuffer(assignBody))
	assignReq.Header.Set("Authorization", "Bearer "+token)
	assignReq.Header.Set("Content-Type", "application/json")
	assignRes := httptest.NewRecorder()
	router.ServeHTTP(assignRes, assignReq)
	assert.Equal(t, http.StatusNoContent, assignRes.Code)

	// Cleanup: delete class
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/classes/"+classID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRes := httptest.NewRecorder()
	router.ServeHTTP(deleteRes, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRes.Code)
}

func (r *testUserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *testUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *testUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *testUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *testUserRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) GetWithRolesAndDormitories(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Preload("Dormitories").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *testUserRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Create(&entity.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *testUserRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entity.UserRole{}).Error
}

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, service.TokenService, func()) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	testutil.SetTestEnv()

	// Temporarily replace database.DB for repositories
	originalDB := database.DB
	database.DB = testDB

	// Initialize repositories with test database
	userRepo := &testUserRepository{db: testDB}
	roleRepo := infraRepo.NewRoleRepository() // These will use database.DB
	dormitoryRepo := infraRepo.NewDormitoryRepository()
	studentRepo := infraRepo.NewStudentRepository()
	fanRepo := infraRepo.NewFanRepository()
	classRepo := infraRepo.NewClassRepository()
	enrollmentRepo := infraRepo.NewStudentClassEnrollmentRepository()
	classStaffRepo := infraRepo.NewClassStaffRepository()
	teacherRepo := infraRepo.NewTeacherRepository()
	subjectRepo := infraRepo.NewSubjectRepository()
	classScheduleRepo := infraRepo.NewClassScheduleRepository()
	sksDefinitionRepo := infraRepo.NewSKSDefinitionRepository()
	sksExamRepo := infraRepo.NewSKSExamScheduleRepository()
	scheduleSlotRepo := infraRepo.NewScheduleSlotRepository()
	permissionRepo := infraRepo.NewPermissionRepository()
	auditLogRepo := infraRepo.NewAuditLogRepository()
	provinceRepo := infraRepo.NewProvinceRepository()
	regencyRepo := infraRepo.NewRegencyRepository()
	districtRepo := infraRepo.NewDistrictRepository()
	villageRepo := infraRepo.NewVillageRepository()

	// Initialize services
	tokenService := infraService.NewJWTService()
	auditLogger := appService.NewAuditLogger(auditLogRepo)
	ensureRoleExists(t, roleRepo, "teacher")

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, auditLogger)
	dormitoryUseCase := usecase.NewDormitoryUseCase(dormitoryRepo, userRepo, auditLogger)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, dormitoryRepo, auditLogger)
	fanUseCase := usecase.NewFanUseCase(fanRepo, auditLogger)
	classUseCase := usecase.NewClassUseCase(classRepo, fanRepo, studentRepo, enrollmentRepo, classStaffRepo, auditLogger)
	teacherUseCase := usecase.NewTeacherUseCase(teacherRepo, userRepo, roleRepo, auditLogger)
	scheduleSlotUseCase := usecase.NewScheduleSlotUseCase(scheduleSlotRepo, dormitoryRepo, auditLogger)
	classScheduleUseCase := usecase.NewClassScheduleUseCase(classScheduleRepo, classRepo, teacherRepo, subjectRepo, scheduleSlotRepo, dormitoryRepo, auditLogger)
	sksDefinitionUseCase := usecase.NewSKSDefinitionUseCase(sksDefinitionRepo, fanRepo, subjectRepo, auditLogger)
	sksExamUseCase := usecase.NewSKSExamScheduleUseCase(sksExamRepo, sksDefinitionRepo, teacherRepo, auditLogger)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, auditLogger)
	locationUseCase := usecase.NewLocationUseCase(provinceRepo, regencyRepo, districtRepo, villageRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	auditLogUseCase := usecase.NewAuditLogUseCase(auditLogRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	dormitoryHandler := handler.NewDormitoryHandler(dormitoryUseCase)
	studentHandler := handler.NewStudentHandler(studentUseCase)
	fanHandler := handler.NewFanHandler(fanUseCase)
	classHandler := handler.NewClassHandler(classUseCase)
	teacherHandler := handler.NewTeacherHandler(teacherUseCase)
	scheduleSlotHandler := handler.NewScheduleSlotHandler(scheduleSlotUseCase)
	classScheduleHandler := handler.NewClassScheduleHandler(classScheduleUseCase)
	sksDefinitionHandler := handler.NewSKSDefinitionHandler(sksDefinitionUseCase)
	sksExamHandler := handler.NewSKSExamScheduleHandler(sksExamUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	locationHandler := handler.NewLocationHandler(locationUseCase)
	permissionHandler := handler.NewPermissionHandler(permissionUseCase)
	auditLogHandler := handler.NewAuditLogHandler(auditLogUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup router
	r := router.SetupRouter(
		authHandler,
		userHandler,
		dormitoryHandler,
		studentHandler,
		roleHandler,
		locationHandler,
		permissionHandler,
		auditLogHandler,
		fanHandler,
		classHandler,
		teacherHandler,
		classScheduleHandler,
		sksDefinitionHandler,
		sksExamHandler,
		scheduleSlotHandler,
		authMiddleware,
	)

	cleanup := func() {
		database.DB = originalDB // Restore original DB
		testutil.CleanupTestDB(t, testDB)
		testutil.UnsetTestEnv()
	}

	return r, testDB, tokenService, cleanup
}

func seedDormitory(t *testing.T, db *gorm.DB, name string) entity.Dormitory {
	dorm := entity.Dormitory{
		ID:        uuid.New(),
		Name:      name,
		Gender:    "male",
		Level:     "senior",
		Code:      fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&dorm).Error; err != nil {
		t.Fatalf("failed to seed dormitory: %v", err)
	}
	return dorm
}

// createTestUser creates an active user and returns user + bearer token
func createTestUser(t *testing.T, db *gorm.DB, username string, tokenService service.TokenService, permissions ...string) (entity.User, string) {
	user := entity.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  "$2a$10$Q9f1iG.zRV/X9sYt8GvGle6hwzEwA9H9n1tFoZT3zh0TBTtPlqHcC", // bcrypt hash for "password123"
		Name:      "Dorm Staff",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if len(permissions) == 0 {
		permissions = []string{"dorm:update"}
	}

	token, err := tokenService.GenerateAccessToken(user.ID, user.Username, permissions)
	if err != nil {
		t.Fatalf("failed generating access token: %v", err)
	}

	return user, token
}

func assignPermissionsToUser(t *testing.T, db *gorm.DB, userID uuid.UUID, permissionNames []string) {
	assignRoleWithPermissions(t, db, userID, fmt.Sprintf("role-%s", uuid.New().String()), permissionNames)
}

func assignRoleWithPermissions(t *testing.T, db *gorm.DB, userID uuid.UUID, roleName string, permissionNames []string) {
	permissions := make([]entity.Permission, 0, len(permissionNames))
	for _, name := range permissionNames {
		var perm entity.Permission
		if err := db.Where("name = ?", name).First(&perm).Error; err != nil {
			parts := strings.Split(name, ":")
			resource := parts[0]
			action := ""
			if len(parts) > 1 {
				action = parts[1]
			}
			perm = entity.Permission{
				ID:        uuid.New(),
				Name:      name,
				Slug:      fmt.Sprintf("%s-%d", strings.ReplaceAll(name, ":", "-"), time.Now().UnixNano()),
				Resource:  resource,
				Action:    action,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&perm).Error; err != nil {
				t.Fatalf("failed to create permission %s: %v", name, err)
			}
		}
		permissions = append(permissions, perm)
	}

	role := entity.Role{
		ID:        uuid.New(),
		Name:      roleName,
		Slug:      fmt.Sprintf("%s-%d", roleName, time.Now().UnixNano()),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to create role: %v", err)
	}

	for _, perm := range permissions {
		rp := entity.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
		if err := db.Create(&rp).Error; err != nil {
			t.Fatalf("failed to assign permission %s: %v", perm.Name, err)
		}
	}

	userRole := entity.UserRole{UserID: userID, RoleID: role.ID}
	if err := db.Create(&userRole).Error; err != nil {
		t.Fatalf("failed to assign role to user: %v", err)
	}
}

func assignDormitoryAdminRole(t *testing.T, db *gorm.DB, userID uuid.UUID) {
	assignRoleWithPermissions(t, db, userID, "admin", []string{"dorm:update"})
}

func assignStudentAdminRole(t *testing.T, db *gorm.DB, userID uuid.UUID) {
	assignRoleWithPermissions(t, db, userID, "student-admin", []string{"student:read", "student:create", "student:update"})
}

func ensureRoleExists(t *testing.T, roleRepo domainRepo.RoleRepository, slug string) {
	ctx := context.Background()
	if role, _ := roleRepo.GetBySlug(ctx, slug); role != nil {
		return
	}
	role := &entity.Role{
		ID:          uuid.New(),
		Name:        strings.Title(strings.ReplaceAll(slug, "_", " ")),
		Slug:        slug,
		IsActive:    true,
		IsProtected: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := roleRepo.Create(ctx, role); err != nil {
		t.Fatalf("failed to create %s role: %v", slug, err)
	}
}

func TestAuthIntegration_RegisterAndLogin(t *testing.T) {
	router, _, _, cleanup := setupTestRouter(t)
	defer cleanup()

	// Test Register
	registerReq := dto.RegisterRequest{
		Username: "integration",
		Password: "password123",
		Name:     "Integration Test User",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)

	assert.Equal(t, http.StatusCreated, registerW.Code)

	var registerResp map[string]interface{}
	err := json.Unmarshal(registerW.Body.Bytes(), &registerResp)
	require.NoError(t, err)
	assert.True(t, registerResp["success"].(bool))

	// Extract tokens from response
	data := registerResp["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)
	require.NotEmpty(t, accessToken)

	// Test Login
	loginReq := dto.LoginRequest{
		Username: "integration",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusOK, loginW.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	assert.True(t, loginResp["success"].(bool))

	// Test protected endpoint with token
	// protectedReq, _ := http.NewRequest(http.MethodGet, "/api/users/me", nil)
	// protectedReq.Header.Set("Authorization", "Bearer "+accessToken)

	// protectedW := httptest.NewRecorder()
	// router.ServeHTTP(protectedW, protectedReq)

	// Note: This test assumes there's a /api/users/me endpoint
	// Adjust based on your actual routes
	// assert.True(t, protectedW.Code == http.StatusOK || protectedW.Code == http.StatusNotFound)
}

func TestAuthIntegration_InvalidCredentials(t *testing.T) {
	router, _, _, cleanup := setupTestRouter(t)
	defer cleanup()

	// Register user first
	registerReq := dto.RegisterRequest{
		Username: "test",
		Password: "password123",
		Name:     "Test User",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)
	assert.Equal(t, http.StatusCreated, registerW.Code)

	// Try login with wrong password
	loginReq := dto.LoginRequest{
		Username: "test",
		Password: "wrongpassword",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)
}

func TestDormitoryIntegration_AssignAndRemoveUser(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Integration Dorm")
	user, token := createTestUser(t, db, "dorm-staff", tokenService)
	assignDormitoryAdminRole(t, db, user.ID)

	assignPayload := map[string]string{"user_id": user.ID.String()}
	body, _ := json.Marshal(assignPayload)
	assignReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/dormitories/%s/users", dorm.ID.String()), bytes.NewBuffer(body))
	assignReq.Header.Set("Content-Type", "application/json")
	assignReq.Header.Set("Authorization", "Bearer "+token)

	assignW := httptest.NewRecorder()
	router.ServeHTTP(assignW, assignReq)
	assert.Equal(t, http.StatusNoContent, assignW.Code)

	removeReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/dormitories/%s/users/%s", dorm.ID.String(), user.ID.String()), nil)
	removeReq.Header.Set("Authorization", "Bearer "+token)

	removeW := httptest.NewRecorder()
	router.ServeHTTP(removeW, removeReq)
	assert.Equal(t, http.StatusNoContent, removeW.Code)
}

func TestStudentIntegration_CreateStatusMutate(t *testing.T) {
	router, db, tokenService, cleanup := setupTestRouter(t)
	defer cleanup()

	dorm := seedDormitory(t, db, "Student Dorm")
	studentNumber := fmt.Sprintf("STD%d", time.Now().UnixNano())
	user, token := createTestUser(t, db, "student-admin", tokenService, "student:read", "student:create", "student:update")
	assignStudentAdminRole(t, db, user.ID)

	birthDate := time.Now().AddDate(-15, 0, 0).UTC().Format(time.RFC3339)
	createPayload := map[string]interface{}{
		"student_number": studentNumber,
		"full_name":      "Integration Student",
		"birth_date":     birthDate,
		"gender":         "male",
		"parent_name":    "Integration Parent",
	}
	createBody, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/api/students", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	assert.Equal(t, http.StatusCreated, createW.Code)

	type studentResp struct {
		Success bool                `json:"success"`
		Data    dto.StudentResponse `json:"data"`
	}

	var created studentResp
	require.NoError(t, json.Unmarshal(createW.Body.Bytes(), &created))
	require.True(t, created.Success)
	require.NotEmpty(t, created.Data.ID)

	// Update status
	statusPayload := map[string]string{"status": entity.StudentStatusLeave}
	statusBody, _ := json.Marshal(statusPayload)
	statusReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/students/%s/status", created.Data.ID), bytes.NewBuffer(statusBody))
	statusReq.Header.Set("Content-Type", "application/json")
	statusReq.Header.Set("Authorization", "Bearer "+token)

	statusW := httptest.NewRecorder()
	router.ServeHTTP(statusW, statusReq)
	assert.Equal(t, http.StatusOK, statusW.Code)

	var statusResp studentResp
	require.NoError(t, json.Unmarshal(statusW.Body.Bytes(), &statusResp))
	assert.Equal(t, entity.StudentStatusLeave, statusResp.Data.Status)

	// Mutate dormitory
	mutatePayload := map[string]interface{}{
		"dormitory_id": dorm.ID.String(),
		"start_date":   time.Now().UTC().Format(time.RFC3339),
	}
	mutateBody, _ := json.Marshal(mutatePayload)
	mutateReq, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/students/%s/mutate-dormitory", created.Data.ID), bytes.NewBuffer(mutateBody))
	mutateReq.Header.Set("Content-Type", "application/json")
	mutateReq.Header.Set("Authorization", "Bearer "+token)

	mutateW := httptest.NewRecorder()
	router.ServeHTTP(mutateW, mutateReq)
	assert.Equal(t, http.StatusOK, mutateW.Code)

	var mutateResp studentResp
	require.NoError(t, json.Unmarshal(mutateW.Body.Bytes(), &mutateResp))
	require.Greater(t, len(mutateResp.Data.DormitoryHistory), 0)
	assert.Equal(t, dorm.ID.String(), mutateResp.Data.DormitoryHistory[0].DormitoryID)
}
