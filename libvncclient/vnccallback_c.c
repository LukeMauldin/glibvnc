#include <rfb/rfbclient.h>
#include "_cgo_export.h"


void setMallocFrameBufferProc(rfbClient *client) {
  client->MallocFrameBuffer = (MallocFrameBufferProc)onSourceClientResize;
}

void setGotFrameBufferUpdate(rfbClient *client) {
  client->GotFrameBufferUpdate = (GotFrameBufferUpdateProc)onSourceClientUpdate;
}
