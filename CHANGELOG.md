# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/),
и этот проект придерживается [Semantic Versioning](https://semver.org/lang/ru/).

## [Unreleased]

## [0.0.1] - 2025-12-11

### Added
- Инициализация проекта с Go backend и Next.js frontend
- Go HTTP сервер на порту 8080 с базовыми эндпоинтами (/, /api/health)
- CORS middleware для взаимодействия с фронтендом
- Next.js приложение с TypeScript на порту 3000
- Tailwind CSS для стилизации
- Структура компонентов: ui/, shared/, layout/
- API клиент (`lib/api.ts`) для связи с бекендом
- Базовые компоненты: Button с вариантами, Header
- Environment variables конфигурация (.env.local)
- Скрипт `run` для запуска проекта с автоматической очисткой портов
- Документация: README.md и CHANGELOG.md

### Технологии
- Go 1.25.5
- Node.js 25.2.0
- Next.js 16.0.8
- React 19.2.1
- TypeScript 5.x
- Tailwind CSS 4.x
