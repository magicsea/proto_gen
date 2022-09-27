package main

import (
	"hash/crc32"
)

type ProtoProcess interface {
	GetNamespace() string
	Read(content string) error
	GetTables() []string
	GenExternConf() []*StreamDef
}

type ProtoDef struct {
	process         ProtoProcess
	namespaceRegexp string
	messageRegexp   string
	fileNameExt     string
	newTypeMsgGo    string
	newTypeMsgCs    string
}

//支持的协议类型
const (
	ProtoPB  = "pb"
	ProtoFBS = "fbs"
	ProtoSB  = "sb"
)

var protoDefMap = map[string]ProtoDef{
	ProtoFBS: {fileNameExt: ".fbs",
		process:         &HeaderProcess{},
		newTypeMsgGo:    "interface{}",
		newTypeMsgCs:    "IFBSMessage",
		namespaceRegexp: `namespace (?s:(.*?))\;`,
		messageRegexp:   `table (?s:(.*?))\{`,
	},
	ProtoPB: {
		fileNameExt:     ".proto",
		process:         &HeaderProcess{},
		newTypeMsgGo:    "interface{}",
		newTypeMsgCs:    "IPBMessage",
		namespaceRegexp: `package (?s:(.*?))\;`,
		messageRegexp:   `message (?s:(.*?))\{`,
	},
	ProtoSB: {
		fileNameExt:  ".json",
		process:      &StreamProcess{},
		newTypeMsgGo: "IMsg",
		newTypeMsgCs: "IMessage",
	},
}

type Generator interface {
	GenHeader(p ProtoProcess) string
	GenExtern(p ProtoProcess) string
	GetGenFilePath() string
}

type GenIDFunc func(start int, index int, name string) uint32

var currGenIDFunc GenIDFunc

func GenIDByIndex(start int, index int, name string) uint32 {
	return uint32(start + index)
}

func GenIDByHash(start int, index int, name string) uint32 {
	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum([]byte(name), crc32q)
}
