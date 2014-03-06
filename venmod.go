package main

import (
    "flag"
    "log"
    "github.com/howeyc/fsnotify"
    "github.com/visualfc/go-iup/iup"
    "github.com/yofu/st/stlib"
    "github.com/yofu/st/stgui"
)

const (
    version = "0.1.0"
    modified = "LAST CHANGE:06-Mar-2014."
    HOME = "C:/D/CDOCS/Hogan-venhira"
)

func main () {
    hide := flag.Bool("h", false, "Hide Rate<1.0")
    flag.Parse()

    iup.Open()
    defer iup.Close()
    sw := stgui.NewWindow(HOME)
    sw.Version = version
    sw.Modified = modified

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    done := make(chan bool)

    var read bool
    go func () {
        for {
            select {
            case ev := <-watcher.Event:
                if ev.IsModify() {
                    if read {
                        sw.ReadAll()
                        sw.ExecCommand(stgui.ERRORELEM)
                        if *hide {
                            sw.HideNotSelected()
                        } else {
                            sw.LockNotSelected()
                        }
                    }
                    read = !read
                }
            case err := <-watcher.Error:
                log.Println("error: ", err)
            }
        }
    }()

    go func () {
        err = watcher.Watch(st.Ce(flag.Arg(0), ".rlt"))
        if err != nil {
            log.Fatal(err)
        }
    }()

    sw.Dlg.Show()
    sw.OpenFile(st.Ce(flag.Arg(0), ".inp"))
    sw.EscapeAll()
    sw.ShowCenter()
    sw.ReadAll()
    sw.ExecCommand(stgui.ERRORELEM)
    sw.NodeCaptionOff("NC_NUM")
    sw.HideEtype(st.WALL)
    val := iup.MainLoop()
    if val == iup.CLOSE {
        <-done
        watcher.Close()
    }
}
