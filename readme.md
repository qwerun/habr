Билд проекта  docker compose up -d --build

Порядок выполнения команд для полноценного тестирования
1. curl -v -X POST http://localhost:8081/api/v1/register   -H "Content-Type: application/json"   -d '{"email": "abiojppppp@mail.com", "password": "Hewer)8p", "nickname": "ppppps32"}'

2. docker ps

3. Найти контейнер id у контейнера habr-notification-service и выполнить команду
docker logs df6273237d51
В последнем логе будет отображаться сообщение 2025/07/10 16:13:06 Отправка письма по событию: {"abiojppppp@mail.com":446615}
Если сообщение есть, значит оно отправилось из http-auth-go в кафку и было считано http-notification-go
Значение ключа это и есть код верификации

4. Далее нужно верифицировать аккаунт, выполнив команду
curl -v -X POST http://localhost:8081/api/v1/verify-email -H "Content-Type: application/json"   -d '{"email": "abiojppppp@mail.com", "code": "446615"}'
Если в ответе отображается {"success":true} , значит аккаунт успешно верифицирован и доступна операция - логин

5. Также можно сменить пароль через ручку /api/v1/change-password
curl -v -X POST http://localhost:8081/api/v1/change-password -H "Content-Type: application/json"   -d '{"email": "abiojppppp@mail.com", "password": "Hewer)8p", "newPassword": "Hetui(112"}'
Если в ответе возвращается {"success":true}  значит пароль успешно сменился

6. После верификации аккаунта можно произвести авторизацию curl -v -X POST http://localhost:8081/api/v1/login -H "Content-Type: application/json" -d '{"email": "abiojppppp@mail.com", "password": "Hetui(112", "fingerprint": "123"}'
В ответе вернется:
{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZjI2MmMyNS0zMzE4LTQ5ZDgtOTk2OS02N2Y0ZGNmNWJmMDUiLCJleHAiOjE3NTI0ODAyMDIsImlhdCI6MTc1MjQ3OTYwMiwianRpIjoiMTIzIn0.jncFY6vZAGoFKAxzOxw6Mv4pv4tNAFmmhNJPmEuFV9o",
"refresh":"8BIdUvHdxKBqUbWGE8xwywmJFcKrtaYXNu_qes3BOpk","fingerprint":"123"}

7. Также в любой момент, пока действителен refresh токен, можно обновить access и refresh токены, если выполнить
curl -v -X POST http://localhost:8081/api/v1/refresh -H "Content-Type: application/json" -d '{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZjI2MmMyNS0zMzE4LTQ5ZDgtOTk2OS02N2Y0ZGNmNWJmMDUiLCJleHAiOjE3NTI0ODAyMDIsImlhdCI6MTc1MjQ3OTYwMiwianRpIjoiMTIzIn0.jncFY6vZAGoFKAxzOxw6Mv4pv4tNAFmmhNJPmEuFV9o","refresh":"8BIdUvHdxKBqUbWGE8xwywmJFcKrtaYXNu_qes3BOpk","fingerprint":"123"}'
                   .                 .                     .
И в ответ можно получить новые токены
{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZjI2MmMyNS0zMzE4LTQ5ZDgtOTk2OS02N2Y0ZGNmNWJmMDUiLCJleHAiOjE3NTI0ODAyMDIsImlhdCI6MTc1MjQ3OTYwMiwianRpIjoiMTIzIn0.99pvy_IXS5jbsSNm5Saqa19IP3Pz5HVf9loNHQfb_Gw","refresh":"TV_lvSlV57rVUYW5YWqnQL5SZ3uMc2rBwPWN3Rxoiwg","fingerprint":"123"}