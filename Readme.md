# pc-agent

Лёгкий агент мониторинга ПК на Go: собирает CPU, RAM, диски
(далее — сеть, процессы, UPS) и отправляет данные на сервер.

Рассчитан на промышленные ПК, включая Windows 7 32-bit.

## Сборка под Windows 7 32-bit

Нужен Go 1.20 (последняя версия с поддержкой Win7):

    go install golang.org/dl/go1.20.14@latest
    go1.20.14 download
    go env -w GOOS=windows GOARCH=386
    go1.20.14 build -o agent.exe

## Статус

- [x] Сбор CPU / RAM / дисков, запуск на Win7 32-bit
- [ ] Сборщики и JSON-формат
- [ ] Отправка на сервер (Django)
- [ ] Мнемосхема в стиле ISA-101