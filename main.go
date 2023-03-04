package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	// . "github.com/go-git/go-git/v5/_examples"

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

func track() {

	output_file := "rewindMe/" + time.Now().String() + ".mkv"
	capture_cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-pix_fmt", "yuyv422", "-i", "1:0", "-r", "0.5", output_file)
	output, _ := capture_cmd.CombinedOutput()
	fmt.Println(string(output))
}

func pause(capture_cmd chan exec.Cmd) {
	// (<-capture_cmd).Process.Kill()
	cd_cmd := exec.Command("cd", "rewindMe")
	pwd_cmd := exec.Command("pwd")
	dir, _ := pwd_cmd.CombinedOutput()
	fmt.Println(string(dir))
	add_cmd := exec.Command("git", "add", ".")
	output, _ := add_cmd.CombinedOutput()
	fmt.Println(string(output))
	fmt.Println("Added")
	msg := "msg"
	commit_cmd := exec.Command("git", "commit", "-m", msg)
	fmt.Println("committed")
	push_cmd := exec.Command("git", "push")
	out, _ := push_cmd.CombinedOutput()
	fmt.Println(string(out))
	fmt.Println("pushed")

	cd_cmd.Run()

	add_cmd.Run()
	commit_cmd.Run()
	push_cmd.Run()
	fmt.Println("pause")
}

func main() {
	runtime.LockOSThread()
	cocoa.TerminateAfterWindowsClose = false

	mode = PAUSE

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {

		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("▶️")

		itemTrackOrPause := cocoa.NSMenuItem_New()
		nextClicked := make(chan bool)
		capture_cmd := make(chan exec.Cmd)
		go func() {
			for {
				select {
				case <-nextClicked:
					if mode == PAUSE {
						mode = TRACK
					} else {
						mode = PAUSE
					}
					if mode == TRACK {
						core.Dispatch(func() { itemTrackOrPause.SetTitle("Pause") })
						go track()
					} else {
						core.Dispatch(func() { itemTrackOrPause.SetTitle("Track") })
						go pause(capture_cmd)
					}
				}
			}
		}()

		// nextClicked <- true
		itemTrackOrPause.SetTitle("Track")
		itemTrackOrPause.SetAction(objc.Sel("nextClicked:"))
		cocoa.DefaultDelegateClass.AddMethod("nextClicked:", func(_ objc.Object) {
			nextClicked <- true
			fmt.Println("track/pause!")
		})

		itemTrace := cocoa.NSMenuItem_New()
		itemTrace.SetTitle("Trace")
		itemTrace.SetAction(objc.Sel("trace:"))
		cocoa.DefaultDelegateClass.AddMethod("trace:", func(_ objc.Object) {
			mode = TRACE
			fmt.Println("trace!")
		})

		itemQuit := cocoa.NSMenuItem_New()
		itemQuit.SetTitle("Quit")
		itemQuit.SetAction(objc.Sel("terminate:"))
		cocoa.DefaultDelegateClass.AddMethod("terminate:", func(_ objc.Object) {
			fmt.Println("quit!")
		})

		menu := cocoa.NSMenu_New()
		menu.AddItem(itemTrackOrPause)
		menu.AddItem(itemTrace)
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)

		//<-quit
		/*
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
		*/
		// events := make(chan cocoa.NSEvent, 64)
		// go func() {
		// 	<-events
		// 	cocoa.NSApp().Terminate()
		// }()
		// cocoa.NSEvent_GlobalMonitorMatchingMask(cocoa.NSEventMaskAny, events)
	})

	app.ActivateIgnoringOtherApps(true)

	app.Run()

}
