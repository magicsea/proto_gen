package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type GeneratorCsharp struct {
	stdTypes map[string]StdTypeDef
}

func NewGeneratorCsharp() *GeneratorCsharp {
	return &GeneratorCsharp{stdTypes: stdTypesCs}
}

// 是否基础类型
func (gen *GeneratorCsharp) isStdType(t string) bool {
	for k, _ := range gen.stdTypes {
		if k == t {
			return true
		}
	}
	return false
}
func (gen *GeneratorCsharp) anyType() string {
	return "object"
}
func (gen *GeneratorCsharp) GetGenFilePath() string {
	return *outputDir + currProtoConfig.process.GetNamespace() + "Init.cs"
}

func (gen *GeneratorCsharp) GenHeader(p ProtoProcess) string {
	start := int(*startId)
	idname := gen.getIDTypeName()
	tables := p.GetTables()
	if *sortMsg > 0 {
		sort.Slice(tables, func(i, j int) bool { return tables[i] < tables[j] })
	}

	initText := ""
	for i, v := range tables {
		id := currGenIDFunc(start, i, v)
		initText += fmt.Sprintf("            [%v] = typeof(%s)", id, v)
		if i != len(tables)-1 {
			initText += ",\n"
		} else {
			initText += "\n"
		}
	}

	tabExtText := ""
	for i, v := range tables {
		id := currGenIDFunc(start, i, v)
		tabExtText += fmt.Sprintf(cstabExtTemplate, v, currProtoConfig.newTypeMsgCs, idname, id, idname, id, v)
	}

	importText := ""
	confs := p.GenExternConf()

	var imports []string
	if confs != nil {
		imports = getImportcs(confs)

	}

	if len(*importpack) > 0 {
		argsImports := strings.Split(*importpack, ",")
		imports = append(imports, argsImports...)
	}

	if len(imports) > 0 {
		for _, s := range imports {
			importText += fmt.Sprintf("using %s;\n", s)
		}
	}

	msgName := p.GetNamespace()
	if len(*mgrname) > 0 {
		msgName = *mgrname
	}

	caseStr := gen.genNewTypeCase(p)
	vs := strings.Split(*version, "#")
	vstring := fmt.Sprintf("		public static string %s { get {return \"%s\";} }", vs[0], vs[1])
	content := fmt.Sprintf(csfileTemplate, importText, p.GetNamespace(), msgName, vstring, idname, idname, initText, currProtoConfig.newTypeMsgCs, p.GetNamespace(), idname, caseStr, tabExtText)
	return content
}

func (gen *GeneratorCsharp) genNewTypeCase(p ProtoProcess) string {
	start := int(*startId)
	caseContent := ""
	for i, msg := range p.GetTables() {
		id := currGenIDFunc(start, i, msg)
		caseContent += fmt.Sprintf(`
				case %v:
					return new %s();`, id, msg)
	}
	return caseContent
}

func (gen *GeneratorCsharp) GenExtern(p ProtoProcess) string {
	confs := p.GenExternConf()
	if confs == nil {
		return ""
	}

	var content = ""
	for _, conf := range confs {
		content += fmt.Sprintf(`
namespace %s
{`, p.GetNamespace())
		for _, v := range conf.Messages {
			content += gen.genMsg(&v)
		}
		content += "}"
	}

	return content
}

func (gen *GeneratorCsharp) getArrayTypeName(etype string) string {
	if std, ok := gen.stdTypes[etype]; ok {
		return std.TypeName + "[]"
	}
	return etype + "[]"
}

func (gen *GeneratorCsharp) getIDTypeName() string {
	return gen.stdTypes[*idType].TypeName
}

func (gen *GeneratorCsharp) genMsg(msg *MsgDef) string {
	fieldContent := gen.genField(msg)
	readC, writeC, sizeC := gen.genFunc(msg)

	content := fmt.Sprintf(csclassTemplate, msg.About, msg.Name, fieldContent, writeC, readC, sizeC)
	return content
}

func (gen *GeneratorCsharp) genField(msg *MsgDef) string {

	fieldContent := ""
	for _, f := range msg.Fields {
		tname := f.Type
		if b, etype := isArrayType(f.Type); b {
			tname = gen.getArrayTypeName(etype)
		} else if gen.isStdType(f.Type) {
			tname = gen.stdTypes[f.Type].TypeName
		} else if isAnyType(f.Type) {
			tname = gen.anyType()
		}
		fieldContent += fmt.Sprintf("		public %s %s;\n", tname, f.Name)
	}
	return fieldContent
}

