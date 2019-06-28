# github.com/domgoer/election

主从选举

## 依赖

 `etcd` 或者 `zookeeper` 

## Example

```go
package main

import (
    "github.com/domgoer/election"
    "github.com/domgoer/election/etcd"
)

func main(){
    cli , _ := etcd.New([]string{"addr"},nil)
    id := "uuid"
    elec := github.com/domgoer/election.New(cli,id)
    stopCh := make(chan struct{}{})
    idCh , _ := elec.Elect(stopCh)

    leaderID := <- idCh
    if leaderID == id {
        fmt.Println("leader")
    }
}
```

## 异常处理

如果etcd或者zk不可用， 会让所有服务都升级为leader， 直到etcd或者zk恢复会重新进行选举。 