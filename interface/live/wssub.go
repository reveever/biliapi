package live

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	WsBodyProtocolVersionNormal = 0
	WsBodyProtocolVersionZlib   = 2
	WsBodyProtocolVersionBrotli = 3

	WsOpHeartbeat          = 2
	WsOpHeartbeatReply     = 3
	WsOpMessage            = 5
	WsOpUserAuthentication = 7
	WsOpConnectSuccess     = 8

	wsSubPktHdrLen = 16
)

type WsSubPkt []byte

func NewWsSubPkt(ver uint16, op uint32, payload []byte) WsSubPkt {
	length := uint32(wsSubPktHdrLen + len(payload))
	pkt := make([]byte, length)
	binary.BigEndian.PutUint32(pkt[0:4], length)
	binary.BigEndian.PutUint16(pkt[4:6], wsSubPktHdrLen)
	binary.BigEndian.PutUint16(pkt[6:8], ver)
	binary.BigEndian.PutUint32(pkt[8:12], op)
	// binary.BigEndian.PutUint32(v[12:16], seq)
	copy(pkt[wsSubPktHdrLen:], payload)
	return pkt
}

func ReadWsSubPkt(r io.Reader) (WsSubPkt, error) {
	pkt := make(WsSubPkt, wsSubPktHdrLen)
	if n, err := r.Read(pkt); err != nil {
		return nil, fmt.Errorf("read pkt len: %v", err)
	} else if n != wsSubPktHdrLen {
		return nil, fmt.Errorf("insufficient read")
	} else if pkt.HdrLen() != wsSubPktHdrLen {
		return nil, fmt.Errorf("unexpected header len: %d", pkt.HdrLen())
	}

	toRead := int(pkt.PktLen() - wsSubPktHdrLen)
	pkt = append(pkt, make([]byte, toRead)...)
	if n, err := r.Read(pkt[wsSubPktHdrLen:]); err != nil {
		return nil, fmt.Errorf("read ws sub pkt: %v", err)
	} else if n != toRead {
		return nil, fmt.Errorf("insufficient read")
	}
	return pkt, nil
}

func (m WsSubPkt) PktLen() uint32 {
	return binary.BigEndian.Uint32(m[0:4])
}

func (m WsSubPkt) HdrLen() uint16 {
	return binary.BigEndian.Uint16(m[4:6])
}

func (m WsSubPkt) Version() uint16 {
	return binary.BigEndian.Uint16(m[6:8])
}

func (m WsSubPkt) Operation() uint32 {
	return binary.BigEndian.Uint32(m[8:12])
}

func (m WsSubPkt) Sequence() uint32 {
	return binary.BigEndian.Uint32(m[12:16])
}

func (m WsSubPkt) Body() []byte {
	if len(m) <= wsSubPktHdrLen {
		return nil
	}
	body := make([]byte, len(m)-wsSubPktHdrLen)
	copy(body, m[wsSubPktHdrLen:])
	return body
}
