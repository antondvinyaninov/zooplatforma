# @pet/shared - Общие компоненты

Пакет с переиспользуемыми компонентами, типами и утилитами для всех сервисов проекта.

## ⚠️ ВАЖНО: Tailwind CSS v4 и компоненты

**Проблема:** Tailwind CSS v4 НЕ сканирует файлы из `node_modules` по умолчанию, даже если они указаны в `content` конфигурации.

**Решение:** Копируйте компоненты напрямую в каждый сервис вместо импорта из `@pet/shared`.

### Правильный подход

```bash
# Копируем компонент в сервис
cp shared/src/components/AuthForm.tsx main/frontend/app/components/
cp shared/src/components/AuthForm.tsx admin/frontend/app/components/

# Используем локальный импорт
import AuthForm from '../components/AuthForm';
```

### ❌ Неправильный подход (не работает с Tailwind v4)

```typescript
// Это НЕ РАБОТАЕТ - Tailwind не сгенерирует стили
import { AuthForm } from '@pet/shared';
```

## Структура

```
shared/
├── src/
│   ├── components/     # React компоненты (копировать в сервисы!)
│   │   ├── AuthForm.tsx
│   │   └── admin/      # Шаблон админ-панели
│   │       ├── AdminLayout.tsx
│   │       ├── widgets/
│   │       │   ├── StatsWidget.tsx
│   │       │   ├── TableWidget.tsx
│   │       │   └── ChartWidget.tsx
│   │       ├── index.ts
│   │       └── README.md
│   ├── api/           # API клиенты (можно импортировать)
│   ├── types/         # TypeScript типы (можно импортировать)
│   └── hooks/         # React хуки (можно импортировать)
├── package.json
└── README.md
```

## Что можно импортировать из @pet/shared

✅ **Безопасно импортировать:**
- Типы TypeScript (`types/`)
- API клиенты (`api/`)
- Хуки без JSX (`hooks/`)
- Утилиты (`utils/`)

❌ **НЕ импортировать (копировать):**
- React компоненты с Tailwind классами (`components/`)

## Использование

### Для компонентов с Tailwind

```bash
# 1. Копируем компонент
cp shared/src/components/AuthForm.tsx your-service/app/components/

# 2. Для админ-шаблона копируем всю папку
cp -r shared/src/components/admin your-service/app/components/

# 3. Импортируем локально
import AuthForm from '../components/AuthForm';
import AdminLayout from '../components/admin/AdminLayout';
```

### Для типов и API

```typescript
// Это работает нормально
import { User, Post } from '@pet/shared/types';
import { apiClient } from '@pet/shared/api';
```

## Разработка

```bash
# Сборка TypeScript
npm run build

# Watch mode
npm run dev

# Проверка типов
npm run type-check
```

## Добавление нового компонента

1. Создайте компонент в `src/components/`
2. Экспортируйте в `src/index.ts`
3. **Скопируйте** в каждый сервис, который его использует
4. Обновите документацию

## Почему не монорепозиторий?

Мы используем микросервисную архитектуру, где каждый сервис независим. Shared пакет - это просто удобное место для хранения общего кода, но компоненты копируются в сервисы для корректной работы Tailwind CSS v4.

## Альтернативы (для будущего)

Если проект вырастет, можно рассмотреть:
- Переход на Tailwind CSS v3 (поддерживает сканирование node_modules)
- Использование CSS-in-JS решений
- Создание отдельной UI библиотеки с pre-compiled стилями
- Монорепозиторий с Turborepo/Nx
