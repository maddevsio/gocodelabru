## Шаг 10. Имплементируем LRU (часть 1)
Для того. чтобы сделать LRU хранилище, сначала нужно рассказать что это такое.

**LRU (least recently used)** — это алгоритм кеширования данных, при котором вытесняются значения, которые дольше всего не запрашивались. Этот механизм удобен тем, что мы например инициализируем кеш на 20 элементов. И как только мы стараемся добавить 21й, то самый долго неиспользуемый элемент удалится.

Этот механизм удобен еще тем, что эту логику мы реализуем только в одном месте.
Его будем делать в пакете `storage/lru` создайте новую папку, создайте два файла `storage/lru/lru.go` и `storage/lru/lru_test.go` с контентом `package lru`

## Структуры данных 

Нам нужно подумать, как же реализовывать кеш. Нам нужен какой-то список из элементов, в котором мы бы могли двигать элементы как вверх стопки, так и забирать последний. Плюс ко всему нам нужна возможность удалять значения из этого списка. Кажется, что [container/list](https://golang.org/pkg/container/list/) нам подойдет.

Получается кеш мы сможем описать следующим образом

```Go
type LRU struct {
  size int
  evictList *list.List
  items map[interface{}]*list.Element
}
```
Также нам нужна структура для того, чтобы хранить данные в списке и карте.
```Go
// entry used to store value in evictList
type entry struct {
	key   interface{}
	value interface{}
}
```

в файле `lru/lru.go` следующий код
```Go
package lru

import (
	"container/list"
)

type (
	LRU struct {
		size      int
		evictList *list.List
		items     map[interface{}]*list.Element
	}
	// entry used to store value in evictList
	entry struct {
		key   interface{}
		value interface{}
	}
)

```

## New
Для того, чтобы создать кеш нам нужно передать его размер в параметрах и проинициализировать все структуры хранения.
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
У нас Add делает две вещи. Добавляет и обновляет значение по ключу. При этом с помощью API `container/list` он управляет положением элемента в списке. Не хватает лишь удаления элемента, если у нас элементов в нашем кеше больше чем его размер

## Удаление наименее используемого элемента
```Go
func (l *LRU) removeOldest() {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
	}
}
func (l *LRU) removeElement(e *list.Element) {
	l.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(l.items, kv.key)
}
```
По закладываемой логике, нам нужно удалить всего-лишь последний элемент из списка. Удалять нужно будет его и из нашего списка и из карты.

Вернемся к добавлению элемента и внедрим удаление самого старого элемента, если мы превышаем размер хранилища. Получится следующий код

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
	evict := l.evictList.Len() > l.size
	if evict {
		l.removeOldest()
	}
	return evict
}
```
Напишем на него тест в `lru/lru_test.go`
```Go
// Test that Add returns true/false if an eviction occurred
func TestLRU_Add(t *testing.T) {

	l, err := New(1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add(1, 1) == true {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false {
		t.Errorf("should have an eviction")
	}
}

```
## Len
С Len() все просто. Нам нужно вернуть только длину списка, чтобы узнать сколько элементов у нас сейчас в кеше
```Go
// Len returns the number of items in cache
func (l *LRU) Len() int {
	return l.evictList.Len()
}
```

# Поздравляю!

Вы реализовали создание, добавление и удаление самого старого элемента из кеша и написали тест на добавление элемента. Мы продолжим работу в [следующей части](../step11/README.md)
