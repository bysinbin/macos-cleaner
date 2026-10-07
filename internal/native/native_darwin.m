#import "native_darwin.h"

// Forward declarations of Go callbacks from _cgo_export.h
extern void goQuickFreeRAM(void);
extern void goQuickSmartCare(void);
extern void goOpenAppWindow(void);

@interface DiskCleanerAppDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate, WKNavigationDelegate>
@property (strong, nonatomic) NSWindow *window;
@property (strong, nonatomic) WKWebView *webView;
@property (strong, nonatomic) NSStatusItem *statusItem;
@property (copy, nonatomic) NSString *serverURL;
@property (copy, nonatomic) NSString *appTitle;
@end

static DiskCleanerAppDelegate *globalDelegate = nil;

@implementation DiskCleanerAppDelegate

- (instancetype)initWithURL:(NSString *)url title:(NSString *)title {
    self = [super init];
    if (self) {
        _serverURL = [url copy];
        _appTitle = [title copy];
    }
    return self;
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    // 1. Create Native macOS Window
    NSRect screenRect = [[NSScreen mainScreen] visibleFrame];
    CGFloat winWidth = 1220;
    CGFloat winHeight = 820;
    if (screenRect.size.width < winWidth) winWidth = screenRect.size.width - 40;
    if (screenRect.size.height < winHeight) winHeight = screenRect.size.height - 60;

    NSRect frame = NSMakeRect((screenRect.size.width - winWidth) / 2 + screenRect.origin.x,
                              (screenRect.size.height - winHeight) / 2 + screenRect.origin.y,
                              winWidth, winHeight);

    NSUInteger style = NSWindowStyleMaskTitled |
                       NSWindowStyleMaskClosable |
                       NSWindowStyleMaskMiniaturizable |
                       NSWindowStyleMaskResizable |
                       NSWindowStyleMaskFullSizeContentView;

    self.window = [[NSWindow alloc] initWithContentRect:frame
                                              styleMask:style
                                                backing:NSBackingStoreBuffered
                                                  defer:NO];

    self.window.title = self.appTitle;
    self.window.titleVisibility = NSWindowTitleHidden;
    self.window.titlebarAppearsTransparent = YES;
    self.window.backgroundColor = [NSColor colorWithSRGBRed:0.06 green:0.09 blue:0.16 alpha:1.0];
    [self.window setDelegate:self];

    // 2. Create and configure WKWebView
    WKWebViewConfiguration *config = [[WKWebViewConfiguration alloc] init];
    [config.preferences setValue:@YES forKey:@"developerExtrasEnabled"];

    self.webView = [[WKWebView alloc] initWithFrame:self.window.contentView.bounds configuration:config];
    self.webView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    self.webView.navigationDelegate = self;
    self.window.contentView = self.webView;

    // Load local server URL
    NSURL *targetURL = [NSURL URLWithString:self.serverURL];
    NSURLRequest *req = [NSURLRequest requestWithURL:targetURL cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:15.0];
    [self.webView loadRequest:req];

    // 3. Setup macOS Menubar Status Item (Tray Agent)
    self.statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    self.statusItem.button.title = @"🍏 DiskCleaner";

    NSMenu *menu = [[NSMenu alloc] initWithTitle:@"DiskCleaner"];

    NSMenuItem *openItem = [[NSMenuItem alloc] initWithTitle:@"🍏 DiskCleaner Pro'yu Aç" action:@selector(actionOpenWindow:) keyEquivalent:@"o"];
    [openItem setTarget:self];
    [menu addItem:openItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *actHeader = [[NSMenuItem alloc] initWithTitle:@"Hızlı Bakım İşlemleri:" action:nil keyEquivalent:@""];
    [actHeader setEnabled:NO];
    [menu addItem:actHeader];

    NSMenuItem *ramItem = [[NSMenuItem alloc] initWithTitle:@"⚡ Inaktif RAM Belleği Boşalt" action:@selector(actionFreeRAM:) keyEquivalent:@"r"];
    [ramItem setTarget:self];
    [menu addItem:ramItem];

    NSMenuItem *smartCareItem = [[NSMenuItem alloc] initWithTitle:@"✨ Smart Care Hızlı Temizle" action:@selector(actionSmartCare:) keyEquivalent:@"s"];
    [smartCareItem setTarget:self];
    [menu addItem:smartCareItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *watcherItem = [[NSMenuItem alloc] initWithTitle:@"🗑️ Çöp Kutusu Artık Takipçisi (Aktif)" action:nil keyEquivalent:@""];
    [watcherItem setEnabled:NO];
    [menu addItem:watcherItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:@"Çıkış (Quit)" action:@selector(terminate:) keyEquivalent:@"q"];
    [quitItem setTarget:NSApp];
    [menu addItem:quitItem];

    self.statusItem.menu = menu;

    // Show window and bring app to foreground
    [self.window makeKeyAndOrderFront:nil];
    [self.window makeFirstResponder:self.webView];
    [NSApp activateIgnoringOtherApps:YES];
}

- (BOOL)windowShouldClose:(NSWindow *)sender {
    if ([self.window styleMask] & NSWindowStyleMaskFullScreen) {
        [self.window toggleFullScreen:nil];
    }
    [NSApp terminate:nil];
    return YES;
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender hasVisibleWindows:(BOOL)flag {
    [self.window makeKeyAndOrderFront:nil];
    [self.window makeFirstResponder:self.webView];
    [NSApp activateIgnoringOtherApps:YES];
    return YES;
}

- (void)webView:(WKWebView *)webView didFailProvisionalNavigation:(WKNavigation *)navigation withError:(NSError *)error {
    // Retry loading if local server is binding
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(300 * NSEC_PER_MSEC)), dispatch_get_main_queue(), ^{
        NSURL *targetURL = [NSURL URLWithString:self.serverURL];
        NSURLRequest *req = [NSURLRequest requestWithURL:targetURL cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:15.0];
        [self.webView loadRequest:req];
    });
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    exit(0);
}

- (void)actionOpenWindow:(id)sender {
    [self.window makeKeyAndOrderFront:nil];
    [self.window makeFirstResponder:self.webView];
    [NSApp activateIgnoringOtherApps:YES];
    goOpenAppWindow();
}

- (void)actionFreeRAM:(id)sender {
    goQuickFreeRAM();
}

- (void)actionSmartCare:(id)sender {
    goQuickSmartCare();
}

- (void)updateStatusTitle:(NSString *)title {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (self.statusItem && self.statusItem.button) {
            self.statusItem.button.title = title;
        }
    });
}

