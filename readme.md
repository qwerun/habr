Порядок выполнения команд для полноценного тестирования
1. curl -v -X POST http://localhost:8081/api/v1/register   -H "Content-Type: application/json"   -d '{"email": "abiojppppp@mail.com", "password": "Hewer)8p", "nickname": "ppppps32"}'

2. docker ps

3. Найти контейнер id у контейнера habr-notification-service и выполнить команду
docker logs df6273237d51
В последнем логе будет отображаться сообщение 2025/07/10 16:13:06 Отправка письма по событию: {"abiojppppp@mail.com":446615}
Это и есть код верификации

4. Далее нужно верифицировать аккаунт, выполнив команду
curl -v -X POST http://localhost:8081/api/v1/verify-email -H "Content-Type: application/json"   -d '{"email": "abiojppppp@mail.com", "code": "446615"}'
Если в ответе отображается {"check":true} , значит аккаунт успешно верифицирован и доступна операция - логин





