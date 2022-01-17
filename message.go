package sockunx

import (
	"fmt"
	"reflect"
	"sockunx/pkg/helper"
)

// Message struct
type Message struct {
	Data   []byte `json:"data"`
	Length int    `json:"length"`
}

func (o Message) String() string {
	return helper.ToJSONString(o)
}

// Bytes of message
func (o Message) Bytes() []byte {
	return []byte(helper.ToJSONString(o))
}

// ParseMessage from bytes
func ParseMessage(data interface{}) (*Message, error) {
	var e error
	var result Message
	if v, ok := data.([]byte); ok {
		e = helper.FromJSON(v, &result)
	} else if v, ok := data.(string); ok {
		e = helper.FromJSONString(v, &result)
	} else {
		e = fmt.Errorf("invalid data type %s", reflect.TypeOf(data))
	}
	return &result, e
}
