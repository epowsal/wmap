// wmap project wmap.go by suirosu exgaya epowsal wlb iwlb@outlook.com exgaya@gmail.com 20241122;
package wmap

import (
	"encoding/binary"
	"errors"
	"math"
	"quickset"
	"reflect"
	"strings"
)

type Wmap[K comparable, V any] struct {
	root *quickset.QuickSet[*Item[K, V]]
	K2B  func(k K) []byte
	B2K  func([]byte) K
}

type Item[K comparable, V any] struct {
	ek string
	v  V
	bv bool
	i  *quickset.QuickSet[*Item[K, V]]
}

func cmp[K comparable, V any](a, b *Item[K, V]) int {
	if a.ek[0] < b.ek[0] {
		return -1
	} else if a.ek[0] > b.ek[0] {
		return 1
	} else {
		return 0
	}
}

func scmp(a, b string) (nsame, cmpr int) {
	al := 0
	bl := 0
	for true {
		if al >= len(a) {
			if len(a) < len(b) {
				return nsame, -1
			} else if len(a) > len(b) {
				return nsame, 1
			} else {
				return nsame, 0
			}
		}
		if bl >= len(b) {
			if len(a) < len(b) {
				return nsame, -1
			} else if len(a) > len(b) {
				return nsame, 1
			} else {
				return nsame, 0
			}
		}
		if a[al] == b[bl] {
			nsame += 1
			al += 1
			bl += 1
		} else {
			return nsame, int(a[al]) - int(b[bl])
		}
	}
	return nsame, strings.Compare(a, b)
}

func ssame(a, b string) (nsame int) {
	al := 0
	bl := 0
	for true {
		if al >= len(a) {
			return nsame
		}
		if bl >= len(b) {
			return nsame
		}
		if a[al] == b[bl] {
			nsame += 1
			al += 1
			bl += 1
		} else {
			return nsame
		}
	}
	return nsame
}

func New[K comparable, V any]() *Wmap[K, V] {
	w := &Wmap[K, V]{}
	w.root = quickset.New[*Item[K, V]]()
	return w
}

func (w *Wmap[K, V]) Put(key K, val V) error {
	kb := w.KeyToBytes(key)
	if len(kb) == 0 {
		return errors.New("key empty error")
	}
	cur := w.root
	ke := string(kb)
	for true {
		i := cur.Index(cmp[K, V], &Item[K, V]{ek: string([]byte{ke[0]})})
		if i == -1 {
			cur.Insert(cmp[K, V], &Item[K, V]{ek: ke, v: val, bv: true, i: nil})
			return nil
		} else {
			sn := ssame(string(ke), cur.It[i].ek)
			if sn == len(cur.It[i].ek) {
				if len(ke) == len(cur.It[i].ek) {
					cur.It[i].v = val
					cur.It[i].bv = true
					return nil
				} else {
					ke = ke[sn:]
					cur = cur.It[i].i
				}
			} else { //be sn<len(ek);
				// if sn >= len(cur.It[i].ek) {
				// 	e.P("sn >= len(cur.It[i].ek)", sn, []byte(ke), []byte(cur.It[i].ek))
				// 	panic("value sn is >= ek error")
				// }
				if sn == len(ke) {
					on := &Item[K, V]{ek: cur.It[i].ek[sn:], v: cur.It[i].v, bv: cur.It[i].bv, i: cur.It[i].i}
					cur.It[i].ek = cur.It[i].ek[:sn]
					cur.It[i].v = val
					cur.It[i].bv = true
					cur.It[i].i = quickset.NewN(cmp[K, V], on)
					return nil
				} else if sn < len(ke) {
					on := &Item[K, V]{ek: cur.It[i].ek[sn:], v: cur.It[i].v, bv: cur.It[i].bv, i: cur.It[i].i}
					nn := &Item[K, V]{ek: ke[sn:], v: val, bv: true, i: nil}
					cur.It[i].ek = cur.It[i].ek[:sn]
					var df V
					cur.It[i].v = df
					cur.It[i].bv = false
					cur.It[i].i = quickset.NewN(cmp[K, V], on, nn)
					return nil
				}
			}
		}
	}
	return errors.New("put data fail")
}

func (w *Wmap[K, V]) Get(key K) (val V, er error) {
	kb := w.KeyToBytes(key)
	if len(kb) == 0 {
		return val, errors.New("key empty error")
	}
	cur := w.root
	ke := string(kb)
	for true {
		i := cur.Index(cmp[K, V], &Item[K, V]{ek: string([]byte{ke[0]})})
		if i == -1 {
			return val, errors.New("not found")
		} else {
			sn := ssame(string(ke), cur.It[i].ek)
			if sn == len(cur.It[i].ek) {
				if len(ke) == len(cur.It[i].ek) {
					if cur.It[i].bv {
						return cur.It[i].v, nil
					} else {
						return val, errors.New("not found")
					}
				} else {
					if sn >= 1 {
						ke = ke[sn:]
						cur = cur.It[i].i
					} else {
						return val, errors.New("not found")
					}
				}
			} else {
				return val, errors.New("not found")
			}
		}
	}
	return val, errors.New("not found")
}

