package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/go-backend-starter/internal/application/dto"
	"github.com/your-org/go-backend-starter/internal/application/usecase/mocks"
	"github.com/your-org/go-backend-starter/internal/domain/entity"
	domainErrors "github.com/your-org/go-backend-starter/internal/domain/errors"
	"github.com/your-org/go-backend-starter/internal/domain/repository"
)

type classScheduleNoopAuditLogger struct{}

func (n *classScheduleNoopAuditLogger) Log(ctx context.Context, resource, action, targetID string, metadata map[string]string) error {
	return nil
}

func newClassScheduleUCForTest() (*ClassScheduleUseCase, *mocks.ClassScheduleRepositoryMock, *mocks.ClassRepositoryMock, *mocks.MockTeacherRepository, *mocks.SubjectRepositoryMock, *mocks.MockScheduleSlotRepository, *mocks.MockDormitoryRepository) {
	scheduleRepo := new(mocks.ClassScheduleRepositoryMock)
	classRepo := new(mocks.ClassRepositoryMock)
	teacherRepo := new(mocks.MockTeacherRepository)
	subjectRepo := new(mocks.SubjectRepositoryMock)
	slotRepo := new(mocks.MockScheduleSlotRepository)
	dormRepo := new(mocks.MockDormitoryRepository)
	uc := NewClassScheduleUseCase(scheduleRepo, classRepo, teacherRepo, subjectRepo, slotRepo, dormRepo, &classScheduleNoopAuditLogger{})
	return uc, scheduleRepo, classRepo, teacherRepo, subjectRepo, slotRepo, dormRepo
}

func TestClassScheduleUseCase_Create_WithSlot(t *testing.T) {
	uc, scheduleRepo, classRepo, teacherRepo, _, slotRepo, dormRepo := newClassScheduleUCForTest()
	classID := uuid.New()
	teacherID := uuid.New()
	dormID := uuid.New()
	slotID := uuid.New()

	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, IsActive: true}, nil)
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)
	slotRepo.On("GetByID", mock.Anything, slotID).Return(&entity.ScheduleSlot{ID: slotID, DormitoryID: dormID, IsActive: true, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)}, nil)
	scheduleRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *entity.ClassSchedule) bool {
		return s.ClassID == classID && s.TeacherID == teacherID && s.SlotID != nil && *s.SlotID == slotID
	})).Return(nil)

	req := dto.CreateClassScheduleRequest{
		ClassID:     classID.String(),
		DormitoryID: dormID.String(),
		TeacherID:   teacherID.String(),
		SlotID:      stringPtrCS(slotID.String()),
		DayOfWeek:   "mon",
	}

	resp, err := uc.CreateClassSchedule(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, classID.String(), resp.ClassID)
	scheduleRepo.AssertExpectations(t)
}

