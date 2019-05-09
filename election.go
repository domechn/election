/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : selection.go
#   Created       : 2019-04-29 11:06
#   Last Modified : 2019-04-29 11:06
#   Describe      :
#
# ====================================================*/
package election

import (
	"context"
	"time"
)

// Election automatic voting
type Election struct {
	sw        Elector
	id        string
	key       string
	isLeader  bool
	succeedCh chan string
}

// New selection
func New(sw Elector, id string) *Election {
	return &Election{
		sw:        sw,
		id:        id,
		key:       leaderSelectionKey,
		succeedCh: make(chan string),
	}
}

// Elect monitor election results. If any new process is successful, all processes will be notified through chan
// will not block
func (s *Election) Elect(stopCh <-chan struct{}) (chan string, error) {
	watchRes, err := s.sw.Watch(context.Background(), s.key, stopCh)
	if err != nil {
		return nil, err
	}

	go s.doHeartBeat()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			case wr := <-watchRes:
				// fix nil pointer when server is close
				if wr == nil {
					continue
				}
				switch wr.Type {
				case OptionTypeDelete:
					s.doSelect()
				case OptionTypeNew, OptionTypeUpdate:
					newID := string(wr.KV.Value)
					s.succeedCh <- newID
				}
			}
		}
	}()

	return s.succeedCh, nil
}

// Resign give up the leader status
func (s *Election) Resign() error {
	if !s.isLeader {
		return ErrIsNotLeader
	}
	err := s.sw.Delete(context.Background(), s.key)
	if err == nil {
		s.isLeader = false
	}
	return err
}

func (s *Election) doSelect() error {
	wo := &WriteOptions{
		TTL:       defaultTimeout,
		KeepAlive: true,
	}
	err := s.sw.SetNx(context.Background(), leaderSelectionKey, []byte(s.id), wo)
	if err == nil {
		s.isLeader = true
	}
	return err
}

func (s *Election) doHeartBeat() {
	t := time.NewTicker(heartbeatInterval)
	for {
		<-t.C
		if err := s.sw.HeartBeat(); err != nil {
			s.succeedCh <- s.id
		}
	}
}

// DoSelector put value to store, if successful return nil
func (s *Election) DoSelect() error {
	return s.doSelect()
}

func (s *Election) ID() string {
	return s.id
}
