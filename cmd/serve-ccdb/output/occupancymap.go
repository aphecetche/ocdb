package output

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/alice-go/ocdb"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go-hep.org/x/hep/groot"
)

func JSONOccupancyMap(w io.Writer, run int, deid int) error {

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

	cdbEntry, err := froot.Get("AliCDBEntry")
	if err != nil {
		return errors.Wrap(err, "Could not get AliCDBEntry from file")
	}
	e := cdbEntry.(*ocdb.Entry)

	obj := e.Object()
	bimap := obj.(*ocdb.AliMUON2DMap)

	manus := bimap.GetManusForDE(deid)

	type DS struct {
		DSId int
		Mean float64
	}
	data := struct {
		DEId       int
		DualSampas []DS
	}{
		DEId: deid,
	}

	for _, m := range manus {
		o := bimap.GetObject(m.DeID, m.ID)
		c := o.(*ocdb.AliMUONCalibParamND)
		data.DualSampas = append(data.DualSampas, DS{
			DSId: int(c.ID1()),
			Mean: c.Value(0, 0) / c.Value(3, 0) / c.Value(4, 0),
		})
	}
	b, _ := json.MarshalIndent(data, "", " ")
	fmt.Fprintln(w, string(b[:]))
	return nil

}
