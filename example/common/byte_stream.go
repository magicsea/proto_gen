package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// ByteStream 包装一些字节流的读写方法，方便二进制消息的序列和反序列
type ByteStream struct {
	data     []byte
	readPos  uint32
	writePos uint32
}

// NewByteStream 创建一个新的字节流
func NewByteStream(data []byte) *ByteStream {
	return &ByteStream{
		data:     data,
		readPos:  0,
		writePos: 0,
	}
}

func (bs *ByteStream) Reset() {
	bs.readPos = 0
	bs.writePos = 0
}

func (bs *ByteStream) readCheck(c uint32) error {
	if bs.data == nil {
		return errors.New("data is nil")
	}

	if bs.readPos+c > uint32(len(bs.data)) {
		return errors.New("no enough read space")
	}

	return nil
}

// ReadEnd 是否读完
// func (bs *ByteStream) ReadEnd() bool {
// 	return bs.readPos != uint32(len(bs.data))
// }

// ReadByte 读取一个字节
func (bs *ByteStream) ReadUint8() (byte, error) {
	if err := bs.readCheck(1); err != nil {
		return 0, err
	}

	v := bs.data[bs.readPos]
	bs.readPos = bs.readPos + 1

	return v, nil
}

// ReadBool 读bool
func (bs *ByteStream) ReadBool() (bool, error) {
	v, err := bs.ReadUint8()
	if err != nil {
		return false, err
	}

	return (v != 0), nil
}

// ReadInt8 读取一个int8
func (bs *ByteStream) ReadInt8() (int8, error) {
	v, err := bs.ReadUint8()
	return int8(v), err
}

// ReadInt16 读取一个int16
func (bs *ByteStream) ReadInt16() (int16, error) {
	v, err := bs.ReadUint16()
	return int16(v), err
}

// ReadInt32 读取一个int32
func (bs *ByteStream) ReadInt32() (int32, error) {
	v, err := bs.ReadUint32()
	return int32(v), err
}

// ReadInt64 读取一个int64
func (bs *ByteStream) ReadInt64() (int64, error) {
	v, err := bs.ReadUint64()
	return int64(v), err
}

// ReadUint16 读取一个Uint16
func (bs *ByteStream) ReadUint16() (uint16, error) {
	if err := bs.readCheck(2); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint16(bs.data[bs.readPos : bs.readPos+2])
	bs.readPos = bs.readPos + 2

	return v, nil
}

// ReadUint32 读取一个Int
func (bs *ByteStream) ReadUint32() (uint32, error) {
	if err := bs.readCheck(4); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint32(bs.data[bs.readPos : bs.readPos+4])
	bs.readPos = bs.readPos + 4

	return v, nil
}

// ReadUint64 读取一个Uint64
func (bs *ByteStream) ReadUint64() (uint64, error) {
	if err := bs.readCheck(8); err != nil {
		return 0, err
	}

	v := binary.LittleEndian.Uint64(bs.data[bs.readPos : bs.readPos+8])
	bs.readPos = bs.readPos + 8

	return v, nil
}

// ReadStr 读取一个string
func (bs *ByteStream) ReadString() (string, error) {

	len, err := bs.ReadUint16()
	if err != nil {
		return "", err
	}

	if err = bs.readCheck(uint32(len)); err != nil {
		return "", err
	}

	v := string(bs.data[bs.readPos : bs.readPos+uint32(len)])
	bs.readPos = bs.readPos + uint32(len)

	return v, nil
}

// ReadBytes 读取一个byte[]
func (bs *ByteStream) ReadBytes() ([]byte, error) {
	len, err := bs.ReadUint32()
	if err != nil {
		return nil, err
	}

	if len == 0 {
		return nil, nil
	}

	if err := bs.readCheck(uint32(len)); err != nil {
		return nil, err
	}

	b := make([]byte, len)

	copy(b, bs.data[bs.readPos:bs.readPos+uint32(len)])
	bs.readPos = bs.readPos + uint32(len)

	return b, nil
}

// ReadFloat32 读取float32
func (bs *ByteStream) ReadFloat32() (float32, error) {
	u, err := bs.ReadUint32()
	if err != nil {
		return 0, err
	}

	f := math.Float32frombits(u)
	return f, nil
}

// ReadFloat64 读取float64
func (bs *ByteStream) ReadFloat64() (float64, error) {
	u, err := bs.ReadUint64()
	if err != nil {
		return 0, err
	}

	f := math.Float64frombits(u)
	return f, nil
}

//TODO:auto grow
func (bs *ByteStream) writeCheck(c uint32) error {
	if bs.data == nil {
		return errors.New("data is nil")
	}

	if bs.writePos+c > uint32(len(bs.data)) {
		return errors.New("no enough write space")
	}

	return nil
}

// WriteByte 写字节
func (bs *ByteStream) WriteUint8(v byte) error {
	if err := bs.writeCheck(1); err != nil {
		return err
	}

	bs.data[bs.writePos] = v
	bs.writePos = bs.writePos + 1

	return nil
}

// WriteBool 写bool, 1代表true, 0代表false
func (bs *ByteStream) WriteBool(v bool) error {
	if v {
		return bs.WriteUint8(1)
	}
	return bs.WriteUint8(0)
}

// WriteInt8 写Int8
func (bs *ByteStream) WriteInt8(v int8) error {
	return bs.WriteUint8(byte(v))
}

