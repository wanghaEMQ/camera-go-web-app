package main

import (
  "fmt"
  "log"
  "net/http"
  "io"
  "os"
  "path"
  "encoding/json"
  "time"
  "bytes"
)

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

  byteArray := ask_camera_preview()

  /*
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
  */
  rw.Write(byteArray)
}

func handler_startrecord(rw http.ResponseWriter, r *http.Request) {
  ask_camera_on()

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
  ask_camera_off()

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

func handler_status(rw http.ResponseWriter, r *http.Request) {
  ans := ask_camera_status()

  res := Record {
    Txt: string(ans),
  }
  byteArray, err := json.Marshal(res)
  if err != nil {
    fmt.Println(err)
  }
  rw.Write(byteArray)
}

func ask_camera_on() {
	res := http_send("start-record")
	log.Println("start-record => ", string(res))
}

func ask_camera_off() {
	res := http_send("stop-record")
	log.Println("stop-record => ", string(res))
}

func ask_camera_preview() []byte {
	res := http_send("preview")
	log.Println("preview => ", string(res)[:10])
	return res
}

func ask_camera_status() []byte {
	res := http_send("status")
	log.Println("status => ", string(res))
	return res
}

func http_send(data string) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", "http://127.0.0.1:9999", bytes.NewReader([]byte(data)))
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	return body
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
  http.HandleFunc("/startrecord", handler_startrecord)
  http.HandleFunc("/stoprecord", handler_stoprecord)
  http.HandleFunc("/status", handler_status)
  http.ListenAndServe(":11111", nil)
}
