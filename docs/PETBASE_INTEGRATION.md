# Интеграция ЗооБазы с основным сайтом

## Обзор

ЗооБаза (PetID) - это отдельный микросервис, который является единым источником правды (Single Source of Truth) для всех данных о питомцах. Основной сайт запрашивает данные через API.

## Архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                      Основной сайт                          │
│                   (localhost:3000)                          │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Профиль    │  │    Посты     │  │   Питомцы    │    │
│  │ пользователя │  │              │  │              │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│         │                 │                  │             │
│         └─────────────────┴──────────────────┘             │
│                           │                                │
│                    API запросы                             │
└───────────────────────────┼────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                        ЗооБаза                              │
│                   (localhost:8100)                          │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   REST API                           │  │
│  │  /api/pets, /api/species, /api/breeds              │  │
│  └──────────────────────────────────────────────────────┘  │
│                           │                                │
│                           ▼                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              База данных SQLite                      │  │
│  │         (единый источник правды)                     │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Этапы интеграции

### Этап 1: Настройка CORS

ЗооБаза должна разрешать запросы с основного сайта.

**Файл:** `petbase/backend/main.go`

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Разрешаем запросы с основного сайта
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        // Обработка preflight запросов
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    // ... существующий код ...
    
    // Применяем CORS middleware
    handler := corsMiddleware(mux)
    
    log.Printf("PetBase API server starting on :8100")
    log.Fatal(http.ListenAndServe(":8100", handler))
}
```

### Этап 2: API клиент на основном сайте

Создаём клиент для работы с API ЗооБазы.

**Файл:** `main/frontend/lib/petbase-api.ts`

```typescript
const PETBASE_API_URL = process.env.NEXT_PUBLIC_PETBASE_API_URL || 'http://localhost:8100/api';

export interface Pet {
  id: number;
  user_id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  size?: string;
  weight?: number;
  chip_number?: string;
  tattoo_number?: string;
  ear_tag_number?: string;
  passport_number?: string;
  is_sterilized: boolean;
  sterilization_date?: string;
  is_vaccinated: boolean;
  health_notes?: string;
  character_traits?: string;
  special_needs?: string;
  status: string;
  status_updated_at?: string;
  photo?: string;
  photos?: string;
  story?: string;
  created_at: string;
  updated_at: string;
  // Поля паспорта
  distinctive_marks?: string;
  owner_name?: string;
  owner_address?: string;
  owner_phone?: string;
  owner_email?: string;
  blood_type?: string;
  allergies?: string;
  chronic_diseases?: string;
  current_medications?: string;
  pedigree_number?: string;
  registration_org?: string;
  // Куратор и локация
  curator_id?: number;
  curator_name?: string;
  curator_phone?: string;
  location?: string;
  foster_address?: string;
  shelter_name?: string;
}

export interface Species {
  id: number;
  name: string;
  name_en: string;
  description?: string;
  icon?: string;
  created_at: string;
}

export interface Breed {
  id: number;
  species_id: number;
  species_name: string;
  name: string;
  name_en: string;
  description?: string;
  origin?: string;
  size?: string;
  weight_min?: number;
  weight_max?: number;
  lifespan_min?: number;
  lifespan_max?: number;
  temperament?: string;
  care_level?: string;
  photo?: string;
  created_at: string;
}

class PetBaseAPI {
  private baseURL: string;

  constructor(baseURL: string = PETBASE_API_URL) {
    this.baseURL = baseURL;
  }

  // Питомцы
  async getPets(): Promise<Pet[]> {
    const response = await fetch(`${this.baseURL}/pets`);
    const data = await response.json();
    return data.success ? data.data : [];
  }

  async getPet(id: number): Promise<Pet | null> {
    const response = await fetch(`${this.baseURL}/pets/${id}`);
    const data = await response.json();
    return data.success ? data.data : null;
  }

  async getPetsByUser(userId: number): Promise<Pet[]> {
    const response = await fetch(`${this.baseURL}/pets/user/${userId}`);
    const data = await response.json();
    return data.success ? data.data : [];
  }

