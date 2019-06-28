/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : common.go
#   Created       : 2019-04-29 12:02
#   Last Modified : 2019-04-29 12:02
#   Describe      :
#
# ====================================================*/
package election

import (
	"errors"
	"time"
)

const (
	// write timeout
	defaultTimeout     = time.Second
	heartbeatInterval  = time.Second * 2
	leaderSelectionKey = "/api-gateway-controller/leader/selection"
	// DefaultDialTimeout dial timeout
	DefaultDialTimeout = time.Second * 10
)

var (
	// ErrKeyIsExist key exists
	ErrKeyIsExist = errors.New("selection: key is already exist")

	// ErrKeyIsNotExist key is not exist at the beginning
	ErrKeyIsNotExist = errors.New("selection watch: key is not exist")

	// ErrHeartBeatFail connection is not available
	ErrHeartBeatFail = errors.New("heart beat failed")

	// ErrIsNotLeader node is not leader
	ErrIsNotLeader = errors.New("selection: is not leader")
)
