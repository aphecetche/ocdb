// Copyright 2019 The Alice-Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	manus := bimap.GetManusForDE(706)
	// manus := bimap.GetManus()

	for _, m := range manus {
		o := bimap.GetObject(m.DeID, m.ID)
		fmt.Printf("DE %4d MANU %4d", m.DeID, m.ID)
		c := o.(*ocdb.AliMUONCalibParamND)
		c.Dump(os.Stdout)
	}
}
