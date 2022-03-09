package manager

import (
	"encoding/json"
	"fmt"
	"huangjihui511/p2p_discover/entity/network_define"
	"sort"
)

type Connection struct {
	MyNodeId      network_define.NodeId8
	IP            network_define.IP
	NodeIdManager *NodeIdManager
	Cmd           chan *ConnectionCmdWithData
}

type ConnectionCmd int16

const (
	ConnectionCmdEnd ConnectionCmd = iota
	ConnectionCmdPing
	ConnectionCmdRegister
)

type ConnectionCmdWithData struct {
	Cmd      ConnectionCmd
	Data     []byte
	BackData []byte
	BackSync chan bool
}

type IPAndNodeId struct {
	MyNodeId network_define.NodeId8
	IP       network_define.IP
}

func (self *Connection) Name() string {
	return fmt.Sprintf("conn[ip=%v,node=%v]", self.IP, self.MyNodeId)
}

func (self *Connection) Run(rootIp network_define.IP) {
	fmt.Printf("%v starting\n", self.Name())
	SingleIpManager().SaveConnection(self.IP, self)
	rootConn, ok := SingleIpManager().GetConnection(rootIp)
	if ok {
		bytes, _ := json.Marshal(IPAndNodeId{
			MyNodeId: self.MyNodeId,
			IP:       self.IP,
		})
		rCmd := &ConnectionCmdWithData{
			Cmd:  ConnectionCmdRegister,
			Data: bytes,
		}
		rootConn.Execute(rCmd)
		backData := rCmd.BackData
		iPAndNodeIds := make([]IPAndNodeId, 0)
		_ = json.Unmarshal(backData, &iPAndNodeIds)
		self.LoadIpAndIds(iPAndNodeIds)
		self.NodeIdManager.TrySaveNodeId(self.MyNodeId, rootConn.MyNodeId, rootConn.IP)
	} else {
		fmt.Printf("%v register failed: conn not found root_ip=%v \n", self.Name(), rootIp)
	}

	go func() {
		defer func() {
			fmt.Printf("%v closed\n", self.Name())
			SingleIpManager().DelConnection(self.IP)
		}()
		for {
			select {
			case cmd := <-self.Cmd:
				{
					switch cmd.Cmd {
					case ConnectionCmdEnd:
						// cmd.BackSync <- true
						return
					case ConnectionCmdPing:
						self.Ping(cmd)
					case ConnectionCmdRegister:
						self.Register(cmd)
					}
				}
				cmd.BackSync <- true
			}
		}
	}()
}

func (self *Connection) Execute(cmd *ConnectionCmdWithData) {
	cmd.BackSync = make(chan bool)
	select {
	case self.Cmd <- cmd:
		<-cmd.BackSync
	}

	return
}

func (self *Connection) Register(cmd *ConnectionCmdWithData) {
	var c IPAndNodeId
	json.Unmarshal(cmd.Data, &c)
	_ = self.NodeIdManager.TrySaveNodeId(self.MyNodeId, c.MyNodeId, c.IP)
	fmt.Printf("%v Register node_id=%v\n", self.Name(), c.MyNodeId)
	fmt.Printf("%v Register %v ids\n", self.Name(), len(self.NodeIdManager.NodeIds))
	iPAndNodeIds := make([]IPAndNodeId, 0)
	for nodeId, ip := range self.NodeIdManager.NodeIds {
		iPAndNodeIds = append(iPAndNodeIds, IPAndNodeId{
			IP:       ip,
			MyNodeId: nodeId,
		})
	}
	iPAndNodeIds = append(iPAndNodeIds, IPAndNodeId{
		IP:       self.IP,
		MyNodeId: self.MyNodeId,
	})
	bytes, _ := json.Marshal(iPAndNodeIds)
	cmd.BackData = bytes
}

func (self *Connection) LoadIpAndIds(iPAndNodeIds []IPAndNodeId) {
	for _, v := range iPAndNodeIds {
		_ = self.NodeIdManager.TrySaveNodeId(self.MyNodeId, v.MyNodeId, v.IP)
	}
	fmt.Printf("%v load %v ids\n", self.Name(), len(self.NodeIdManager.NodeIds))
}

func (self *Connection) Ping(cmd *ConnectionCmdWithData) {
	var nodeId network_define.NodeId8
	json.Unmarshal(cmd.Data, &nodeId)
	if ip, ok := self.NodeIdManager.NodeIds[nodeId]; !ok {
		iPAndNodeIds := make([]IPAndNodeId, 0)
		for nodeId, ip := range self.NodeIdManager.NodeIds {
			iPAndNodeIds = append(iPAndNodeIds, IPAndNodeId{
				IP:       ip,
				MyNodeId: nodeId,
			})
		}
		sort.Slice(iPAndNodeIds, func(i, j int) bool {
			disI, _ := iPAndNodeIds[i].MyNodeId.GetDistance(nodeId)
			disJ, _ := iPAndNodeIds[j].MyNodeId.GetDistance(nodeId)
			return disI < disJ
		})
		for _, v := range iPAndNodeIds {
			conn, ok := SingleIpManager().GetConnection(v.IP)
			if !ok {
				continue
			}
			bytes, _ := json.Marshal(nodeId)
			cmdNew := &ConnectionCmdWithData{
				Cmd:  ConnectionCmdPing,
				Data: bytes,
			}
			conn.Execute(cmdNew)
			if cmdNew.BackData == nil {
				continue
			}
			var c IPAndNodeId
			json.Unmarshal(cmdNew.BackData, &c)
			self.NodeIdManager.TrySaveNodeId(self.MyNodeId, c.MyNodeId, c.IP)
			cmd.BackData = cmdNew.BackData
			fmt.Printf("%v ping success remote node_id=%v\n", self.Name(), nodeId)
		}
		fmt.Printf("%v ping not found node_id=%v\n", self.Name(), nodeId)
		return
	} else {
		if _, okIp := SingleIpManager().GetConnection(ip); !okIp {
			fmt.Printf("%v ping not found ip=%v\n", self.Name(), ip)
		}
		fmt.Printf("%v ping success node_id=%v\n", self.Name(), nodeId)
		result := IPAndNodeId{
			IP:       ip,
			MyNodeId: nodeId,
		}
		bytes, _ := json.Marshal(result)
		cmd.BackData = bytes
	}
}
