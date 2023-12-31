// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/unbeman/av-prac-task/internal/database (interfaces: IDatabase)

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "github.com/unbeman/av-prac-task/internal/model"
)

// MockIDatabase is a mock of IDatabase interface.
type MockIDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockIDatabaseMockRecorder
}

// MockIDatabaseMockRecorder is the mock recorder for MockIDatabase.
type MockIDatabaseMockRecorder struct {
	mock *MockIDatabase
}

// NewMockIDatabase creates a new mock instance.
func NewMockIDatabase(ctrl *gomock.Controller) *MockIDatabase {
	mock := &MockIDatabase{ctrl: ctrl}
	mock.recorder = &MockIDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDatabase) EXPECT() *MockIDatabaseMockRecorder {
	return m.recorder
}

// AddSegmentToRandomUsers mocks base method.
func (m *MockIDatabase) AddSegmentToRandomUsers(arg0 context.Context, arg1 *model.Segment, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSegmentToRandomUsers", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSegmentToRandomUsers indicates an expected call of AddSegmentToRandomUsers.
func (mr *MockIDatabaseMockRecorder) AddSegmentToRandomUsers(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSegmentToRandomUsers", reflect.TypeOf((*MockIDatabase)(nil).AddSegmentToRandomUsers), arg0, arg1, arg2)
}

// CreateDeleteUserSegments mocks base method.
func (m *MockIDatabase) CreateDeleteUserSegments(arg0 context.Context, arg1 *model.User, arg2, arg3 []model.Slug) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeleteUserSegments", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDeleteUserSegments indicates an expected call of CreateDeleteUserSegments.
func (mr *MockIDatabaseMockRecorder) CreateDeleteUserSegments(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeleteUserSegments", reflect.TypeOf((*MockIDatabase)(nil).CreateDeleteUserSegments), arg0, arg1, arg2, arg3)
}

// CreateSegment mocks base method.
func (m *MockIDatabase) CreateSegment(arg0 context.Context, arg1 *model.Segment) (*model.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSegment", arg0, arg1)
	ret0, _ := ret[0].(*model.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSegment indicates an expected call of CreateSegment.
func (mr *MockIDatabaseMockRecorder) CreateSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSegment", reflect.TypeOf((*MockIDatabase)(nil).CreateSegment), arg0, arg1)
}

// DeleteSegment mocks base method.
func (m *MockIDatabase) DeleteSegment(arg0 context.Context, arg1 *model.Segment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSegment indicates an expected call of DeleteSegment.
func (mr *MockIDatabaseMockRecorder) DeleteSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegment", reflect.TypeOf((*MockIDatabase)(nil).DeleteSegment), arg0, arg1)
}

// GetSegment mocks base method.
func (m *MockIDatabase) GetSegment(arg0 context.Context, arg1 *model.Segment) (*model.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegment", arg0, arg1)
	ret0, _ := ret[0].(*model.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegment indicates an expected call of GetSegment.
func (mr *MockIDatabaseMockRecorder) GetSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegment", reflect.TypeOf((*MockIDatabase)(nil).GetSegment), arg0, arg1)
}

// GetSegments mocks base method.
func (m *MockIDatabase) GetSegments(arg0 context.Context, arg1 []model.Slug) ([]*model.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegments", arg0, arg1)
	ret0, _ := ret[0].([]*model.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegments indicates an expected call of GetSegments.
func (mr *MockIDatabaseMockRecorder) GetSegments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegments", reflect.TypeOf((*MockIDatabase)(nil).GetSegments), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockIDatabase) GetUser(arg0 context.Context, arg1 *model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockIDatabaseMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockIDatabase)(nil).GetUser), arg0, arg1)
}

// GetUserSegmentsHistory mocks base method.
func (m *MockIDatabase) GetUserSegmentsHistory(arg0 context.Context, arg1 *model.User, arg2, arg3 time.Time) ([]model.UserSegment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSegmentsHistory", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]model.UserSegment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSegmentsHistory indicates an expected call of GetUserSegmentsHistory.
func (mr *MockIDatabaseMockRecorder) GetUserSegmentsHistory(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSegmentsHistory", reflect.TypeOf((*MockIDatabase)(nil).GetUserSegmentsHistory), arg0, arg1, arg2, arg3)
}

// GetUserWithActiveSegments mocks base method.
func (m *MockIDatabase) GetUserWithActiveSegments(arg0 context.Context, arg1 *model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithActiveSegments", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithActiveSegments indicates an expected call of GetUserWithActiveSegments.
func (mr *MockIDatabaseMockRecorder) GetUserWithActiveSegments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithActiveSegments", reflect.TypeOf((*MockIDatabase)(nil).GetUserWithActiveSegments), arg0, arg1)
}
