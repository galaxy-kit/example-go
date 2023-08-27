package main

import (
	"context"
	"fmt"
	"kit.golaxy.org/plugins/gtp"
	"kit.golaxy.org/plugins/gtp_client"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing endpoint")
	}

	cli, err := gtp_client.Connect(context.Background(), os.Args[1],
		gtp_client.Option{}.RecvDataHandlers(func(client *gtp_client.Client, data []byte) error {
			fmt.Println(string(data))
			return nil
		}),
		gtp_client.Option{}.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_XChaCha20,
			BlockCipherMode:     gtp.BlockCipherMode_None,
			PaddingMode:         gtp.PaddingMode_None,
			MACHash:             gtp.Hash_Fnv1a64,
		}),
		gtp_client.Option{}.CompressedSize(128),
		gtp_client.Option{}.IOTimeout(3*time.Second),
		gtp_client.Option{}.IOBufferCap(1024*1024*5),
		gtp_client.Option{}.AutoReconnect(true),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Println("this console is", cli.GetSessionId())

	for {
		var text string
		fmt.Scanln(&text)
		if err := cli.SendData([]byte(text)); err != nil {
			fmt.Println("send data err:", err)
		}
	}
}
