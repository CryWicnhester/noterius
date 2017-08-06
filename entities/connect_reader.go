package entities

import (
	"encoding/binary"
	"net"

	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/Nyarum/noterius/network"
)

type ConnectReader struct {
	Client       net.Conn
	PacketReader *actor.PID
	Network      network.INetwork
	Logger       *zap.SugaredLogger
}

func (state *ConnectReader) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case ReadPacket:
		if msg.Len == 0 {
			state.Client.Write([]byte{0x00, 0x02})
			return
		}

		uniqueID := binary.LittleEndian.Uint32(msg.Buf[0:4])
		opcode := binary.BigEndian.Uint16(msg.Buf[4:6])

		state.Logger.Debugw("Received a new packet", "len", msg.Len, "uniqueID", uniqueID, "opcode", opcode)

		if msg.Len >= 6 {
			packet, err := state.Network.Unmarshal(opcode, msg.Buf[6:])
			if err != nil {
				state.Logger.Errorw("Error unmarshal packet", "error", err)
				return
			}

			state.PacketReader.Tell(packet)
		}
	}
}
