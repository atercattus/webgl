// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten_jsutil

import (
	"syscall/js"
)

var (
	object       = js.Global().Get("Object")
	arrayBuffer  = js.Global().Get("ArrayBuffer")
	uint8Array   = js.Global().Get("Uint8Array")
	float32Array = js.Global().Get("Float32Array")
)

var (
	// temporaryArrayBuffer is a temporary buffer used at gl.readPixels or gl.texSubImage2D.
	// The read data is converted to Go's byte slice as soon as possible.
	// To avoid often allocating ArrayBuffer, reuse the buffer whenever possible.
	temporaryArrayBuffer = arrayBuffer.New(16)

	// temporaryUint8Array is a Uint8ArrayBuffer whose underlying buffer is always temporaryArrayBuffer.
	temporaryUint8Array = uint8Array.New(temporaryArrayBuffer)

	// temporaryFloat32Array is a Float32ArrayBuffer whose underlying buffer is always temporaryArrayBuffer.
	temporaryFloat32Array = float32Array.New(temporaryArrayBuffer)
)

var (
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

// TemporaryUint8ArrayFromUint8Slice returns a Uint8Array.
// data must be a slice of a numeric type for initialization, or nil if you don't need initialization.
func TemporaryUint8ArrayFromUint8Slice(minLength int, data []uint8) js.Value {
	ensureTemporaryArrayBufferSize(minLength)
	copyUint8SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array.Call("subarray", 0, minLength)
}

// TemporaryUint8ArrayFromUint16Slice returns a Uint8Array.
// data must be a slice of a numeric type for initialization, or nil if you don't need initialization.
func TemporaryUint8ArrayFromUint16Slice(minLength int, data []uint16) js.Value {
	ensureTemporaryArrayBufferSize(minLength * 2)
	copyUint16SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array.Call("subarray", 0, minLength*2)
}

// TemporaryUint8ArrayFromFloat32Slice returns a Uint8Array.
// data must be a slice of a numeric type for initialization, or nil if you don't need initialization.
func TemporaryUint8ArrayFromFloat32Slice(minLength int, data []float32) js.Value {
	ensureTemporaryArrayBufferSize(minLength * 4)
	copyFloat32SliceToTemporaryArrayBuffer(data)
	return temporaryUint8Array.Call("subarray", 0, minLength*4)
}

// TemporaryFloat32Array returns a Float32Array.
// data must be a slice of a numeric type for initialization, or nil if you don't need initialization.
func TemporaryFloat32Array(minLength int, data []float32) js.Value {
	ensureTemporaryArrayBufferSize(minLength * 4)
	copyFloat32SliceToTemporaryArrayBuffer(data)
	return temporaryFloat32Array.Call("subarray", 0, minLength)
}
