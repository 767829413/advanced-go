# 老生常谈3: Go module

## 说明

Go module 和 package 管理一向都是'精妙'的设计, 最终导致对 go module 掌握的程度而出现 Go 小白工程师和 Go 老鸟工程师，一般 Go 的教学文档会花费大量篇幅来介绍说明 go module,然后剩下留给 Go 的语法和特性.

关于 go module 在哪些版本支持哪些特性都是很废话的说明,咱就简单明了的直接从 Go 1.23.0 说一下,没必要追溯某个版本 go module 光荣的加入了牛逼的新特性，而仅仅是基于当前 Go 1.23.0 版本来将 Go module 的新知和困惑以及其中出现的版本定义做个稍微的介绍和说明.

基于以上的种种,这肯定不是一个官方应该呈现的教程，只是单纯的方便在吹牛逼或者使用过程中记住 Go module 中的对于 version 的定义。尝试自己做个归纳而已, 只是方便个人的理解和记忆，如果把这个当作茴香豆的茴有几种写法来恶心为难别人那就有点小丑了.

## Go module 中几处定义 version 的地方
---------------------------

在 Go 编程语言中，`go.mod` 文件和环境变量用于配置和管理 Go 项目的编译和运行环境。特别是，`go directive`、`toolchain directive`、`go env` 中的 `GOTOOLCHAIN` 以及环境变量中的 `GOTOOLCHAIN` 都有各自的作用和相互之间的关系。以下是详细描述：

### 1 `go.mod`中的`go directive`

* **描述**: `go directive` 用于指定当前模块使用的 Go 语言版本。这有助于确保模块在不同开发环境中具有一致的行为。你不应该在代码中使用比这个版本更高的语法和 API。
    
* **语法**: `go 1.xx`，例如 `go 1.18`, 虽然你也可以使用 `patch` 号，比如 `go 1.23.0`，但是是没必要的，因为 `patch` 不会有 API 的改变，我们只需要定义 `MAJOR.MINOR`。
    
* **作用**:
    
* 确定编译器应使用的语言特性和标准库的**最低版本**。
    
* 影响依赖解析和模块的语义版本控制。
    
* **其他** - 当你使用 `go mod tidy`, 你可能会发现你的 go. mod 文件中 go directive 这一行中的 go 版本自动升高到一个莫名其妙的版本上了。不用惊慌，那是因为它会检查依赖，可能会被设置为依赖的最小版本。
    
### 2 `go.mod`中的`toolchain directive`

* **描述**: `toolchain directive`用于指定编译和运行该模块所需的 Go 工具链版本。
    
* **语法**: `toolchain go1.xx.xx`，例如 `toolchain go1.23.0`
    
* **作用**:
    
* 强制使用指定版本的 Go 工具链进行编译和运行。
    
* 确保编译和运行环境的一致性，尤其是对于特定版本的工具链特性或修复。
    
* **其他** - 执行 `go mod tidy`, 有时候 toolchain directive 会消失，有时候会被修改为一个特定的值，基本属于薛定谔的猫的状态。当然背后有一套复杂的机制在运作。对于个人或者公司来说，指定为当前支持的最新版本吧。 - 而且如果你的 GOTOOLCHAIN 如果设置为 auto, 而莫得 toolchain directive 又指定了一个你机器没有的 go 版本，系统可能会自动下载。 - 注意以下三种，版本好前面都带 `go` 前缀，和 go directive 是不一致的。
    
### 3 `go env`中的`GOTOOLCHAIN`

* **描述**: `GOTOOLCHAIN`是一个环境变量，用于指定 Go 工具链的版本。
    
* **设置方式**: `go env -w GOTOOLCHAIN=go1.xx`，例如 `go env -w GOTOOLCHAIN=go1.18`
    
* **作用**:
    
* 设置全局工具链版本，影响所有 Go 命令的执行。
    
* 可以覆盖默认的 Go 版本，提供一个一致的开发和编译环境。
    
* **其他** - 运行 `go env`, 你会看到这个全局变量，比如 `GOTOOLCHAIN='auto'`, 正如上面所说，你可以修改它，不过默认是 `auto`。 - 当设置为 `auto` 时，如果 `go.mod` 中 go directive 或者 toolchain 声明的 Go 版本更新，那么系统会自动下载最新的 Go 。 - 另外一个格式是 `GOTOOLCHAIN=go1.21.3+auto`, 默认选择 `go1.21.3`, 如果 `go1.21.3` 不满足则下载最新的 Go 。 - 还可以设置 `GOTOOLCHAIN=local`, 工具链使用本机默认安装的 Go 。 - 还可以组合 `GOTOOLCHAIN=local+auto`，等价于 `GOTOOLCHAIN=auto`。 - 还可以设置 `GOTOOLCHAIN=path, 应该是从` $PATH \`,官方文档语焉不详。
    
### 4 环境变量中的`GOTOOLCHAIN`

* **描述**: 这是系统环境变量，可以在 shell 或其他环境配置文件中设置，用于指定 Go 工具链的版本。
    
* **设置方式**:
    
* 在 Unix/Linux/macOS 中：`export GOTOOLCHAIN=go1.xx.xx`
    荣幸
* 在 Windows 中：`set GOTOOLCHAIN=go1.xx.xx`
    
* **作用**:
    
* 类似于`go env`中的`GOTOOLCHAIN`，但优先级更高，因为它是在操作系统级别设置的。
    
* 影响所有 Go 命令和工具的执行。
    
## 相互之间的关系和优先级
---------------------------

1. **优先级**（从高到低）:
    
* 环境变量中的`GOTOOLCHAIN`
    
* `go env`中的`GOTOOLCHAIN`
    
* `go.mod`中的`toolchain directive`
    
* `go.mod`中的`go directive`
    
3. **相互抑制关系**:
    
* 环境变量中的`GOTOOLCHAIN`会覆盖`go env`中的`GOTOOLCHAIN`设置。
    
* `go env`中的`GOTOOLCHAIN`会覆盖`go.mod`中的`toolchain directive`。
    
* `go.mod`中的`toolchain directive`会覆盖`go directive`关于工具链版本的影响，但不影响语言层面的版本控制。
    
5. **选择与使用**:
    
* **开发阶段**: 可以使用`go.mod`中的`go directive`和`toolchain directive`来确保团队使用一致的 Go 语言版本和工具链版本。
    
* **部署和 CI/CD**: 可以使用环境变量中的`GOTOOLCHAIN`来强制指定工具链版本，确保编译和运行环境的一致性。
    
通过理解这些配置和设置之间的关系，可以更有效地管理和控制 Go 项目的编译和运行环境，从而确保项目的稳定性和一致性。

其实 `go.mod` 还有一个新加的 `godebug` directive, 前面已经足够学习的了，这个很少会使用，就不介绍了。

在加上 `go.work` 中的配置，就更复杂了。