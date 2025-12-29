# NPS (Net Promoter Score) - Техническая реализация

**Версия:** 1.0  
**Дата:** 28 декабря 2025  
**Статус:** Планируется для MVP 1.0.0

---

## 📊 Что такое NPS?

**Net Promoter Score (NPS)** - метрика лояльности пользователей, основанная на одном вопросе:

> "Насколько вероятно, что вы порекомендуете ЗооПлатформу своим друзьям или коллегам?"

### Шкала оценки (0-10):

- **0-6**: Критики (Detractors) - недовольны, могут навредить репутации
- **7-8**: Нейтралы (Passives) - удовлетворены, но не лояльны
- **9-10**: Промоутеры (Promoters) - восторженные адвокаты бренда

### Формула расчёта:

```
NPS = % Промоутеров - % Критиков
```

**Диапазон:** от -100 (все критики) до +100 (все промоутеры)

### Интерпретация:

- **< 0**: Критическая ситуация, больше критиков чем промоутеров
- **0-30**: Нужны улучшения
- **30-70**: Хороший результат
- **> 70**: Отличный результат, мировой класс

---

## 🎯 Цели для ЗооПлатформы

### MVP 1.0.0 (Январь 2026):
- **NPS > 40** - хороший результат для нового продукта
- **Response rate > 30%** - процент ответивших на опрос
- **< 5% отписок** - пользователи не раздражены опросами
- **Сегментация по ролям** - понимание удовлетворённости разных групп

### Версия 1.1.0+:
- **NPS > 50** - рост лояльности
- **Тренд вверх** - улучшение со временем
- **Корреляция с retention** - связь NPS с удержанием пользователей

---

## 🗄️ Структура базы данных

### Таблица: nps_surveys

Хранит ответы пользователей на NPS опросы.

```sql
CREATE TABLE nps_surveys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    score INTEGER NOT NULL CHECK(score >= 0 AND score <= 10),
    feedback TEXT,
    trigger_type TEXT NOT NULL, -- 'day_7', 'success_story', 'day_30', 'quarterly'
    user_role TEXT, -- 'user', 'volunteer', 'shelter_admin', 'clinic_admin'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_nps_surveys_user_id ON nps_surveys(user_id);
CREATE INDEX idx_nps_surveys_score ON nps_surveys(score);
CREATE INDEX idx_nps_surveys_trigger_type ON nps_surveys(trigger_type);
CREATE INDEX idx_nps_surveys_created_at ON nps_surveys(created_at);
```

### Таблица: nps_survey_shows

Отслеживает показы опросов (для контроля частоты).

```sql
CREATE TABLE nps_survey_shows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    trigger_type TEXT NOT NULL,
    shown_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed BOOLEAN DEFAULT FALSE,
    dismissed BOOLEAN DEFAULT FALSE, -- пользователь закрыл опрос
    opted_out BOOLEAN DEFAULT FALSE, -- пользователь отписался от опросов
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_nps_shows_user_id ON nps_survey_shows(user_id);
CREATE INDEX idx_nps_shows_shown_at ON nps_survey_shows(shown_at);
```

---

## 📊 Зависимость от аналитики активности

**ВАЖНО:** NPS опросы работают на основе данных активности пользователей!

Для корректной работы триггеров необходимо:
- Отслеживать дату регистрации пользователя
- Считать количество входов и сессий
- Логировать все действия пользователя
- Отслеживать события (пристройство животных)

**Подробная документация:** [USER_ANALYTICS.md](USER_ANALYTICS.md)

Без системы аналитики NPS опросы не смогут определить:
- Когда показывать опрос (7 дней, 30 дней, квартал)
- Достаточно ли пользователь активен
- Был ли момент успеха (пристройство животного)

---

## 🔔 Триггеры показа опроса

### 1. После 7 дней использования (day_7)

**Цель:** Первое впечатление от платформы

**Условия:**
- Пользователь зарегистрирован 7 дней назад
- Был активен хотя бы 3 раза за эти 7 дней
- Не видел NPS опрос ранее

**Когда показывать:**
- При входе на платформу на 8-й день

### 2. После успешного пристройства (success_story)

**Цель:** Момент максимального удовлетворения

**Условия:**
- Пользователь создал пост с меткой "ищет дом"
- Пост был обновлён на статус "дома" (пристроен)
- Прошло 1-3 дня после пристройства

**Когда показывать:**
- Сразу после обновления статуса животного на "дома"

### 3. После 30 дней активности (day_30)

