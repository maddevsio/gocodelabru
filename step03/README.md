## Шаг 3. Имплементируем LRU (часть 1)
И вот теперь мы подобрались к интересным моментам нашей кодлабы. У вас уже есть структура проекта и теперь осталось написать  код :)
Поэтому открываем пакет [storage](storage) в вашей IDE или редакторе и начнем с пакета `storage/lru`

## Структуры данных 

Нам нужно подумать, как же реализовывать кеш. С одной стороны нам нужно key-value хралилище и стандартный `map[interface{}]interface{}` нам подойдет. В этом случае мы сможем и добавлять и удалять элементы в кеше быстро. С другой стороны нам нужен какой-то список из элементов, в котором мы бы могли двигать элементы как вверх стопки, так и забирать последний. Плюс ко всему нам нужна возможность удалять значения из этого списка. Кажется, что [container/list](https://golang.org/pkg/container/list/) нам подойдет.

Получается кеш мы сможем описать следующим образом

```Go
type LRU struct {
  size int
  evictList *list.List
  items map[interface{}]*list.Element
}
```
в `lru/lru.go` у нас теперь следующий код

```Go
package lru

import "container/list"

type LRU struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
}
```
## New
Для того, чтобы создать кеш нам нужно передать его размер в параметрах и проинициализировать все структуры храненния.
Получается следующее
```Go
// New initialized a new LRU with fixed size
func New(size int) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Size must be greater than 0")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
	}
	return c, nil
}
```

## Add

```Go
// Add adds a value to the cache. Return true if eviction occured
func (l *LRU) Add(key, value interface{}) bool {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	}
	ent := &entry{key, value}
	entry := l.evictList.PushFront(ent)
	l.items[key] = entry
}

```
