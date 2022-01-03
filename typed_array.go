package webgl

// It's almost copy&paste of https://github.com/hajimehoshi/ebiten/tree/main/internal/jsutil

import (
	"fmt"
	"runtime"
	"syscall/js"
	"unsafe"
)

// SliceHeader defined here, because Len & Cap fields type differs at go-wasm and tinygo-wasm
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

var (
	object       = js.Global().Get("Object")
	arrayBuffer  = js.Global().Get("ArrayBuffer")
	uint8Array   = js.Global().Get("Uint8Array")
	float32Array = js.Global().Get("Float32Array")

	temporaryArrayBuffer  = arrayBuffer.New(16)
	temporaryUint8Array   = uint8Array.New(temporaryArrayBuffer)
	temporaryFloat32Array = float32Array.New(temporaryArrayBuffer)

	temporaryArrayBufferByteLengthFunc  js.Value
	temporaryUint8ArrayByteLengthFunc   js.Value
	temporaryFloat32ArrayByteLengthFunc js.Value
)

func init() {
	temporaryArrayBufferByteLengthFunc = object.Call("getOwnPropertyDescriptor", arrayBuffer.Get("prototype"), "byteLength").Get("get").Call("bind", temporaryArrayBuffer)
	temporaryUint8ArrayByteLengthFunc = object.Call("getOwnPropertyDescriptor", object.Call("getPrototypeOf", uint8Array).Get("prototype"), "byteLength").Get("get").Call("bind", temporaryUint8Array)
	temporaryFloat32ArrayByteLengthFunc = object.Call("getOwnPropertyDescriptor", object.Call("getPrototypeOf", float32Array).Get("prototype"), "byteLength").Get("get").Call("bind", temporaryFloat32Array)
}

func temporaryArrayBufferByteLength() int {
	return temporaryArrayBufferByteLengthFunc.Invoke().Int()
}

func temporaryUint8ArrayByteLength() int {
	return temporaryUint8ArrayByteLengthFunc.Invoke().Int()
}

func temporaryFloat32ArrayByteLength() int {
	return temporaryFloat32ArrayByteLengthFunc.Invoke().Int()
}

func ensureTemporaryArrayBufferSize(byteLength int) {
	bufl := temporaryArrayBufferByteLength()
	if bufl < byteLength {
		for bufl < byteLength {
			bufl *= 2
		}
		temporaryArrayBuffer = arrayBuffer.New(bufl)
	}
	if temporaryUint8ArrayByteLength() < bufl {
		temporaryUint8Array = uint8Array.New(temporaryArrayBuffer)
	}
	if temporaryFloat32ArrayByteLength() < bufl {
		temporaryFloat32Array = float32Array.New(temporaryArrayBuffer)
	}
}

func copySliceHeaderToTemporaryArrayBuffer(h *SliceHeader, itemSize int) {
	h.Len *= itemSize
	h.Cap *= itemSize
	bs := *(*[]byte)(unsafe.Pointer(h))
	js.CopyBytesToJS(temporaryUint8Array, bs)
}

func copyUint8SliceToTemporaryArrayBuffer(src []uint8) {
	if len(src) == 0 {
		return
	}
	js.CopyBytesToJS(temporaryUint8Array, src)
}

func copyFloat32SliceToTemporaryArrayBuffer(src []float32) {
	if len(src) == 0 {
		return
	}
	h := (*SliceHeader)(unsafe.Pointer(&src))
	copySliceHeaderToTemporaryArrayBuffer(h, 4)
	runtime.KeepAlive(src)
}

func copyUint32SliceToTemporaryArrayBuffer(src []uint32) {
	if len(src) == 0 {
		return
	}
	h := (*SliceHeader)(unsafe.Pointer(&src))
	copySliceHeaderToTemporaryArrayBuffer(h, 4)
	runtime.KeepAlive(src)
}

func copyUint16SliceToTemporaryArrayBuffer(src []uint16) {
	if len(src) == 0 {
		return
	}
	h := (*SliceHeader)(unsafe.Pointer(&src))
	h.Len *= 2
	h.Cap *= 2
	bs := *(*[]byte)(unsafe.Pointer(h))
	runtime.KeepAlive(src)
	js.CopyBytesToJS(temporaryUint8Array, bs)
}

func getTemporaryFloat32Array(data []float32) js.Value {
	ensureTemporaryArrayBufferSize(len(data) * 4)
	copyFloat32SliceToTemporaryArrayBuffer(data)
	return temporaryFloat32Array
}

func getTemporaryUint8ArrayFromUint32Slice(data []uint32) js.Value {
	ensureTemporaryArrayBufferSize(len(data) * 4)
	copyUint32SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array
}

func getTemporaryUint8ArrayFromUint16Slice(data []uint16) js.Value {
	ensureTemporaryArrayBufferSize(len(data) * 2)
	copyUint16SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array
}

func getTemporaryUint8Array(data []uint8) js.Value {
	ensureTemporaryArrayBufferSize(len(data))
	copyUint8SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array
}

func jsTypedArrayOf(srcData interface{}) js.Value {
	switch s := srcData.(type) {
	case []float32:
		return getTemporaryFloat32Array(s)
	case []uint32:
		return getTemporaryUint8ArrayFromUint32Slice(s)
	case []uint16:
		return getTemporaryUint8ArrayFromUint16Slice(s)
	case []uint8:
		return getTemporaryUint8Array(s)
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at sliceToBytesSlice: %T %v", s, s))
	}
}
