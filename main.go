package main

import (
	"huangjihui511/p2p_discover/entity/client"
)

func main() {
	rootC := client.NewClientDiscover(10)
	rootC.Start(0)

	cn := make([]*client.ClientDiscover, 0)
	for i := 0; i < 20; i++ {
		c1 := client.NewClientDiscover(10)
		c1.Start(rootC.IP)
		cn = append(cn, c1)
		// c1.Ping(rootC.MyNodeId)
	}
	// for i := 0; i < 20; i++ {
	// 	cn[0].Ping(cn[i].MyNodeId)
	// }
	// for i := 0; i < 20; i++ {
	// 	cn[10].Ping(cn[i].MyNodeId)
	// }
	// for i := 0; i < 20; i++ {
	// 	cn[19].Ping(cn[i].MyNodeId)
	// }
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			cn[i].Ping(cn[j].MyNodeId)
		}
	}
	select {}
}
