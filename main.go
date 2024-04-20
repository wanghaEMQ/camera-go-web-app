package main

import (
  "fmt"
  "net/http"
  "os"
  "path"
)


func handler(w http.ResponseWriter, r *http.Request) {
  fileName := "testfile.jpg"
  fmt.Fprintf(w, "<html></br><img src='/images/" + fileName + "' ></html>")
}

func main() {
  rootdir, err := os.Getwd()
  if err != nil {
    rootdir = "No dice"
  }

  // Handler for anything pointing to /images/
  http.Handle("/images/", http.StripPrefix("/images",
        http.FileServer(http.Dir(path.Join(rootdir, "images/")))))
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
