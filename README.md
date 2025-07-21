Простое REST API для условного маркетплейса, реализованное на языке Go. Поддерживаются следующие функции:

- Регистрация и авторизация пользователей
- Размещение новых объявлений (авторизованные пользователи)
- Лента объявлений с фильтрами, сортировкой и пагинацией
- JWT-токены для авторизации

### Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/plusha-fullstack/marketplaceAPI.git
```

Запустите приложение:
```bash
docker compose up --build
```

### Примеры запросов  

```bash
#регистрация пользователя
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"password123"}'
```
```bash
# логин в систему (при успехе возврат токена)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"password123"}'
```
```bash
# создание объявления
curl -X POST http://localhost:8080/api/ads \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer PASTE_YOUR_TOKEN" \
  -d '{"title":"Test Ad","description":"Test description","image_url":"https://example.com/image.jpg","price":99.99}'
```
```bash
# получение списка объявлений с пагинацией, сортировкой, фильтрацией
curl "http://localhost:8080/api/ads?page=1&limit=10&sort_by=date&sort_order=desc&min_price=0&max_price=1000"
```


### Хостинг
Приложение так же как минимум до 21.08.2025 доступно по адресу http://62.133.61.149

### Примеры запросов на хостинг

```bash
#регистрация пользователя
curl -X POST http://62.133.61.149/api/register \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"password123"}'
```
```bash
# логин в систему (при успехе возврат токена)
curl -X POST http://62.133.61.149/api/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"password123"}'
```
```bash
# создание объявления
curl -X POST http://62.133.61.149/api/ads \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer PASTE_YOUR_TOKEN" \
  -d '{"title":"Test Ad","description":"Test description","image_url":"https://example.com/image.jpg","price":99.99}'
```
```bash
# получение списка объявлений с пагинацией, сортировкой, фильтрацией
curl "http://62.133.61.149/api/ads?page=1&limit=10&sort_by=date&sort_order=desc&min_price=0&max_price=1000"
```
