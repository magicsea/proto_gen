# proto_gen_ext
- 根据协议（flatbuffers,protobuff），生成相应的init文件。   
- 对消息进行排序，编号，方便消息序列化。   
- 扩展消息的方法

### 示例
```
bin/gen_example.bat
proto_gen_ext.exe --startid=100 --proto=pb --in=./../example/ --out=./../output/ --type=go
```
### 参数
```
var startId = flag.Uint64("startid", 1, "start id")                  //起始id
var protoType = flag.String("proto", "pb", "proto type:fbs or pb")   //协议
var inputDir = flag.String("in", "./", "input dir")                  //输入目录
var outputDir = flag.String("out", "./output/", "output dir")        //输出目录
var ouputType = flag.String("type", "go", "ouput type:go or csharp") //输出文件类型
```

### 类型
```
- 基础类型:和golang语法一样,例如:int8,uint32,float32,bool,string
- 数组类型：[]类型，例如:[]int32
- 字节流：bytes
- 任意类型：any
```