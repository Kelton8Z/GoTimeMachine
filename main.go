package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	// . "github.com/go-git/go-git/v5/_examples"
	"github.com/go-errors/errors"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func Crash() error {
	return errors.Errorf("this function is supposed to crash")
}

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

	stop_cmd := exec.Command("Ctrl-C")
	cd_cmd := exec.Command("cd", "rewindMe")
	add_cmd := exec.Command("git", "add", ".")
	msg := "msg"
	commit_cmd := exec.Command("git", "commit", " -m", msg)
	push_cmd := exec.Command("git", "push")

	cd_cmd.Run()

	output_file := time.Now().String() + ".mkv"
	capture_cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-pix_fmt", "yuyv422", "-i", "1:0", "-r", "0.5", output_file)
	// output, _ := capture_cmd.CombinedOutput()
	// fmt.Println(string(output))
	capture_cmd.Run()
	stop_cmd.Run()
	add_cmd.Run()
	commit_cmd.Run()
	push_cmd.Run()

	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
