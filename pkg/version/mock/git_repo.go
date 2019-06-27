// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/arnaud-deprez/gsemver/pkg/version (interfaces: GitRepo)

// Package mock_version is a generated GoMock package.
package mock_version

import (
	git "github.com/arnaud-deprez/gsemver/pkg/git"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGitRepo is a mock of GitRepo interface
type MockGitRepo struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepoMockRecorder
}

// MockGitRepoMockRecorder is the mock recorder for MockGitRepo
type MockGitRepoMockRecorder struct {
	mock *MockGitRepo
}

// NewMockGitRepo creates a new mock instance
func NewMockGitRepo(ctrl *gomock.Controller) *MockGitRepo {
	mock := &MockGitRepo{ctrl: ctrl}
	mock.recorder = &MockGitRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitRepo) EXPECT() *MockGitRepoMockRecorder {
	return m.recorder
}

// CountCommits mocks base method
func (m *MockGitRepo) CountCommits(arg0, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountCommits", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountCommits indicates an expected call of CountCommits
func (mr *MockGitRepoMockRecorder) CountCommits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountCommits", reflect.TypeOf((*MockGitRepo)(nil).CountCommits), arg0, arg1)
}

// GetCommits mocks base method
func (m *MockGitRepo) GetCommits(arg0, arg1 string) ([]git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommits", arg0, arg1)
	ret0, _ := ret[0].([]git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommits indicates an expected call of GetCommits
func (mr *MockGitRepoMockRecorder) GetCommits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommits", reflect.TypeOf((*MockGitRepo)(nil).GetCommits), arg0, arg1)
}

// GetCurrentBranch mocks base method
func (m *MockGitRepo) GetCurrentBranch() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentBranch")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentBranch indicates an expected call of GetCurrentBranch
func (mr *MockGitRepoMockRecorder) GetCurrentBranch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentBranch", reflect.TypeOf((*MockGitRepo)(nil).GetCurrentBranch))
}

// GetLastRelativeTag mocks base method
func (m *MockGitRepo) GetLastRelativeTag(arg0 string) (git.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastRelativeTag", arg0)
	ret0, _ := ret[0].(git.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastRelativeTag indicates an expected call of GetLastRelativeTag
func (mr *MockGitRepoMockRecorder) GetLastRelativeTag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastRelativeTag", reflect.TypeOf((*MockGitRepo)(nil).GetLastRelativeTag), arg0)
}
