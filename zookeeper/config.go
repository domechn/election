/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : config.go
#   Created       : 2019-05-07 19:33
#   Last Modified : 2019-05-07 19:33
#   Describe      :
#
# ====================================================*/
package zookeeper

import (
	"time"
)

// Config is the config use to connection zk
type Config struct {
	sessionTimeout time.Duration
}
