package nats

//func RequestError(status int, msg *nats.Msg) {
//	out, err := proto.Marshal(&bookpb.ErrorResponse{
//		StatusCode: int32(status),
//	})
//	if err != nil {
//		_ = msg.Nak()
//	}
//
//	_ = msg.Respond(out)
//}
