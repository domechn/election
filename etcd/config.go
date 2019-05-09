/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : config.go
#   Created       : 2019-04-29 10:54
#   Last Modified : 2019-04-29 10:54
#   Describe      :
#
# ====================================================*/
package etcd

import (
	"time"
)

// Config contains the options for a storage client
type Config struct {
	ConnectionTimeout time.Duration
	Username          string
	Password          string

	CAFile string

	CertFile string
	KeyFile  string
}

// ClientTLSConfig contains data for a Client TLS configuration in the form
// the etcd client wants it.  Eventually we'll adapt it for ZK and Consul.
type ClientTLSConfig struct {
	CertFile   string
	KeyFile    string
	CACertFile string
}