func (w *Wmap[K, V]) Del(key K) error {
	kb := w.KeyToBytes(key)
	if len(kb) == 0 {
		return errors.New("key empty error")
	}
	cur := w.root
	ke := string(kb)
	for true {
		i := cur.Index(cmp[K, V], &Item[K, V]{ek: string([]byte{ke[0]})})
		if i == -1 {
			return errors.New("not found")
		} else {
			sn := ssame(string(ke), cur.It[i].ek)
			if sn == len(cur.It[i].ek) {
				if len(ke) == len(cur.It[i].ek) {
					if cur.It[i].bv {
						var df V
						cur.It[i].v = df
						cur.It[i].bv = false
						return nil
					} else {
						return errors.New("not found")
					}
				} else {
					ke = ke[sn:]
					cur = cur.It[i].i
				}
			}
		}
	}
	return errors.New("not found")
}

func (w *Wmap[K, V]) KeyToBytes(k K) []byte {
	tpn := reflect.TypeOf(k).String()
	switch tpn {
	case "int":
		v := reflect.ValueOf(k).Int()
		vb := make([]byte, reflect.TypeOf(k).Bits()/8)
		if len(vb) == 8 {
			binary.BigEndian.PutUint64(vb, uint64(v))
		} else {
			binary.BigEndian.PutUint32(vb, uint32(v))
		}
		return vb
	case "uint":
		v := reflect.ValueOf(k).Uint()
		vb := make([]byte, reflect.TypeOf(k).Bits()/8)
		if len(vb) == 8 {
			binary.BigEndian.PutUint64(vb, uint64(v))
		} else {
			binary.BigEndian.PutUint32(vb, uint32(v))
		}
		return vb
	case "int8":
		v := reflect.ValueOf(k).Int()
		vb := make([]byte, 1)
		vb[0] = byte(v)
		return vb
	case "uint8", "byte":
		v := reflect.ValueOf(k).Uint()
		vb := make([]byte, 1)
		vb[0] = byte(v)
		return vb
	case "int16":
		v := reflect.ValueOf(k).Int()
		vb := make([]byte, 2)
		binary.BigEndian.PutUint16(vb, uint16(v))
		return vb
	case "uint16":
		v := reflect.ValueOf(k).Uint()
		vb := make([]byte, 2)
		binary.BigEndian.PutUint16(vb, uint16(v))
		return vb
	case "int32":
		v := reflect.ValueOf(k).Int()
		vb := make([]byte, 4)
		binary.BigEndian.PutUint32(vb, uint32(v))
		return vb
	case "uint32":
		v := reflect.ValueOf(k).Uint()
		vb := make([]byte, 4)
		binary.BigEndian.PutUint32(vb, uint32(v))
		return vb
	case "int64":
		v := reflect.ValueOf(k).Int()
		vb := make([]byte, 8)
		binary.BigEndian.PutUint64(vb, uint64(v))
		return vb
	case "uint64":
		v := reflect.ValueOf(k).Uint()
		vb := make([]byte, 8)
		binary.BigEndian.PutUint64(vb, uint64(v))
		return vb
	case "float32":
		v := reflect.ValueOf(k).Float()
		vb := make([]byte, 4)
		binary.BigEndian.PutUint32(vb, math.Float32bits(float32(v)))
		return vb
	case "float64":
		v := reflect.ValueOf(k).Float()
		vb := make([]byte, 8)
		binary.BigEndian.PutUint64(vb, math.Float64bits(float64(v)))
		return vb
	case "string":
		v := reflect.ValueOf(k).String()
		return []byte(v)
	default:
		return w.K2B(k)
	}
}

func (w *Wmap[K, V]) KeyFromBytes(b []byte) (k K) {
	tpn := reflect.TypeOf(k).String()
	switch tpn {
	case "int":
		if len(b) == 8 {
			u := binary.BigEndian.Uint64(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(u))
			return k
		} else if len(b) == 4 {
			u := binary.BigEndian.Uint32(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(u))
			return k
		}
	case "uint":
		if len(b) == 8 {
			u := binary.BigEndian.Uint64(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(u))
			return k
		} else if len(b) == 4 {
			u := binary.BigEndian.Uint32(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(u))
			return k
		}
	case "int8":
		if len(b) == 1 {
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(b[0]))
			return k
		}
	case "uint8":
		if len(b) == 1 {
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(b[0]))
			return k
		}
	case "int16":
		if len(b) == 2 {
			u := binary.BigEndian.Uint16(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(u))
			return k
		}
	case "uint16":
		if len(b) == 2 {
			u := binary.BigEndian.Uint16(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(u))
			return k
		}
	case "int32":
		if len(b) == 4 {
			u := binary.BigEndian.Uint32(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(u))
			return k
		}
	case "uint32":
		if len(b) == 4 {
			u := binary.BigEndian.Uint32(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(u))
			return k
		}
	case "int64":
		if len(b) == 8 {
			u := binary.BigEndian.Uint64(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetInt(int64(u))
			return k
		}
	case "uint64":
		if len(b) == 8 {
			u := binary.BigEndian.Uint64(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetUint(uint64(u))
			return k
		}
	case "float32":
		if len(b) == 4 {
			u := binary.BigEndian.Uint32(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetFloat(float64(math.Float32frombits(u)))
			return k
		}
	case "float64":
		if len(b) == 8 {
			u := binary.BigEndian.Uint64(b)
			v := reflect.ValueOf(&k).Elem()
			v.SetFloat(math.Float64frombits(u))
			return k
		}
	case "string":
		v := reflect.ValueOf(&k).Elem()
		v.SetString(string(b))
		return k
	default:
		return w.B2K(b)
	}
	return k
}
