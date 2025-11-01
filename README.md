# Finance TUI — отчет по архитектуре

## 1. Общая идея решения
- Приложение — текстовый интерфейс для управления банковскими счетами, категориями и финансовыми операциями. Пользователь может создавать, редактировать, удалять и просматривать сущности, фильтровать операции, импортировать и экспортировать данные в JSON.
- UI построен на Bubble Tea: корневой экран `main_menu.go` маршрутизирует в экранные модули для разных подсистем (`accounts`, `categories`, `operations`, `files`).
- На уровне бизнес-логики приложение разбито на слои:
  - `domain`: агрегаты и фабрики, интерфейсы репозиториев;
  - `application`: фасады, команды, сервисы экспорт/импорт;
  - `infrastructure`: in-memory репозитории, генератор ULID, файловые импортеры/экспортеры, DI.
- Структура повторяет принципы DDD: доменный слой не зависит от инфраструктуры, границы выражены через интерфейсы (`internal/domain/repository`, `internal/domain/factory`), а модуль `application` управляет сценариями использования и выступает антикоррупционным слоем между TUI и доменом.
- Сборка зависимостей централизована через собственный DI-контейнер (`internal/infrastructure/di`) и bootstrap (`internal/infrastructure/di/bootstrap`). Точка входа `cmd/finance/main.go` управляет средовыми ресурсами (лог таймингов) и вызывает `bootstrap.Build`, который собирает граф сервисов и возвращает готовую модель TUI.

## 2. SOLID и GRASP
- **Single Responsibility Principle**
  - `internal/application/files/export/service.go` и `internal/application/files/import/service.go` отвечают только за экспорт/импорт.
  - `internal/application/command/account/service.go`, `.../category/service.go`, `.../operation/service.go` собирают команды строго для своего агрегата.
- **Open/Closed Principle**
  - Сервис экспорта и импорта расширяются новыми форматами через интерфейсы `fileexport.Exporter` и `fileimport.Importer` (`internal/application/files/export/service.go:38`, `.../import/service.go:32`) без правок существующего кода.
- **Liskov Substitution Principle**
  - Клиенты работают с интерфейсами из `internal/domain/repository/*.go`; in-memory реализации (`internal/infrastructure/repository/memory/*.go`) подставляются без изменений.
  - Фасады `internal/application/facade/*.go` предоставляют единый контракт, который использует TUI.
- **Interface Segregation Principle**
  - Разделённые интерфейсы репозиториев (`AccountRepository`, `CategoryRepository`, `OperationRepository`) позволяют слоям зависеть только от нужных методов.
- **Dependency Inversion Principle**
  - Слой `application` получает зависимости через абстракции `domainfactory.*` и `repository.*`, а конкретные реализации предоставляет DI (`internal/infrastructure/di/bootstrap/application.go`).
- **GRASP Creator**
  - `internal/domain/factory/*.go` создают агрегаты, инкапсулируя логику генерации идентификаторов и валидации.
- **GRASP Controller**
  - Фасады (`internal/application/facade`) координируют сценарии использования; TUI обращается только к ним.
- **High Cohesion / Low Coupling**
  - Пакетная структура и DI-контейнер (`internal/infrastructure/di`) отделяют слои и снижают связность.
- **Pure Fabrication**
  - Слой команд и декораторов (`internal/application/command/*`) — искусственные объекты, обеспечивающие удобный API для TUI и таймингов.

## 3. Паттерны GoF
- **Abstract Factory / Factory Method**
  - `internal/domain/factory/bank_account_factory.go`, `category_factory.go`, `operation_factory.go` инкапсулируют создание агрегатов и проверку входных данных.
- **Command**
  - `internal/application/command/command.go` задаёт интерфейс, а сервисы команд (`.../account/service.go`, `.../category/service.go`, `.../operation/service.go`, `.../export/service.go`, `.../import/service.go`) формируют конкретные команды для TUI.
- **Decorator**
  - `internal/application/command/decorator/timed.go` добавляет измерение времени выполнения команд без изменения базовой логики.
- **Template Method**
  - `internal/application/files/import/service.go` реализует алгоритм импорта: открытие источника, парсинг через выбранный формат и применение payload к фасаду. Разные форматы переопределяют только `Importer.Parse`.
- **Visitor**
  - Экспорт использует посетителей (`internal/application/files/export/visitor.go`, `internal/infrastructure/files/export/json_exporter.go`): сервис экспортирует агрегаты, а конкретный визитор знает, как сериализовать сущность (например, в JSON).
- **Facade**
  - `internal/application/facade/*.go` скрывают сложность доменного слоя за простыми методами CRUD.
- **Strategy**
  - Экспортеры и импортеры выбираются по ключу (`internal/application/files/export/service.go:58`, `.../import/service.go:56`), что позволяет подставлять разные реализации сериализации.
- **Singleton (через DI)**
  - Контейнер (`internal/infrastructure/di/container.go`) кеширует созданные экземпляры, гарантируя единичность репозиториев и сервисов на время запуска.

## 4. Bootstrap и DI
- Контейнер (`internal/infrastructure/di/container.go`) поддерживает регистрация/возврат зависимостей с дженериками, потокобезопасность и обработку ошибок.
- Bootstrap (`internal/infrastructure/di/bootstrap/*.go`) разбит на слои:
  - `registerInfrastructure` — генератор ID, репозитории;
  - `registerDomain` — фабрики агрегатов;
  - `registerApplication` — фасады и сервисы импорт/экспорт;
  - `registerCommands` — сервисы команд с декораторами;
  - `registerUI` — корневой экран Bubble Tea.
- `cmd/finance/main.go` корпус: открывает лог таймингов, подготавливает наборы импортеров/экспортеров и вызывает `bootstrap.Build`. UI запускается через Bubble Tea.

## 5. Запуск и ресурсы
- Требования: Go ≥ 1.24.
- Сборка и запуск:
  ```bash
  go run ./cmd/finance
  ```
- Логи таймингов пишутся в `cmd/finance/logs/timings.log`.
- Экспортированные файлы — JSON, импорт из аналогичного формата.