**Цель:** Устоявшееся мнение о платформе

**Условия:**
- Пользователь зарегистрирован 30 дней назад
- Был активен хотя бы 10 раз за эти 30 дней
- Не проходил опрос day_7 или прошло >21 день с последнего опроса

**Когда показывать:**
- При входе на платформу на 31-й день

### 4. Квартальный опрос (quarterly)

**Цель:** Отслеживание динамики для активных пользователей

**Условия:**
- Пользователь активен >90 дней
- Прошло >90 дней с последнего NPS опроса
- Активность >20 действий за последние 90 дней

**Когда показывать:**
- При входе на платформу после 90 дней с последнего опроса

---

## 🎨 UI компонент

### Модальное окно (Desktop)

```
┌─────────────────────────────────────────────┐
│  ✕                                          │
│                                             │
│  🐾 Помогите нам стать лучше!               │
│                                             │
│  Насколько вероятно, что вы порекомендуете  │
│  ЗооПлатформу своим друзьям или коллегам?   │
│                                             │
│  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
│  │ 0 │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │ 9 │10 │
│  └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
│  Не рекомендую                  Точно рекомендую
│                                             │
│  Расскажите, почему? (опционально)          │
│  ┌─────────────────────────────────────┐   │
│  │                                     │   │
│  │                                     │   │
│  └─────────────────────────────────────┘   │
│                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │ Отправить│  │  Позже   │  │ Больше не│ │
│  │          │  │          │  │спрашивать│ │
│  └──────────┘  └──────────┘  └──────────┘ │
└─────────────────────────────────────────────┘
```

### Slide-in панель (Mobile)

```
┌─────────────────────────┐
│  🐾 Помогите нам!       │
│                         │
│  Порекомендуете нас?    │
│                         │
│  [0][1][2][3][4][5]     │
│  [6][7][8][9][10]       │
│                         │
│  Почему? (опционально)  │
│  ┌─────────────────┐   │
│  │                 │   │
│  └─────────────────┘   │
│                         │
│  [Отправить] [Позже]   │
│  [Больше не спрашивать]│
└─────────────────────────┘
```

### Состояния UI:

1. **Начальное** - показ вопроса и шкалы
2. **Выбрана оценка** - подсветка выбранной цифры
3. **Отправка** - loader, кнопка неактивна
4. **Успех** - "Спасибо за отзыв! 🙏"
5. **Ошибка** - "Что-то пошло не так, попробуйте позже"

---

## 🔌 API Endpoints

### POST /api/nps/submit

Отправка ответа на NPS опрос.

**Request:**
```json
{
  "score": 9,
  "feedback": "Отличная платформа! Помогла пристроить собаку.",
  "trigger_type": "success_story"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Спасибо за отзыв!"
}
```

**Errors:**
- 400: Invalid score (не в диапазоне 0-10)
- 401: Unauthorized
- 429: Too many requests (rate limiting)

### GET /api/nps/should-show

Проверка, нужно ли показывать NPS опрос пользователю.

**Response:**
```json
{
  "should_show": true,
  "trigger_type": "day_7",
  "message": "Вы с нами уже неделю! Поделитесь впечатлениями?"
}
```

### POST /api/nps/dismiss

Пользователь закрыл опрос (показать позже).

**Request:**
```json
{
  "trigger_type": "day_7"
}
```

**Response:**
```json
{
  "success": true,
  "next_show_in_days": 7
}
```

### POST /api/nps/opt-out

Пользователь отписался от NPS опросов.

**Response:**
```json
{
  "success": true,
  "message": "Вы больше не будете получать опросы"
}
```

---

## 📊 Дашборд в Admin панели

### Страница: /admin/analytics/nps

#### Главные метрики (карточки):

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  NPS Score   │  │ Промоутеры   │  │  Нейтралы    │  │   Критики    │
│              │  │              │  │              │  │              │
│      42      │  │     45%      │  │     35%      │  │     20%      │
│   ↑ +5      │  │   (90 чел)   │  │   (70 чел)   │  │   (40 чел)   │
└──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘
```

#### График динамики:

```
NPS по неделям
 60 │                                    ●
 50 │                          ●───●───●
 40 │                ●───●───●
 30 │      ●───●───●
 20 │●───●
    └─────────────────────────────────────
     W1  W2  W3  W4  W5  W6  W7  W8  W9
