/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : interface.go
#   Created       : 2019-04-29 11:07
#   Last Modified : 2019-04-29 11:07
#   Describe      :
#
# ====================================================*/
package election

import (
	"context"
	"time"
)

// OptType key change type
type OptType int

const (
	// OptionTypeNew event type new
	OptionTypeNew OptType = iota
	// OptionTypeUpdate event type update
	OptionTypeUpdate
	// OptionTypeDelete event type delete
	OptionTypeDelete
)

// KVPair encapsulate the results of the query
type KVPair struct {
	Key       string
	Value     []byte
	LastIndex int64
}

// WatchRes encapsulate the values that change when you watch
type WatchRes struct {
	KV   *KVPair
	Type OptType
}

// WriteOptions configuration when writing values
type WriteOptions struct {
	// expire time of key
	TTL time.Duration
	// Key does not disappear until the connection is disconnected
	KeepAlive bool
}

// Elector use set-watch mechanism to realize main-subordinate election
type Elector interface {
	// Watch watch what happens to key
	Watch(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan *WatchRes, error)
	// Put change the key-value
	SetNx(ctx context.Context, key string, value []byte, opt *WriteOptions) error

	// Delete key-value
	Delete(ctx context.Context, key string) error

	// HeartBeat checks if the connection is available
	HeartBeat() error
}
