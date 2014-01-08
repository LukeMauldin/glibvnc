#include <rfb/rfbclient.h>
#include <stdio.h>
#include <stdarg.h>
#include "_cgo_export.h"


void setMallocFrameBufferProc(rfbClient *client) {
  client->MallocFrameBuffer = (MallocFrameBufferProc)onSourceClientResize;
}

void setGotFrameBufferUpdate(rfbClient *client) {
  client->GotFrameBufferUpdate = (GotFrameBufferUpdateProc)onSourceClientUpdate;
}

void rfbLogInfoToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyLogInfo(buffer, n);
}

void rfbLogErrToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyLogErr(buffer, n);
}

void setRfbLog() {
	rfbClientLog = rfbLogInfoToString;
	rfbClientErr = rfbLogErrToString;
}