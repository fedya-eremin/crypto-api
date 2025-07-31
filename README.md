# Crypto API
Использовано API CoinMarketCap

## Запуск
Необходимо скопировать .env.example в .env и заменить переменную CMC_API_KEY на актуальное значение
```bash
cp .env.example .env
```

Далее запуск с помощью Docker Compose
```bash
docker compose up
```

API работает на localhost:8001, сваггер доступен на эндпоинте /docs
