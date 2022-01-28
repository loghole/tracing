// DO NOT EDIT.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	attribute "go.opentelemetry.io/otel/attribute"
	instrumentation "go.opentelemetry.io/otel/sdk/instrumentation"
	resource "go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	trace0 "go.opentelemetry.io/otel/trace"
)

// MockReadOnlySpan is a mock of ReadOnlySpan interface.
type MockReadOnlySpan struct {
	tracesdk.ReadWriteSpan // DO NOT EDIT.

	ctrl     *gomock.Controller
	recorder *MockReadOnlySpanMockRecorder
}

// MockReadOnlySpanMockRecorder is the mock recorder for MockReadOnlySpan.
type MockReadOnlySpanMockRecorder struct {
	tracesdk.ReadWriteSpan // DO NOT EDIT.

	mock *MockReadOnlySpan
}

// NewMockReadOnlySpan creates a new mock instance.
func NewMockReadOnlySpan(ctrl *gomock.Controller) *MockReadOnlySpan {
	mock := &MockReadOnlySpan{ctrl: ctrl}
	mock.recorder = &MockReadOnlySpanMockRecorder{mock: mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadOnlySpan) EXPECT() *MockReadOnlySpanMockRecorder {
	return m.recorder
}

// Attributes mocks base method.
func (m *MockReadOnlySpan) Attributes() []attribute.KeyValue {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].([]attribute.KeyValue)
	return ret0
}

// Attributes indicates an expected call of Attributes.
func (mr *MockReadOnlySpanMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockReadOnlySpan)(nil).Attributes))
}

// ChildSpanCount mocks base method.
func (m *MockReadOnlySpan) ChildSpanCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChildSpanCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// ChildSpanCount indicates an expected call of ChildSpanCount.
func (mr *MockReadOnlySpanMockRecorder) ChildSpanCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChildSpanCount", reflect.TypeOf((*MockReadOnlySpan)(nil).ChildSpanCount))
}

// DroppedAttributes mocks base method.
func (m *MockReadOnlySpan) DroppedAttributes() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedAttributes")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedAttributes indicates an expected call of DroppedAttributes.
func (mr *MockReadOnlySpanMockRecorder) DroppedAttributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedAttributes", reflect.TypeOf((*MockReadOnlySpan)(nil).DroppedAttributes))
}

// DroppedEvents mocks base method.
func (m *MockReadOnlySpan) DroppedEvents() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedEvents")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedEvents indicates an expected call of DroppedEvents.
func (mr *MockReadOnlySpanMockRecorder) DroppedEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedEvents", reflect.TypeOf((*MockReadOnlySpan)(nil).DroppedEvents))
}

// DroppedLinks mocks base method.
func (m *MockReadOnlySpan) DroppedLinks() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedLinks")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedLinks indicates an expected call of DroppedLinks.
func (mr *MockReadOnlySpanMockRecorder) DroppedLinks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedLinks", reflect.TypeOf((*MockReadOnlySpan)(nil).DroppedLinks))
}

// EndTime mocks base method.
func (m *MockReadOnlySpan) EndTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EndTime indicates an expected call of EndTime.
func (mr *MockReadOnlySpanMockRecorder) EndTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTime", reflect.TypeOf((*MockReadOnlySpan)(nil).EndTime))
}

// Events mocks base method.
func (m *MockReadOnlySpan) Events() []trace.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Events")
	ret0, _ := ret[0].([]trace.Event)
	return ret0
}

// Events indicates an expected call of Events.
func (mr *MockReadOnlySpanMockRecorder) Events() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Events", reflect.TypeOf((*MockReadOnlySpan)(nil).Events))
}

// InstrumentationLibrary mocks base method.
func (m *MockReadOnlySpan) InstrumentationLibrary() instrumentation.Library {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstrumentationLibrary")
	ret0, _ := ret[0].(instrumentation.Library)
	return ret0
}

// InstrumentationLibrary indicates an expected call of InstrumentationLibrary.
func (mr *MockReadOnlySpanMockRecorder) InstrumentationLibrary() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstrumentationLibrary", reflect.TypeOf((*MockReadOnlySpan)(nil).InstrumentationLibrary))
}

// Links mocks base method.
func (m *MockReadOnlySpan) Links() []trace.Link {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Links")
	ret0, _ := ret[0].([]trace.Link)
	return ret0
}

// Links indicates an expected call of Links.
func (mr *MockReadOnlySpanMockRecorder) Links() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Links", reflect.TypeOf((*MockReadOnlySpan)(nil).Links))
}

// Name mocks base method.
func (m *MockReadOnlySpan) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockReadOnlySpanMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockReadOnlySpan)(nil).Name))
}

// Parent mocks base method.
func (m *MockReadOnlySpan) Parent() trace0.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parent")
	ret0, _ := ret[0].(trace0.SpanContext)
	return ret0
}

// Parent indicates an expected call of Parent.
func (mr *MockReadOnlySpanMockRecorder) Parent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parent", reflect.TypeOf((*MockReadOnlySpan)(nil).Parent))
}

// Resource mocks base method.
func (m *MockReadOnlySpan) Resource() *resource.Resource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resource")
	ret0, _ := ret[0].(*resource.Resource)
	return ret0
}

// Resource indicates an expected call of Resource.
func (mr *MockReadOnlySpanMockRecorder) Resource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resource", reflect.TypeOf((*MockReadOnlySpan)(nil).Resource))
}

// SpanContext mocks base method.
func (m *MockReadOnlySpan) SpanContext() trace0.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpanContext")
	ret0, _ := ret[0].(trace0.SpanContext)
	return ret0
}

// SpanContext indicates an expected call of SpanContext.
func (mr *MockReadOnlySpanMockRecorder) SpanContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpanContext", reflect.TypeOf((*MockReadOnlySpan)(nil).SpanContext))
}

// SpanKind mocks base method.
func (m *MockReadOnlySpan) SpanKind() trace0.SpanKind {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpanKind")
	ret0, _ := ret[0].(trace0.SpanKind)
	return ret0
}

// SpanKind indicates an expected call of SpanKind.
func (mr *MockReadOnlySpanMockRecorder) SpanKind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpanKind", reflect.TypeOf((*MockReadOnlySpan)(nil).SpanKind))
}

// StartTime mocks base method.
func (m *MockReadOnlySpan) StartTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// StartTime indicates an expected call of StartTime.
func (mr *MockReadOnlySpanMockRecorder) StartTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTime", reflect.TypeOf((*MockReadOnlySpan)(nil).StartTime))
}

// Status mocks base method.
func (m *MockReadOnlySpan) Status() trace.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(trace.Status)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockReadOnlySpanMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockReadOnlySpan)(nil).Status))
}

// private mocks base method.
func (m *MockReadOnlySpan) private() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "private")
}

// private indicates an expected call of private.
func (mr *MockReadOnlySpanMockRecorder) private() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "private", reflect.TypeOf((*MockReadOnlySpan)(nil).private))
}
