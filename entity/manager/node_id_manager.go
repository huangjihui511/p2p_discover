package manager

import "huangjihui511/p2p_discover/entity/network_define"

type NodeIdManager struct {
	NodeIds map[network_define.NodeId8]network_define.IP
	Size    int
}

func NewNodeIdManager(size int) *NodeIdManager {
	return &NodeIdManager{
		NodeIds: make(map[network_define.NodeId8]network_define.IP),
		Size:    size,
	}
}

func (self *NodeIdManager) TrySaveNodeId(p, nodeId network_define.NodeId8, ip network_define.IP) error {
	if len(self.NodeIds) < self.Size {
		self.NodeIds[nodeId] = ip
		return nil
	}
	far, _, err := self.getFarNearNodeId(p)
	if err != nil {
		return err
	}
	delete(self.NodeIds, far)
	self.NodeIds[nodeId] = ip
	return nil
}

func (self *NodeIdManager) getFarNearNodeId(p network_define.NodeId8) (far, near network_define.NodeId8, err error) {
	farDistance, nearDistance := int64(0), int64(9999999)
	for nodeId := range self.NodeIds {
		distance, err := nodeId.GetDistance(p)
		if err != nil {
			return far, near, err
		}
		if distance > farDistance {
			far = nodeId
		}
		if distance < nearDistance {
			near = nodeId
		}
	}
	return far, near, nil
}
