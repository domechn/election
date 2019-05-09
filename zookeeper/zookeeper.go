/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : zookeeper.go
#   Created       : 2019-05-07 19:28
#   Last Modified : 2019-05-07 19:28
#   Describe      :
#
# ====================================================*/
package zookeeper

import (
	"bytes"
	"context"
	"strings"
	"time"

	"goelection"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	// SOH control character
	SOH = "\x01"
)

type Election struct {
	client *zk.Conn
}

// New init a election based on zk
// addrs => host:port
func New(addrs []string, cfg *Config) (*Election, error) {
	// If the client is disconnected, zk deletes the temporary node created by the client after timeout
	timeout := time.Second

	if cfg != nil {
		if cfg.sessionTimeout > 0 {
			timeout = cfg.sessionTimeout
		}
	}

	conn, _, err := zk.Connect(addrs, timeout)
	if err != nil {
		return nil, err
	}
	return &Election{
		client: conn,
	}, nil
}

func (s *Election) get(key string) (*election.KVPair, error) {
	resp, meta, err := s.client.Get(key)

	if err != nil {
		if err == zk.ErrNoNode {
			return nil, election.ErrKeyIsNotExist
		}
		return nil, err
	}

	// conflict
	if string(resp) == SOH {
		return s.get(key)
	}
	pair := &election.KVPair{
		Key:       key,
		Value:     resp,
		LastIndex: int64(meta.Version),
	}
	return pair, nil
}

func (s *Election) Delete(ctx context.Context, key string) error {
	err := s.client.Delete(key, -1)
	if err == zk.ErrNoNode {
		return election.ErrKeyIsNotExist
	}
	return err
}

// Watch watch the key
func (s *Election) Watch(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan *election.WatchRes, error) {
	watchCh := make(chan *election.WatchRes)

	go func() {
		defer close(watchCh)

		for {
			_, _, eventCh, err := s.client.ExistsW(key)
			if err != nil {
				return
			}
			select {
			case <-stopCh:
				return
			case e := <-eventCh:
				switch e.Type {
				case zk.EventNodeCreated, zk.EventNodeDataChanged:
					tp := election.OptionTypeUpdate
					if e.Type == zk.EventNodeCreated {
						tp = election.OptionTypeNew
					}
					if entry, err := s.get(key); err == nil {
						// FIXME: sometimes zk can not get the value
						if bytes.Compare(entry.Value, []byte{}) == 0 {
							if entry, err = s.get(key); err == nil {
								watchCh <- &election.WatchRes{
									KV:   entry,
									Type: tp,
								}
							}
						}
					}
				case zk.EventNodeDeleted:
					watchCh <- &election.WatchRes{
						Type: election.OptionTypeDelete,
					}
				}
			}
		}
	}()

	return watchCh, nil
}

func (s *Election) createFullPath(path []string, ephemeral bool) error {
	for i := 1; i <= len(path); i++ {
		newpath := "/" + strings.Join(path[:i], "/")
		if i == len(path) && ephemeral {
			_, err := s.client.Create(newpath, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			return err
		}
		_, err := s.client.Create(newpath, []byte{}, 1, zk.WorldACL(zk.PermAll))
		if err != nil {
			// Skip if node already exists
			if err != zk.ErrNodeExists {
				return err
			}
		}
	}
	return nil
}

// SetNx set the key if key is not exist
// because zk is special, you cannot set TTL for non-persistent nodes
func (s *Election) SetNx(ctx context.Context, key string, value []byte, opt *election.WriteOptions) error {
	flag := 0

	exsits, _, err := s.client.Exists(key)

	if err != nil {
		return err
	}

	if opt != nil {
		if opt.KeepAlive {
			flag = zk.FlagEphemeral
		}
	}

	if !exsits {
		fullPath := strings.Split(strings.TrimPrefix(strings.TrimSuffix(key, "/"), "/"), "/")
		if flag == 0 {
			err = s.createFullPath(fullPath, false)
		} else {
			err = s.createFullPath(fullPath, true)
		}
	} else {
		return election.ErrKeyIsExist
	}
	if err == zk.ErrNodeExists {
		return election.ErrKeyIsExist
	}

	_, err = s.client.Set(key, value, -1)
	return err
}

// HeartBeat check the connection weather is active
func (s *Election) HeartBeat() error {
	if s.client.State() != zk.StateHasSession {
		return election.ErrHeartBeatFail
	}
	return nil
}

// Close shutdown the connection
func (s *Election) Close() error {
	s.client.Close()
	return nil
}
