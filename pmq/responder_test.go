package pmq_test

import (
	"fmt"
	"github.com/nidhhoggr/xipc/pmq"
	"reflect"
	"syscall"
	"testing"
	"time"
)

const wired = "Narwhals and ice cream"

func TestRecvErrwithNonblocking(t *testing.T) {

	config := pmq.QueueConfig{
		Name:  "pmq_testing_recverrwnblk",
		Flags: pmq.O_RDWR | pmq.O_CREAT | pmq.O_NONBLOCK,
	}
	mqr := pmq.NewResponder(&config, nil)
	assert(t, !mqr.HasErrors())
	assertNotNil(t, mqr)

	msg, err := mqr.Read()

	assertNil(t, msg)
	assertNotNil(t, err)
	assertEqual(t, syscall.EAGAIN, err.(syscall.Errno))
	err = mqr.CloseResponder()
	assertNil(t, err)
}

func TestRecvwithNonblocking(t *testing.T) {

	config := pmq.QueueConfig{
		Name:  "pmq_testing_recvwnblk",
		Flags: pmq.O_RDWR | pmq.O_CREAT | pmq.O_NONBLOCK,
	}
	mqr := pmq.NewResponder(&config, nil)
	assert(t, !mqr.HasErrors())
	assertNotNil(t, mqr)

	err := mqr.HandleRequest(func(request []byte) (processed []byte, err error) {
		return []byte(fmt.Sprintf("I recieved request: %s\n", request)), nil
	})
	assertNil(t, err)
	err = mqr.CloseResponder()
	assertNil(t, err)
}

func TestRecvErrwithBlocking(t *testing.T) {

	config := pmq.QueueConfig{
		Name:  "pmq_testing_recverrwblk",
		Flags: pmq.O_RDWR | pmq.O_CREAT,
	}
	mqr := pmq.NewResponder(&config, nil)
	assert(t, !mqr.HasErrors())
	assertNotNil(t, mqr)

	msg, err := mqr.ReadTimed(time.Second)

	assertNil(t, msg)
	assertNotNil(t, err)
	assertEqual(t, syscall.ETIMEDOUT, err.(syscall.Errno))
	err = mqr.CloseResponder()
	assertNil(t, err)
}

func assert(t *testing.T, i bool) {
	if !i {
		t.Errorf("expected %-v to be true", i)
	}
}

func assertNil(t *testing.T, i interface{}) {
	if !isNil(i) {
		t.Errorf("expected %-v to be nil", i)
	}
}

func assertNotNil(t *testing.T, i interface{}) {
	if isNil(i) {
		t.Errorf("expected %-v to not be nil", i)
	}
}

func assertEqual[T any](t *testing.T, ptr T, ptr2 T) {
	if !reflect.ValueOf(ptr).Equal(reflect.ValueOf(ptr2)) {
		t.Errorf("expected %-v to equal %-v", ptr, ptr2)
	}
}

// containsKind checks if a specified kind in the slice of kinds.
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

// isNil checks if a specified object is nil or not, without Failing.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}
