package cachekey

//// NewCollectionKey 构建一个带有前缀的缓存key 使用 ":" 分隔
//func (p *KeyPrefix) NewCollectionKey(cc *collectioncache.Cache) *CollectionKey {
//	return &CollectionKey{
//		keyPrefix: p,
//		cc:        cc,
//	}
//}
//
//// CollectionKey 实际key参数
//type CollectionKey struct {
//	keyPrefix *KeyPrefix
//	cc        *collectioncache.Cache
//}
//
//// BuildKey  获取key
//func (p *CollectionKey) BuildKey(keys ...any) string {
//	keyStr := make([]string, 0)
//	for _, v := range keys {
//		keyStr = append(keyStr, conv.String(v))
//	}
//	return strings.Join(keyStr, ":")
//}
//
//// FinalKey 获取实际key
//func (p *CollectionKey) FinalKey(key string) string {
//	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, key}, ":")
//}
//
//// TTL 获取缓存key的过期时间time.Duration
//func (p *CollectionKey) TTL() time.Duration {
//	return p.keyPrefix.ExpirationTime
//}
//
//// TTLSecond 获取缓存key的过期时间 Second
//func (p *CollectionKey) TTLSecond() int {
//	return int(p.keyPrefix.ExpirationTime / time.Second)
//}
//
//// CollectionCache 进程内缓存生成
//func (p *CollectionKey) CollectionCache(key string, fn func() (string, error)) (string, error) {
//	take, err := p.cc.TakeWithExpire(p.FinalKey(key), p.TTL(), func() (interface{}, error) {
//		return fn()
//	})
//	if err != nil {
//		return "", err
//	}
//	return take.(string), nil
//}
//
//// CollectionCacheDel 进程内缓存删除
//func (p *CollectionKey) CollectionCacheDel(key string) error {
//	p.cc.Del(p.FinalKey(key))
//	return nil
//}
