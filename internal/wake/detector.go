package wake

/*
#cgo CFLAGS: -I${SRCDIR}/porcupine/include
#cgo LDFLAGS: -L${SRCDIR}/porcupine/lib/windows/amd64 -lporcupine
#include "porcupine.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Detector struct {
	handle C.pv_porcupine_t
}

func New(sensitivity float32) (*Detector, error) {
	var handle C.pv_porcupine_t
	status := C.pv_porcupine_init(
		C.CString("./porcupine/resources/keyword_windows.ppn"),
		C.float(sensitivity),
		&handle,
	)
	if status != 0 {
		return nil, fmt.Errorf("init failed")
	}
	return &Detector{handle: handle}, nil
}

func (d *Detector) Process(pcm []int16) bool {
	// Porcupine принимает int16*
	cPtr := unsafe.Pointer(&pcm[0])
	return C.pv_porcupine_process(d.handle, (*C.int16_t)(cPtr)) > 0
}

func (d *Detector) Delete() {
	C.pv_porcupine_delete(d.handle)
}
