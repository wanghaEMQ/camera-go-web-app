package main

import (
  "fmt"
  "net/http"
  "os"
  "path"
  "encoding/json"
  "time"

  "go.nanomsg.org/mangos/v3"
  "go.nanomsg.org/mangos/v3/protocol/push"
  _ "go.nanomsg.org/mangos/v3/transport/all"
)

var ipcsock mangos.Socket
var record_running int = 0

var preview_idx int = 0
var preview_cap int = 20

type CameraPreview struct {
    Path string
}

type Record struct {
    Txt string
}

func handlertest(w http.ResponseWriter, r *http.Request) {
  fileName := "testfile.jpg"
  fmt.Fprintf(w, "<html></br><img src='/images/" + fileName + "' ></html>")
}

func handler_camerapreview(rw http.ResponseWriter, r *http.Request) {
  // Wait 100ms for writing
  time.Sleep(100 * time.Millisecond)

  if record_running != 1 {
    path := CameraPreview {
      Path: "/images/preview.jpg",
    }
    byteArray, err := json.Marshal(path)
    if err != nil {
      fmt.Println(err)
    }
    rw.Write(byteArray)
    return
  }

  mangos_send_preview()

  var str string
  str = fmt.Sprintf("%s%s%d.%s", "/images/", "preview", preview_idx, "jpg")
  fmt.Println(str)
  preview_idx = (preview_idx + 1) % preview_cap

  path := CameraPreview {
      Path: str,
  }
  byteArray, err := json.Marshal(path)
  if err != nil {
      fmt.Println(err)
  }
  rw.Write(byteArray)
}

func handler_startrecord(rw http.ResponseWriter, r *http.Request) {
  mangos_send_start()

  res := Record {
    Txt: "Successfully start",
  }
  byteArray, err := json.Marshal(res)
  if err != nil {
    fmt.Println(err)
  }
  rw.Write(byteArray)

  record_running = 1
}

func handler_stoprecord(rw http.ResponseWriter, r *http.Request) {
  mangos_send_stop()

  res := Record {
    Txt: "Successfully stop",
  }
  byteArray, err := json.Marshal(res)
  if err != nil {
    fmt.Println(err)
  }
  rw.Write(byteArray)

  record_running = 0
}

var ipcurl string = "ipc:///tmp/camerarecord.ipc"

func mangos_start() {
	var err error

	if ipcsock, err = push.NewSocket(); err != nil {
		fmt.Println("can't get new push socket: %s", err.Error())
	}
	if err = ipcsock.Dial(ipcurl); err != nil {
		fmt.Println("can't dial on push socket: %s", err.Error())
	}
}

func mangos_send_start() {
	mangos_send("start-record")
}

func mangos_send_stop() {
	mangos_send("stop-record")
}

func mangos_send_preview() {
	mangos_send("preview")
}

func mangos_send(data string) {
	// data := "IPC://EXTERNAL2NANO:{\"key\":1000,\"offset\":100}"
	// for {
	fmt.Printf("CLIENT: PUBLISHING DATA %s\n", data)
	if err := ipcsock.Send([]byte(data)); err != nil {
		fmt.Println("Failed publishing: %s", err.Error())
	}
	//time.Sleep(time.Millisecond * 200)
}

func main() {
  rootdir, err := os.Getwd()
  if err != nil {
    rootdir = "No dice"
  }
  mangos_start()

  http.Handle("/", http.FileServer(http.Dir("web")))
  // Handler for anything pointing to /images/
  http.Handle("/images/", http.StripPrefix("/images",
        http.FileServer(http.Dir(path.Join(rootdir, "images/")))))
  http.HandleFunc("/test", handlertest)
  http.HandleFunc("/camerapreview", handler_camerapreview)
  http.HandleFunc("/startrecord", handler_startrecord)
  http.HandleFunc("/stoprecord", handler_stoprecord)
  http.ListenAndServe(":11111", nil)
}
