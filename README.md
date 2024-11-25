# Реализация онлайн библиотеки песен 🎶

**Задание**
Необходимо реализовать следующее:

1. **Выставить rest методы**:
- Получение данных библиотеки с фильтрацией по всем полям и пагинацией
- Получение текста песни с пагинацией по куплетам
- Удаление песни
- Изменение данных песни
- Добавление новой песни в формате

```json
{
 "group": "Muse",
 "song": "Supermassive Black Hole"
}
```

2. **При добавлении сделать запрос в АПИ, описанного сваггером**

openapi: 3.0.3
info:
  title: Music info
  version: 0.0.1
paths:
  /info:
    get:
      parameters:
        - name: group
          in: query
          required: true
          schema:
            type: string
        - name: song
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SongDetail'
        '400':
          description: Bad request
        '500':
          description: Internal server error
components:
  schemas:
    SongDetail:
      required:
        - releaseDate
        - text
        - link
      type: object
      properties:
        releaseDate:
          type: string
          example: 16.07.2006
        text:
          type: string
          example: Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight
        link:
          type: string
          example: https://www.youtube.com/watch?v=Xsp3_a-PMTw


3. **Обогащенную информацию положить в БД postgres (структура БД должна быть создана путем миграций при старте сервиса)**
4. **Покрыть код debug- и info-логами**
5. **Вынести конфигурационные данные в .env-файл**
6. **Сгенерировать сваггер на реализованное АПИ**

**Инструкция для работы с **

# Реализация онлайн библиотеки песен 🎶

## API Documentation

### 1. **GET /search**

**Описание:**  
Получение списка песен с фильтрацией по полям и пагинацией.

**Параметры запроса:**
- `group` (необязательный) — Название музыкальной группы.
- `song` (необязательный) — Название песни.

**Ответ:**
- **200 OK**: Возвращает список песен в формате JSON.
- **400 Bad Request**: Ошибка в запросе.
- **500 Internal Server Error**: Ошибка поиска в базе данных.
  
**Пример запроса:**
```bash
GET /search?group=Muse
```
```json
[
  {
    "group": "Muse",
    "song": "Supermassive Black Hole",
  },
  {
    "group": "Muse",
    "song": "Hysteria",
  }
]
```
### 2. **GET /info**

**Описание:**  
Получение информации о песне по названию группы и песни.

**Параметры запроса:**
- `group` (обязательный) — Название музыкальной группы.
- `song` (обязательный) — Название песни.

**Ответ:**
- **200 OK**: Возвращает подробную информацию о песне.
- **400 Bad Request**: Ошибка в запросе.
- **500 Internal Server Error**: Ошибка поиска в базе данных.
  
**Пример запроса:**
```bash
GET /info?group=Muse&song=Supermassive
```
```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole",
  "releaseDate": "16.07.2006",
  "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses",
  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
}
```
### 3. **DELETE /remove**

**Описание:**  
Удаление песни по её названию группы и песни.

**Параметры запроса:**
- `group` (обязательный) — Название музыкальной группы.
- `song` (обязательный) — Название песни.

**Ответ:**
- **204 No Content:**: Песня успешно удалена.
- **404 Not Found:**: Песня не найдена.
- **500 Internal Server Error:**: Ошибка при удалении.
**Пример запроса:**
```bash
DELETE /remove?group=Muse&song=Supermassive
```
### 4. **PATCH /update**

**Описание:**  
Обновление информации о песне.

**Параметры запроса:**
- `group` (обязательный) — Название музыкальной группы.
- `song` (обязательный) — Название песни.

**Тело запроса**
В теле запроса необходимо передать обновленную информацию о песне.

**Ответ:**
- **200 OK**: Песня успешно обновлена.
- **404 Not Found**: Песня не найдена.
- **500 Internal Server Error**: Ошибка при обновлении.
  
**Пример запроса:**
```bash
PATCH /update?group=Muse&song=Supermassive
Content-Type: application/json
```
```json
{
  {
  "releaseDate": "17.07.2006",
  "text": "Updated text of the song",
  "link": "https://www.youtube.com/watch?v=updatedLink"
}
}
```
### 5. **POST /new**

**Описание:**  
Добавление новой песни в библиотеку.

**Тело запроса**
```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole",
}
```

**Ответ:**
- **201 Created**: Песня успешно добавлена.
- **400 Bad Request**: Неверный формат данных.
- **500 Internal Server Error**: Ошибка при добавлении песни.
  
**Пример запроса:**
```bash
POST /new
Content-Type: application/json
```
```json
{
  {
  "group": "Muse",
  "song": "Supermassive Black Hole"
  }
}
```
## Configuration
Для работы с сервисом необходим файл конфигурации **.env**, который должен содержать следующие параметры:
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres
DB_USERNAME=your_username
DB_PASSWORD=your_password
LISTEN_IS_DEBUG=true
LISTEN_TYPE=port
LISTEN_BIND_IP=0.0.0.0
LISTEN_PORT=10000
```
## Swagger Documentation
Swagger документация для вашего API доступна по следующему адресу:
```bash
http://localhost:10000/swagger/index.html
```
Swagger автоматически генерирует и отображает документацию на основе комментариев в коде, помеченных с помощью аннотаций, таких как **@Summary**, **@Param**, **@Success**, и других.
## Запуск и Миграции
1. **Запуск приложения**: Для запуска сервиса используйте команду:
```bash
go run main.go
```