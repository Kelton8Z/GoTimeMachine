package main

import (
	"context"
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

type Mode int64

var mode Mode

const (
	TRACK Mode = 0
	PAUSE Mode = 1
	TRACE Mode = 2
)

func Crash() error {
	return errors.Errorf("this function is supposed to crash")
}

func main() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		fontName := flag.String("font", "Helvetica", "font to use")
		flag.Parse()

		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("▶️")

		itemTrack := cocoa.NSMenuItem_New()
		itemTrack.SetTitle("Track")
		itemTrack.SetAction(objc.Sel("track:"))
		cocoa.DefaultDelegateClass.AddMethod("track:", func(_ objc.Object) {
			mode = TRACK
			fmt.Println("track!")
		})

		itemPause := cocoa.NSMenuItem_New()
		itemPause.SetTitle("Pause")
		itemPause.SetAction(objc.Sel("pause:"))
		cocoa.DefaultDelegateClass.AddMethod("trace:", func(_ objc.Object) {
			mode = PAUSE
			fmt.Println("pause!")
		})

		itemTrace := cocoa.NSMenuItem_New()
		itemTrace.SetTitle("Trace")
		itemTrace.SetAction(objc.Sel("trace:"))
		cocoa.DefaultDelegateClass.AddMethod("trace:", func(_ objc.Object) {
			mode = TRACE
			fmt.Println("trace!")
		})

		menu := cocoa.NSMenu_New()
		menu.AddItem(itemTrack)
		menu.AddItem(itemPause)
		menu.AddItem(itemTrace)
		obj.SetMenu(menu)

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

	app.ActivateIgnoringOtherApps(true)
	app.Run()

	cd_cmd := exec.Command("cd", "rewindMe")
	add_cmd := exec.Command("git", "add", ".")
	msg := "msg"
	commit_cmd := exec.Command("git", "commit", " -m", msg)
	push_cmd := exec.Command("git", "push")

	cd_cmd.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Second)
	defer cancel()

	output_file := time.Now().String() + ".mkv"
	capture_cmd := exec.CommandContext(ctx, "ffmpeg", "-f", "avfoundation", "-pix_fmt", "yuyv422", "-i", "1:0", "-r", "0.5", output_file)
	output, _ := capture_cmd.CombinedOutput()
	fmt.Println(string(output))
	for {
		capture_cmd.Run()
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Paused")
			break
		}
	}
	capture_cmd.Process.Kill()

	add_cmd.Run()
	commit_cmd.Run()
	push_cmd.Run()

}
