# goelection

go语言实现的主从选举

## 依赖

`etcd`或者`zookeeper`

## Feature

1. 主从选举
2. 依赖奔溃/恢复 状态自动切换

## Example

```go
package main

import (
    "goelection"
    "goelection/etcd"
)

func main(){
    cli , _ := etcd.New([]string{"addr"},nil)
    elec := goelection.New(cli,"id")
    stopCh := make(chan struct{}{})
    idCh , _ := elec.Elect(stopCh)

    id := <- idCh
    if id == "id" {
        fmt.Println("leader")
    }
}
```