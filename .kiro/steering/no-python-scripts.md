---
inclusion: manual
---

# ⛔️ КРИТИЧЕСКОЕ ПРАВИЛО: НЕ ИСПОЛЬЗУЙ PYTHON СКРИПТЫ

## ЗАПРЕЩЕНО

**НИКОГДА не используй Python скрипты для модификации файлов проекта!**

### Почему:
- Python скрипты часто ломают файлы
- AWK скрипты удаляют содержимое файлов
- Sed может неправильно обработать многострочные блоки
- Автоматические замены могут сломать синтаксис

### Что случилось:
- AWK скрипт удалил содержимое `comments.go`, `pets.go`, `organizations.go`, `friends.go`
- Пришлось восстанавливать из git
- Потеряно время на восстановление

## ✅ ПРАВИЛЬНО

### Для массовых изменений:
1. **Используй `strReplace` tool** - безопасная замена строк
2. **Делай по одному файлу** - контролируемо и безопасно
3. **Проверяй компиляцию** после каждого изменения
4. **Коммить часто** - чтобы можно было откатить

### Для сложных изменений:
1. **Читай файл** - `readFile`
2. **Анализируй структуру** - понимай что меняешь
3. **Используй `strReplace`** - точечная замена
4. **Проверяй результат** - компиляция, тесты

### Пример ПРАВИЛЬНОГО подхода:

```typescript
// 1. Читаем файл
readFile('main/backend/handlers/comments.go')

// 2. Находим нужное место
// Видим: result, err := db.Exec(...)
//        commentID, _ := result.LastInsertId()

// 3. Используем strReplace
strReplace({
  path: 'main/backend/handlers/comments.go',
  oldStr: `result, err := db.Exec(ConvertPlaceholders(\`
    INSERT INTO comments (...)
  \`), ...)
  
  if err != nil {
    ...
  }
  
  commentID, _ := result.LastInsertId()`,
  newStr: `var commentID int64
  err = db.QueryRow(ConvertPlaceholders(\`
    INSERT INTO comments (...)
    RETURNING id
  \`), ...).Scan(&commentID)
  
  if err != nil {
    ...
  }`
})

// 4. Проверяем компиляцию
executeBash('go build', cwd: 'main/backend')

// 5. Коммитим
executeBash('git add -A && git commit -m "Fix comments.go"')
```

## ❌ НЕПРАВИЛЬНО

```bash
# НЕ ДЕЛАЙ ТАК!
cat > /tmp/fix.py << 'EOF'
import re
with open('file.go', 'r') as f:
    content = f.read()
content = re.sub(r'pattern', 'replacement', content)
with open('file.go', 'w') as f:
    f.write(content)
EOF
python3 /tmp/fix.py
```

```bash
# НЕ ДЕЛАЙ ТАК!
awk '/pattern/ { ... }' file.go > temp && mv temp file.go
```

```bash
# НЕ ДЕЛАЙ ТАК!
sed -i 's/pattern/replacement/g' file.go
```

## Исключения

**Единственный допустимый случай использования скриптов:**
- Простые операции чтения (не записи!)
- Анализ структуры проекта
- Генерация отчетов

**Но даже в этих случаях предпочитай встроенные tools!**

## Checklist перед массовым изменением

- [ ] Можно ли сделать через `strReplace`?
- [ ] Можно ли сделать по одному файлу?
- [ ] Есть ли резервная копия (git commit)?
- [ ] Проверена ли компиляция после изменения?
- [ ] Точно ли нужен скрипт, или можно вручную?

**Помни:** Лучше потратить 10 минут на ручное исправление, чем 30 минут на восстановление сломанных файлов!

---

**ПРАВИЛО:** Если хочешь использовать Python/AWK/Sed скрипт - ОСТАНОВИСЬ и используй `strReplace` вместо этого!
