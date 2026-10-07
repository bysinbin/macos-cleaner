#ifndef NATIVE_DARWIN_H
#define NATIVE_DARWIN_H

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

void runCocoaApp(const char *cUrl, const char *cTitle);
void setMenubarTitle(const char *cTitle);
void bringWindowToFront(void);
void showNativeNotification(const char *title, const char *subtitle, const char *message);

#endif
