package oneway

import "sync"


type Getter func(key interface{}) (value interface{}, err error)


type Cache interface {

	Get(key interface{}) (value interface{})
	Set(key interface{}, value interface{})

}

type MapCache struct {
    cache          map[interface{}] interface{}
}

func (m *MapCache) Get(key interface{}) (value interface{}){
	return m.cache[key]
}

func (m *MapCache) Set(key interface{}, value interface{}){
	m.cache[key] = value
}

func NewMapCache() *MapCache{
	return &MapCache{
		cache: make(map[interface{}] interface{}),
	}
}


type group struct {
    key          interface{}
    value        interface{}
    wait         sync.WaitGroup
    err          error
}



type OneWay struct {
    cache          Cache
    lock           sync.Mutex
}


func (one *OneWay) Get(key interface{}, getter Getter) (value interface{}, err error){
	one.lock.Lock()

	var g *group

	if v := one.cache.Get(key); v != nil{
		one.lock.Unlock()
		g = v.(*group)

		g.wait.Wait()
		return g.value, g.err
	}

	g = &group{
		key:   key,
		value: nil,
		wait:  sync.WaitGroup{},
		err:   nil,
	}

	g.wait.Add(1)
	one.cache.Set(key, g)
	one.lock.Unlock()

	g.value, g.err = getter(key)
	g.wait.Done()
	return g.value, g.err

}


func NewOneway(cache Cache) *OneWay{
	return &OneWay{
		cache: cache,
		lock:  sync.Mutex{},
	}
}


