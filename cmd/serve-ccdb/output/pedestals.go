package output

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go-hep.org/x/hep/groot"
)

func dumpRequest(w io.Writer, r *http.Request) {
	output, err := httputil.DumpRequest(r, false)
	if err != nil {
		fmt.Fprintln(w, "Error dumping request:", err)
		return
	}
	fmt.Fprintln(w, string(output))
}

func dumpResponse(w io.Writer, r *http.Response) {
	output, err := httputil.DumpResponse(r, false)
	if err != nil {
		fmt.Fprintln(w, "Error dumping response:", err)
		return
	}
	fmt.Fprintln(w, string(output))
}

func getLocation(run int) (string, error) {
	serverURL := viper.GetString("ccdb")

	u, err := url.Parse(serverURL + "QcTaskMCH/QcMuonChambers_Pedestals_DS39" + strconv.Itoa(run))
	if err != nil {
		return "", errors.Wrap(err, "Could not parse 1st url")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", errors.Wrap(err, "Could not create 1st request")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not make 1st request")
	}

	return resp.Header.Get("Location"), nil
}

func doRequest(run int) (*http.Response, error) {

	location, err := getLocation(run)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get location for run %d", run)
	}

	serverURL := viper.GetString("ccdb")
	u, err := url.Parse(serverURL + "/" + location)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse 2nd url")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create 2nd request")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func getRootFileName(run int) (string, error) {
	resp, err := doRequest(run)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// dump the response body into a temporary (root) file
	// that we'll read from afterwards
	f, err := ioutil.TempFile("", "alice-go-ocdb")
	if err != nil {
		return "", errors.Wrap(err, "Could not open temporary file")
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Could not copy body response to temporary file")
	}
	err = f.Close()
	if err != nil {
		return "", errors.Wrap(err, "Could not close temporary file")
	}
	return f.Name(), err
}

func JSONPedestals(w io.Writer, run int, deid int) error {

	fname, err := getRootFileName(run)
	defer os.Remove(fname)
	if err != nil {
		return errors.Wrap(err, "Could not get Root file name")
	}

	froot, err := groot.Open(fname)
	if err != nil {
		return errors.Wrap(err, "Could not get Root file from CCDB")
	}
	defer froot.Close()

	// cdbEntry, err := froot.Get("AliCDBEntry")
	// if err != nil {
	// 	return errors.Wrap(err, "Could not get AliCDBEntry from file")
	// }
	// e := cdbEntry.(*ocdb.Entry)
	//
	// obj := e.Object()
	// bimap := obj.(*ocdb.AliMUON2DMap)
	//
	// manus := bimap.GetManusForDE(deid)
	//
	// type DS struct {
	// 	DSId int
	// 	Mean float64
	// }
	// data := struct {
	// 	DEId       int
	// 	DualSampas []DS
	// }{
	// 	DEId: deid,
	// }
	//
	// for _, m := range manus {
	// 	o := bimap.GetObject(m.DeID, m.ID)
	// 	c := o.(*ocdb.AliMUONCalibParamND)
	// 	data.DualSampas = append(data.DualSampas, DS{
	// 		DSId: int(c.ID1()),
	// 		Mean: c.Value(0, 0) / c.Value(3, 0) / c.Value(4, 0),
	// 	})
	// }
	// b, _ := json.MarshalIndent(data, "", " ")
	// fmt.Fprintln(w, string(b[:]))
	return nil

}
