# Запуск
Запуск проекта производится через docker-compose файл в корневой директории.

```bash
docker-compose up
```
- Пользовательский интерфейс доступен по адресу <http://frontend.localhost>
- API сервиса доступен по адресу <http://api.localhost>
- Swagger UI доступен по адресу <http://api.localhost/swagger>

# Сидирование
Для наполнения топика тестовыми данными существует отдельное приложение `producer`. \
 Для отправки `N` сообщений с рандомными данными можно воспользоваться командой ниже.
```bash
docker-compose exec producer produce --count=N
```
