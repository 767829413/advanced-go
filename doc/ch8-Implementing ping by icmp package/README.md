# 通过 icmp 软件包实现 ping

## 实现基础

golang.org/x/net/icmp 很容易实现基于 ICMP 的工具, 但是 Matt Layher (Go 网络编程的专家) 刚推出一个新的 ICMP 库：[mdlayher/icmpx](https://github.com/mdlayher/icmpx), 这个库的使用也非常简单, ReadFrom 用来读, WriteTo 用来发, Close 用来关闭, SetTOS 可以设置 TOS 值.

## 实现方式
