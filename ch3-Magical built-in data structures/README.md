# 神奇内置数据结构

## 内置数据结构⼀览

`runtime`

- channel
- timer
- semaphore
- map
- iface
- eface
- slice
- string

`sync`

- mutex
- cond
- pool
- once
- map
- waitgroup

`container`

- heap
- list
- ring

`os`

- os related

`context`

- context

`memory`

- allocation related
- gc related

`netpoll`

- netpoll related

## Channel

`演示动画: 基本执⾏流程`

<https://www.figma.com/proto/vfhlrTqsKicCO5ZbQZXgD4/runtime-structs?node-id=25-2&starting-point-node-id=25%3A2>

`发送流程示意图`

![channel_send_flow.png](https://s2.loli.net/2023/06/06/sq3X5KtMDkLCSmz.png)

`接收流程示意图`

![channel_recv_flow.png](https://s2.loli.net/2023/06/06/clidBXekPQjRoVS.png)

`并发安全`

`挂起和唤醒`

## Timer

## Map

## Context