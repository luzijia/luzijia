package cache



import (
    "sync"
    "time"
)

type myCache struct {
    value        string
    expiredTime int64
}

const delChannelCap = 100

type ExpiredMap struct {
    m       map[interface{}]*myCache
    timeMap map[int64][]interface{}
    lck     *sync.Mutex
    stop    chan struct{}
}

func NewExpiredMap() *ExpiredMap {
    e := ExpiredMap{
        m:       make(map[interface{}]*myCache),
        lck:     new(sync.Mutex),
        timeMap: make(map[int64][]interface{}),
        stop:    make(chan struct{}),
    }
    go e.run(time.Now().Unix())
    return &e
}

type delMsg struct {
    keys []interface{}
    t    int64
}

func (e *ExpiredMap) run(now int64) {
    t := time.NewTicker(time.Second * 1)
    defer t.Stop()
    delCh := make(chan *delMsg, delChannelCap)
    go func() {
        for v := range delCh {
            e.multiDelete(v.keys, v.t)
        }
    }()
    for {
        select {
        case <-t.C:
            now++
            e.lck.Lock()
            if keys, found := e.timeMap[now]; found {
                e.lck.Unlock()
                delCh <- &delMsg{keys: keys, t: now}
            } else {
                e.lck.Unlock()
            }
        case <-e.stop:
            close(delCh)
            return
        }
    }
}

func (e *ExpiredMap) Set(key int , value string, expireSeconds int64) {
    if expireSeconds <= 0 {
        return
    }
    e.lck.Lock()
    defer e.lck.Unlock()
    expiredTime := time.Now().Unix() + expireSeconds
    e.m[key] = &myCache{
        value:        value,
        expiredTime: expiredTime,
    }
    e.timeMap[expiredTime] = append(e.timeMap[expiredTime], key)
}

func (e *ExpiredMap) Get(key interface{}) (found bool, value string) {
    e.lck.Lock()
    defer e.lck.Unlock()
    if found = e.checkDeleteKey(key); !found {
        return
    }
    value = e.m[key].value
    return
}
func (e *ExpiredMap) checkDeleteKey(key interface{}) bool {
    if val, found := e.m[key]; found {
        if val.expiredTime <= time.Now().Unix() {
            delete(e.m, key)
            return false
        }
        return true
    }
    return false
}
func (e *ExpiredMap) multiDelete(keys []interface{}, t int64) {
    e.lck.Lock()
    defer e.lck.Unlock()
    delete(e.timeMap, t)
    for _, key := range keys {
        delete(e.m, key)
    }
}

func (e *ExpiredMap) Delete(key interface{}) {
    e.lck.Lock()
    delete(e.m, key)
    e.lck.Unlock()
}

func (e *ExpiredMap) Close() {
    e.lck.Lock()
    defer e.lck.Unlock()
    e.stop <- struct{}{}
}

