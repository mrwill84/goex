<div align="center">
<img width="409" heigth="205" src="https://upload-images.jianshu.io/upload_images/6760989-dec7dc747846880e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240"  alt="goex">
<img width="198" height="205" src="https://upload-images.jianshu.io/upload_images/6760989-81f29f7a5dbd9bb6.jpg" >
</div>

### goex 目标

goex 项目是为了统一并标准化各个数字资产交易平台的接口而设计，同一个策略可以随时切换到任意一个交易平台，而不需要更改任何代码。

[English](https://github.com/mrwill84/goex/blob/dev/README_en.md)

### wiki 文档

[文档](https://github.com/mrwill84/goex/wiki)

### goex 已支持交易所 `23+`

| 交易所                 | 行情接口      | 交易接口 | 版本号 |
| ---------------------- | ------------- | -------- | ------ |
| huobi.pro              | Y             | Y        | 1      |
| hbdm.com               | Y (REST / WS) | Y        | 1      |
| okex.com (spot/future) | Y (REST / WS) | Y        | 1      |
| okex.com (swap future) | Y             | Y        | 2      |
| binance.com            | Y             | Y        | 1      |
| kucoin.com             | Y             | Y        | 1      |
| bitstamp.net           | Y             | Y        | 1      |
| bitfinex.com           | Y             | Y        | 1      |
| zb.com                 | Y             | Y        | 1      |
| kraken.com             | Y             | Y        | \*     |
| poloniex.com           | Y             | Y        | \*     |
| big.one                | Y             | Y        | 2\|3   |
| hitbtc.com             | Y             | Y        | \*     |
| coinex.com             | Y             | Y        | 1      |
| exx.com                | Y             | Y        | 1      |
| bithumb.com            | Y             | Y        | \*     |
| gate.io                | Y             | N        | 1      |
| bittrex.com            | Y             | N        | 1.1    |

### 安装 goex 库

> go get

`go get github.com/mrwill84/goex`

> 建议 go mod 管理依赖

```
require (
          github.com/mrwill84/goex latest
)
```

### 注意事项

1. 推荐使用 GoLand 开发。
2. 推荐关闭自动格式化功能,代码请使用 go fmt 格式化.
3. 不建议对现已存在的文件进行重新格式化，这样会导致 commit 特别糟糕。
4. 请用 OrderID2 这个字段代替 OrderID
5. 请不要使用 deprecated 关键字标注的方法和字段，后面版本可能随时删除的

---
