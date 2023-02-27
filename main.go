package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func main() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		fontName := flag.String("font", "Helvetica", "font to use")
		flag.Parse()

		screen := cocoa.NSScreen_Main().Frame().Size
		text := fmt.Sprintf(" %s ", strings.Join(flag.Args(), " "))
		tr, fontSize := func() (rect core.NSRect, size float64) {
			t := cocoa.NSTextView_Init(core.Rect(0, 0, 0, 0))
			t.SetString(text)
			for s := 70.0; s <= 550; s += 12 {
				t.SetFont(cocoa.Font(*fontName, s))
				t.LayoutManager().EnsureLayoutForTextContainer(t.TextContainer())
				rect = t.LayoutManager().UsedRectForTextContainer(t.TextContainer())
				size = s
				if rect.Size.Width >= screen.Width {
					break
				}
			}
			return rect, size
		}()

		t := cocoa.NSTextView_Init(tr)
		t.SetString(text)
		t.SetFont(cocoa.Font(*fontName, fontSize))
		t.SetEditable(false)
		t.SetImportsGraphics(false)
		t.SetDrawsBackground(false)

		c := cocoa.NSView_Init(core.Rect(0, 0, 0, 0))
		c.SetBackgroundColor(cocoa.Color(0, 0, 0, 0.75))
		c.SetWantsLayer(true)
		c.Layer().SetCornerRadius(32.0)
		c.AddSubviewPositionedRelativeTo(t, cocoa.NSWindowAbove, nil)

		tr.Size.Height = screen.Height
		tr.Size.Width = screen.Width
		tr.Origin.X = (screen.Width / 2) - (tr.Size.Width / 2)
		tr.Origin.Y = (screen.Height / 2) - (tr.Size.Height / 2)

		w := cocoa.NSWindow_Init(core.Rect(0, 0, 0, 0),
			cocoa.NSWindowStyleMaskFullScreen, cocoa.NSBackingStoreBuffered, false)
		w.SetContentView(c)
		w.SetTitlebarAppearsTransparent(true)
		w.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		w.SetOpaque(false)
		w.SetBackgroundColor(cocoa.NSColor_Clear())
		w.SetLevel(cocoa.NSMainMenuWindowLevel + 2)
		w.SetFrameDisplay(tr, true)
		cls := objc.NewClass("AppDelegate", "NSWindow")
		cls.AddMethod("applicationDidFinishLaunching:", func(app objc.Object) {
			w.ToggleFullScreen(nil)
			fmt.Println("Launched!")
		})

		w.MakeKeyAndOrderFront(nil)
		objc.RegisterClass(cls)
		fmt.Println(screenshot.NumActiveDisplays())

		delegate := objc.Get("AppDelegate").Alloc().Init()
		app := objc.Get("NSApplication").Get("sharedApplication")
		app.Set("delegate:", delegate)
		app.Send("run")

		// events := make(chan cocoa.NSEvent, 64)
		// go func() {
		// 	<-events
		// 	cocoa.NSApp().Terminate()
		// }()
		// cocoa.NSEvent_GlobalMonitorMatchingMask(cocoa.NSEventMaskAny, events)
	})

	// go full screen
	// cls := objc.NewClass("AppDelegate", "NSObject")

	// objc.RegisterClass(cls)

	// delegate := objc.Get("AppDelegate").Alloc().Init()
	// app.Set("delegate:", delegate)
	// app.NSApplicationPresentationFullScreen
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