  async getPetsByStatus(status: string): Promise<Pet[]> {
    const response = await fetch(`${this.baseURL}/pets/status/${status}`);
    const data = await response.json();
    return data.success ? data.data : [];
  }

  async createPet(pet: Partial<Pet>): Promise<Pet | null> {
    const response = await fetch(`${this.baseURL}/pets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(pet),
    });
    const data = await response.json();
    return data.success ? data.data : null;
  }

  async updatePet(id: number, pet: Partial<Pet>): Promise<Pet | null> {
    const response = await fetch(`${this.baseURL}/pets/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(pet),
    });
    const data = await response.json();
    return data.success ? data.data : null;
  }

  async deletePet(id: number): Promise<boolean> {
    const response = await fetch(`${this.baseURL}/pets/${id}`, {
      method: 'DELETE',
    });
    const data = await response.json();
    return data.success;
  }

  // Виды
  async getSpecies(): Promise<Species[]> {
    const response = await fetch(`${this.baseURL}/species`);
    const data = await response.json();
    return data.success ? data.data : [];
  }

  // Породы
  async getBreeds(): Promise<Breed[]> {
    const response = await fetch(`${this.baseURL}/breeds`);
    const data = await response.json();
    return data.success ? data.data : [];
  }

  async getBreedsBySpecies(speciesId: number): Promise<Breed[]> {
    const response = await fetch(`${this.baseURL}/breeds/species/${speciesId}`);
    const data = await response.json();
    return data.success ? data.data : [];
  }
}

export const petBaseAPI = new PetBaseAPI();
```

### Этап 3: Использование на основном сайте

#### 3.1. Профиль пользователя - список питомцев

**Файл:** `main/frontend/app/profile/page.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { petBaseAPI, Pet } from '@/lib/petbase-api';

