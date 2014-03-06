package main

import (
    "os"
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

    go func () {
        for {
            select {
            case ev := <-watcher.Event:
                if ev.IsModify() {
                    sw.ReadAll()
                    sw.Redraw()
                }
            case err := <-watcher.Error:
                log.Println("error: ", err)
            }
        }
    }()

    go func () {
        err = watcher.Watch(st.Ce(os.Args[1], ".rlt"))
        if err != nil {
            log.Fatal(err)
        }
    }()

    sw.Dlg.Show()
    sw.OpenFile(st.Ce(os.Args[1], ".inp"))
    sw.EscapeAll()
    sw.ShowCenter()
    sw.ExecCommand(stgui.ERRORELEM)
    val := iup.MainLoop()
    if val == iup.CLOSE {
        <-done
        watcher.Close()
    }
}
