# Проект умного дома

# Задание 1. Анализ и планирование

### 1. Описание функциональности монолитного приложения

Нынешнее приложение компании позволяет только управлять отоплением в доме и проверять температуру.
Каждая установка сопровождается выездом специалиста по подключению системы отопления в доме к текущей версии системы.
Система поддерживает только специальные датчики и реле, самостоятельно подключить свой датчик к системе пользователь не может.

**Управление отоплением:**

- Пользователи могут удалённо включать/выключать отопление в своих домах.

**Мониторинг температуры:**

- Пользователи могут просматривать текущую температуру в своих домах через веб-интерфейс
- Система получает данные о температуре с датчиков, установленных в домах, через запрос от сервера к датчику.

### 2. Анализ архитектуры монолитного приложения

- Язык программирования: Go
- База данных: PostgreSQL
- Архитектура: Монолитная, все компоненты системы (обработка запросов, бизнес-логика, работа с данными) находятся в рамках одного приложения.
- Взаимодействие: Синхронное, запросы обрабатываются последовательно.
- Масштабируемость: Ограничена, так как монолит сложно масштабировать по частям.
- Развертывание: Требует остановки всего приложения.

### 3. Определение доменов и границы контекстов

**Управленине домами**
- Справочные данные: назавние, адрес, часовой пояс

**Управленине пользователями**
- Учетные записи пользователей
- Аутентификация

**Управленине устройствами**
- Справочник поддерживаемых моделей (описание, тип, протокол, поддерживаемые функции)
- Экземпляры устройств, привязанных к дому
- Регистрация устройств

**Управленине автоматизацией**
 - Создание сценариев для автоматизации (триггер, условие, действие)
 - Расписание выполнения сценариев

**Телеметрия**
 - Данные, получаемые от устройств

### **4. Проблемы монолитного решения**

- Высокий риск ошибок. Изменения в одной части приложения могут непредсказуемо влиять на другие части. Например, изменения в обработки телеметрии может повлять доступность всей системы. Это увеличивает риск ошибок и требует дополнительных ресурсов на тестирование и отладку всей системы.
- Длительные циклы разработки и развёртывания. При каждом изменении приходится тестировать всё приложение целиком. Это замедляет выпуск новых функций.
- Трудно управлять командой. Сейчас все компоненты приложения объединены в единый блок. Чтобы внести изменение или добавить функциональность, нужно внести изменения во всё приложение. Задачи, которыми занимаются разные команды, могут блокировать друг друга.
- Трудно масштабировать отдельные компоненты системы. Например, часть системы, которая отвечает за обработку телеметрии, со временем обабатывать все больший объем данных. С монолитной архитектурой не получится масштабировать только эту часть — придётся масштабировать приложение целиком.

### 5. Визуализация контекста системы — диаграмма С4

[Диаграмма контекста](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/context/Context-SmartHome_Context_Diagram.png)

# Задание 2. Проектирование микросервисной архитектуры

**Диаграмма контейнеров (Containers)**

[Диаграмма контейнеров](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/container/Container-SmartHouse_Container_Diagram.png)


**Диаграмма компонентов (Components)**

[Диаграмма DeviceManagementService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentDeviceManagement-SmartHouse_Device_Management_Service_Component_Diagram.png)

[Диаграмма DeviceCatalogService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentDeviceCatalog-SmartHouse_Device_Catalog_Service_Component_Diagram.png)

[Диаграмма AutomationService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentAutomation-SmartHouse_Automation_Service_Component_Diagram.png)

[Диаграмма TelemetryService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentTelemetry-SmartHouse_Telemetry_Service_Component_Diagram.png)

[Диаграмма UserService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentUser-SmartHouse_User_Service_Component_Diagram.png)

[Диаграмма WebApplication](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentWeb-SmartHouse_Web_Application_Component_Diagram.png)

[Диаграмма HouseService](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/component/ComponentHouse-SmartHouse_House_Service_Component_Diagram.png)



**Диаграмма кода (Code)**

[Диаграмма кода](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/code/Code-SmartHouse_Device_Management_Service_Code_Diagram.png)

# Задание 3. Разработка ER-диаграммы

[Диаграмма ER](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/schemas/ER/ER-SmartHouse_Device_Management_Service_ER_Diagram.png)

# Задание 4. Создание и документирование API

### 1. Тип API

Для взаимодейтвие между микросервисами будет использоваться REST API, там где не требуется асинхронные вызовы
Там где требуются асинхронные вызовы, там будет использоваться брокер сообщений

### 2. Документация API

[Swagger](https://github.com/reverdp/architecture-pro-warmhouse/blob/warmhouse/swagger/api-gateway.yaml)

# Задание 5. Работа с docker и docker-compose

Настройки выполнены

# **Задание 6. Разработка MVP**

Созданы новые микросервисы device-api и telemetry-api и обеспечена их интеграции с существующим монолитом
