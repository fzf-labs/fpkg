# repo 代码生成

## 增:
```
CreateOne
```
## 删:
### 唯一索引
```
DeleteOneByField
DeleteMultiByFieldComplex
```
## 改:
```
UpdateOne
```
## 查
### 唯一索引
#### 单个字段
```
FindOneCacheByField   带缓存
FindMultiCacheByFieldComplex 带缓存
```
#### 多个字段
```
FindOneCacheByFields 带缓存
```
### 非唯一索引：

#### 单个字段
```
FindMultiByFieldComplex 不带缓存
```
#### 多个字段
```
FindMultiByFields 不带缓存
```
