#!/bin/bash
# Установка пароля "123" для пользователя

# Генерируем bcrypt хеш для пароля "123"
HASH='$2a$10$Z31MaRaEMTDXiQhiyljoue/skO.uhY0U/uHbdHiTjwnYR85pP/uz6'

# Обновляем пароль в Auth Service базе
sqlite3 auth.db "UPDATE users SET password_hash = '$HASH' WHERE email = 'anton@dvinyaninov.ru'"

echo "✅ Password updated for anton@dvinyaninov.ru"
