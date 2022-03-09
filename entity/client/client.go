package client

import (
	"encoding/json"
	"huangjihui511/p2p_discover/entity/manager"
	"huangjihui511/p2p_discover/entity/network_define"
)

type ClientDiscover struct {
	MyNodeId   network_define.NodeId8
	IP         network_define.IP
	Connection *manager.Connection
	Size       int
	ConnCmd    chan manager.ConnectionCmdWithData
}

func NewClientDiscover(size int) *ClientDiscover {
	return &ClientDiscover{
		MyNodeId: network_define.NewNodeId8(),
		Size:     size,
		IP:       network_define.NewIP(),
	}
}

func (self *ClientDiscover) Start(rootIp network_define.IP) {
	self.ConnCmd = make(chan manager.ConnectionCmdWithData)
	self.Connection = &manager.Connection{
		MyNodeId:      self.MyNodeId,
		NodeIdManager: manager.NewNodeIdManager(self.Size),
		IP:            self.IP,
		Cmd:           make(chan *manager.ConnectionCmdWithData),
	}
	self.Connection.Run(rootIp)
}

func (self *ClientDiscover) End() {
	self.Connection.Execute(&manager.ConnectionCmdWithData{
		Cmd: manager.ConnectionCmdEnd,
	})
}

func (self *ClientDiscover) Ping(nodeId network_define.NodeId8) {
	bytes, _ := json.Marshal(nodeId)
	self.Connection.Execute(&manager.ConnectionCmdWithData{
		Cmd:  manager.ConnectionCmdPing,
		Data: bytes,
	})
}
