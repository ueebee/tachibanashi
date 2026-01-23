package event

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestDecodeFrameBasic(t *testing.T) {
	raw := "p_no\x021\x01p_date\x022018.12.03-13:11:22.122\x01p_cmd\x02ST\x01p_X\x02a\x03b\x03"
	frame, err := DecodeFrame([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeFrame() error = %v", err)
	}
	if frame.Command != CommandST {
		t.Fatalf("command = %s", frame.Command)
	}
	if frame.No != 1 {
		t.Fatalf("no = %d", frame.No)
	}
	if frame.Date != "2018.12.03-13:11:22.122" {
		t.Fatalf("date = %s", frame.Date)
	}
	values := frame.Values("p_X")
	if len(values) != 3 || values[0] != "a" || values[1] != "b" || values[2] != "" {
		t.Fatalf("values = %#v", values)
	}
}

func TestDecodeEventRequiresCommand(t *testing.T) {
	_, err := DecodeEvent([]byte("p_no\x021"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeEventST(t *testing.T) {
	raw := "p_no\x02208\x01p_date\x022018.12.03-13:11:22.122\x01p_errno\x022\x01p_err\x02session inactive.\x01p_cmd\x02ST"
	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	st, ok := event.(ST)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if st.ErrNo != "2" {
		t.Fatalf("p_errno = %s", st.ErrNo)
	}
	if st.Err != "session inactive." {
		t.Fatalf("p_err = %s", st.Err)
	}
}

func TestDecodeEventKP(t *testing.T) {
	raw := "p_no\x0220\x01p_date\x022018.12.03-11:34:59.138\x01p_cmd\x02KP"
	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	if _, ok := event.(KP); !ok {
		t.Fatalf("event type mismatch")
	}
}

func TestDecodeEventBase64(t *testing.T) {
	titleBytes := []byte{0x83, 0x65, 0x83, 0x58, 0x83, 0x67}
	encoded := base64.StdEncoding.EncodeToString(titleBytes)

	want, err := decodeShiftJIS(titleBytes)
	if err != nil {
		t.Fatalf("decodeShiftJIS() error = %v", err)
	}

	raw := fmt.Sprintf("p_no\x021\x01p_date\x022018.12.03-13:11:22.122\x01p_cmd\x02NS\x01p_HDL\x02%s\x01p_TX\x02%s", encoded, encoded)
	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}

	frame, ok := event.(Frame)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if got := frame.Value("p_HDL"); got != want {
		t.Fatalf("p_HDL = %s", got)
	}
	if got := frame.Value("p_TX"); got != want {
		t.Fatalf("p_TX = %s", got)
	}
}
