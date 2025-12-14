# school-21-workshop

## 09.11.2025 Архитектура BFF (Backend-for-Frontend) микросервиса

**Темы по Go**: интерфейсы, встраивание структур, HTTP, error.

Реализуем API роутер, который по IP адресу будет возвращать температуру в городе. Познакомимся с архитектурой организации директорий и 
адаптированным подходом vertical slice для разделения логики на слои:
- handler: транспортный слой
- usecase: слой бизнес логики
- gateway: слой абстракции для сетевых запросов

Ссылки на код:
- [Пример программы без применения подхода](https://github.com/yasonofriychuk/school-21-workshop/tree/2025-11-09_bff_service_simple)
- [С применением адаптированного Vertical Slice](https://github.com/yasonofriychuk/school-21-workshop/tree/2025-11-09_bff_service_vertical_slice)

## 23.11.2025 Пирамида тестирования, тесты в Go

**Темы по Go**: mock, виды тестов, инструментарий `go test`, testcontainers для интеграционных тестов, table tests.

Поговорим про пирамиду тестирования, виды тестов. Реализуем интеграционный тест с применением testcontainers и подхода table tests,
используя gomock для моков зависимостей. Рассмотрим флаги запуска тестов команды `go test` и механизм тегирования для разделения
интеграционных от unit тестов.

Ссылки на код:
- [Проект в состоянии до покрытия тестами](https://github.com/yasonofriychuk/school-21-workshop/tree/2025-11-23_test)
- [Проект с тестами](https://github.com/yasonofriychuk/school-21-workshop/tree/2025-11-23_test_finish)

## 14.12.2025 Обработка ошибок в Go

**Темы по Go**: error, пакет errors, кастомные ошибки

Рассмотрим приёмы работы с ошибками в Golang: обёртка, определение типа ошибки, объявление собственных видов ошибок.

Ссылки на код:
- [Код в конце занятия](https://github.com/yasonofriychuk/school-21-workshop/tree/2025-12-14-handle-errors)

