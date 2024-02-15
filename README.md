# 节点服务

集成多方免费和付费的区块链节点服务，开发一套涵盖多链多协议的扫快节点服务。核心功能是实时、高效率、低延迟的加密货币出账入账通知。

## 项目结构

```

    ├── server
        ├── api             (api层)
        │   └── v1          (v1版本接口)
        ├── config          (配置包)
        ├── core            (核心文件)
        ├── docs            (swagger文档目录)
        ├── global          (全局对象)
        ├── initialize      (初始化)
        │   └── internal    (初始化内部函数)
        ├── middleware      (中间件层)
        ├── model           (模型层)
        │   ├── request     (入参结构体)
        │   └── response    (出参结构体)
        ├── packfile        (静态文件打包)
        ├── resource        (静态资源文件夹)
        │   ├── excel       (excel导入导出默认路径)
        │   ├── page        (表单生成器)
        │   └── template    (模板)
        ├── router          (路由层)
        ├── service         (service层)
        ├── source          (source层)
        └── utils           (工具包)
            ├── timer       (定时器接口封装)
            └── upload      (oss接口封装)

```