```

#### Распределение оценок:

```
Распределение оценок (0-10)
10 ████████████████████ 25%
 9 ████████████████ 20%
 8 ████████████ 15%
 7 ████████████ 15%
 6 ████ 5%
 5 ████ 5%
 4 ██ 3%
 3 ██ 3%
 2 ██ 3%
 1 ██ 3%
 0 ██ 3%
```

#### Сегментация по ролям:

```
┌─────────────────┬─────────┬──────────────┬──────────┬──────────┐
│ Роль            │ NPS     │ Промоутеры   │ Нейтралы │ Критики  │
├─────────────────┼─────────┼──────────────┼──────────┼──────────┤
│ Владельцы       │ 45      │ 50%          │ 35%      │ 15%      │
│ Волонтёры       │ 55      │ 60%          │ 30%      │ 10%      │
│ Приюты          │ 30      │ 35%          │ 40%      │ 25%      │
│ Ветклиники      │ 40      │ 45%          │ 35%      │ 20%      │
└─────────────────┴─────────┴──────────────┴──────────┴──────────┘
```

#### Комментарии критиков (0-6):

```
┌────────────────────────────────────────────────────────────┐
│ 🔴 Оценка 3 | Владелец | 25.12.2025                        │
│ "Мало функций для приютов, нет отчётности"                 │
│ [Ответить] [Отметить как решённое]                         │
├────────────────────────────────────────────────────────────┤
│ 🔴 Оценка 5 | Волонтёр | 24.12.2025                        │
│ "Долго загружаются видео, иногда вылетает"                 │
│ [Ответить] [Отметить как решённое]                         │
└────────────────────────────────────────────────────────────┘
```

#### Комментарии промоутеров (9-10):

```
┌────────────────────────────────────────────────────────────┐
│ 🟢 Оценка 10 | Владелец | 26.12.2025                       │
│ "Помогли пристроить собаку за 2 недели! Спасибо!"          │
│ [Использовать как отзыв] [Поделиться]                      │
├────────────────────────────────────────────────────────────┤
│ 🟢 Оценка 9 | Волонтёр | 25.12.2025                        │
│ "Удобный интерфейс, всё интуитивно понятно"                │
│ [Использовать как отзыв] [Поделиться]                      │
└────────────────────────────────────────────────────────────┘
```

#### Фильтры:

- По периоду (неделя, месяц, квартал, всё время)
- По ролям (все, владельцы, волонтёры, приюты, клиники)
- По триггерам (все, day_7, success_story, day_30, quarterly)
- По оценкам (все, критики, нейтралы, промоутеры)

#### Экспорт:

- CSV (все ответы)
- PDF (отчёт с графиками)
- Excel (детальная таблица)

---

## 🔄 Логика работы

### Backend (Go)

#### 1. Проверка, нужно ли показывать опрос

```go
// handlers/nps.go
func ShouldShowNPSSurvey(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    // Проверяем, не отписался ли пользователь
    var optedOut bool
    db.QueryRow(`
        SELECT opted_out FROM nps_survey_shows 
        WHERE user_id = ? AND opted_out = TRUE 
        LIMIT 1
    `, userID).Scan(&optedOut)
    
    if optedOut {
        c.JSON(200, gin.H{"should_show": false})
        return
    }
    
    // Проверяем триггеры
    trigger := checkNPSTriggers(userID)
    
    if trigger != "" {
        c.JSON(200, gin.H{
            "should_show": true,
            "trigger_type": trigger,
            "message": getTriggerMessage(trigger),
        })
    } else {
        c.JSON(200, gin.H{"should_show": false})
    }
}

