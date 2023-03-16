package Znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
)

type DataPack struct {
}

var DefaultDataPack *DataPack

func init()  {
	DefaultDataPack=&DataPack{}
}

func (p *DataPack) Pack(msg Zinterface.MessageI)([]byte,error) {
	buf:=bytes.NewBuffer([]byte{})
	//写dataLen
	//log.Printf("prepare to pack,message len and id are %d,%d,%x",msg.GetDataLen(),msg.GetMessageId(),msg.GetData())
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		//log.Println("[Pack] Message error:"+err.Error())
		return nil, err
	}
	//fmt.Println("Pack finished",buf)
	//写msgID
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMessageId()); err != nil {
		//log.Println("[Pack] Message error:"+err.Error())
		return nil, err
	}
	//fmt.Println("Pack finished",buf)
	//写data数据
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		//log.Println("[Pack] Message error:"+err.Error())
		return nil ,err
	}
	//fmt.Printf("Pack finished % x",buf.Bytes())
	return buf.Bytes(), nil

}

func (p *DataPack) GetHeadLen()uint32  {
	/*|DataLen**|MessageID|Data|*/
	/*|***4Byte*|**4Byte**|*/
	return 8
}

func (p *DataPack) UnPack(data []byte)(Zinterface.MessageI,error)  {
	reader:=bytes.NewReader(data)
	msg:=&Message{}
	if err:=binary.Read(reader,binary.LittleEndian,&msg.DataLen);err!=nil{
		//log.Println("[Pack]UnPack Message error:"+err.Error())
		return nil, err
	}
	if err:=binary.Read(reader,binary.LittleEndian,&msg.MessageId);err!=nil{
		//log.Println("[Pack]UnPack Message error:"+err.Error())
		return nil, err
	}
	if utils.GlobalConfig.MaxPackageSize>0&&utils.GlobalConfig.MaxPackageSize<int(msg.DataLen){
		return nil,errors.New("[Pack]pack length of data is too long")
	}
	//根据数据长度设立一个缓冲池获取数据
	buf:=make([]byte,msg.GetDataLen())
	if err:=binary.Read(reader,binary.LittleEndian,&buf);err!=nil{
		//log.Println("[Pack]UnPack Message error:"+err.Error())
		return nil, err
	}
	msg.SetData(buf)
	return msg,nil
}
