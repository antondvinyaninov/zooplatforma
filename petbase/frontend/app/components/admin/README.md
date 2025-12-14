# Admin Dashboard Template

Переиспользуемый шаблон для создания админ-панелей с единым дизайном.

## ⚠️ ВАЖНО: Tailwind CSS v4

Эти компоненты используют Tailwind классы, поэтому их нужно **копировать** в каждый сервис, а не импортировать из @pet/shared!

```bash
# Копируем компоненты админки в новый сервис
cp -r shared/src/components/admin your-service/frontend/app/components/
```

## Компоненты

### AdminLayout

Основной layout с header, sidebar и content area.

**Пример использования:**

```typescript
import AdminLayout, { AdminTab } from './components/admin/AdminLayout';
import { UsersIcon, ChartBarIcon } from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'users', label: 'Пользователи', icon: <UsersIcon className="w-5 h-5" /> },
  { id: 'stats', label: 'Статистика', icon: <ChartBarIcon className="w-5 h-5" /> },
];

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState('users');
  
  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText="Моя Админка"
      logoAlt="Админка"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      adminUser={{ email: 'admin@example.com', role: 'superadmin' }}
      onLogout={() => router.push('/auth')}
      mainSiteUrl="http://localhost:3000"
    >
      {/* Ваш контент */}
    </AdminLayout>
  );
}
```

### Виджеты

#### StatsWidget

Карточка со статистикой.

```typescript
<StatsWidget
  title="Всего пользователей"
  value={1234}
  icon={<UsersIcon className="w-6 h-6" />}
  color="blue"
  trend={{ value: '+12% за месяц', isPositive: true }}
/>
```

**Props:**
- `title` - заголовок
- `value` - значение (число или строка)
- `icon` - иконка (опционально)
- `color` - цвет: 'blue' | 'green' | 'orange' | 'purple'
- `trend` - тренд (опционально): `{ value: string, isPositive: boolean }`

#### TableWidget

Виджет с таблицей.

```typescript
<TableWidget
  title="Пользователи (100)"
  actions={
    <>
      <input type="text" placeholder="Поиск..." />
      <button>Создать</button>
    </>
  }
>
  <table>...</table>
</TableWidget>
```

**Props:**
- `title` - заголовок
- `children` - содержимое (обычно таблица)
- `actions` - кнопки/фильтры в header (опционально)

#### ChartWidget

Виджет для графиков.

```typescript
<ChartWidget
  title="Активность пользователей"
  subtitle="За последние 7 дней"
>
  <YourChart />
</ChartWidget>
```

**Props:**
- `title` - заголовок
- `subtitle` - подзаголовок (опционально)
- `children` - содержимое (график)

## Адаптивность

Шаблон полностью адаптивен:
- На экранах < 1024px sidebar скрывается
- Появляется кнопка-гамбургер для открытия меню
- Padding контента уменьшается на мобильных
- Все виджеты адаптируются под размер экрана

## Пример полной страницы

```typescript
'use client';

import { useState } from 'react';
import AdminLayout, { AdminTab } from './components/admin/AdminLayout';
import StatsWidget from './components/admin/widgets/StatsWidget';
import TableWidget from './components/admin/widgets/TableWidget';
import { UsersIcon, ChartBarIcon } from '@heroicons/react/24/outline';

const tabs: AdminTab[] = [
  { id: 'stats', label: 'Статистика', icon: <ChartBarIcon className="w-5 h-5" /> },
  { id: 'users', label: 'Пользователи', icon: <UsersIcon className="w-5 h-5" /> },
];

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState('stats');

  return (
    <AdminLayout
      logoSrc="/favicon.svg"
      logoText="Админка"
      logoAlt="Админка"
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      adminUser={{ email: 'admin@example.com', role: 'admin' }}
      onLogout={() => console.log('logout')}
    >
      {activeTab === 'stats' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Статистика</h2>
            <p className="text-gray-600">Основные метрики</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatsWidget
              title="Всего пользователей"
              value={1234}
              icon={<UsersIcon className="w-6 h-6" />}
              color="blue"
            />
          </div>
        </div>
      )}

      {activeTab === 'users' && (
        <div className="space-y-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Пользователи</h2>
            <p className="text-gray-600">Управление пользователями</p>
          </div>

          <TableWidget title="Пользователи">
            <table className="w-full">
              {/* Ваша таблица */}
            </table>
          </TableWidget>
        </div>
      )}
    </AdminLayout>
  );
}
```

## Стилизация

Все компоненты используют Tailwind CSS v4. Цветовая схема:
- Фон: `#f5f7fa`
- Карточки: белые с `shadow-sm` и `border-gray-100`
- Акцент: синий `#2787f5`
- Текст: `text-gray-900` для заголовков, `text-gray-600` для описаний

## Создание нового сервиса

1. Скопируйте компоненты:
   ```bash
   cp -r shared/src/components/admin your-service/frontend/app/components/
   ```

2. Создайте страницу с AdminLayout

3. Добавьте свои вкладки и контент

4. Готово! 🎉
