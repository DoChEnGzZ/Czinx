package Znet

/*|DataLen**|MessageID|Data|*/
/*|***4Byte*|**4Byte**|*/
type Message struct {
	Data []byte
	DataLen uint32
	MessageId uint32
}

func NewMessage(data []byte,id uint32)*Message  {
	return &Message{
		Data:      data,
		DataLen: uint32(len(data)),
		MessageId: id,
	}
}

func (m *Message) GetData()[]byte  {
	return m.Data
}
func (m *Message) GetMessageId()uint32  {
	return m.MessageId
}
func (m *Message) GetDataLen()uint32  {
	return m.DataLen
}
func (m *Message) SetData(data []byte)  {
	m.Data=data
}
func (m *Message) SetMessageId(id uint32)  {
	m.MessageId=id
}
func (m *Message) SetDataLen(len uint32)  {
	m.DataLen=len
}