// WriteInt16 写Int16
func (bs *ByteStream) WriteInt16(v int16) error {
	return bs.WriteUint16(uint16(v))
}

// WriteInt32 写Int32
func (bs *ByteStream) WriteInt32(v int32) error {
	return bs.WriteUint32(uint32(v))
}

// WriteInt64 写Int64
func (bs *ByteStream) WriteInt64(v int64) error {
	return bs.WriteUint64(uint64(v))
}

// WriteUint16 写Uint16
func (bs *ByteStream) WriteUint16(v uint16) error {
	if err := bs.writeCheck(2); err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(bs.data[bs.writePos:bs.writePos+2], v)
	bs.writePos = bs.writePos + 2

	return nil
}

// WriteUint32 写Uint32
func (bs *ByteStream) WriteUint32(v uint32) error {

	if err := bs.writeCheck(4); err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(bs.data[bs.writePos:bs.writePos+4], v)
	bs.writePos = bs.writePos + 4

	return nil
}

// WriteUint64 写Uint64
func (bs *ByteStream) WriteUint64(v uint64) error {

	if err := bs.writeCheck(8); err != nil {
		return err
	}

	binary.LittleEndian.PutUint64(bs.data[bs.writePos:bs.writePos+8], v)
	bs.writePos = bs.writePos + 8

	return nil
}

// WriteStr 写string
func (bs *ByteStream) WriteString(v string) error {

	if err := bs.writeCheck(uint32(len(v) + 2)); err != nil {
		return err
	}

	bs.WriteUint16(uint16(len(v)))

	if len(v) != 0 {
		copy(bs.data[bs.writePos:bs.writePos+uint32(len(v))], v)
		bs.writePos = bs.writePos + uint32(len(v))
	}

	return nil
}

// WriteBytes 写[]byte
func (bs *ByteStream) WriteBytes(v []byte) error {

	if v == nil {
		bs.WriteUint32(0)
		return nil
	}

	if err := bs.writeCheck(uint32(len(v) + 4)); err != nil {
		return err
	}

	bs.WriteUint32(uint32(len(v)))

	copy(bs.data[bs.writePos:bs.writePos+uint32(len(v))], v)

	bs.writePos = bs.writePos + uint32(len(v))
	return nil
}

// WriteFloat32 写Float32
func (bs *ByteStream) WriteFloat32(f float32) error {
	u := math.Float32bits(f)
	return bs.WriteUint32(u)
}

// WriteFloat64 写Float64
func (bs *ByteStream) WriteFloat64(f float64) error {
	u := math.Float64bits(f)
	return bs.WriteUint64(u)
}

// CalcSize 计算序列化所需长度
func CalcSizeReflect(content interface{}) int {
	size := 0

	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		switch field.(type) {
		case uint8, int8, bool:
			size++
		case uint16, int16:
			size += 2
		case uint32, int32, float32:
			size += 4
		case uint64, int64, float64:
			size += 8
		case string:
			size += 2
			size += len(field.(string))
		case []byte:
			size += 4
			size += len(field.([]byte))
		default:
			panic(fmt.Sprintf("不支持的类型 %+v", field))
		}
	}

	return size
}

// Marshal 序列化content
func (bs *ByteStream) MarshalReflect(content interface{}) error {
	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		var err error
		switch field.(type) {
		case bool:
			err = bs.WriteBool(v.Field(i).Interface().(bool))
		case uint8:
			err = bs.WriteUint8(v.Field(i).Interface().(uint8))
		case uint16:
			err = bs.WriteUint16(v.Field(i).Interface().(uint16))
		case uint32:
			err = bs.WriteUint32(v.Field(i).Interface().(uint32))
		case uint64:
			err = bs.WriteUint64(v.Field(i).Interface().(uint64))
		case string:
			err = bs.WriteString(v.Field(i).Interface().(string))
		case float32:
			err = bs.WriteFloat32(v.Field(i).Interface().(float32))
		case float64:
			err = bs.WriteFloat64(v.Field(i).Interface().(float64))
		case []uint8:
			err = bs.WriteBytes(v.Field(i).Interface().([]byte))
		default:
			panic(fmt.Sprintf("不支持的类型 %+v", field))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal 反序列化
func (bs *ByteStream) UnmarshalReflect(content interface{}) error {
	v := reflect.ValueOf(content).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()

		var err error
		var value interface{}
		switch field.(type) {
		case bool:
			value, err = bs.ReadBool()
		case uint8:
			value, err = bs.ReadUint8()
		case uint16:
			value, err = bs.ReadUint16()
		case uint32:
			value, err = bs.ReadUint32()
		case uint64:
			value, err = bs.ReadUint64()
		case string:
			value, err = bs.ReadString()
		case float32:
			value, err = bs.ReadFloat32()
		case float64:
			value, err = bs.ReadFloat64()
		case []byte:
			value, err = bs.ReadBytes()
		default:
			panic(fmt.Sprintf("不支持的类型 %+v", field))
		}

		if err != nil {
			return err
		}
		v.Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}

// GetUsedSlice 获取已经写入的部分Slice
func (bs *ByteStream) GetUsedSlice() []byte {
	if bs.data == nil || bs.writePos == 0 {
		return nil
	}

	return bs.data[0:bs.writePos]
}

func (bs *ByteStream) GetRWPos() (int, int) {
	return int(bs.readPos), int(bs.writePos)
}
