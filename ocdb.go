// Copyright 2019 The Alice-Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocdb exposes types and functions to read and write OCDB files.
package ocdb

import (
	"fmt"
	"io"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

// Entry represents a single entry in an OCDB data file.
type Entry struct {
	base  rbase.Object
	obj   root.Object `groot:"fObject"`
	id    ID          `groot:"fId"`
	meta  *MetaData   `groot:"fMetaData"`
	owner bool        `groot:"fIsOwner"`
}

func (*Entry) Class() string   { return "AliCDBEntry" }
func (*Entry) RVersion() int16 { return 1 }

func (entry *Entry) Object() root.Object { return entry.obj }

func (entry *Entry) Display(w io.Writer) {
	fmt.Fprintf(w, `=== Entry ===
ID: %v
Owner: %v
`,
		entry.id, entry.owner,
	)
	if entry.meta != nil {
		fmt.Fprintf(w, "MetaData:\n")
		entry.meta.Display(w)
	}
	if entry.obj != nil {
		fmt.Fprintf(w, "Object: %T\n%v\n===\n", entry.obj, entry.obj)
	}
}

// MarshalROOT implements rbytes.Marshaler
func (entry *Entry) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(entry.RVersion())

	entry.base.MarshalROOT(w)
	w.WriteObjectAny(entry.obj)
	entry.id.MarshalROOT(w)
	w.WriteObjectAny(entry.meta)
	w.WriteBool(entry.owner)

	return w.SetByteCount(pos, entry.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *Entry) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	o.obj = r.ReadObjectAny()
	if err := o.id.UnmarshalROOT(r); err != nil {
		return err
	}
	o.meta = nil
	if obj := r.ReadObjectAny(); obj != nil {
		o.meta = obj.(*MetaData)
	}
	o.owner = r.ReadBool()

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

// ID uniquely identifies an entry in an OCDB file.
type ID struct {
	base    rbase.Object
	path    Path     `groot:"fPath"`
	runs    RunRange `groot:"fRunRange"`
	vers    int32    `groot:"fVersion"`
	subvers int32    `groot:"fSubVersion"`
	last    string   `groot:"fLastStorage"`
}

func (*ID) Class() string   { return "AliCDBId" }
func (*ID) RVersion() int16 { return 1 }

func (id ID) String() string {
	return fmt.Sprintf("AliCDBId{Path: %v, RunRange: %v, Version: 0x%x, SubVersion: 0x%x, Last: %q}",
		id.path, id.runs, id.vers, id.subvers, id.last,
	)
}

// MarshalROOT implements rbytes.Marshaler
func (id *ID) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(id.RVersion())

	id.base.MarshalROOT(w)
	id.path.MarshalROOT(w)
	id.runs.MarshalROOT(w)
	w.WriteI32(id.vers)
	w.WriteI32(id.subvers)
	w.WriteString(id.last)

	return w.SetByteCount(pos, id.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (id *ID) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := id.base.UnmarshalROOT(r); err != nil {
		return err
	}
	if err := id.path.UnmarshalROOT(r); err != nil {
		return err
	}
	if err := id.runs.UnmarshalROOT(r); err != nil {
		return err
	}
	id.vers = r.ReadI32()
	id.subvers = r.ReadI32()
	id.last = r.ReadString()

	r.CheckByteCount(pos, bcnt, start, id.Class())
	return r.Err()
}

// Path represents a provenance path in an OCDB file.
type Path struct {
	base     rbase.Object
	path     string `groot:"fPath"`
	lvl0     string `groot:"fLevel0"`
	lvl1     string `groot:"fLevel1"`
	lvl2     string `groot:"fLevel2"`
	valid    bool   `groot:"fIsValid"`
	wildcard bool   `groot:"fIsWildCard"`
}

func (p Path) String() string {
	return fmt.Sprintf("Path{Path: %q, Level0: %q, Level1: %q, Level2: %q, Valid: %v, WildCard: %v}",
		p.path, p.lvl0, p.lvl1, p.lvl2, p.valid, p.wildcard,
	)
}

func (*Path) Class() string   { return "AliCDBPath" }
func (*Path) RVersion() int16 { return 1 }

// MarshalROOT implements rbytes.Marshaler
func (p *Path) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(p.RVersion())

	p.base.MarshalROOT(w)
	w.WriteString(p.path)
	w.WriteString(p.lvl0)
	w.WriteString(p.lvl1)
	w.WriteString(p.lvl2)
	w.WriteBool(p.valid)
	w.WriteBool(p.wildcard)

	return w.SetByteCount(pos, p.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (p *Path) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := p.base.UnmarshalROOT(r); err != nil {
		return err
	}

	p.path = r.ReadString()
	p.lvl0 = r.ReadString()
	p.lvl1 = r.ReadString()
	p.lvl2 = r.ReadString()
	p.valid = r.ReadBool()
	p.wildcard = r.ReadBool()

	r.CheckByteCount(pos, bcnt, start, p.Class())
	return r.Err()
}

// RunRange represents a [first, last] range of run numbers.
type RunRange struct {
	base  rbase.Object
	First int32 `groot:"fFirstRun"`
	Last  int32 `groot:"fLastRun"`
}

func (*RunRange) Class() string   { return "AliCDBRunRange" }
func (*RunRange) RVersion() int16 { return 1 }

func (rr RunRange) String() string {
	return fmt.Sprintf("RunRange{First: %d, Last: %d}", rr.First, rr.Last)
}

// MarshalROOT implements rbytes.Marshaler
func (rr *RunRange) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(rr.RVersion())

	rr.base.MarshalROOT(w)
	w.WriteI32(rr.First)
	w.WriteI32(rr.Last)

	return w.SetByteCount(pos, rr.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (rr *RunRange) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := rr.base.UnmarshalROOT(r); err != nil {
		return err
	}
	rr.First = r.ReadI32()
	rr.Last = r.ReadI32()

	r.CheckByteCount(pos, bcnt, start, rr.Class())
	return r.Err()
}

// MetaData stores optional metadata associated with an entry in an OCDB file.
type MetaData struct {
	base    rbase.Object
	class   string    `groot:"fObjectClassName"`
	resp    string    `groot:"fResponsible"`
	beam    uint32    `groot:"fBeamPeriod"`
	vers    string    `groot:"fAliRootVersion"`
	comment string    `groot:"fComment"`
	props   rcont.Map `groot:"fProperties"`
}

func (*MetaData) Class() string   { return "AliCDBMetaData" }
func (*MetaData) RVersion() int16 { return 1 }

func (meta *MetaData) Display(w io.Writer) {
	fmt.Fprintf(w, "Class: %q\nResponsible: %q\nBeamPeriod: %d\nAliRoot Version: %q\nComment: %q\nProperties: %d\n",
		meta.class, meta.resp, meta.beam, meta.vers, meta.comment, len(meta.props.Table()),
	)
	for k, v := range meta.props.Table() {
		fmt.Fprintf(w, "  key: %v\n  val: %v\n", k, v)
	}
}

// MarshalROOT implements rbytes.Marshaler
func (meta *MetaData) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(meta.RVersion())

	meta.base.MarshalROOT(w)
	w.WriteString(meta.class)
	w.WriteString(meta.resp)
	w.WriteU32(meta.beam)
	w.WriteString(meta.vers)
	w.WriteString(meta.comment)
	meta.props.MarshalROOT(w)

	return w.SetByteCount(pos, meta.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (meta *MetaData) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := meta.base.UnmarshalROOT(r); err != nil {
		return err
	}

	meta.class = r.ReadString()
	meta.resp = r.ReadString()
	meta.beam = r.ReadU32()
	meta.vers = r.ReadString()
	meta.comment = r.ReadString()

	if err := meta.props.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, start, meta.Class())
	return r.Err()
}
