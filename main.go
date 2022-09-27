package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var startId = flag.Uint64("startid", 1, "start id")                       //起始id
var protoType = flag.String("proto", "pb", "proto type:fbs,pb,sb")        //协议
var inputDir = flag.String("in", "./", "input dir")                       //输入目录
var outputDir = flag.String("out", "./output/", "output dir")             //输出目录
var ouputType = flag.String("type", "go", "ouput type:go or csharp")      //输出文件类型
var sortMsg = flag.Uint64("sort", 1, "sort message by name")              //消息是否按字母排序
var importpack = flag.String("import", "", "import packages")             //导入的包 ，用逗号分割
var outfile = flag.String("outfile", "", "outfile")                       //输出文件名，如果写了会用这个替代
var mgrname = flag.String("mgrname", "", "mgrname")                       //管理对象的命名，默认包名，用作一个包导出多份文件
var packagename = flag.String("package", "", "replace package name")      //替换原来的包名
var genidType = flag.String("idrule", "inc", "how to gen id:hash or inc") //生成id方式
var idType = flag.String("idtype", "uint16", "type of id")                //id类型
var version = flag.String("version", "ProtoVersion#0", "version")         //版本标记

var currProtoConfig *ProtoDef //协议解析器
var currGenerator Generator   //语言生成器

func main() {
	flag.Parse()
	fmt.Println("inputDir:", *inputDir)
	fmt.Println("outputDir:", *outputDir)
	fmt.Println("protoType:", *protoType)
	fmt.Println("startId:", *startId)
	fmt.Println("import", *importpack)
	fmt.Println("outfile", *outfile)
	fmt.Println("mgrname", *mgrname)
	fmt.Println("packagename", *packagename)
	fmt.Println("genidType", *genidType)
	fmt.Println("idType", *idType)
	fmt.Println("version", *version)
	//protoType
	for key, v := range protoDefMap {
		if key == *protoType {
			currProtoConfig = &v
			fmt.Printf("use proto:%+v\n", v)
			break
		}
	}
	if currProtoConfig == nil {
		log.Fatal("not found protoDef,use fbs,pb,bs")
	}
	//ouputType
	if *ouputType == "go" {
		currGenerator = NewGeneratorGolang()
	} else if *ouputType == "csharp" {
		currGenerator = NewGeneratorCsharp()
	} else {
		log.Fatal("not found generator,ouputType need go,csharp")
	}

	//currGenIDFunc
	if *genidType == "hash" {
		currGenIDFunc = GenIDByHash
	} else {
		currGenIDFunc = GenIDByIndex
	}

	//read
	err := filepath.Walk(*inputDir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, currProtoConfig.fileNameExt) {
			return nil
		}
		fmt.Println("do file: ", file.Name())

		errRead := readDSL(path)
		return errRead
	})
	if err != nil {
		log.Fatal("read fail:", err)
	}

	//write
	dir := *outputDir
	errDir := os.MkdirAll(dir, os.ModePerm)
	if errDir != nil {
		log.Fatal("create dir fail:", errDir)
	}
	errWrite := writeGenFile()
	if errWrite != nil {
		log.Fatal("write fail:", errWrite)
	}
}

func readDSL(path string) error {
	if !filepath.IsAbs(path) {
		file, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		path = file
	}

	bts, errFile := ioutil.ReadFile(path)
	if errFile != nil {
		panic(errFile)
	}
	//fmt.Println("fbs:", string(bts))
	str := string(bts)
	err := currProtoConfig.process.Read(str)
	return err
}

func writeGenFile() error {
	filePath := currGenerator.GetGenFilePath()
	if len(*outfile) > 0 {
		filePath = *outputDir + *outfile
	}
	//check same
	m := map[uint32]string{}
	start := int(*startId)
	tables := currProtoConfig.process.GetTables()
	for i, v := range tables {
		id := currGenIDFunc(start, i, v)
		if n, ok := m[id]; ok {
			panic(fmt.Errorf("gen same id(%v),please change proto name:%v,%v", id, n, v))
		}
		m[id] = v
	}

	content := currGenerator.GenHeader(currProtoConfig.process)
	content += currGenerator.GenExtern(currProtoConfig.process)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.WriteString(content)
	f.Sync()
	fmt.Println("writeGenFile end:", filePath)
	return f.Close()
}
