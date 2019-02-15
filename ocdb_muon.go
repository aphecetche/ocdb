package ocdb

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

type AliMUON2DMap struct {
	base  AliMUONVStore
	exmap *AliMpExMap `groot:"fMap"`
	opt   bool        `groot:"fOptimizeForDEManu"`
}

type Manu struct {
	DeID int
	ID   int
}

func sortedManus(manus []Manu) []Manu {
	sort.Slice(manus, func(i, j int) bool {
		if manus[i].DeID == manus[j].DeID {
			return manus[i].ID < manus[j].ID
		}
		return manus[i].DeID < manus[j].DeID
	})
	return manus
}

func (*AliMUON2DMap) Class() string   { return "AliMUON2DMap" }
func (*AliMUON2DMap) RVersion() int16 { return 1 }

func (m *AliMUON2DMap) ExMap() *AliMpExMap { return m.exmap }

func (m *AliMUON2DMap) GetObject(deid, manuid int) root.Object {
	objects := m.exmap.Objects()
	keys := m.exmap.Keys()
	for i := 0; i < objects.Len(); i++ {
		de := keys.At(i)
		if int(de) != deid {
			continue
		}
		om := objects.At(i).(*AliMpExMap)
		k := om.Keys()
		o := om.Objects()
		for j := 0; j < o.Len(); j++ {
			if int(k.At(j)) != manuid {
				continue
			}
			return o.At(j)
		}
	}
	return nil
}

func (m *AliMUON2DMap) GetManusForDE(deid int) []Manu {
	var manus []Manu
	objects := m.exmap.Objects()
	keys := m.exmap.Keys()
	for i := 0; i < objects.Len(); i++ {
		de := keys.At(i)
		if int(de) != deid {
			continue
		}
		om := objects.At(i).(*AliMpExMap)
		k := om.Keys()
		o := om.Objects()
		for j := 0; j < o.Len(); j++ {
			manus = append(manus, Manu{int(de), int(k.At(j))})
		}
	}
	return sortedManus(manus)
}

func (m *AliMUON2DMap) GetManus() []Manu {
	var manus []Manu
	objects := m.exmap.Objects()
	keys := m.exmap.Keys()
	for i := 0; i < objects.Len(); i++ {
		de := keys.At(i)
		manus = append(manus, m.GetManusForDE(int(de))...)
	}
	return sortedManus(manus)
}

func (m *AliMUON2DMap) String() string {
	return fmt.Sprintf("MUON2DMap{Opt: %v, Map: %v}", m.opt, *m.exmap)
}

// MarshalROOT implements rbytes.Marshaler
func (o *AliMUON2DMap) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	o.base.MarshalROOT(w)
	w.WriteObjectAny(o.exmap)
	w.WriteBool(o.opt)

	return w.SetByteCount(pos, o.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *AliMUON2DMap) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	o.exmap = nil
	if obj := r.ReadObjectAny(); obj != nil {
		o.exmap = obj.(*AliMpExMap)
	}
	o.opt = r.ReadBool()

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

type AliMUONVStore struct {
	base rbase.Object
}

func (*AliMUONVStore) Class() string   { return "AliMUONVStore" }
func (*AliMUONVStore) RVersion() int16 { return 1 }

// MarshalROOT implements rbytes.Marshaler
func (o *AliMUONVStore) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	o.base.MarshalROOT(w)

	return w.SetByteCount(pos, o.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *AliMUONVStore) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

type AliMpExMap struct {
	base rbase.Object
	objs rcont.ObjArray `groot:"fObjects"`
	keys rcont.ArrayL64 `groot:"fKeys"`
}

func (exmap AliMpExMap) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "ExMap{Objs: [")
	for i := 0; i < exmap.objs.Len(); i++ {
		if i > 0 {
			fmt.Fprintf(o, ", ")
		}
		fmt.Fprintf(o, "%v", exmap.objs.At(i))
	}
	fmt.Fprintf(o, "], Keys: %v}", exmap.keys.Data)
	return o.String()
}

func (*AliMpExMap) Class() string   { return "AliMpExMap" }
func (*AliMpExMap) RVersion() int16 { return 1 }

