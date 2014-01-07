package libvncserver

/*
 #include <rfb/rfb.h>
 #cgo CFLAGS: -I/usr/local/include
 #cgo LDFLAGS: -L/usr/local/lib -lvncserver
 extern void setRfbLog();
*/
import "C"

import (
	"io"
	"unsafe"
)

var RfbInfoLogger io.Writer
var RfbErrLogger io.Writer

type GoRfbServer struct {
	rfbServer   C.rfbScreenInfoPtr
	frameBuffer []uint8
}

//Bool conversion
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

func (f *GoRfbServer) ShutdownServer(disconnectClients bool) {
	C.rfbShutdownServer(f.rfbServer, toRfbBool(disconnectClients))
}

func (f *GoRfbServer) ScreenCleanup() {
	//Call libvnc ScreenCleanup
	C.rfbScreenCleanup(f.rfbServer)

	//Cleanup frameBuffer
	f.frameBuffer = nil
}

func (f *GoRfbServer) NewFrameBuffer(width int, height int,
	bitsPerSample int, samplesPerPixel int, bytesPerPixel int) {
	//Create frameBuffer
	f.SetFrameBuffer(width, height, bytesPerPixel)

	//Notify rfbServer
	_frameBuffer := (*C.char)(unsafe.Pointer(&f.frameBuffer[0]))
	_width := (C.int)(width)
	_height := (C.int)(height)
	_bitsPerSample := (C.int)(bitsPerSample)
	_samplesPerPixel := (C.int)(samplesPerPixel)
	_bytesPerPixel := (C.int)(bytesPerPixel)
	C.rfbNewFramebuffer(f.rfbServer, _frameBuffer, _width, _height, _bitsPerSample, _samplesPerPixel, _bytesPerPixel)
}

func (f *GoRfbServer) GetFrameBuffer() *[]uint8 {
	return &f.frameBuffer
}

func (f *GoRfbServer) MarkRectAsModified(x1 int, y1 int, x2 int, y2 int) {
	_x1 := (C.int)(x1)
	_y1 := (C.int)(y1)
	_x2 := (C.int)(x2)
	_y2 := (C.int)(y2)
	C.rfbMarkRectAsModified(f.rfbServer, _x1, _y1, _x2, _y2)
}

func GetScreen(width int, height int, bitsPerSample int, samplesPerPixel int, bytesPerPixel int, argc int, argv *int8) (ret *GoRfbServer) {
	_argc := (*C.int)(unsafe.Pointer(&argc))
	_argv := (**C.char)(unsafe.Pointer(&argv))
	_width := C.int(width)
	_height := C.int(height)
	_bitsPerSample := C.int(bitsPerSample)
	_samplesPerPixel := C.int(samplesPerPixel)
	_bytesPerPixel := C.int(bytesPerPixel)
	_ret := C.rfbGetScreen(_argc, _argv, _width, _height, _bitsPerSample, _samplesPerPixel, _bytesPerPixel)
	if _ret == nil {
		panic("Error in rfbGetScreen")
	}
	ret = &GoRfbServer{
		rfbServer: _ret,
	}
	return
}

func (f *GoRfbServer) ProcessEvents(usec int64) (ret bool) {
	_usec := C.long(usec)
	_ret := C.rfbProcessEvents(f.rfbServer, _usec)
	ret = toGoBool(_ret)
	return
}

func (f *GoRfbServer) SetFrameBuffer(width, height, bytesPerPixel int) {
	f.frameBuffer = make([]uint8, width*height*bytesPerPixel)
	f.rfbServer.frameBuffer = (*C.char)(unsafe.Pointer(&f.frameBuffer[0]))
}

func (f *GoRfbServer) InitServer() {
	C.rfbInitServer(f.rfbServer)
}

func (f *GoRfbServer) SetConfiguration(desktopName string, tcpPort int, httpPort int, httpDir string) {
	//Assign parameters to rfbServer
	f.rfbServer.desktopName = C.CString(desktopName)
	f.rfbServer.alwaysShared = C.TRUE
	f.rfbServer.port = C.int(tcpPort)
	f.rfbServer.httpPort = C.int(httpPort)
	f.rfbServer.httpDir = C.CString(httpDir)
}

func (f *GoRfbServer) IsActive() (ret bool) {
	_ret := C.rfbIsActive(f.rfbServer)
	ret = toGoBool(_ret)
	return
}

//export notifyLogInfo
func notifyLogInfo(str *C.char, n C.int) {
	if RfbInfoLogger != nil {
		goStr := C.GoBytes(unsafe.Pointer(str), n)
		RfbInfoLogger.Write(goStr)
	}
}

//export notifyLogErr
func notifyLogErr(str *C.char, n C.int) {
	if RfbErrLogger != nil {
		goStr := C.GoBytes(unsafe.Pointer(str), n)
		RfbErrLogger.Write(goStr)
	}
}
