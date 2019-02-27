package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alice-go/ocdb"
	"go-hep.org/x/hep/groot"
)

var ccdbServer = "http://localhost:6464"

func dumpRequest(r *http.Request) {
	output, err := httputil.DumpRequest(r, false)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

func dumpResponse(r *http.Response) {
	output, err := httputil.DumpResponse(r, false)
	if err != nil {
		fmt.Println("Error dumping response:", err)
		return
	}
	fmt.Println(string(output))
}

func process(client *http.Client, path string) {
	f, err := groot.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	key := "AliCDBEntry"
	o, err := f.Get(key)
	if err != nil {
		fmt.Printf("Could not get key %s from file %s\n", key, path)
	}
	v := o.(*ocdb.Entry)
	// using run range as timestamp range for the moment
	// FIXME: should read the corresponding GRP/GRP/Data object to
	// get the run->timestamp relationship and use timestamps
	// as validity range for the put
	r1 := v.Id().Runs().First
	r2 := v.Id().Runs().Last
	fmt.Printf("%T path=%s R=%d,%d\n", v, path, r1, r2)
	url := ccdbServer + "/OccupancyMap/MUON/" + strconv.Itoa(int(r1))

	r, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open file %s", path)
	}
	var requestBody bytes.Buffer
	mpw := multipart.NewWriter(&requestBody)

	w, err := mpw.CreateFormFile("data", path)
	if err != nil {
		log.Fatal("Cannot create form file %s", err.Error())
	}

	_, err = io.Copy(w, r)
	if err != nil {
		log.Fatal("Cannot copy file to request body %s", err.Error())
	}
	mpw.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	req.Header.Set("Content-Type", mpw.FormDataContentType())

	if err != nil {
		log.Fatal("Could not create request %s", err.Error())
	}
	dumpRequest(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request did not go well %s", err.Error())
	}
	defer resp.Body.Close()
	dumpResponse(resp)
}

func main() {
	dir := "/Users/laurent/cernbox/ocdbs/2018/OCDB/MUON/Calib/OccupancyMap"
	imax := 0
	client := &http.Client{Timeout: 2 * time.Second}

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), "Run") &&
			filepath.Ext(path) == ".root" {
			imax--
			if imax == 0 {
				return fmt.Errorf("toto")
			}
			process(client, path)
		}
		return nil
	})
}
