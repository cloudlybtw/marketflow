
# 🐰 marketflow — Real-Time Market Data Processing System 📈🚀

> *“Fast, fluffy & concurrent — just like a bunny hopping between exchanges!”* 🐇✨

---

## 🌟 О проекте

**marketflow** — это система для обработки рыночных данных в реальном времени.
Она собирает цены криптовалют с разных бирж, обрабатывает их параллельно 🧵, кэширует ⚡️ и сохраняет в базу данных 📦.

Система поддерживает:

* 🔴 **Live Mode** — подключение к эмуляторам криптобирж (BTCUSDT, DOGEUSDT, TONUSDT, SOLUSDT, ETHUSDT).
* 🟢 **Test Mode** — генерацию тестовых данных через встроенный генератор 🎲.

Ты получаешь **реальный опыт**, как строятся трейдинговые платформы, используя:

* Go’s concurrency (goroutines + channels) 🦫
* Hexagonal Architecture 🛠️
* Redis (кэширование) 🍩
* PostgreSQL (хранение данных) 📊
* REST API (для запросов) 🌐

---

## 🧩 Архитектура

```
         🟢 Source 1     🟢 Source 2     🟢 Source 3
          (Exchange)      (Exchange)      (Exchange)
                │              │              │
                ▼              ▼              ▼
             Fan-Out        Fan-Out        Fan-Out
                │              │              │
        ┌───────┼───────┐─────┼───────┐─────┼───────┐
        ▼       ▼       ▼     ▼       ▼     ▼       ▼
      Worker   Worker  Worker Worker Worker Worker Worker
        │        │       │      │       │     │      │
        └────────┴───────┴──────┴───────┴─────┘
                         │
                         ▼
                       Fan-In
                         │
                         ▼
                   PostgreSQL + Redis
                         │
                         ▼
                    🌐 REST API
```

---

## ⚡️ Основные возможности

* 🐇 **Fan-in / Fan-out / Worker pool** — многопоточная обработка данных
* 🔄 **Автоматический реконнект** при падении источников
* 🛡 **Failover**: если Redis упал → PostgreSQL всё равно сохраняет
* ⏱ **Ежеминутная агрегация**: avg / min / max для каждой пары
* 🗄 **PostgreSQL таблица** для хранения истории:

| pair\_name | exchange | timestamp | average\_price | min\_price | max\_price |
| ---------- | -------- | --------- | -------------- | ---------- | ---------- |

---

## 🐙 API

📡 **Market Data API**

* `GET /prices/latest/{symbol}` → Последняя цена
* `GET /prices/highest/{symbol}?period=1m` → Самая высокая цена
* `GET /prices/lowest/{symbol}` → Самая низкая цена
* `GET /prices/average/{exchange}/{symbol}?period=5m` → Средняя цена

🔀 **Data Mode API**

* `POST /mode/test` → Включить Test Mode 🧪
* `POST /mode/live` → Включить Live Mode 🔴

💓 **Healthcheck**

* `GET /health` → Состояние системы (Redis, PostgreSQL, Sources)

---

## 🛠 Установка и запуск

### 📦 Зависимости

* Go >= 1.22
* PostgreSQL 🐘
* Redis 🍩
* Docker 🐳 (для эмуляторов бирж)

### 🚀 Запуск

```bash
# Сборка проекта
go build -o marketflow .

# Запуск в Live Mode
./marketflow --port 8080

# Запуск в Test Mode
./marketflow --port 8080 --mode test
```

### 🐳 Запуск биржевых эмуляторов

```bash
docker load -i exchange1_amd64.tar
docker run -p 40101:40101 --name exchange1 -d exchange1-amd64
```

---

## 🧸 Пример

```bash
curl http://localhost:8080/prices/latest/BTCUSDT
```

Ответ:

```json
{
  "pair": "BTCUSDT",
  "exchange": "exchange1",
  "price": 65720.54,
  "timestamp": "2025-08-28T10:10:00Z"
}
```

---

## 🐾 Полезные команды

```bash
./marketflow --help

Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number
```

---

## 💖 Credits

👩‍💻 Сделано с любовью и Go 🦫
🎓 Основано на реальных проектах из финтеха
