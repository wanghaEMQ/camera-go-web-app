package main

import (
  "fmt"
  "net/http"
  "os"
  "path"
  "encoding/json"
)

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
  path := CameraPreview {
      Path: "/images/preview.jpg",
  }
  byteArray, err := json.Marshal(path)
  if err != nil {
      fmt.Println(err)
  }
  rw.Write(byteArray)
}

func main() {
  rootdir, err := os.Getwd()
  if err != nil {
    rootdir = "No dice"
  }

  http.Handle("/", http.FileServer(http.Dir("web")))
  // Handler for anything pointing to /images/
  http.Handle("/images/", http.StripPrefix("/images",
        http.FileServer(http.Dir(path.Join(rootdir, "images/")))))
  http.HandleFunc("/test", handlertest)
  http.HandleFunc("/camerapreview", handler_camerapreview)
  http.ListenAndServe(":8080", nil)
}
