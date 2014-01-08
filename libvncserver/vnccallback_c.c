#include <rfb/rfb.h>
#include <stdio.h>
#include <stdarg.h>
#include "_cgo_export.h"

void rfbServerLogInfoToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyServerLogInfo(buffer, n);
}

void rfbServerLogErrToString(const char *format, ...) {
	int bufferSize = 4096;
	char buffer[bufferSize];
	va_list argptr;
    va_start(argptr, format);
    int n = vsnprintf(buffer, bufferSize, format, argptr);
    va_end(argptr);
    notifyServerLogErr(buffer, n);
}

void setServerRfbLog() {
	rfbLog = rfbServerLogInfoToString;
	rfbErr = rfbServerLogErrToString;
}