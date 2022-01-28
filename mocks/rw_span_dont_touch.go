// DO NOT EDIT.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	attribute "go.opentelemetry.io/otel/attribute"
	codes "go.opentelemetry.io/otel/codes"
	instrumentation "go.opentelemetry.io/otel/sdk/instrumentation"
	resource "go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	trace0 "go.opentelemetry.io/otel/trace"
)

// MockReadWriteSpan is a mock of ReadWriteSpan interface.
type MockReadWriteSpan struct {
	tracesdk.ReadWriteSpan // DO NOT EDIT.

	ctrl     *gomock.Controller
	recorder *MockReadWriteSpanMockRecorder
}

// MockReadWriteSpanMockRecorder is the mock recorder for MockReadWriteSpan.
type MockReadWriteSpanMockRecorder struct {
	tracesdk.ReadWriteSpan // DO NOT EDIT.

	mock *MockReadWriteSpan
}

// NewMockReadWriteSpan creates a new mock instance.
func NewMockReadWriteSpan(ctrl *gomock.Controller) *MockReadWriteSpan {
	mock := &MockReadWriteSpan{ctrl: ctrl}
	mock.recorder = &MockReadWriteSpanMockRecorder{mock: mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadWriteSpan) EXPECT() *MockReadWriteSpanMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockReadWriteSpan) AddEvent(arg0 string, arg1 ...trace0.EventOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AddEvent", varargs...)
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockReadWriteSpanMockRecorder) AddEvent(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockReadWriteSpan)(nil).AddEvent), varargs...)
}

// Attributes mocks base method.
func (m *MockReadWriteSpan) Attributes() []attribute.KeyValue {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].([]attribute.KeyValue)
	return ret0
}

// Attributes indicates an expected call of Attributes.
func (mr *MockReadWriteSpanMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockReadWriteSpan)(nil).Attributes))
}

// ChildSpanCount mocks base method.
func (m *MockReadWriteSpan) ChildSpanCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChildSpanCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// ChildSpanCount indicates an expected call of ChildSpanCount.
func (mr *MockReadWriteSpanMockRecorder) ChildSpanCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChildSpanCount", reflect.TypeOf((*MockReadWriteSpan)(nil).ChildSpanCount))
}

// DroppedAttributes mocks base method.
func (m *MockReadWriteSpan) DroppedAttributes() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedAttributes")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedAttributes indicates an expected call of DroppedAttributes.
func (mr *MockReadWriteSpanMockRecorder) DroppedAttributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedAttributes", reflect.TypeOf((*MockReadWriteSpan)(nil).DroppedAttributes))
}

// DroppedEvents mocks base method.
func (m *MockReadWriteSpan) DroppedEvents() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedEvents")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedEvents indicates an expected call of DroppedEvents.
func (mr *MockReadWriteSpanMockRecorder) DroppedEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedEvents", reflect.TypeOf((*MockReadWriteSpan)(nil).DroppedEvents))
}

// DroppedLinks mocks base method.
func (m *MockReadWriteSpan) DroppedLinks() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DroppedLinks")
	ret0, _ := ret[0].(int)
	return ret0
}

// DroppedLinks indicates an expected call of DroppedLinks.
func (mr *MockReadWriteSpanMockRecorder) DroppedLinks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DroppedLinks", reflect.TypeOf((*MockReadWriteSpan)(nil).DroppedLinks))
}

// End mocks base method.
func (m *MockReadWriteSpan) End(arg0 ...trace0.SpanEndOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "End", varargs...)
}

// End indicates an expected call of End.
func (mr *MockReadWriteSpanMockRecorder) End(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "End", reflect.TypeOf((*MockReadWriteSpan)(nil).End), arg0...)
}

// EndTime mocks base method.
func (m *MockReadWriteSpan) EndTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EndTime indicates an expected call of EndTime.
func (mr *MockReadWriteSpanMockRecorder) EndTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTime", reflect.TypeOf((*MockReadWriteSpan)(nil).EndTime))
}

// Events mocks base method.
func (m *MockReadWriteSpan) Events() []trace.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Events")
	ret0, _ := ret[0].([]trace.Event)
	return ret0
}

// Events indicates an expected call of Events.
func (mr *MockReadWriteSpanMockRecorder) Events() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Events", reflect.TypeOf((*MockReadWriteSpan)(nil).Events))
}

// InstrumentationLibrary mocks base method.
func (m *MockReadWriteSpan) InstrumentationLibrary() instrumentation.Library {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstrumentationLibrary")
	ret0, _ := ret[0].(instrumentation.Library)
	return ret0
}

// InstrumentationLibrary indicates an expected call of InstrumentationLibrary.
func (mr *MockReadWriteSpanMockRecorder) InstrumentationLibrary() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstrumentationLibrary", reflect.TypeOf((*MockReadWriteSpan)(nil).InstrumentationLibrary))
}

