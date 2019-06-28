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
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/domgoer/election"
)

var etcdAddr = "127.0.0.1:2379"

var cli *clientv3.Client

func init() {
	addr := os.Getenv("ETCD_PORT_2380_TCP_ADDR")
	if addr != "" {
		etcdAddr = addr + ":2379"
	}
	cli, _ = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 10 * time.Second,
		// dont do auto sync
		// AutoSyncInterval: periodicSync,
	})
}

func TestNew(t *testing.T) {
	type args struct {
		addrs   []string
		options *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Election
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "new_case1",
			args: args{
				addrs: []string{etcdAddr},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.addrs, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Close()
		})
	}
}

func TestElection_Put(t *testing.T) {
	type fields struct {
		cli *clientv3.Client
	}
	type args struct {
		ctx   context.Context
		key   string
		value []byte
		opts  *election.WriteOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "put_case",
			fields: fields{
				cli: cli,
			},
			args: args{
				ctx:   context.Background(),
				key:   "test_put",
				value: []byte("1"),
				opts: &election.WriteOptions{
					TTL:       time.Second,
					KeepAlive: true,
				},
			},
		},
		{
			name: "putnx_case",
			fields: fields{
				cli: cli,
			},
			args: args{
				ctx:   context.Background(),
				key:   "test_put",
				value: []byte("1"),
				opts: &election.WriteOptions{
					TTL: time.Second,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				cli: tt.fields.cli,
			}
			if err := s.SetNx(tt.args.ctx, tt.args.key, tt.args.value, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Election.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestElection_Watch(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	type fields struct {
		cli *clientv3.Client
	}
	type args struct {
		ctx    context.Context
		key    string
		stopCh <-chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *election.WatchRes
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "watch_case",
			fields: fields{
				cli: cli,
			},
			args: args{
				ctx:    context.Background(),
				key:    "test_watch",
				stopCh: ctx.Done(),
			},
			want: &election.WatchRes{
				KV: &election.KVPair{
					Key:       "test_watch",
					Value:     []byte("1"),
					LastIndex: 0,
				},
			},
		},
	}
	go func() {
		time.Sleep(time.Millisecond * 50)
		cli.Put(context.Background(), "test_watch", "1")
		defer cli.Delete(context.Background(), "test_watch")
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				cli: tt.fields.cli,
			}
			got, err := s.Watch(tt.args.ctx, tt.args.key, tt.args.stopCh)
			if (err != nil) != tt.wantErr {
				t.Errorf("Election.Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual((<-got).KV.Key, tt.want.KV.Key) {
				t.Errorf("Election.Watch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElection_Delete(t *testing.T) {
	cli.Put(context.Background(), "test_delete", "")
	type fields struct {
		cli *clientv3.Client
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				cli: cli,
			},
			args: args{
				ctx: context.Background(),
				key: "/test_delete",
			},
		}, {
			name: "case2",
			fields: fields{
				cli: cli,
			},
			args: args{
				ctx: context.Background(),
				key: "/test_deletesdsgfsdgfd1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				cli: tt.fields.cli,
			}
			if err := s.Delete(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Election.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClose(t *testing.T) {
	cli.Close()
}
