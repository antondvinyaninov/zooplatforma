# Тест SSO для админ-панели

## Проблема
Старые токены не содержат поле `roles`, поэтому Admin Backend их не принимает.

## Решение
Нужно перелогиниться, чтобы получить новый токен с ролями.

## Шаги для проверки:

### 1. Выйдите из системы
Откройте http://localhost:3000 и нажмите "Выйти"

### 2. Войдите заново
Войдите с теми же credentials: anton@dvinyaninov.ru

### 3. Проверьте токен в DevTools
```javascript
// Откройте Console (F12) на http://localhost:3000
// Вставьте этот код:

fetch('http://localhost:8080/api/auth/verify', {
  credentials: 'include'
})
.then(r => r.json())
.then(data => {
  console.log('Токен:', data);
  if (data.data && data.data.roles) {
    console.log('✅ Роли найдены:', data.data.roles);
  } else {
    console.log('❌ Роли отсутствуют - нужно перелогиниться');
  }
});
```

### 4. Проверьте доступ к админке
```javascript
// В той же Console:

fetch('http://localhost:8083/api/admin/auth/me', {
  credentials: 'include'
})
.then(r => r.json())
.then(data => {
  console.log('Admin API:', data);
  if (data.success) {
    console.log('✅ Доступ к админке разрешен!');
  } else {
    console.log('❌ Доступ запрещен:', data.error);
  }
});
```

### 5. Получите список пользователей
```javascript
fetch('http://localhost:8083/api/admin/users', {
  credentials: 'include'
})
.then(r => r.json())
.then(data => console.log('Пользователи:', data));
```

### 6. Получите статистику
```javascript
fetch('http://localhost:8083/api/admin/stats/overview', {
  credentials: 'include'
})
.then(r => r.json())
.then(data => console.log('Статистика:', data));
```

## Альтернатива: Тест через curl

Если хотите протестировать через curl, нужно сначала получить токен:

```bash
# 1. Войдите и сохраните cookie
curl -c /tmp/cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"anton@dvinyaninov.ru","password":"ВАШ_ПАРОЛЬ"}'

# 2. Проверьте токен
curl -b /tmp/cookies.txt http://localhost:8080/api/auth/verify

# 3. Проверьте доступ к админке
curl -b /tmp/cookies.txt http://localhost:8083/api/admin/auth/me

# 4. Получите список пользователей
curl -b /tmp/cookies.txt http://localhost:8083/api/admin/users
```

## Ожидаемый результат

После перелогина вы должны увидеть:

```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "anton@dvinyaninov.ru",
    "roles": ["user", "superadmin"],
    "valid": true
  }
}
```

И доступ к админке должен работать! 🎉
