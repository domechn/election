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
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"goelection"
	"github.com/samuel/go-zookeeper/zk"
)

var addrt = "127.0.0.1:2181"
var testKey = "/api-gateway-controller/leader/selection"
var cli *zk.Conn

func init() {
	addr := os.Getenv("ETCD_PORT_2380_TCP_ADDR")
	if addr != "" {
		addrt = addr + ":2181"
	}
	var err error
	cli, _, err = zk.Connect([]string{addrt}, 10*time.Second)
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	type args struct {
		addrs []string
		cfg   *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Election
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				addrs: []string{addrt},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.addrs, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Close()
		})
	}
}

func TestElection_get(t *testing.T) {
	type fields struct {
		client *zk.Conn
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *election.KVPair
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				client: cli,
			},
			args: args{
				key: testKey,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				client: tt.fields.client,
			}
			got, err := s.get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Election.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Election.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElection_Watch(t *testing.T) {
	type fields struct {
		client *zk.Conn
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
			name: "case1",
			fields: fields{
				client: cli,
			},
			args: args{
				ctx: context.Background(),
				key: "/test_key",
			},
			want: &election.WatchRes{
				KV: &election.KVPair{
					Key: "/test_key",
				},
				Type: election.OptionTypeNew,
			},
		},
	}
	go func() {
		time.Sleep(time.Millisecond * 50)
		cli.Create("/test_key", []byte{}, 1, zk.WorldACL(zk.PermAll))
	}()
	for _, tt := range tests {
		stCh := make(chan struct{})
		tt.args.stopCh = stCh
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				client: tt.fields.client,
			}
			got, err := s.Watch(tt.args.ctx, tt.args.key, tt.args.stopCh)
			if (err != nil) != tt.wantErr {
				t.Errorf("Election.Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			g := <-got
			fmt.Println(g)
			stCh <- struct{}{}
			if !reflect.DeepEqual(g.KV.Key, tt.want.KV.Key) {
				t.Errorf("Election.Watch() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(g.Type, tt.want.Type) {
				t.Errorf("Election.Watch() = %v, want %v", got, tt.want)
			}
			// tt.args.stopCh <- struct{}{}
		})
	}
}

func TestElection_SetNx(t *testing.T) {
	type fields struct {
		client *zk.Conn
	}
	type args struct {
		ctx   context.Context
		key   string
		value []byte
		opt   *election.WriteOptions
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
				client: cli,
			},
			args: args{
				ctx:   context.Background(),
				key:   testKey,
				value: []byte{},
				opt: &election.WriteOptions{
					KeepAlive: true,
				},
			},
		}, {
			name: "case2",
			fields: fields{
				client: cli,
			},
			args: args{
				ctx:   context.Background(),
				key:   testKey,
				value: []byte{},
				opt: &election.WriteOptions{
					KeepAlive: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				client: tt.fields.client,
			}
			if err := s.SetNx(tt.args.ctx, tt.args.key, tt.args.value, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Election.SetNx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestElection_HeartBeat(t *testing.T) {
	type fields struct {
		client *zk.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				client: cli,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				client: tt.fields.client,
			}
			if err := s.HeartBeat(); (err != nil) != tt.wantErr {
				t.Errorf("Election.HeartBeat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestElection_Delete(t *testing.T) {
	cli.Create("/test_delete", []byte{}, 1, zk.WorldACL(zk.PermAll))
	type fields struct {
		client *zk.Conn
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
				client: cli,
			},
			args: args{
				ctx: context.Background(),
				key: "/test_delete",
			},
		}, {
			name: "case2",
			fields: fields{
				client: cli,
			},
			args: args{
				ctx: context.Background(),
				key: "/test_delete1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Election{
				client: tt.fields.client,
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