// IsRecording mocks base method.
func (m *MockReadWriteSpan) IsRecording() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRecording")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRecording indicates an expected call of IsRecording.
func (mr *MockReadWriteSpanMockRecorder) IsRecording() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRecording", reflect.TypeOf((*MockReadWriteSpan)(nil).IsRecording))
}

// Links mocks base method.
func (m *MockReadWriteSpan) Links() []trace.Link {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Links")
	ret0, _ := ret[0].([]trace.Link)
	return ret0
}

// Links indicates an expected call of Links.
func (mr *MockReadWriteSpanMockRecorder) Links() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Links", reflect.TypeOf((*MockReadWriteSpan)(nil).Links))
}

// Name mocks base method.
func (m *MockReadWriteSpan) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockReadWriteSpanMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockReadWriteSpan)(nil).Name))
}

// Parent mocks base method.
func (m *MockReadWriteSpan) Parent() trace0.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parent")
	ret0, _ := ret[0].(trace0.SpanContext)
	return ret0
}

// Parent indicates an expected call of Parent.
func (mr *MockReadWriteSpanMockRecorder) Parent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parent", reflect.TypeOf((*MockReadWriteSpan)(nil).Parent))
}

// RecordError mocks base method.
func (m *MockReadWriteSpan) RecordError(arg0 error, arg1 ...trace0.EventOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "RecordError", varargs...)
}

// RecordError indicates an expected call of RecordError.
func (mr *MockReadWriteSpanMockRecorder) RecordError(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordError", reflect.TypeOf((*MockReadWriteSpan)(nil).RecordError), varargs...)
}

// Resource mocks base method.
func (m *MockReadWriteSpan) Resource() *resource.Resource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resource")
	ret0, _ := ret[0].(*resource.Resource)
	return ret0
}

// Resource indicates an expected call of Resource.
func (mr *MockReadWriteSpanMockRecorder) Resource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resource", reflect.TypeOf((*MockReadWriteSpan)(nil).Resource))
}

// SetAttributes mocks base method.
func (m *MockReadWriteSpan) SetAttributes(arg0 ...attribute.KeyValue) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetAttributes", varargs...)
}

// SetAttributes indicates an expected call of SetAttributes.
func (mr *MockReadWriteSpanMockRecorder) SetAttributes(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAttributes", reflect.TypeOf((*MockReadWriteSpan)(nil).SetAttributes), arg0...)
}

// SetName mocks base method.
func (m *MockReadWriteSpan) SetName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetName", arg0)
}

// SetName indicates an expected call of SetName.
func (mr *MockReadWriteSpanMockRecorder) SetName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetName", reflect.TypeOf((*MockReadWriteSpan)(nil).SetName), arg0)
}

// SetStatus mocks base method.
func (m *MockReadWriteSpan) SetStatus(arg0 codes.Code, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStatus", arg0, arg1)
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockReadWriteSpanMockRecorder) SetStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockReadWriteSpan)(nil).SetStatus), arg0, arg1)
}

// SpanContext mocks base method.
func (m *MockReadWriteSpan) SpanContext() trace0.SpanContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpanContext")
	ret0, _ := ret[0].(trace0.SpanContext)
	return ret0
}

// SpanContext indicates an expected call of SpanContext.
func (mr *MockReadWriteSpanMockRecorder) SpanContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpanContext", reflect.TypeOf((*MockReadWriteSpan)(nil).SpanContext))
}

// SpanKind mocks base method.
func (m *MockReadWriteSpan) SpanKind() trace0.SpanKind {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpanKind")
	ret0, _ := ret[0].(trace0.SpanKind)
	return ret0
}

// SpanKind indicates an expected call of SpanKind.
func (mr *MockReadWriteSpanMockRecorder) SpanKind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpanKind", reflect.TypeOf((*MockReadWriteSpan)(nil).SpanKind))
}

// StartTime mocks base method.
func (m *MockReadWriteSpan) StartTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// StartTime indicates an expected call of StartTime.
func (mr *MockReadWriteSpanMockRecorder) StartTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTime", reflect.TypeOf((*MockReadWriteSpan)(nil).StartTime))
}

// Status mocks base method.
func (m *MockReadWriteSpan) Status() trace.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(trace.Status)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockReadWriteSpanMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockReadWriteSpan)(nil).Status))
}

// TracerProvider mocks base method.
func (m *MockReadWriteSpan) TracerProvider() trace0.TracerProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TracerProvider")
	ret0, _ := ret[0].(trace0.TracerProvider)
	return ret0
}

// TracerProvider indicates an expected call of TracerProvider.
func (mr *MockReadWriteSpanMockRecorder) TracerProvider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TracerProvider", reflect.TypeOf((*MockReadWriteSpan)(nil).TracerProvider))
}

// private mocks base method.
func (m *MockReadWriteSpan) private() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "private")
}

// private indicates an expected call of private.
func (mr *MockReadWriteSpanMockRecorder) private() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "private", reflect.TypeOf((*MockReadWriteSpan)(nil).private))
}
