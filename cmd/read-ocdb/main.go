// Copyright 2019 The Alice-Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/alice-go/ocdb"
	"go-hep.org/x/hep/groot"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func main() {
	flag.Parse()

	f, err := groot.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, k := range f.Keys() {
		log.Printf("key: %v", k.Name())

		o, err := f.Get(k.Name())
		if err != nil {
			log.Fatal(err)
		}

		v := o.(*ocdb.Entry)

		dumpEntry(v)
	}
}

func dumpEntry(v *ocdb.Entry) {

	obj := v.Object()
	bimap := obj.(*ocdb.AliMUON2DMap)

	deid := 706

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
	fmt.Println(string(b[:]))
}