func TestClassScheduleUseCase_Create_ManualTime(t *testing.T) {
	uc, scheduleRepo, classRepo, teacherRepo, _, slotRepo, dormRepo := newClassScheduleUCForTest()
	classID := uuid.New()
	teacherID := uuid.New()
	dormID := uuid.New()

	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, IsActive: true}, nil)
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)
	scheduleRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *entity.ClassSchedule) bool {
		return s.SlotID == nil && s.StartTime != nil && s.EndTime != nil
	})).Return(nil)

	start := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
	end := time.Now().Add(3 * time.Hour).Format(time.RFC3339)
	req := dto.CreateClassScheduleRequest{
		ClassID:     classID.String(),
		DormitoryID: dormID.String(),
		TeacherID:   teacherID.String(),
		DayOfWeek:   "tue",
		StartTime:   &start,
		EndTime:     &end,
	}

	resp, err := uc.CreateClassSchedule(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "tue", resp.DayOfWeek)
	scheduleRepo.AssertExpectations(t)
	slotRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestClassScheduleUseCase_UpdateSchedule(t *testing.T) {
	uc, scheduleRepo, _, teacherRepo, subjectRepo, slotRepo, _ := newClassScheduleUCForTest()
	scheduleID := uuid.New()
	teacherID := uuid.New()
	subjectID := uuid.New()
	slotID := uuid.New()
	dormID := uuid.New()
	start := time.Now()
	end := time.Now().Add(time.Hour)

	scheduleRepo.On("GetByID", mock.Anything, scheduleID).Return(&entity.ClassSchedule{
		ID:          scheduleID,
		ClassID:     uuid.New(),
		DormitoryID: dormID,
		TeacherID:   uuid.New(),
		DayOfWeek:   "wed",
		StartTime:   &start,
		EndTime:     &end,
		IsActive:    true,
	}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(&entity.Teacher{ID: teacherID, IsActive: true}, nil)
	subjectRepo.On("GetByID", mock.Anything, subjectID).Return(&entity.Subject{ID: subjectID}, nil)
	slotRepo.On("GetByID", mock.Anything, slotID).Return(&entity.ScheduleSlot{ID: slotID, DormitoryID: dormID, IsActive: true, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)}, nil)
	scheduleRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	req := dto.UpdateClassScheduleRequest{
		TeacherID: stringPtrCS(teacherID.String()),
		SubjectID: stringPtrCS(subjectID.String()),
		SlotID:    stringPtrCS(slotID.String()),
		DayOfWeek: stringPtrCS("thu"),
		IsActive:  boolPtrCS(false),
	}

	resp, err := uc.UpdateClassSchedule(context.Background(), scheduleID, req)
	assert.NoError(t, err)
	assert.Equal(t, "thu", resp.DayOfWeek)
	assert.Equal(t, false, resp.IsActive)
}

func TestClassScheduleUseCase_ListSchedules(t *testing.T) {
	uc, scheduleRepo, _, _, _, _, _ := newClassScheduleUCForTest()
	scheduleRepo.On("List", mock.Anything, mock.MatchedBy(func(filter repository.ClassScheduleFilter) bool {
		return filter.Page == 0 && filter.PageSize == 0 && filter.DayOfWeek == "fri"
	})).Return([]*entity.ClassSchedule{{ID: uuid.New(), ClassID: uuid.New(), TeacherID: uuid.New(), DormitoryID: uuid.New(), DayOfWeek: "fri"}}, int64(1), nil)

	resp, err := uc.ListClassSchedules(context.Background(), "", "", "", "fri", 0, 0, nil)
	assert.NoError(t, err)
	assert.Len(t, resp.Schedules, 1)
}

func TestClassScheduleUseCase_DeleteSchedule(t *testing.T) {
	uc, scheduleRepo, _, _, _, _, _ := newClassScheduleUCForTest()
	scheduleID := uuid.New()
	scheduleRepo.On("GetByID", mock.Anything, scheduleID).Return(&entity.ClassSchedule{ID: scheduleID}, nil)
	scheduleRepo.On("Delete", mock.Anything, scheduleID).Return(nil)

	assert.NoError(t, uc.DeleteClassSchedule(context.Background(), scheduleID))
	scheduleRepo.AssertExpectations(t)
}

func TestClassScheduleUseCase_Create_InvalidTeacher(t *testing.T) {
	uc, _, classRepo, teacherRepo, _, _, dormRepo := newClassScheduleUCForTest()
	classID := uuid.New()
	teacherID := uuid.New()
	dormID := uuid.New()

	classRepo.On("GetByID", mock.Anything, classID).Return(&entity.Class{ID: classID}, nil)
	teacherRepo.On("GetByID", mock.Anything, teacherID).Return(nil, assert.AnError)
	dormRepo.On("GetByID", mock.Anything, dormID).Return(&entity.Dormitory{ID: dormID}, nil)

	_, err := uc.CreateClassSchedule(context.Background(), dto.CreateClassScheduleRequest{
		ClassID:     classID.String(),
		DormitoryID: dormID.String(),
		TeacherID:   teacherID.String(),
		DayOfWeek:   "mon",
		StartTime:   stringPtrCS(time.Now().Format(time.RFC3339)),
		EndTime:     stringPtrCS(time.Now().Add(time.Hour).Format(time.RFC3339)),
	})
	assert.ErrorIs(t, err, domainErrors.ErrTeacherNotFound)
}

func stringPtrCS(val string) *string {
	return &val
}

func boolPtrCS(val bool) *bool {
	return &val
}
