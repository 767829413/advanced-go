通用分析工具

readelf -h ./advanced-go | grep -E "Entry point address:"
dlv exec ./advanced-go

```dlv
b *0x453880
c
si
b *runqput
c
dt
```

dlv debug ./main.go

```dlv
b *main.main
help
disass
b *runtime.closechan
c
n
```

go tool compile -S ./hello.go | grep “hello.go:5”
go tool objdump ./advanced-go | grep  

google calling conven golang