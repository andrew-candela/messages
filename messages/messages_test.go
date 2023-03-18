package messages

import (
	"testing"
)

func TestEncode(t *testing.T) {
	pack := Packet{
		SenderName: "Andrew",
		Content:    "Hello",
	}
	pack_bytes := pack.ToBytes()
	new_packet := FromBytes(pack_bytes)
	if new_packet != pack {
		t.Errorf("Oops!")
	}
}
