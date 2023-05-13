// Code generated by MockGen. DO NOT EDIT.
// Source: ./handler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	context "context"
	reflect "reflect"
	comment "task-manager-api/internal/comment"
	profile "task-manager-api/internal/profile"
	taskmanager "task-manager-api/internal/taskmanager"

	gomock "github.com/golang/mock/gomock"
)

// MockITasks is a mock of ITasks interface.
type MockITasks struct {
	ctrl     *gomock.Controller
	recorder *MockITasksMockRecorder
}

// MockITasksMockRecorder is the mock recorder for MockITasks.
type MockITasksMockRecorder struct {
	mock *MockITasks
}

// NewMockITasks creates a new mock instance.
func NewMockITasks(ctrl *gomock.Controller) *MockITasks {
	mock := &MockITasks{ctrl: ctrl}
	mock.recorder = &MockITasksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITasks) EXPECT() *MockITasksMockRecorder {
	return m.recorder
}

// ArchiveTask mocks base method.
func (m *MockITasks) ArchiveTask(ctx context.Context, ownerId, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveTask", ctx, ownerId, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveTask indicates an expected call of ArchiveTask.
func (mr *MockITasksMockRecorder) ArchiveTask(ctx, ownerId, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveTask", reflect.TypeOf((*MockITasks)(nil).ArchiveTask), ctx, ownerId, id)
}

// CreateTask mocks base method.
func (m *MockITasks) CreateTask(ctx context.Context, ownerId, topic, desc string) (*taskmanager.TaskDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, ownerId, topic, desc)
	ret0, _ := ret[0].(*taskmanager.TaskDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockITasksMockRecorder) CreateTask(ctx, ownerId, topic, desc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockITasks)(nil).CreateTask), ctx, ownerId, topic, desc)
}

// GetAllTask mocks base method.
func (m *MockITasks) GetAllTask(ctx context.Context, page, limit int) ([]taskmanager.TaskDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTask", ctx, page, limit)
	ret0, _ := ret[0].([]taskmanager.TaskDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTask indicates an expected call of GetAllTask.
func (mr *MockITasksMockRecorder) GetAllTask(ctx, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTask", reflect.TypeOf((*MockITasks)(nil).GetAllTask), ctx, page, limit)
}

// GetTask mocks base method.
func (m *MockITasks) GetTask(ctx context.Context, id string) (*taskmanager.TaskDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", ctx, id)
	ret0, _ := ret[0].(*taskmanager.TaskDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockITasksMockRecorder) GetTask(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockITasks)(nil).GetTask), ctx, id)
}

// UpdateTaskStatus mocks base method.
func (m *MockITasks) UpdateTaskStatus(ctx context.Context, ownerId, id string, status int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskStatus", ctx, ownerId, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskStatus indicates an expected call of UpdateTaskStatus.
func (mr *MockITasksMockRecorder) UpdateTaskStatus(ctx, ownerId, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskStatus", reflect.TypeOf((*MockITasks)(nil).UpdateTaskStatus), ctx, ownerId, id, status)
}

// MockIComments is a mock of IComments interface.
type MockIComments struct {
	ctrl     *gomock.Controller
	recorder *MockICommentsMockRecorder
}

// MockICommentsMockRecorder is the mock recorder for MockIComments.
type MockICommentsMockRecorder struct {
	mock *MockIComments
}

// NewMockIComments creates a new mock instance.
func NewMockIComments(ctrl *gomock.Controller) *MockIComments {
	mock := &MockIComments{ctrl: ctrl}
	mock.recorder = &MockICommentsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIComments) EXPECT() *MockICommentsMockRecorder {
	return m.recorder
}

// CreateComment mocks base method.
func (m *MockIComments) CreateComment(ctx context.Context, ownerId, topicId, content string) (*comment.CommentDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", ctx, ownerId, topicId, content)
	ret0, _ := ret[0].(*comment.CommentDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockICommentsMockRecorder) CreateComment(ctx, ownerId, topicId, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockIComments)(nil).CreateComment), ctx, ownerId, topicId, content)
}

// GetTopicComments mocks base method.
func (m *MockIComments) GetTopicComments(ctx context.Context, topicId string, page, limit int) ([]comment.CommentDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopicComments", ctx, topicId, page, limit)
	ret0, _ := ret[0].([]comment.CommentDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopicComments indicates an expected call of GetTopicComments.
func (mr *MockICommentsMockRecorder) GetTopicComments(ctx, topicId, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopicComments", reflect.TypeOf((*MockIComments)(nil).GetTopicComments), ctx, topicId, page, limit)
}

// MockIProfile is a mock of IProfile interface.
type MockIProfile struct {
	ctrl     *gomock.Controller
	recorder *MockIProfileMockRecorder
}

// MockIProfileMockRecorder is the mock recorder for MockIProfile.
type MockIProfileMockRecorder struct {
	mock *MockIProfile
}

// NewMockIProfile creates a new mock instance.
func NewMockIProfile(ctrl *gomock.Controller) *MockIProfile {
	mock := &MockIProfile{ctrl: ctrl}
	mock.recorder = &MockIProfileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProfile) EXPECT() *MockIProfileMockRecorder {
	return m.recorder
}

// GetProfile mocks base method.
func (m *MockIProfile) GetProfile(ctx context.Context, ownerId string) (*profile.ProfileDoc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", ctx, ownerId)
	ret0, _ := ret[0].(*profile.ProfileDoc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockIProfileMockRecorder) GetProfile(ctx, ownerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockIProfile)(nil).GetProfile), ctx, ownerId)
}