export default function ProfilePage() {
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Получаем ID текущего пользователя из сессии
  const userId = 1; // TODO: получить из сессии

  useEffect(() => {
    const loadPets = async () => {
      try {
        const userPets = await petBaseAPI.getPetsByUser(userId);
        setPets(userPets);
      } catch (error) {
        console.error('Error loading pets:', error);
      } finally {
        setLoading(false);
      }
    };

    loadPets();
  }, [userId]);

  return (
    <div>
      <h1>Мои питомцы</h1>
      {loading ? (
        <p>Загрузка...</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {pets.map(pet => (
            <div key={pet.id} className="border rounded-lg p-4">
              <h3>{pet.name}</h3>
              <p>{pet.species} {pet.breed && `• ${pet.breed}`}</p>
              <a href={`http://localhost:4100/pets/${pet.id}`} target="_blank">
                Открыть карточку
              </a>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

#### 3.2. Создание поста с привязкой питомца

**Файл:** `main/frontend/app/posts/create/page.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { petBaseAPI, Pet } from '@/lib/petbase-api';

export default function CreatePostPage() {
  const [pets, setPets] = useState<Pet[]>([]);
  const [selectedPetId, setSelectedPetId] = useState<number | null>(null);
  const userId = 1; // TODO: получить из сессии

  useEffect(() => {
    const loadPets = async () => {
      const userPets = await petBaseAPI.getPetsByUser(userId);
      setPets(userPets);
    };
    loadPets();
  }, [userId]);

  return (
    <form>
      <h1>Создать пост</h1>
      
      {/* Выбор питомца */}
      <div>
        <label>Выберите питомца</label>
        <select 
          value={selectedPetId || ''} 
          onChange={(e) => setSelectedPetId(Number(e.target.value))}
        >
          <option value="">Без питомца</option>
          {pets.map(pet => (
            <option key={pet.id} value={pet.id}>
              {pet.name} ({pet.species})
            </option>
          ))}
        </select>
      </div>

      {/* Остальные поля поста */}
      {/* ... */}
    </form>
  );
}
```

#### 3.3. Отображение питомца в посте

**Файл:** `main/frontend/app/components/PetCard.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { petBaseAPI, Pet } from '@/lib/petbase-api';

interface PetCardProps {
  petId: number;
}

export default function PetCard({ petId }: PetCardProps) {
  const [pet, setPet] = useState<Pet | null>(null);

  useEffect(() => {
    const loadPet = async () => {
      const petData = await petBaseAPI.getPet(petId);
      setPet(petData);
    };
    loadPet();
  }, [petId]);

  if (!pet) return <div>Загрузка...</div>;

  return (
    <div className="border rounded-lg p-4 bg-gray-50">
      <div className="flex items-center gap-3">
        <div className="text-3xl">🐾</div>
        <div>
          <h4 className="font-semibold">{pet.name}</h4>
          <p className="text-sm text-gray-600">
            {pet.species} {pet.breed && `• ${pet.breed}`}
          </p>
        </div>
      </div>
      <a 
        href={`http://localhost:4100/pets/${pet.id}`} 
        target="_blank"
        className="text-blue-600 hover:text-blue-700 text-sm mt-2 inline-block"
      >
        Открыть карточку →
      </a>
    </div>
  );
}
```

### Этап 4: Переменные окружения

**Файл:** `main/frontend/.env.local`

```env
NEXT_PUBLIC_PETBASE_API_URL=http://localhost:8100/api
NEXT_PUBLIC_PETBASE_FRONTEND_URL=http://localhost:4100
```

**Файл:** `petbase/frontend/.env.local`

```env
NEXT_PUBLIC_MAIN_SITE_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8100/api
```

## Сценарии использования

### 1. Пользователь создаёт питомца

```
Основной сайт → POST /api/pets → ЗооБаза
                                    ↓
                              Сохранение в БД
                                    ↓
                              Возврат данных
```

### 2. Пользователь создаёт пост о питомце

```
Основной сайт:
1. GET /api/pets/user/{userId} → получить список питомцев
2. Пользователь выбирает питомца
3. POST /api/posts (с pet_id) → сохранить пост в основной БД
```

### 3. Отображение поста с питомцем

```
Основной сайт:
1. GET /api/posts → получить посты из основной БД
2. Для каждого поста с pet_id:
   GET /api/pets/{pet_id} → получить данные питомца из ЗооБазы
3. Отобразить пост с карточкой питомца
```

### 4. Редактирование питомца

```
Пользователь переходит на http://localhost:4100/pets/{id}
                                    ↓
                            Админка ЗооБазы
                                    ↓
                        PUT /api/pets/{id} → обновление
                                    ↓
                        Изменения сразу видны на основном сайте
```

## Безопасность

### Аутентификация

1. **Токены JWT**: Основной сайт передаёт JWT токен в заголовке Authorization
2. **Проверка владельца**: ЗооБаза проверяет, что пользователь может редактировать только своих питомцев

```go
// petbase/backend/middleware/auth.go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        // Проверка токена
        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Добавляем userID в контекст
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### CORS

- Разрешаем запросы только с `http://localhost:3000` (основной сайт)
- В продакшене: только с домена основного сайта

## Развёртывание

### Development

```bash
# Терминал 1: ЗооБаза Backend
cd petbase/backend
air

# Терминал 2: ЗооБаза Frontend
cd petbase/frontend
npm run dev

# Терминал 3: Основной сайт
cd main/frontend
npm run dev
```

### Production

```bash
# ЗооБаза Backend
cd petbase/backend
go build -o petbase-server
./petbase-server

# ЗооБаза Frontend
cd petbase/frontend
npm run build
npm start

# Основной сайт
cd main/frontend
npm run build
npm start
```

## Мониторинг

### Логирование

- Все API запросы логируются
- Ошибки отправляются в систему мониторинга

### Метрики

- Количество запросов к API
- Время ответа
- Ошибки

## Следующие шаги

1. ✅ Создать CORS middleware в ЗооБазе
2. ✅ Создать API клиент на основном сайте
3. ⏳ Добавить аутентификацию через JWT
4. ⏳ Интегрировать в профиль пользователя
5. ⏳ Интегрировать в создание постов
6. ⏳ Добавить кэширование данных
7. ⏳ Настроить мониторинг

---

**Последнее обновление:** 29 декабря 2024
