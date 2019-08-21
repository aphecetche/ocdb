package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alice-go/ocdb/cmd/serve-ccdb/output"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func serveOccupancyMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		q := r.URL.Query()
		run, err := strconv.Atoi(q.Get("run"))
		if err != nil {
			w.Write([]byte("Malformed run number"))
			return
		}
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed detection element id"))
			return
		}
		err = output.JSONOccupancyMap(w, run, deid)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

func servePedestals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		q := r.URL.Query()
		run, err := strconv.Atoi(q.Get("run"))
		if err != nil {
			w.Write([]byte("Malformed run number"))
			return
		}
		deid, err := strconv.Atoi(q.Get("deid"))
		if err != nil {
			w.Write([]byte("Malformed detection element id"))
			return
		}
		err = output.JSONPedestals(w, run, deid)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

func main() {
	pflag.Int("port", 4242, "port to listen to")
	pflag.String("ccdb", "http://localhost:6464", "url to contact the CCDB")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	http.HandleFunc("/occupancymap", serveOccupancyMap())
	http.HandleFunc("/pedestals", servePedestals())
	fmt.Printf("Serving on port %d and using CCDB %s\n", viper.GetInt("port"), viper.GetString("ccdb"))
	if err := http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), nil); err != nil {
		panic(err)
	}
}
