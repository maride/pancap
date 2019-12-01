package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

func Test_analyzePCAP(t *testing.T) {
	type args struct {
		source   *gopacket.PacketSource
		linkType layers.LinkType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Faulty link type",
			args: args{
				source: &gopacket.PacketSource{
					DecodeOptions: gopacket.DecodeOptions{
						Lazy:                     false,
						NoCopy:                   false,
						SkipDecodeRecovery:       false,
						DecodeStreamsAsDatagrams: false,
					},
				},
				linkType: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := analyzePCAP(tt.args.source, tt.args.linkType); (err != nil) != tt.wantErr {
				t.Errorf("analyzePCAP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}