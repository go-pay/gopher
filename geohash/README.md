# geohash

<a href="https://www.igoogle.ink" target="_blank"><img src="https://img.shields.io/badge/Author-Jerry-blue.svg"/></a>
<a href="https://golang.org" target="_blank"><img src="https://img.shields.io/badge/Golang-1.11+-brightgreen.svg"/></a>
<img src="https://img.shields.io/badge/Build-passing-brightgreen.svg"/>

# 使用手册

## 安装
```bash
$ go get -u github.com/iGoogle-ink/geohash
```

## 计算geohash
```go
//计算geohash
//    lat：纬度
//    lng：纬度
//    precision：精度值
geohash := geohash.Encode(39.928167, 116.389550, 6)
fmt.Println(geohash)
```

# geohash原理

geohash基本原理是将地球理解为一个二维平面，将平面递归分解成更小的子块，每个子块在一定经纬度范围内拥有相同的编码，这种方式简单粗暴，可以满足对小规模的数据进行经纬度的检索

### geohash算法

以经纬度值：（116.389550， 39.928167）进行算法说明，对纬度39.928167进行逼近编码 (地球纬度区间是[-90,90]）

- 1、区间[-90,90]进行二分为[-90,0),[0,90]，称为左右区间，可以确定39.928167属于右区间[0,90]，给标记为1
- 2、接着将区间[0,90]进行二分为 [0,45),[45,90]，可以确定39.928167属于左区间 [0,45)，给标记为0
- 3、递归上述过程39.928167总是属于某个区间[a,b]。随着每次迭代区间[a,b]总在缩小，并越来越逼近39.928167
- 4、同理，地球经度区间是[-180,180]，可以对经度116.389550进行编码

- 通过上述计算，纬度产生的编码为1 1 0 1 0 0 1 0 1 1 0 0 0 1 0，经度产生的编码为1 0 1 1 1 0 0 0 1 1 0 0 1 0 0
- 合并：偶数位放经度，奇数位放纬度，把2串编码组合生成新串如下：1 1 1 0 0 1 1 1 0 1 0 0 1 0 0 0 1 1 1 1 0 0 0 0 0 1 1 0 0 0
- 将合并后的二进制转成十进制，并对转后的十进制做base32编码，如下：
    
[1 1 1 0 0] => 28 => w

[1 1 1 0 1] => 29 => x

[0 0 1 0 0] => 4 => 4

[0 1 1 1 1] => 15 => g

[0 0 0 0 0] => 0 => 0

[1 1 0 0 0] => 24 => s

- 最后得到结果：wx4g0s

### geohash长度和精度对照表

<img src="https://raw.githubusercontent.com/iGoogle-ink/gotil/main/geohash/table.jpg"/>