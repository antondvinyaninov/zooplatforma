---
title: Работа с @pet/shared компонентами
inclusion: always
---

# Правила работы с @pet/shared

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Tailwind CSS v4

**Tailwind CSS v4 НЕ сканирует файлы из `node_modules`**, даже если они указаны в `content` конфигурации.

### Что делать с React компонентами

**ВСЕГДА копируйте компоненты с Tailwind классами напрямую в сервис:**

```bash
# Правильно
cp shared/src/components/AuthForm.tsx main/frontend/app/components/
cp shared/src/components/AuthForm.tsx admin/frontend/app/components/

# Импорт
import AuthForm from '../components/AuthForm';
```

**НИКОГДА не импортируйте компоненты из @pet/shared:**

```typescript
// ❌ НЕ ДЕЛАЙТЕ ТАК - стили не применятся
import { AuthForm } from '@pet/shared';
```

### Что можно импортировать

✅ **Безопасно:**
- TypeScript типы
- API клиенты
- Хуки без JSX
- Утилиты

❌ **Копировать:**
- React компоненты с Tailwind классами

## Процесс добавления компонента

1. Создайте компонент в `shared/src/components/`
2. Протестируйте локально
3. **Скопируйте** в каждый сервис:
   ```bash
   cp shared/src/components/NewComponent.tsx main/frontend/app/components/
   cp shared/src/components/NewComponent.tsx admin/frontend/app/components/
   ```
4. Используйте локальный импорт в каждом сервисе

## Почему так?

- Tailwind v4 изменил механизм сканирования файлов
- `node_modules` исключены из сканирования по умолчанию
- Это известное ограничение v4
- Копирование компонентов - самое простое и надежное решение

## Структура Tailwind конфигурации

```typescript
// tailwind.config.ts
const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    // ❌ Это НЕ работает в v4:
    // './node_modules/@pet/shared/**/*.{js,ts,jsx,tsx}',
  ],
}
```

## Синхронизация компонентов

Если обновляете компонент в `shared/`:

1. Обновите файл в `shared/src/components/`
2. Скопируйте во все сервисы:
   ```bash
   cp shared/src/components/AuthForm.tsx main/frontend/app/components/
   cp shared/src/components/AuthForm.tsx admin/frontend/app/components/
   ```
3. Проверьте работу в каждом сервисе

## Альтернативы (не рекомендуется сейчас)

- Downgrade на Tailwind v3 (теряем новые фичи)
- CSS-in-JS (усложняет проект)
- Pre-compiled UI библиотека (overhead)
- Монорепозиторий (меняет архитектуру)

## Checklist при создании нового сервиса

- [ ] Скопировать нужные компоненты из `shared/src/components/`
- [ ] Настроить Tailwind config с правильными путями
- [ ] Использовать локальные импорты
- [ ] Проверить что все стили применяются
- [ ] Не добавлять `node_modules/@pet/shared` в Tailwind content
