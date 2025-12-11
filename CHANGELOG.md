# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/),
и этот проект придерживается [Semantic Versioning](https://semver.org/lang/ru/).

## [Unreleased]

## [0.0.3] - 2025-12-11

### Added
- Heroicons библиотека для иконок
- Минималистичный дизайн UI в стиле VK/ЗооПлатформа
- Боковое меню с иконками Heroicons
- Разделители между секциями меню
- Badge для мессенджера (43 непрочитанных)
- Обновленный хедер с поиском и уведомлениями

### Changed
- Полностью переработан дизайн интерфейса
- Sidebar: чистый минималистичный стиль без фоновых подложек
- Hover эффекты с закругленными углами
- Header: выровнен с контентом страницы
- Логотип в хедере выровнен с боковым меню
- Layout: центральная колонка 580px, правая панель расширена
- Меню включает: Главная, Лента, Объявления, Профиль, Мессенджер, Зоожурнал, Афиша, Организации, Зоомаркет, Учебный центр, Сервисы, Мои питомцы, Горячая линия
- Дополнительные ссылки: О платформе, Фонд "ЗооПомощь", Статистика, Техподдержка, Команда

### Технологии
- Go 1.25.5
- Node.js 25.2.0
- Next.js 16.0.8
- React 19.2.1
- React Native 0.81.5
- Expo 54.0.27
- TypeScript 5.x
- Tailwind CSS 4.x
- Heroicons 2.x
- Air 1.63.4
- SQLite 3
- JWT (golang-jwt/jwt/v5)
- Bcrypt (golang.org/x/crypto/bcrypt)

## [0.0.2] - 2025-12-11

### Added
- React Native мобильное приложение (Expo)
- Папка `shared/` для общей логики между платформами
- Air для hot reload Go бекенда
- Metro config для React Native
- Пример интеграции mobile с backend API
- SQLite база данных
- Отдельный модуль `database/` для работы с БД
- Таблица `users` с полем password
- Автоматическая инициализация БД при старте сервера
- REST API структура (handlers, models)
- JWT авторизация
- Endpoints для регистрации и входа (`/api/auth/register`, `/api/auth/login`)
- Middleware для защиты endpoints
- Bcrypt хеширование паролей
- CRUD endpoints для users (GET, POST, PUT, DELETE)
- TypeScript типы для всех API запросов/ответов
- Стандартный формат ответов API (success, data, error)
- Документация структуры проекта (STRUCTURE.md)

### Changed
- CORS настроен на `*` для поддержки всех клиентов
- CORS headers включают Authorization
- Скрипт `run` теперь запускает три сервиса (backend, web, mobile)
- API клиент вынесен в shared (для web) и продублирован в mobile
- Обновлен Go до версии 1.25.5
- База данных вынесена в отдельный модуль в корне проекта
- Endpoints `/api/users` теперь защищены JWT токеном

### Технологии
- Go 1.25.5
- Node.js 25.2.0
- Next.js 16.0.8
- React 19.2.1
- React Native 0.81.5
- Expo 54.0.27
- TypeScript 5.x
- Tailwind CSS 4.x
- Air 1.63.4
- SQLite 3
- JWT (golang-jwt/jwt/v5)
- Bcrypt (golang.org/x/crypto/bcrypt)

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
- Git репозиторий с тегом v0.0.1

### Технологии
- Go 1.25.5
- Node.js 25.2.0
- Next.js 16.0.8
- React 19.2.1
- React Native 0.81.5
- Expo 54.0.27
- TypeScript 5.x
- Tailwind CSS 4.x
- Air 1.63.4
