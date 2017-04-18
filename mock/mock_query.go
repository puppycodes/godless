// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/johnny-morrice/godless (interfaces: QueryVisitor)

package mock_godless

import (
	gomock "github.com/golang/mock/gomock"
	godless "github.com/johnny-morrice/godless"
)

// Mock of QueryVisitor interface
type MockQueryVisitor struct {
	ctrl     *gomock.Controller
	recorder *_MockQueryVisitorRecorder
}

// Recorder for MockQueryVisitor (not exported)
type _MockQueryVisitorRecorder struct {
	mock *MockQueryVisitor
}

func NewMockQueryVisitor(ctrl *gomock.Controller) *MockQueryVisitor {
	mock := &MockQueryVisitor{ctrl: ctrl}
	mock.recorder = &_MockQueryVisitorRecorder{mock}
	return mock
}

func (_m *MockQueryVisitor) EXPECT() *_MockQueryVisitorRecorder {
	return _m.recorder
}

func (_m *MockQueryVisitor) LeaveWhere(_param0 *godless.QueryWhere) {
	_m.ctrl.Call(_m, "LeaveWhere", _param0)
}

func (_mr *_MockQueryVisitorRecorder) LeaveWhere(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "LeaveWhere", arg0)
}

func (_m *MockQueryVisitor) VisitAST(_param0 *godless.QueryAST) {
	_m.ctrl.Call(_m, "VisitAST", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitAST(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitAST", arg0)
}

func (_m *MockQueryVisitor) VisitJoin(_param0 *godless.QueryJoin) {
	_m.ctrl.Call(_m, "VisitJoin", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitJoin(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitJoin", arg0)
}

func (_m *MockQueryVisitor) VisitOpCode(_param0 godless.QueryOpCode) {
	_m.ctrl.Call(_m, "VisitOpCode", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitOpCode(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitOpCode", arg0)
}

func (_m *MockQueryVisitor) VisitParser(_param0 *godless.QueryParser) {
	_m.ctrl.Call(_m, "VisitParser", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitParser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitParser", arg0)
}

func (_m *MockQueryVisitor) VisitPredicate(_param0 *godless.QueryPredicate) {
	_m.ctrl.Call(_m, "VisitPredicate", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitPredicate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitPredicate", arg0)
}

func (_m *MockQueryVisitor) VisitRowJoin(_param0 int, _param1 *godless.QueryRowJoin) {
	_m.ctrl.Call(_m, "VisitRowJoin", _param0, _param1)
}

func (_mr *_MockQueryVisitorRecorder) VisitRowJoin(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitRowJoin", arg0, arg1)
}

func (_m *MockQueryVisitor) VisitSelect(_param0 *godless.QuerySelect) {
	_m.ctrl.Call(_m, "VisitSelect", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitSelect(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitSelect", arg0)
}

func (_m *MockQueryVisitor) VisitTableKey(_param0 string) {
	_m.ctrl.Call(_m, "VisitTableKey", _param0)
}

func (_mr *_MockQueryVisitorRecorder) VisitTableKey(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitTableKey", arg0)
}

func (_m *MockQueryVisitor) VisitWhere(_param0 int, _param1 *godless.QueryWhere) {
	_m.ctrl.Call(_m, "VisitWhere", _param0, _param1)
}

func (_mr *_MockQueryVisitorRecorder) VisitWhere(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VisitWhere", arg0, arg1)
}