func checkNPSTriggers(userID int) string {
    // Проверяем day_7
    if shouldShowDay7(userID) {
        return "day_7"
    }
    
    // Проверяем success_story
    if shouldShowSuccessStory(userID) {
        return "success_story"
    }
    
    // Проверяем day_30
    if shouldShowDay30(userID) {
        return "day_30"
    }
    
    // Проверяем quarterly
    if shouldShowQuarterly(userID) {
        return "quarterly"
    }
    
    return ""
}
```

#### 2. Сохранение ответа

```go
func SubmitNPSSurvey(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    var req struct {
        Score       int    `json:"score" binding:"required,min=0,max=10"`
        Feedback    string `json:"feedback"`
        TriggerType string `json:"trigger_type" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Получаем роль пользователя
    var userRole string
    db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&userRole)
    
    // Сохраняем ответ
    _, err := db.Exec(`
        INSERT INTO nps_surveys (user_id, score, feedback, trigger_type, user_role)
        VALUES (?, ?, ?, ?, ?)
    `, userID, req.Score, req.Feedback, req.TriggerType, userRole)
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to save survey"})
        return
    }
    
    // Обновляем статус показа
    db.Exec(`
        UPDATE nps_survey_shows 
        SET completed = TRUE 
        WHERE user_id = ? AND trigger_type = ?
    `, userID, req.TriggerType)
    
    c.JSON(200, gin.H{
        "success": true,
        "message": "Спасибо за отзыв!",
    })
}
```

### Frontend (React/Next.js)

#### Компонент NPSSurvey

```typescript
// components/NPSSurvey.tsx
import { useState, useEffect } from 'react';

export default function NPSSurvey() {
  const [show, setShow] = useState(false);
  const [score, setScore] = useState<number | null>(null);
  const [feedback, setFeedback] = useState('');
  const [triggerType, setTriggerType] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    checkShouldShow();
  }, []);

  const checkShouldShow = async () => {
    const res = await fetch('/api/nps/should-show');
    const data = await res.json();
    
    if (data.should_show) {
      setShow(true);
      setTriggerType(data.trigger_type);
    }
  };

  const handleSubmit = async () => {
    if (score === null) return;
    
    setLoading(true);
    
    await fetch('/api/nps/submit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ score, feedback, trigger_type: triggerType }),
    });
    
    setShow(false);
    // Показать toast "Спасибо за отзыв!"
  };

  const handleDismiss = async () => {
    await fetch('/api/nps/dismiss', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ trigger_type: triggerType }),
    });
    
    setShow(false);
  };

  const handleOptOut = async () => {
    await fetch('/api/nps/opt-out', { method: 'POST' });
    setShow(false);
  };

  if (!show) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full">
        <h3 className="text-xl font-bold mb-4">
          🐾 Помогите нам стать лучше!
        </h3>
        
        <p className="mb-4">
          Насколько вероятно, что вы порекомендуете ЗооПлатформу?
        </p>
        
        <div className="flex gap-2 mb-4">
          {[0,1,2,3,4,5,6,7,8,9,10].map(n => (
            <button
              key={n}
              onClick={() => setScore(n)}
              className={`w-10 h-10 rounded ${
                score === n ? 'bg-blue-500 text-white' : 'bg-gray-200'
              }`}
            >
              {n}
            </button>
          ))}
        </div>
        
        <textarea
          placeholder="Расскажите, почему? (опционально)"
          value={feedback}
          onChange={(e) => setFeedback(e.target.value)}
          className="w-full border rounded p-2 mb-4"
          rows={3}
        />
        
        <div className="flex gap-2">
          <button
            onClick={handleSubmit}
            disabled={score === null || loading}
            className="bg-blue-500 text-white px-4 py-2 rounded"
          >
            Отправить
          </button>
          
          <button
            onClick={handleDismiss}
            className="bg-gray-200 px-4 py-2 rounded"
          >
            Позже
          </button>
          
          <button
            onClick={handleOptOut}
            className="text-sm text-gray-500 underline"
          >
            Больше не спрашивать
          </button>
        </div>
      </div>
    </div>
  );
}
```

---

## 📈 Аналитика и отчёты

### Расчёт NPS

```go
// handlers/nps_analytics.go
func CalculateNPS(startDate, endDate time.Time, role string) NPSResult {
    query := `
        SELECT 
            COUNT(CASE WHEN score >= 9 THEN 1 END) as promoters,
            COUNT(CASE WHEN score >= 7 AND score <= 8 THEN 1 END) as passives,
            COUNT(CASE WHEN score <= 6 THEN 1 END) as detractors,
            COUNT(*) as total
        FROM nps_surveys
        WHERE created_at BETWEEN ? AND ?
    `
    
    if role != "" {
        query += " AND user_role = ?"
    }
    
    var result NPSResult
    // Execute query...
    
    result.PromoterPercent = float64(result.Promoters) / float64(result.Total) * 100
    result.DetractorPercent = float64(result.Detractors) / float64(result.Total) * 100
    result.NPS = result.PromoterPercent - result.DetractorPercent
    
    return result
}
```

### Сегментация

```go
func GetNPSBySegment() []SegmentNPS {
    roles := []string{"user", "volunteer", "shelter_admin", "clinic_admin"}
    results := []SegmentNPS{}
    
    for _, role := range roles {
        nps := CalculateNPS(startDate, endDate, role)
        results = append(results, SegmentNPS{
            Role: role,
            NPS: nps,
        })
    }
    
    return results
}
```

### Экспорт комментариев

```go
func ExportDetractorComments() []Comment {
    rows, _ := db.Query(`
        SELECT u.name, ns.score, ns.feedback, ns.created_at
        FROM nps_surveys ns
        JOIN users u ON ns.user_id = u.id
        WHERE ns.score <= 6 AND ns.feedback != ''
        ORDER BY ns.created_at DESC
    `)
    
    // Parse and return...
}
```

---

## ✅ Checklist реализации

### ⚠️ Зависимости (КРИТИЧНО - делать в первую очередь):
- [ ] **Реализовать систему аналитики активности** (см. [USER_ANALYTICS.md](USER_ANALYTICS.md))
  - [ ] Добавить поля в users: registered_at, last_login_at, login_count, total_sessions
  - [ ] Создать таблицу user_sessions
  - [ ] Создать таблицу user_activity_log
  - [ ] Создать таблицу user_stats
  - [ ] Реализовать middleware для отслеживания сессий
  - [ ] Реализовать функцию LogActivity()
  - [ ] Добавить логирование во все handlers
  - [ ] Реализовать cron jobs (закрытие сессий, агрегация)

### Backend:
- [ ] Создать миграцию для таблиц `nps_surveys` и `nps_survey_shows`
- [ ] Реализовать API endpoints (submit, should-show, dismiss, opt-out)
- [ ] Реализовать логику триггеров (day_7, success_story, day_30, quarterly) **используя данные аналитики**
- [ ] Реализовать расчёт NPS
- [ ] Реализовать сегментацию по ролям
- [ ] Добавить rate limiting для защиты от спама

### Frontend:
- [ ] Создать компонент NPSSurvey (модальное окно)
- [ ] Адаптировать для мобильных устройств
- [ ] Интегрировать в layout (проверка при загрузке страницы)
- [ ] Добавить анимации (появление/исчезновение)
- [ ] Добавить toast уведомление после отправки

### Admin панель:
- [ ] Создать страницу /admin/analytics/nps
- [ ] Реализовать главные метрики (карточки)
- [ ] Реализовать график динамики
- [ ] Реализовать распределение оценок
- [ ] Реализовать сегментацию по ролям
- [ ] Реализовать список комментариев (критики и промоутеры)
- [ ] Добавить фильтры (период, роль, триггер, оценка)
- [ ] Реализовать экспорт (CSV, PDF, Excel)

### Тестирование:
- [ ] Unit тесты для расчёта NPS
- [ ] Unit тесты для логики триггеров
- [ ] Integration тесты для API endpoints
- [ ] E2E тест: показ опроса и отправка ответа
- [ ] Тестирование на разных ролях
- [ ] Тестирование opt-out функции

### Документация:
- [ ] Обновить API документацию
- [ ] Добавить примеры использования
- [ ] Документировать триггеры
- [ ] Создать руководство для админов

---

## 🚀 План внедрения

### Этап 1: Базовая реализация (3 дня)
- День 1: База данных + API endpoints
- День 2: Frontend компонент
- День 3: Тестирование

### Этап 2: Аналитика (2 дня)
- День 4: Дашборд в Admin панели
- День 5: Графики и сегментация

### Этап 3: Полировка (1 день)
- День 6: Экспорт, фильтры, документация

**Итого:** 6 дней (0.5 недели с учётом параллельной работы)

---

## 📊 Метрики успеха реализации

- [ ] NPS опрос показывается корректно всем триггерам
- [ ] Response rate >30% (пользователи отвечают)
- [ ] <5% opt-out (мало отписок)
- [ ] Дашборд работает без ошибок
- [ ] Экспорт данных работает корректно
- [ ] Нет критичных багов

---

## 🔮 Будущие улучшения (v1.1.0+)

### Расширенная аналитика:
- Корреляция NPS с retention rate
- Корреляция NPS с активностью пользователя
- Предсказание churn на основе NPS
- A/B тестирование формулировок вопроса

### Автоматизация:
- Автоматические email-рассылки критикам (предложение помощи)
- Автоматические благодарности промоутерам
- Автоматическое создание задач на основе комментариев критиков

### Интеграция:
- Интеграция с CRM (если появится)
- Интеграция с системой поддержки
- Webhook для отправки NPS в внешние системы

### Геймификация:
- Бейдж "Помог улучшить платформу" за прохождение NPS
- Показ влияния: "Ваш отзыв помог нам добавить функцию X"

---

**Документ создан:** 28 декабря 2025  
**Автор:** ЗооПлатформа Team  
**Версия:** 1.0
