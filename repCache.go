package wecms

type repCache map[ID]interface{}

func (cache repCache) clear() {
	if len(cache) > 0 {
		for key := range cache {
			delete(cache, key)
		}
	}
}