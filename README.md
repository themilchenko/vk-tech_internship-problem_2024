# Профильное задание для прохождения на стажировку в VK Tech

---

:warning: Более стабильная и доработанная версия лежит в ветке **hotfix**.

:notebook_with_decorative_cover: Докумментация доступна по ссылке: https://milchenko.online

:rocket: Деплой API доступен по пути: https://milchenko.online/api

---

Чтобы запустить сервер, нужно прописать:

```bash
go run cmd/main.go
```

При помощи docker-compose:

```bash
docker compose up -d
```

Также можно запустить с помощью конфига, используя флаг `-ConfigPath`, значение которого будет путь к файлу. Сами файлы конфигурации можно посмотреть в `configs/`.
