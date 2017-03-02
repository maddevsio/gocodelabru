# Делаем простую базу для гео данных

Привет, гофер. Ну если ты не гофер и хочешь им стать, тоже привет.  Я предлагаю в этой кодлабе совместить две вещи. Изучить как язык Go и может быть освоить для себя пару новых вешей. 

# Поднимаем окружение
Тебе понадобится следующее:

1. Установленный язык [Go](https://golang.org/)
2. Настроенный `GOPATH` :trollface: (Для 1.8 не актуально)
3. Ты знаком с базовыми вещами в Go. [Тур по Go](https://tour.golang.org/) может хорошо в этом помочь

# Цель лабораторной

У этой лабораторной работы две цели:

1. Получить опыт в Go
2. Научиться понимать как примерно работают key-value хранилища(redis, memcached) 
3. Как работают некоторые индексы.

По итогу БД будет уметь следующие вещи:

* Быстрый поиск по ключу;
* Поиск мест, рядом с вами;
* HTTP интерфейс к БД;
* LRU/expire механизмы для хранения данных;

По Go получите следующие знания:

* Как работает concurrency;
* Поработаете с базовыми синтаксическими вещами;
* Опыт тестирования в go;
* Базовые вещи с Makefile;

# Содержание

Этот воркшоп разделен на несколько частей.

* [Шаг 0. Что такое LRU, R-tree и постановка задачи](step00/README.md)
* [Шаг 1. Бутстрапим проект](step01/README.md)
* [Шаг 2. Что нужно знать о тестировании и написании тестов в Go.](step02/README.md)
* [Шаг 3. Имплементируем LRU (часть 1)](step03/README.md)
* [Шаг 4. Имплементируем LRU (часть 2)](step04/README.md)
* [Шаг 5. Строим сторадж](step05/README.md)
* [Шаг 6. Имплементируем сторадж](step06/README.md)
* [Шаг 7. Проектируем HTTP API](step07/README.md)
* [Шаг 8. Делаем HTTP API](step08/README.md)
* [Шаг 9. Пару штук о автоматизации и Makefile](step09/README.md)
* [Шаг 10. Поздравления!](step10/README.md)

## Комьюнити и ресурсы

Есть несколько мест, где вы можете найти информацию про Go:

- [golang.org](https://golang.org)
- [godoc.org](https://godoc.org) тут вы можете найти документацию по любому пакету
- [Блог языка Go](https://blog.golang.org)

Одно из самых замечательных качеств языка Go - это его сообщество. 
### Сообщества и каналы в телеграм

1. [@bishkekgophers](https://telegram.me/bishkekgophers) - Гоферы Бишкека
2. [@devkg](https://telegram.me/devkg) - Программисты Кыргызстана
3. [@maddevsio](https://telegram.me/maddevsio) - канал нашей компании, где мы делимся всякими интересными штуками. Очень часто говорим про Go

### Сообщества в Slack

1. [golang-ru.slack.com](golang-ru.slack.com) - Рускоязычное сообщество гоферов
2. [gophers.slack.com](gophers.slack.com) - Англоязычное сообщество гоферов


### Подкасты

1. [GolangShow](https://golangshow.com) - Русскоязычный подкаст о языке Go
2. [Gotime](http://gotime.fm) - Англоязычный подкаст о языке Go

### Остальное
- [Go Форум](https://forum.golangbridge.org/)
- [@golang](https://twitter.com/golang) and [#golang](https://twitter.com/search?q=%23golang) on Twitter.
- [Go+ community](https://plus.google.com/u/1/communities/114112804251407510571) on Google Plus.

### Благодарности

1. Francesc Campoy за его воркшоп [Building Web Applications with Go](https://github.com/campoy/go-web-workshop/)
2. Ashley McNamara за картинку в 10м шаге. Вы можете посмотреть и другие работы в [репо](https://github.com/ashleymcnamara/gophers)
3. [Елене Граховац](https://twitter.com/webdeva) за фидбек