func (gen *GeneratorCsharp) genFunc(msg *MsgDef) (string, string, string) {
	writeContent := ""
	for _, f := range msg.Fields {
		if gen.isStdType(f.Type) { //isStdType
			writeContent += fmt.Sprintf("			s.%s(%s);\n", gen.stdTypes[f.Type].Write, f.Name)
		} else if b, etype := isArrayType(f.Type); b { //isArrayType
			if etype == "bytes" {
				log.Fatal("not support []bytes")
			}
			if gen.isStdType(etype) {
				writeContent += fmt.Sprintf(cswriteArrayStdTemplate, f.Name, f.Name, gen.stdTypes[etype].Write)
			} else { //struct
				writeContent += fmt.Sprintf(cswriteArrayStructTemplate, f.Name, f.Name)
			}
		} else if isAnyType(f.Type) {
			writeContent += fmt.Sprintf("			MessagePacker.PackAny(%s,s);\n", f.Name)
		} else {
			writeContent += fmt.Sprintf("			%s.Write(s);\n", f.Name)
		}
	}

	readContent := ""
	for _, f := range msg.Fields {
		if gen.isStdType(f.Type) { //isStdType
			readContent += fmt.Sprintf("			%s = s.%s();\n", f.Name, gen.stdTypes[f.Type].Read)
		} else if b, etype := isArrayType(f.Type); b { //isArrayType
			if gen.isStdType(etype) {
				readContent += fmt.Sprintf(csreadArrayStdTemplate, f.Name, gen.stdTypes[etype].TypeName, f.Name, gen.stdTypes[etype].Read, f.Name)
			} else {
				readContent += fmt.Sprintf(csreadArrayStructTemplate, f.Name, etype, f.Name, etype, f.Name)
			}
		} else if isAnyType(f.Type) {
			readContent += fmt.Sprintf(csreadAnyTemplate, f.Name)
		} else { //struct
			readContent += fmt.Sprintf(csreadStructTemplate, f.Name, f.Type, f.Name)
		}
	}

	sizeContent := ""
	for _, f := range msg.Fields {
		if gen.isStdType(f.Type) {
			sizeContent += fmt.Sprintf("			size += %v;\n", gen.stdTypes[f.Type].Size)
			if gen.stdTypes[f.Type].ElemSize > 0 {
				sizeContent += fmt.Sprintf("			size += %s.Length*%v;\n", f.Name, gen.stdTypes[f.Type].ElemSize)
			}

		} else if b, etype := isArrayType(f.Type); b { //isArrayType
			sizeContent += "			size += 2;\n"
			if etype == "string" {
				sizeContent += fmt.Sprintf(csgetStreamSizeStringListTemplate, f.Name)
			} else if gen.isStdType(etype) {
				sizeContent += fmt.Sprintf("			size += %s.Length * %v;\n", f.Name, gen.stdTypes[etype].Size)
			} else {
				sizeContent += fmt.Sprintf(csgetStreamSizeArrayStructTemplate, f.Name)
			}
		} else if isAnyType(f.Type) {
			sizeContent += fmt.Sprintf("			size += MessagePacker.GetSizeAny(%s);\n", f.Name)
		} else {
			sizeContent += fmt.Sprintf("			size += %s.GetStreamSize();\n", f.Name)
		}
	}

	return readContent, writeContent, sizeContent
}

// 总文件模板
var csfileTemplate = `// Code generated by the proto_gen_ext. DO NOT EDIT.
using System;
using System.Collections.Generic;
using System.Text;
%s

namespace %s
{
    public class %s
    {
%s		
        public static Dictionary<%s, Type> typeIndex = new Dictionary<%s, Type>() {
%s
        };
        public static %s New%s(%s id) {
            switch (id) {
%s
            }
            return null;
        }
    }
%s
}
`

// 注入函数模板
var cstabExtTemplate = `
	public partial class %s : %s
	{
		public const %s CMsgID = %v;
		public %s MsgID { get {return %v;} }
		public string MsgName { get {return "%s";} }
	}
`

const csclassTemplate = `
	// %s
    public partial class %s
    {
%s

        public void Write(BytesStream s) 
        {
%s
        }

        public void Read(BytesStream s)
        {
%s
        }
		public int GetStreamSize() {
			int size = 0;
%s
			return size;
		}
    }
`
const csreadStructTemplate = `
			%s = new %s();
			%s.Read(s);
`

// fname,read,fname
const csreadArrayStdTemplate = `
			%s = new %s[s.ReadUInt16()];
			for (int i = 0; i < %s.Length; i++)
            {
				var v = s.%s();
				%s[i] = v;
            }
`

// fname,ftype,fname,fname
const csreadArrayStructTemplate = `
			%s = new %s[s.ReadUInt16()];
			for (int i = 0; i < %s.Length; i++)
			{
				var v = new %s();
				v.Read(s);
				%s[i] = v;
			}	
`
const cswriteArrayStdTemplate = `
			s.WriteUInt16((UInt16)(%s.Length));
			foreach (var v in %s) {
				s.%s(v);
			}
`
const cswriteArrayStructTemplate = `
			s.WriteUInt16((UInt16)(%s.Length));
			foreach (var v in %s) {
				v.Write(s);
			}
`
const csgetStreamSizeArrayStructTemplate = `
			foreach (var v in %s)
			{
				size += v.GetStreamSize();
			}
`
const csgetStreamSizeStringListTemplate = `
			foreach (var v in %s) 
			{
				size += 2;
				size += v.Length;
			}
`

const csreadAnyTemplate = `
			%s = MessagePacker.UnPackAny(s);
`
