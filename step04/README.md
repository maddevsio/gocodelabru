## Шаг 4. Имплементируем LRU (часть 2)
В прошлой части мы сделали реализацию методов `New`, `Add`, `removeOldest`, `removeElement`, `Len` и написали тест на работу `Add` метода.
В этой части мы продолжим строить  `LRU` кеш.

## Purge
Задача этого метода - полностью удалить все данные в хранилище. Для этого нам нужно просто пройти по всем элементам в карте и удалить их
```Go
// Purge completely clears cache
func (l *LRU) Purge() {
	for k := range l.items {
		delete(l.items, k)
	}
	l.evictList.Init()
}

```

## Get
Get - для того, чтобы получить элемент по его ключу. Забирать его лучше из карты и возвращать, есть ли этот элемент в хранилище или нет
```Go
// Get looks up a key's value from the cache
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}
```

## Contains

```Go
// Contains check if key is in cache without updating
// recent-ness or deleting it for being state.
func (l *LRU) Contains(key interface{}) (ok bool) {
	_, ok = l.items[key]
	return ok
}
```

## Remove

```Go
// Remove removes prodided key from the cache, returning if the
// key was contained
func (l *LRU) Remove(key interface{}) bool {
	if ent, ok := l.items[key]; ok {
		l.removeElement(ent)
		return true
	}
	return false
}
```

## GetOldest и RemoveOldest
```Go
// RemoveOldest removes oldest item from cache
func (l *LRU) RemoveOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// GetOldest returns oldest item from cache
func (l *LRU) GetOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}
```
