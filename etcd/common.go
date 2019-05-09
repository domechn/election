/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : common.go
#   Created       : 2019-04-29 10:56
#   Last Modified : 2019-04-29 10:56
#   Describe      :
#
# ====================================================*/
package etcd

// createEndpoints creates a list of endpoints given the right scheme
func createEndpoints(addrs []string, scheme string) (entries []string) {
	for _, addr := range addrs {
		entries = append(entries, scheme+"://"+addr)
	}
	return entries
}
