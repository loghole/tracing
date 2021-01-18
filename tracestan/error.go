package tracestan

import (
	"errors"
)

var (
	ErrUnsubscribe           = errors.New("unsubscribe")
	ErrEmptyProtoMsg         = errors.New("empty protobuf proto.Message object")
	ErrInvalidProtoMsgEncode = errors.New("invalid protobuf proto.Message object passed to encode")
	ErrInvalidProtoMsgDecode = errors.New("invalid protobuf proto.Message object passed to decode")
	ErrInvalidArgs           = errors.New("invalid callback args")
	ErrInvalidResult         = errors.New("invalid callback result")
)
