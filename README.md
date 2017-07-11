# 二维码 api

使用方式

```bash
./app -port=要监听的端口，默认为 9999
```

# 调用

```bash
http://127.0.0.1:9999/?size=256x256&data=HelloWorld
```

size: 参数大小(xxx，只支持正方形的大小, 比如 size=100 或 size=100x100，这种只是为了兼容网上的 qrserver.com 的api，方便迁移)
data: 二维码的内容