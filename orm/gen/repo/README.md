# repo 代码生成

## 增:

```
CreateOne
```

## 删:

### 唯一索引

#### 单个字段

```
DeleteOneCacheByField
DeleteMultiByCacheFieldPlural
```

#### 多个字段

```
DeleteOneCacheByFields
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
FindMultiCacheByFieldPlural 带缓存
```

#### 多个字段

```
FindOneCacheByFields 带缓存
```

### 非唯一索引：

#### 单个字段

```
FindMultiByFieldPlural 不带缓存
```

#### 多个字段

```
FindMultiByFields 不带缓存
```
