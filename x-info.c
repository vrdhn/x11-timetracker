#include <X11/Xlib.h>
#include <X11/Xmu/WinUtil.h>
#include <X11/Xutil.h>
#include <X11/extensions/scrnsaver.h>
#include <locale.h>
#include <stdlib.h>
#include <string.h>

const char * XInfoFocussedWindowClass();
const char * XInfoFocussedWindowName();


static  char errorMessageBuffer[1024];
static  char appNameBuffer[1024];
static  char classNameBuffer[1024];
static  char titleNameBuffer[1024];

static  Display *openDisplay = NULL;
static  Window prevWindow = 0;

int handle_error(Display* display, XErrorEvent* error)
{
	XGetErrorText(display,
		      error->error_code,
		      errorMessageBuffer,
		      sizeof(errorMessageBuffer));
	return 1;
}

const char * XInfoInitialize() {
	setlocale(LC_ALL, ""); // see man locale
	XSetErrorHandler(handle_error);
	errorMessageBuffer[0] = 0;
	openDisplay = XOpenDisplay(NULL);
	if (!openDisplay) {
		return errorMessageBuffer;
	} else {
		errorMessageBuffer[0] = 0;
		return errorMessageBuffer;
	}
}

int  XInfoIdleTime() {
	XScreenSaverInfo* info = XScreenSaverAllocInfo();
	XScreenSaverQueryInfo(openDisplay, DefaultRootWindow(openDisplay), info);
	int ret =  (info->idle + 999) / 1000;
	XFree(info);
	return ret;
}
const char * XInfoFocussedWindowApp() {
	return appNameBuffer;
}
const char * XInfoFocussedWindowClass() {
	return classNameBuffer;
}
const char * XInfoFocussedWindowTitle() {
	return titleNameBuffer;
}

void XInfoCalculate() {
	// 1. Window
	Window w;
	int revert_to;
	XGetInputFocus(openDisplay, &w, &revert_to); //

	// Optimize...
	if ( w == prevWindow )
		return;
	else
		prevWindow = w;

	Window parent = w;
	Window root = None;
	Window* children;
	unsigned nchildren;
	while (parent != root) {
		w = parent;
		if (XQueryTree(openDisplay, w, &root, &parent, &children, &nchildren))
			XFree(children);
	}

	w =  XmuClientWindow(openDisplay, w);

	// 2. Title
	XTextProperty prop;
	Status s = XGetWMName(openDisplay, w, &prop);
	if (!s)
		return;
	int count = 0, result;
	char** list = NULL;
	result = XmbTextPropertyToTextList(openDisplay, &prop, &list, &count); // see man
	if (result == Success) {
		strncpy(titleNameBuffer,list[0],sizeof(titleNameBuffer)-1);
		titleNameBuffer[sizeof(titleNameBuffer)-1] = 0;
		XFreeStringList(list);
	} else {
		titleNameBuffer[0] = 0;
	}

	// 3. app, class.
	XClassHint* kls;
	kls = XAllocClassHint(); // see man
	if (XGetClassHint(openDisplay, w, kls)) {
		strncpy(appNameBuffer,kls->res_name,sizeof(appNameBuffer)-1);
		appNameBuffer[sizeof(appNameBuffer)-1]=0;
		strncpy(classNameBuffer,kls->res_class,sizeof(classNameBuffer)-1);
		classNameBuffer[sizeof(classNameBuffer)-1]=0;
	} else {
		appNameBuffer[0]=0;
		classNameBuffer[0]=0;
	}
	XFree(kls);
}
