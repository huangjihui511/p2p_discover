package manager

import (
	"huangjihui511/p2p_discover/entity/network_define"
	"sync"
)

type IpManager struct {
	ips *sync.Map
}

var singleIpManager *IpManager
var once sync.Once

func SingleIpManager() *IpManager {
	once.Do(func() {
		singleIpManager = &IpManager{
			ips: &sync.Map{},
		}
	})
	return singleIpManager
}

func (self *IpManager) GetConnection(ip network_define.IP) (*Connection, bool) {
	conn, ok := self.ips.Load(ip)
	if !ok {
		return nil, ok
	}
	return conn.(*Connection), true
}

func (self *IpManager) SaveConnection(ip network_define.IP, conn *Connection) {
	self.ips.Store(ip, conn)
}

func (self *IpManager) DelConnection(ip network_define.IP) {
	self.ips.Delete(ip)
}
