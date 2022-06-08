package starkex

import (
	"math/big"
	"testing"

	"github.com/reddio-com/red-adapter/types"
)

func Test_Mint(t *testing.T) {
	type args struct {
		m *types.L2MintData
	}
	vaultId1, _ := new(big.Int).SetString("1654615998", 10)
	vaultId2, _ := new(big.Int).SetString("878225224", 10)
	vaultId3, _ := new(big.Int).SetString("2147483647", 10)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				m: &types.L2MintData{
					VaultId:  vaultId1,
					Amount:   "6",
					TokenId:  "0x400de4b5a92118719c78df48f4ff31e78de58575487ce1eaf19922ad9b8a714",
					StarkKey: "0x2deb04eb807be0ec943e08d8f666521edb3f12833922fbbc7f93e1434ae810e",
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				m: &types.L2MintData{
					VaultId:  vaultId2,
					Amount:   "4",
					TokenId:  "0x400203fded733e8b421eaeb534097cabaf3897a3e70f16a55485822de1b372a",
					StarkKey: "0x2f116d013fb6ecae90765a876a5bfcf66cd6a6be1f85c9841629cd0bd080ed3",
				},
			},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				m: &types.L2MintData{
					VaultId:  vaultId3,
					Amount:   "4",
					TokenId:  "0x400203fded733e8b421eaeb534097cabaf3897a3e70f16a55485822de1b372a",
					StarkKey: "0x2f116d013fb6ecae90765a876a5bfcf66cd6a6be1f85c9841629cd0bd080ed3",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Mint(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Mint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
