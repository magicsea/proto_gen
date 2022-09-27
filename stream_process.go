package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type FieldDef struct {
	Name string
	Type string
}

type MsgDef struct {
	Name   string
	About  string
	Fields []FieldDef
}
type StreamDef struct {
	Package  string
	Importcs []string
	Importgo []string
	Messages []MsgDef
}

type StreamProcess struct {
	HeaderProcess
	//conf StreamDef
	confs []*StreamDef
}

type StdTypeDef struct {
	Write    string
	Read     string
	Size     int
	ElemSize int    //数组类型的元素大小
	TypeName string //默认等于type名字，如果非“”则替换
}

//任意类型，类似interface{}
const AnyType = "any"

func isAnyType(t string) bool {
	return t == AnyType
}

var stdTypesGo = map[string]StdTypeDef{
	"int8":    {"WriteInt8", "ReadInt8", 1, 0, "int8"},
	"uint8":   {"WriteUint8", "ReadUint8", 1, 0, "uint8"},
	"int16":   {"WriteInt16", "ReadInt16", 2, 0, "int16"},
	"uint16":  {"WriteUint16", "ReadUint16", 2, 0, "uint16"},
	"int32":   {"WriteInt32", "ReadInt32", 4, 0, "int32"},
	"uint32":  {"WriteUint32", "ReadUint32", 4, 0, "uint32"},
	"int64":   {"WriteInt64", "ReadInt64", 8, 0, "int64"},
	"uint64":  {"WriteUint64", "ReadUint64", 8, 0, "uint64"},
	"float32": {"WriteFloat32", "ReadFloat32", 4, 0, "float32"},
	"float64": {"WriteFloat64", "ReadFloat64", 8, 0, "float64"},
	"bool":    {"WriteBool", "ReadBool", 1, 0, "bool"},
	"bytes":   {"WriteBytes", "ReadBytes", 4, 1, "[]byte"},   //长度类型uint32。
	"string":  {"WriteString", "ReadString", 2, 1, "string"}, //长度类型uint16。
}

var stdTypesCs = map[string]StdTypeDef{
	"int8":    {"WriteByte", "ReadByte", 1, 0, "Byte"},
	"uint8":   {"WriteByte", "ReadByte", 1, 0, "Byte"},
	"int16":   {"WriteInt16", "ReadInt16", 2, 0, "Int16"},
	"uint16":  {"WriteUInt16", "ReadUInt16", 2, 0, "UInt16"},
	"int32":   {"WriteInt32", "ReadInt32", 4, 0, "Int32"},
	"uint32":  {"WriteUInt32", "ReadUInt32", 4, 0, "UInt32"},
	"int64":   {"WriteInt64", "ReadInt64", 8, 0, "Int64"},
	"uint64":  {"WriteUInt64", "ReadUInt64", 8, 0, "UInt64"},
	"float32": {"WriteFloat", "ReadFloat", 4, 0, "float"},
	"float64": {"WriteDouble", "ReadDouble", 8, 0, "double"},
	"bool":    {"WriteBool", "ReadBool", 1, 0, "bool"},
	"bytes":   {"WriteBytes", "ReadBytes", 4, 1, "byte[]"},   //长度类型uint32。
	"string":  {"WriteString", "ReadString", 2, 1, "string"}, //长度类型uint16。
}

/*
是否数组类型，返回是否，元素类型
数组长度类型uint16，最大65535。
*/
func isArrayType(t string) (bool, string) {
	if strings.HasPrefix(t, "[]") {
		return true, strings.Replace(t, "[]", "", 1)
	}
	return false, ""
}

func (p *StreamProcess) Read(content string) error {
	var conf StreamDef
	err := json.Unmarshal([]byte(content), &conf)
	if err != nil {
		return err
	}
	namespace := conf.Package
	if namespace == "" {
		return nil
	}
	if p.namespace != "" && p.namespace != namespace {
		return errors.New("must same namespace")
	}
	p.namespace = namespace

	fmt.Println("make ", namespace)
	p.confs = append(p.confs, &conf)
	for _, v := range conf.Messages {
		p.tables = append(p.tables, v.Name)
	}
	return nil
}

func (p *StreamProcess) GenExternConf() []*StreamDef {
	return p.confs
}

func getImportcs(confs []*StreamDef) []string {
	var importMap = map[string]string{}
	for _, conf := range confs {
		for _, s := range conf.Importcs {
			importMap[s] = s
		}
	}
	var ilist []string
	for _, v := range importMap {
		ilist = append(ilist, v)
	}
	sort.Slice(ilist, func(i, j int) bool { return ilist[i] < ilist[2] })
	return ilist
}
func getImportgo(confs []*StreamDef) []string {
	var importMap = map[string]string{}
	for _, conf := range confs {
		for _, s := range conf.Importgo {
			importMap[s] = s
		}
	}
	var ilist []string
	for _, v := range importMap {
		ilist = append(ilist, v)
	}
	sort.Slice(ilist, func(i, j int) bool { return ilist[i] < ilist[2] })
	return ilist
}
