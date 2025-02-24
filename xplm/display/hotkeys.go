//Copyright (c) 2015. The goplane AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package display

/*
#cgo CFLAGS: -I../../SDK/CHeaders -fPIC -DSIMDATA_EXPORTS -DXPLM200=1 -DXPLM210=1 -DXPLM300=1 -DXPLM301=1 -DXPLM302=1 -DXPLM303=1
#include <XPLM/XPLMDisplay.h>
#include <stdlib.h>
#include <string.h>

extern void handleHotKey(void* inRefcon);

*/
import "C"
import (
	"unsafe"

	"github.com/xairline/goplane"
	"github.com/xairline/goplane/xplm"
	"github.com/xairline/goplane/xplm/plugins"
)

type HotKeyCallback func(ref interface{})
type HotKeyID unsafe.Pointer

type hotKeyRegData struct {
	id       HotKeyID
	callback HotKeyCallback
	ref      interface{}
}

var hotkeys = make(map[*C.char]*hotKeyRegData)

//export handleHotKey
func handleHotKey(ref unsafe.Pointer) {
	id := (*C.char)(ref)
	regData, _ := hotkeys[id]
	regData.callback(regData.ref)
}

func RegisterHotKey(key xplm.VirtualKeyCode, flags xplm.KeyFlags, description string, callback HotKeyCallback, ref interface{}) HotKeyID {
	cDescription := C.CString(description)
	defer C.free(unsafe.Pointer(cDescription))
	cId := C.CString(goplane.IdGenerator())
	regData := &hotKeyRegData{nil, callback, ref}
	hotkeys[cId] = regData
	hotkeyId := HotKeyID(C.XPLMRegisterHotKey(C.char(key), C.XPLMKeyFlags(flags), cDescription, C.XPLMHotKey_f(unsafe.Pointer(C.handleHotKey)), unsafe.Pointer(cId)))
	regData.id = hotkeyId
	return hotkeyId
}

func XPLMUnregisterHotKey(hotkeyId HotKeyID) {
	for key, regInfo := range hotkeys {
		if regInfo.id == hotkeyId {
			C.free(unsafe.Pointer(key))
			delete(hotkeys, key)
		}
	}
	C.XPLMUnregisterHotKey(C.XPLMHotKeyID(hotkeyId))
}

func CountHotKeys() int {
	return int(C.XPLMCountHotKeys())
}

func GetNthHotKey(index int) HotKeyID {
	return HotKeyID(C.XPLMGetNthHotKey(C.int(index)))
}

func GetHotKeyInfo(hotkeyId HotKeyID) (virtualKey xplm.VirtualKeyCode, flags xplm.KeyFlags, description string, pluginId plugins.PluginId) {
	descBuf := (*C.char)(C.malloc(255))
	defer C.free(unsafe.Pointer(descBuf))
	C.XPLMGetHotKeyInfo(
		C.XPLMHotKeyID(hotkeyId),
		(*C.char)(unsafe.Pointer(&virtualKey)),
		(*C.XPLMKeyFlags)(unsafe.Pointer(&flags)),
		descBuf,
		(*C.XPLMPluginID)(unsafe.Pointer(&pluginId)))
	description = C.GoString(descBuf)
	return
}

func SetHotKeyCombination(hotkeyId HotKeyID, virtualKey xplm.VirtualKeyCode, flags xplm.KeyFlags) {
	C.XPLMSetHotKeyCombination(C.XPLMHotKeyID(hotkeyId), C.char(virtualKey), C.XPLMKeyFlags(flags))
}
