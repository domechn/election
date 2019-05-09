/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : etcd.go
#   Created       : 2019-04-29 10:52
#   Last Modified : 2019-04-29 10:52
#   Describe      :
#
# ====================================================*/
package etcd

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"goelection"
	tlsc "goelection/pkg/tls"
	"google.golang.org/grpc/connectivity"
)

const periodicSync = 5 * time.Minute // 自动同步间隔

// Election etcd client
type Election struct {
	cli *clientv3.Client
}

// New Create a new etcd client with the given address list and TLS configuration
func New(addrs []string, options *Config) (*Election, error) {
	s := &Election{}

	var (
		entries []string
		err     error
	)

	entries = createEndpoints(addrs, "http")

	cfg := &clientv3.Config{
		Endpoints:        entries,
		DialTimeout:      election.DefaultDialTimeout,
		AutoSyncInterval: periodicSync,
	}
	// set options
	if options != nil {
		if options.CertFile != "" || options.KeyFile != "" {
			tlsConfig, err := tlsc.Config(tlsc.Options{
				CAFile:   options.CAFile,
				CertFile: options.CertFile,
				KeyFile:  options.KeyFile,
			})
			if err != nil {
				return nil, err
			}
			setTLS(cfg, tlsConfig, addrs)
		}
		if options.ConnectionTimeout != 0 {
			setTimeout(cfg, options.ConnectionTimeout)
		}
		if options.Username != "" {
			setCredentials(cfg, options.Username, options.Password)
		}
	}

	c, err := clientv3.New(*cfg)
	if err != nil {
		return nil, err
	}
	s.cli = c

	return s, nil
}

// SetTLS set tls config
func setTLS(cfg *clientv3.Config, tls *tls.Config, addrs []string) {
	entries := createEndpoints(addrs, "https")
	cfg.Endpoints = entries

	cfg.TLS = tls
}

// setTimeout Set the etcd connection timeout
func setTimeout(cfg *clientv3.Config, time time.Duration) {
	cfg.DialTimeout = time
}

// setCredentials Used to set authentication
func setCredentials(cfg *clientv3.Config, username, password string) {
	cfg.Username = username
	cfg.Password = password
}

// SetNx set key value to etcd if it is not exist, you can set writeoptions to change the default status
func (s *Election) SetNx(ctx context.Context, key string, value []byte, opts *election.WriteOptions) error {
	var op []clientv3.OpOption
	if opts != nil {
		if opts.TTL > 0 {
			lease, err := s.cli.Grant(ctx, int64(opts.TTL.Seconds()))
			if err != nil {
				return err
			}
			op = append(op, clientv3.WithLease(lease.ID))
			if opts.KeepAlive {
				// 如果keepalive没有挂，那么key就一直存在，如果keepalie挂了，超过ttl，key就消失了
				_, err := s.cli.KeepAlive(context.Background(), lease.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	err := s.putNx(ctx, key, value, op...)
	return err
}

func (s *Election) putNx(ctx context.Context, key string, value []byte, ops ...clientv3.OpOption) error {
	resp, err := s.cli.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, string(value), ops...)).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return election.ErrKeyIsExist
	}
	return nil
}

// Delete key-value
func (s *Election) Delete(ctx context.Context, key string) error {
	_, err := s.cli.Delete(ctx, key)
	return err
}

// Watch a key change
// it returns a change that chan can use to receive the value of the key
// passing a value to stopch can terminate watch
func (s *Election) Watch(ctx context.Context, key string, stopCh <-chan struct{}) (<-chan *election.WatchRes, error) {
	rch := s.cli.Watch(ctx, key)
	watchCh := make(chan *election.WatchRes)
	go func() {
		defer close(watchCh)
		for {
			select {
			case <-stopCh:
				return
			case wresp := <-rch:
				for _, ev := range wresp.Events {
					var optType election.OptType
					// 判断操作类型
					switch ev.Type {
					case mvccpb.DELETE:
						optType = election.OptionTypeDelete
					case mvccpb.PUT:
						if ev.IsCreate() {
							optType = election.OptionTypeNew
						} else if ev.IsModify() {
							optType = election.OptionTypeUpdate
						}
					}
					watchCh <- &election.WatchRes{
						KV: &election.KVPair{
							Key:       string(ev.Kv.Key),
							Value:     ev.Kv.Value,
							LastIndex: ev.Kv.ModRevision,
						},
						Type: optType,
					}
				}
			}
		}
	}()

	return watchCh, nil
}

// HeartBeat check whether connection is available
func (s *Election) HeartBeat() error {
	conn := s.cli.ActiveConnection()
	if conn.GetState() == connectivity.TransientFailure {
		return election.ErrHeartBeatFail
	}
	return nil
}

// Close 释放连接
func (s *Election) Close() error {
	return s.cli.Close()
}

// keyNotFound checks on the error returned by the KeysAPI
// to verify if the key exists in the store or not
func keyNotFound(err error) bool {
	return false
}
