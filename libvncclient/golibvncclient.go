package libvncclient

/*
 #include <rfb/rfbclient.h>
 #cgo pkg-config: libvncclient
 extern void setMallocFrameBufferProc(rfbClient *client);
 extern void setGotFrameBufferUpdate(rfbClient *client);
*/
import "C"

import (
	"sync"
	"unsafe"
)

//Bool conversion - copied from libvncserver so user would not need to include both packages
func toGoBool(b C.rfbBool) bool {
	if b == C.TRUE {
		return true
	} else {
		return false
	}
}

func toRfbBool(b bool) C.rfbBool {
	if b == true {
		return C.TRUE
	} else {
		return C.FALSE
	}
}

type GoRfbClient struct {
	rfbClient   *C.rfbClient
	frameBuffer []uint8
}

func RfbGetClient(bitsPerSample int, samplesPerPixel int, bytesPerPixel int) (ret *GoRfbClient) {
	_bitsPerSample := C.int(bitsPerSample)
	_samplesPerPixel := C.int(samplesPerPixel)
	_bytesPerPixel := C.int(bytesPerPixel)
	_ret := C.rfbGetClient(_bitsPerSample, _samplesPerPixel, _bytesPerPixel)
	if _ret == nil {
		panic("Error in rfbGetClient")
	}
	ret = &GoRfbClient{
		rfbClient: _ret,
	}
	return
}

func (f *GoRfbClient) SetFrameBuffer(width, height, bytesPerPixel int) {
	f.frameBuffer = make([]uint8, width*height*bytesPerPixel)
	f.rfbClient.frameBuffer = (*C.uint8_t)(&f.frameBuffer[0])
}

func (f *GoRfbClient) GetFrameBuffer() *[]uint8 {
	return &f.frameBuffer
}

func (f *GoRfbClient) SetMallocFramebufferProc() {
	C.setMallocFrameBufferProc(f.rfbClient)
}

func (f *GoRfbClient) SetGotFrameBufferUpdate() {
	C.setGotFrameBufferUpdate(f.rfbClient)
}

func (f *GoRfbClient) ListenForIncomingConnectionsNoFork(usecTimeout int) (ret int) {
	_usecTimeout := C.int(usecTimeout)
	_ret := C.listenForIncomingConnectionsNoFork(f.rfbClient, _usecTimeout)
	ret = int(_ret)
	return
}

func (f *GoRfbClient) InitClient(argc int, argv *int8) (ret bool) {
	_argc := (*C.int)(unsafe.Pointer(&argc))
	var _argv **C.char
	if argv == nil {
		_argv = nil
	} else {
		_argv = (**C.char)(unsafe.Pointer(&argv))
	}
	_ret := C.rfbInitClient(f.rfbClient, _argc, _argv)
	ret = toGoBool(_ret)
	return
}

func (f *GoRfbClient) RfbInitConnection() (ret bool) {
	_ret := C.rfbInitConnection(f.rfbClient)
	ret = toGoBool(_ret)
	return
}

func (f *GoRfbClient) WaitForMessage(usecs int) (ret int) {
	_usecs := C.uint(usecs)
	_ret := C.WaitForMessage(f.rfbClient, _usecs)
	ret = int(_ret)
	return
}

func (f *GoRfbClient) HandleRFBServerMessage() (ret bool) {
	_ret := C.HandleRFBServerMessage(f.rfbClient)
	ret = toGoBool(_ret)
	return
}

func (f *GoRfbClient) ClientCleanup() {
	//Call libvnc ClientCleanup
	C.rfbClientCleanup(f.rfbClient)

	//Cleanup framebuffer
	f.frameBuffer = nil
}

func (f *GoRfbClient) Width() int {
	return int(f.rfbClient.width)
}

func (f *GoRfbClient) Height() int {
	return int(f.rfbClient.height)
}

func (f *GoRfbClient) SetConfiguration(programName string,
	compressLevel int, qualityLevel int, encodingsString string) {

	//Set options
	f.rfbClient.programName = C.CString(programName)
	f.rfbClient.canHandleNewFBSize = C.TRUE
	f.rfbClient.canUseHextile = C.TRUE
	f.rfbClient.appData.compressLevel = C.int(compressLevel)
	f.rfbClient.appData.qualityLevel = C.int(qualityLevel)
	f.rfbClient.appData.encodingsString = C.CString(encodingsString)
}

func (f *GoRfbClient) SetReverseConnectionServer(address string, port int) {
	f.rfbClient.listenSpecified = C.TRUE
	f.rfbClient.listenAddress = C.CString(address)
	f.rfbClient.listenPort = C.int(port)
}

func (f *GoRfbClient) SetServer(address string, port int) {
	f.rfbClient.serverHost = C.CString(address)
	f.rfbClient.serverPort = C.int(port)
}

type RfbCallback interface {
	OnResize()
	OnUpdate(x, y, w, h int)
}

var rfbCallbackTag int = 1000
var (
	rfbCallbackMap      map[uintptr]RfbCallback
	mutexRfbCallbackMap sync.Mutex
)

func init() {
	rfbCallbackMap = make(map[uintptr]RfbCallback)
}

func (f *GoRfbClient) RegisterRfbCallback(callback RfbCallback) {
	mutexRfbCallbackMap.Lock()
	defer mutexRfbCallbackMap.Unlock()

	//Add to map
	rfbCallbackMap[uintptr(unsafe.Pointer(&callback))] = callback

	//Register with libvnc
	_client := f.rfbClient
	_tag := unsafe.Pointer(&rfbCallbackTag)
	_data := unsafe.Pointer(&callback)
	C.rfbClientSetClientData(_client, _tag, _data)

	//Add callbacks
	f.SetMallocFramebufferProc()
	f.SetGotFrameBufferUpdate()
}

func (f *GoRfbClient) UnregisterRfbCallback(ptr uintptr) {
	mutexRfbCallbackMap.Lock()
	defer mutexRfbCallbackMap.Unlock()
	delete(rfbCallbackMap, ptr)
}

func getRfbCallback(ptr uintptr) RfbCallback {
	mutexRfbCallbackMap.Lock()
	defer mutexRfbCallbackMap.Unlock()
	ret, found := rfbCallbackMap[ptr]
	if !found {
		panic("RfbCallback not found")
	}
	return ret
}

//export onSourceClientResize
func onSourceClientResize(client *C.rfbClient) C.rfbBool {
	storedData := C.rfbClientGetClientData(client, unsafe.Pointer(&rfbCallbackTag))
	callback := getRfbCallback(uintptr(storedData))
	callback.OnResize()
	return C.TRUE
}

//export onSourceClientUpdate
func onSourceClientUpdate(client *C.rfbClient, x, y, w, h int) {
	storedData := C.rfbClientGetClientData(client, unsafe.Pointer(&rfbCallbackTag))
	callback := getRfbCallback(uintptr(storedData))
	callback.OnUpdate(x, y, w, h)
}
