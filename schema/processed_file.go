package schema

import (
	"bytes"
	"log"
	"strings"
)

//ProcessedFile collect data in payload and return it in two formats
type ProcessedFile struct {
	FileName   string
	DataSchema *Table

	payload []map[string]interface{}
}

//GetPayload return payload as is
func (pf ProcessedFile) GetPayload() []map[string]interface{} {
	return pf.payload
}

//GetPayloadBytes return marshaling by marshaller func, joined with \n,  bytes
//assume that payload can't be empty
func (pf ProcessedFile) GetPayloadBytes(marshaller Marshaller) []byte {
	var buf *bytes.Buffer

	var fields []string
	//for csv writers using || delimiter
	if marshaller.NeedHeader() {
		fields = pf.DataSchema.Columns.Header()
		buf = bytes.NewBuffer([]byte(strings.Join(fields, "||")))
	}

	for _, object := range pf.payload {
		objectBytes, err := marshaller.Marshal(fields, object)
		if err != nil {
			log.Println("Error marshaling object in processed file:", err)
		} else {
			if buf == nil {
				buf = bytes.NewBuffer(objectBytes)
			} else {
				buf.Write([]byte("\n"))
				buf.Write(objectBytes)
			}
		}
	}

	return buf.Bytes()
}