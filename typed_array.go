package webgl

import (
	"fmt"
	"syscall/js"

	// Copy & paste from https://github.com/hajimehoshi/ebiten/tree/main/internal/jsutil
	"github.com/atercattus/webgl/ebiten_jsutil"
)

func jsTypedArrayOf(srcData interface{}) js.Value {
	switch s := srcData.(type) {
	case []float32:
		arr := ebiten_jsutil.TemporaryFloat32Array(len(s), s)
		return arr.Call("subarray", 0, len(s))
	case []uint16:
		arr := ebiten_jsutil.TemporaryUint8ArrayFromUint16Slice(len(s), s)
		return arr.Call("subarray", 0, len(s))
	case []uint8:
		arr := ebiten_jsutil.TemporaryUint8ArrayFromUint8Slice(len(s), s)
		return arr.Call("subarray", 0, len(s))
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at jsTypedArrayOf: %T %v", s, s))
	}
}
