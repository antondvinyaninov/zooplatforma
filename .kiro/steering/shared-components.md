---
title: Работа с @pet/shared компонентами
inclusion: always
---

---
inclusion: manual
---

# Правила работы с @pet/shared

## Назначение папки shared/

**`shared/` - это единое хранилище переиспользуемых компонентов, стилей, типов и утилит для всех микросервисов!**

### Структура shared/

```
shared/
├── src/
│   ├── components/          # React компоненты (Button, Input, AdminLayout)
│   ├── styles/              # Стили (colors, typography, animations)
│   ├── types/               # TypeScript типы (User, Post, Organization)
│   ├── hooks/               # React хуки (useAuth, useFetch, useDebounce)
│   ├── utils/               # Утилиты (validation, formatting, helpers)
│   ├── api/                 # API клиенты
│   └── constants/           # Константы (routes, config)
└── package.json
```

## Как использовать компоненты из shared/

### 1. Импорт компонентов

**Все компоненты импортируются напрямую из @pet/shared:**

```typescript
// ✅ ПРАВИЛЬНО - импортировать из @pet/shared
import { Button } from '@pet/shared/components/buttons';
import { Input } from '@pet/shared/components/forms';
import { AdminLayout } from '@pet/shared/components/layout';
import { User, Organization } from '@pet/shared/types';
import { isValidEmail, formatDate } from '@pet/shared/utils';

// Использование
<Button variant="primary" onClick={handleClick}>
  Сохранить
</Button>
```

### 2. Настройка Tailwind config

**ОБЯЗАТЕЛЬНО добавьте путь к shared в tailwind.config.ts:**

```typescript
// clinic/frontend/tailwind.config.ts
import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    // ✅ ОБЯЗАТЕЛЬНО: добавить путь к shared компонентам
    '../../../shared/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

export default config;
```

### 3. Установка @pet/shared

**Проверьте что @pet/shared установлен в package.json:**

```json
{
  "dependencies": {
    "@pet/shared": "file:../../../shared"
  }
}
```

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Единый AdminLayout

**Все микросервисы используют ОДИНАКОВЫЙ AdminLayout из shared/!**

### Использование AdminLayout

```typescript
// clinic/frontend/app/(dashboard)/layout.tsx
import { AdminLayout } from '@pet/shared/components/layout';
import { HomeIcon, UsersIcon, CalendarIcon, CogIcon } from '@heroicons/react/24/outline';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const navigationItems = [
    { name: 'Обзор', href: '/overview', icon: HomeIcon },
    { name: 'Пациенты', href: '/patients', icon: UsersIcon },
    { name: 'Записи', href: '/appointments', icon: CalendarIcon },
    { name: 'Настройки', href: '/settings', icon: CogIcon },
  ];

  return (
    <AdminLayout 
      serviceName="Clinic"
      navigationItems={navigationItems}
    >
      {children}
    </AdminLayout>
  );
}
```

### Правила кастомизации AdminLayout

**✅ Можно менять:**
- Список пунктов навигации (navigationItems prop)
- Название сервиса (serviceName prop)

**❌ НЕЛЬЗЯ менять:**
- Структуру layout (sidebar + header + content)
- Размеры (ширина sidebar, высота header)
- Цвета и стили
- Расположение элементов

### Примеры навигации для разных сервисов

**Admin:**
```typescript
const navigationItems = [
  { name: 'Пользователи', href: '/users', icon: UsersIcon },
  { name: 'Посты', href: '/posts', icon: DocumentIcon },
  { name: 'Статистика', href: '/stats', icon: ChartBarIcon },
];
```

**Clinic:**
```typescript
const navigationItems = [
  { name: 'Обзор', href: '/overview', icon: HomeIcon },
  { name: 'Пациенты', href: '/patients', icon: UsersIcon },
  { name: 'Записи', href: '/appointments', icon: CalendarIcon },
];
```

**Shelter:**
```typescript
const navigationItems = [
  { name: 'Обзор', href: '/overview', icon: HomeIcon },
  { name: 'Питомцы', href: '/pets', icon: HeartIcon },
  { name: 'Заявки', href: '/applications', icon: ClipboardIcon },
];
```

## Что ОБЯЗАТЕЛЬНО брать из shared/

1. **AdminLayout** - единый layout для всех сервисов
2. **UI компоненты** - Button, Input, Card, Modal
3. **Стили** - colors, typography, animations
4. **Типы** - User, Post, Organization
5. **Утилиты** - validation, formatting
6. **Хуки** - useAuth, useFetch

**Правило:** Если компонент используется в 2+ сервисах - он должен быть в `shared/`!

## Workflow работы с shared/

### Создание нового компонента

1. **Создать в shared/**
```bash
touch shared/src/components/buttons/LoadingButton.tsx
```

2. **Написать компонент**
```typescript
// shared/src/components/buttons/LoadingButton.tsx
export function LoadingButton({ loading, children, ...props }) {
  return (
    <button disabled={loading} className="px-4 py-2 bg-blue-600 text-white rounded-lg" {...props}>
      {loading ? 'Загрузка...' : children}
    </button>
  );
}
```

3. **Использовать в сервисах**
```typescript
import { LoadingButton } from '@pet/shared/components/buttons';

<LoadingButton loading={isLoading} onClick={handleSave}>
  Сохранить
</LoadingButton>
```

### Обновление существующего компонента

1. Обновить файл в `shared/src/components/`
2. Все сервисы автоматически получат обновление (импортируют из shared)
3. Проверить работу в каждом сервисе
4. Обновить CHANGELOG.md

## Troubleshooting

### Проблема: Стили не применяются

**Решение:**

1. **Проверьте путь к shared в tailwind.config.ts**
   ```typescript
   content: [
     '../../../shared/src/**/*.{js,ts,jsx,tsx}', // Правильный путь
   ]
   ```

2. **Перезапустите dev сервер**
   ```bash
   rm -rf .next
   npm run dev
   ```

3. **Проверьте что @pet/shared установлен**
   ```bash
   npm install
   ```

### Проблема: Компонент не найден

**Решение:**

1. **Проверьте импорт**
   ```typescript
   // ✅ Правильно
   import { Button } from '@pet/shared/components/buttons';
   
   // ❌ Неправильно
   import { Button } from '@pet/shared';
   ```

2. **Проверьте что файл существует**
   ```bash
   ls shared/src/components/buttons/Button.tsx
   ```

## Checklist при создании нового сервиса

- [ ] Добавить `@pet/shared` в package.json
- [ ] Настроить Tailwind config (добавить путь к shared)
- [ ] Импортировать AdminLayout из @pet/shared
- [ ] Настроить navigationItems для сервиса
- [ ] Проверить что все стили применяются
- [ ] Проверить что layout выглядит идентично другим сервисам

## Checklist при добавлении компонента в shared/

- [ ] Компонент переиспользуемый (нужен в 2+ сервисах)
- [ ] Компонент не содержит бизнес-логику конкретного сервиса
- [ ] Компонент хорошо документирован (props, примеры)
- [ ] Компонент добавлен в правильную категорию (buttons/, forms/, layout/)
- [ ] Обновлен CHANGELOG.md

---

**Помни:** `shared/` - это единый источник правды для всех переиспользуемых компонентов!
