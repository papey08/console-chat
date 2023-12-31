// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import model "console-chat/internal/model"

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// RegisterUser provides a mock function with given fields: ctx, nickname, password
func (_m *App) RegisterUser(ctx context.Context, nickname string, password string) (model.User, error) {
	ret := _m.Called(ctx, nickname, password)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) model.User); ok {
		r0 = rf(ctx, nickname, password)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, nickname, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignInUser provides a mock function with given fields: ctx, nickname, password
func (_m *App) SignInUser(ctx context.Context, nickname string, password string) (model.User, error) {
	ret := _m.Called(ctx, nickname, password)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) model.User); ok {
		r0 = rf(ctx, nickname, password)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, nickname, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
