package common

import (
	"reflect"
	"testing"

	. "proto_gen_ext/output"
)

func TestSimple(t *testing.T) {
	data := make([]byte, 128)
	bs := NewByteStream(data)
	if err := bs.WriteUint64(10086); err != nil {
		t.Error(err)
	}
	if err := bs.WriteUint8(1); err != nil {
		t.Error(err)
	}
	if err := bs.WriteString("hello"); err != nil {
		t.Error(err)
	}
	if err := bs.WriteInt32(32); err != nil {
		t.Error(err)
	}

	//write
	if v, err := bs.ReadUint64(); err != nil {
		t.Error(err)
	} else if v != 10086 {
		t.Fatal("v!=10086")
	}

	if v, err := bs.ReadUint8(); err != nil {
		t.Error(err)
	} else if v != 1 {
		t.Fatal("v!=1")
	}

	if v, err := bs.ReadString(); err != nil {
		t.Error(err)

	} else if v != "hello" {
		t.Fatal("v!=hello")
	}

	if v, err := bs.ReadInt32(); err != nil {
		t.Error(err)
	} else if v != 32 {
		t.Fatal("v!=32")
	}

	//reset
	bs.Reset()
	bs.WriteInt16(16)
	if v, err := bs.ReadInt16(); err != nil {
		t.Error(err)
	} else if v != 16 {
		t.Fatal("v!=16")
	}
}

func TestGenSB(t *testing.T) {
	data := make([]byte, 128)
	var bs = NewByteStream(data)
	myv := &MyVec2{
		X: 1,
		Y: 2,
	}
	myv.Write(bs)

	myv2 := &MyVec2{}
	myv2.Read(bs)
	if !reflect.DeepEqual(myv, myv2) {
		t.Fatal("not DeepEqual vec")
	}

	r1 := &MyRole{
		Myname: "aaa",
		Bt:     []byte("bbb"),
		Mypos:  *myv,
		Arri32: []int32{1, 2, 3},
		Plist:  []MyVec2{{1, 2}, {3, 4}},
		SList:  []string{"abc", "123"},
	}

	r1.Write(bs)

	r2 := &MyRole{}
	r2.Read(bs)
	if !reflect.DeepEqual(r1, r2) {
		t.Fatal("not DeepEqual role")
	}
	{
		r, w := bs.GetRWPos()
		t.Log("rw2:", r, w)
	}

	bs.WriteInt16(16)
	i16, err := bs.ReadInt16()
	if err != nil {
		t.Fatal("i16 error:", err)
	}
	if i16 != 16 {
		t.Fatal("i16 != 16", i16)
	}

}
