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

void rfbClientLogInfoToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyClientLogInfo(buffer, n);
}

void rfbClientLogErrToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyClientLogErr(buffer, n);
}

void setClientRfbLog() {
	rfbClientLog = rfbClientLogInfoToString;
	rfbClientErr = rfbClientLogErrToString;
}