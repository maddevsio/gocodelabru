# Делаем простую базу для гео данных

Привет, гофер. Ну если ты не гофер и хочешь им стать, тоже привет.  Я предлагаю в этой кодлабе совместить две вещи. Изучить как язык Go и может быть освоить для себя пару новых вешей.

# Аудитория
Codelab расчитана на людей, у которых есть опыт в программировании и которые хотят попробовать Go. Это может быть люди, пишушие на PHP/Python/Ruby. Для пишуших на C/C++ Codelab будет врятли полезен

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

* [Шаг 0. Постановка задачи](step00/README.md)
* [Шаг 1. Что нужно знать о тестировании и написании тестов в Go.](step01/README.md)
* [Шаг 2. Hello world](step02/README.md)
* [Шаг 3. Проектируем HTTP API](step03/README.md)
* [Шаг 4. Делаем HTTP API](step04/README.md)
* [Шаг 5. Разбиваем main.go на несколько пакетов](step05/README.md)
* [Шаг 6. Makefile, конфигурация и флаги](step06/README.md)
* [Шаг 7. Добавляем хранилище для данных и ищем ближайших водителей наивным путем](step07/README.md)
* [Шаг 8. Пишем первый бенчмарк и зачем он](step08/README.md)
* [Шаг 9. Что такое R-tree и почему оно эффективнее наивной реализации](step09/README.md)
* [Шаг 10. Имплементируем LRU (часть 1)](step10/README.md)
* [Шаг 11. Имплементируем LRU (часть 2)](step11/README.md)
* [Шаг 12. Делаем хранилище консистентным. Внедряем LRU](step12/README.md)
* [Шаг 13. Внедряем хранилище в API](step13/README.md)
* [Шаг 14. Вы прошли курс. Поздравляю](step14/README.md)

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
2. [gophers.slack.com](gophers.slack.com) - Англоязычное сообщество гоферов. Инвайт получить тут [https://invite.slack.golangbridge.org/](https://invite.slack.golangbridge.org/)


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
3. [Елене Граховац](https://twitter.com/webdeva) за ревью и фидбек