- (void)showAppWindow {
    dispatch_async(dispatch_get_main_queue(), ^{
        [self.window makeKeyAndOrderFront:nil];
        [self.window makeFirstResponder:self.webView];
        [NSApp activateIgnoringOtherApps:YES];
    });
}

@end

void runCocoaApp(const char *cUrl, const char *cTitle) {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        NSString *url = [NSString stringWithUTF8String:cUrl];
        NSString *title = [NSString stringWithUTF8String:cTitle];

        globalDelegate = [[DiskCleanerAppDelegate alloc] initWithURL:url title:title];
        [app setDelegate:globalDelegate];
        [app run];
    }
}

void setMenubarTitle(const char *cTitle) {
    if (globalDelegate) {
        NSString *t = [NSString stringWithUTF8String:cTitle];
        [globalDelegate updateStatusTitle:t];
    }
}

void bringWindowToFront(void) {
    if (globalDelegate) {
        [globalDelegate showAppWindow];
    }
}

void showNativeNotification(const char *title, const char *subtitle, const char *message) {
    @autoreleasepool {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
        NSUserNotification *notif = [[NSUserNotification alloc] init];
        notif.title = [NSString stringWithUTF8String:title];
        if (subtitle) {
            notif.subtitle = [NSString stringWithUTF8String:subtitle];
        }
        notif.informativeText = [NSString stringWithUTF8String:message];
        notif.soundName = NSUserNotificationDefaultSoundName;
        [[NSUserNotificationCenter defaultUserNotificationCenter] deliverNotification:notif];
#pragma clang diagnostic pop
    }
}