func (e *AliMpExMap) Objects() rcont.ObjArray { return e.objs }
func (e *AliMpExMap) Keys() rcont.ArrayL64    { return e.keys }

// MarshalROOT implements rbytes.Marshaler
func (o *AliMpExMap) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	o.base.MarshalROOT(w)
	o.objs.MarshalROOT(w)
	o.keys.MarshalROOT(w)

	return w.SetByteCount(pos, o.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *AliMpExMap) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	if err := o.objs.UnmarshalROOT(r); err != nil {
		return err
	}

	if err := o.keys.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

type AliMUONCalibParamND struct {
	base AliMUONVCalibParam
	dim  int32     `groot:"fDimension"`
	size int32     `groot:"fSize"`
	n    int32     `groot:"fN"`
	vs   []float64 `groot:"fValues"`
}

func (*AliMUONCalibParamND) Class() string   { return "AliMUONCalibParamND" }
func (*AliMUONCalibParamND) RVersion() int16 { return 1 }

func (c *AliMUONCalibParamND) index(i, j int) int {
	return i + int(c.size)*j
}

func (c *AliMUONCalibParamND) ID0() uint32 {
	return c.base.base.ID & 0xFFFF
}

func (c *AliMUONCalibParamND) ID1() uint32 {
	return (c.base.base.ID & 0xFFFF0000) >> 16
}

func (c *AliMUONCalibParamND) Value(i, j int) float64 {
	return c.vs[c.index(i, j)]
}

func (c *AliMUONCalibParamND) MeanAndSigma(dim int) (float64, float64) {
	mean := 0.0
	v2 := 0.0
	n := int(c.size)
	for i := 0; i < n; i++ {
		v := c.Value(i, dim)
		mean += v
		v2 += v * v
	}
	mean /= float64(n)
	sigma := 0.0
	if n > 1 {
		sigma = math.Sqrt((v2 - float64(n)*mean*mean) / (float64(n) - 1))
	}
	return mean, sigma
}

func (c *AliMUONCalibParamND) Print(w io.Writer, opt string) {
	fmt.Fprintf(w, "AliMUONCalibParamND Id=(%d,%d) Size=%d Dimension=%d\n",
		c.ID0(), c.ID1(), c.size, c.dim)
	opt = strings.ToUpper(opt)
	if strings.Contains(opt, "FULL") {
		for i := 0; i < int(c.size); i++ {
			fmt.Fprintf(w, "CH %3d", i)
			for j := 0; j < int(c.dim); j++ {
				fmt.Fprintf(w, " %g", c.Value(i, j))
			}
			fmt.Fprint(w, "\n")
		}
	}
	if strings.Contains(opt, "MEAN") {
		var j int
		fmt.Sscanf(opt, "MEAN%d", &j)
		mean, sigma := c.MeanAndSigma(j)
		fmt.Fprintf(w, " Mean(j=%d)=%g Sigma(j=%d)=%g\n", j, mean, j, sigma)
	}
}

// MarshalROOT implements rbytes.Marshaler
func (o *AliMUONCalibParamND) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	o.base.MarshalROOT(w)
	w.WriteI32(o.dim)
	w.WriteI32(o.size)
	w.WriteI32(o.n)
	w.WriteI8(1) // FIXME(sbinet)
	w.WriteFastArrayF64(o.vs)

	return w.SetByteCount(pos, o.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *AliMUONCalibParamND) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	o.dim = r.ReadI32()
	o.size = r.ReadI32()
	o.n = r.ReadI32()
	_ = r.ReadI8() // FIXME(sbinet)
	o.vs = r.ReadFastArrayF64(int(o.n))

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

type AliMUONVCalibParam struct {
	base rbase.Object
}

func (*AliMUONVCalibParam) Class() string   { return "AliMUONVCalibParam" }
func (*AliMUONVCalibParam) RVersion() int16 { return 1 }

// MarshalROOT implements rbytes.Marshaler
func (o *AliMUONVCalibParam) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	o.base.MarshalROOT(w)

	return w.SetByteCount(pos, o.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (o *AliMUONVCalibParam) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := o.base.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}
