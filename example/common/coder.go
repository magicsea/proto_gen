package common

import . "proto_gen_ext/output"

//TODO..处理any类型
func init() {
	SetAnyPacker(msgPacker)
}

var msgPacker = &AnyMsgPacker{}

type AnyMsgPacker struct {
}

func (p *AnyMsgPacker) Pack(msg interface{}, s IStream) error {

	return nil
}

func (p *AnyMsgPacker) UnPack(s IStream) (interface{}, error) {

	return nil, nil
}

//包大小
func (p *AnyMsgPacker) GetPackSize(msg interface{}) (int, bool) {
	var size = 1 //type
	ds, b := p.GetDataSize(msg)
	if !b {
		return 0, false
	}
	size += ds
	return size, true

}

//数据大小
func (p *AnyMsgPacker) GetDataSize(msg interface{}) (int, bool) {

	return 0, false
}
