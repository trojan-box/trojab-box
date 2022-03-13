package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/vo"
	"testing"
)

func TestAuthUseCase_ValidLoginSignature(t *testing.T) {
	type args struct {
		req vo.UserAuthReq
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{req: vo.UserAuthReq{
				Address:   "0x9866c5dB69592ade1F1179D5110fDc7Be46F6f25",
				Timestamp: "1641872317",
				Nonce:     "98768454",
				SignedMsg: "0xf342127e9a5b9d4eb518f62e912d0f88087931cf2ba48b0509064b1f4501b426264f8d905e1502dbdc27a57a1be53fd84bc074986b3b58f2fded5ce33a116ba001",
			}},
			want:    true,
			wantErr: false,
		},
		{
			name: "incorrect address",
			args: args{req: vo.UserAuthReq{
				Address:   "0x9866c5dB69592ade1F1179D5110fDc7Be46F6f25",
				Timestamp: "1641872317",
				Nonce:     "98768454",
				SignedMsg: "0xf342127e9a5b9d4eb518f62e912d0f88087931cf2ba48b0509064b1f4501b426264f8d905e1502dbdc27a57a1be53fd84bc074986b3b58f2fded5ce33a116ba001",
			}},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &AuthUseCase{}
			got, err := u.VerifyLoginSignature(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidLoginSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidLoginSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}
