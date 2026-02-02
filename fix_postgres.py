#!/usr/bin/env python3
"""
Автоматическое исправление PostgreSQL синтаксиса в Go handlers
Оборачивает все SQL запросы с ? в ConvertPlaceholders()
"""

import os
import re
import sys

def fix_sql_queries(content):
    """Находит и оборачивает SQL запросы с ? в ConvertPlaceholders()"""
    
    modified = content
    changes = 0
    
    # Ищем все вызовы DB методов с SQL запросами содержащими ?
    # Паттерн: (database.DB.Method|db.Method)("SQL с ?" или `SQL с ?`)
    
    # Для backtick строк (многострочные)
    pattern_backtick = r'(database\.DB\.(?:Exec|Query|QueryRow)|db\.(?:Exec|Query|QueryRow))\s*\(\s*`([^`]*\?[^`]*)`'
    
    for match in reversed(list(re.finditer(pattern_backtick, modified, re.DOTALL))):
        full_text = match.group(0)
        
        # Пропускаем если уже обернуто
        if 'ConvertPlaceholders' in full_text:
            continue
        
        method = match.group(1)
        sql = match.group(2)
        
        # Заменяем
        new_text = f'{method}(ConvertPlaceholders(`{sql}`)'
        modified = modified[:match.start()] + new_text + modified[match.end():]
        changes += 1
    
    # Для обычных строк (однострочные)
    pattern_quote = r'(database\.DB\.(?:Exec|Query|QueryRow)|db\.(?:Exec|Query|QueryRow))\s*\(\s*"([^"]*\?[^"]*)"'
    
    for match in reversed(list(re.finditer(pattern_quote, modified, re.DOTALL))):
        full_text = match.group(0)
        
        # Пропускаем если уже обернуто
        if 'ConvertPlaceholders' in full_text:
            continue
        
        method = match.group(1)
        sql = match.group(2)
        
        # Заменяем
        new_text = f'{method}(ConvertPlaceholders("{sql}")'
        modified = modified[:match.start()] + new_text + modified[match.end():]
        changes += 1
    
    return modified, changes

def fix_boolean_values(content):
    """Заменяет 1/0 на TRUE/FALSE в SQL запросах"""
    
    patterns = [
        (r'is_read\s*=\s*1', 'is_read = TRUE'),
        (r'is_read\s*=\s*0', 'is_read = FALSE'),
        (r'is_deleted\s*=\s*1', 'is_deleted = TRUE'),
        (r'is_deleted\s*=\s*0', 'is_deleted = FALSE'),
        (r'verified\s*=\s*1', 'verified = TRUE'),
        (r'verified\s*=\s*0', 'verified = FALSE'),
        (r'can_post\s*=\s*1', 'can_post = TRUE'),
        (r'can_post\s*=\s*0', 'can_post = FALSE'),
    ]
    
    modified = content
    changes = 0
    
    for pattern, replacement in patterns:
        new_content = re.sub(pattern, replacement, modified, flags=re.IGNORECASE)
        if new_content != modified:
            changes += len(re.findall(pattern, modified, flags=re.IGNORECASE))
            modified = new_content
    
    return modified, changes

def process_file(filepath):
    """Обрабатывает один файл"""
    
    print(f"Обработка: {filepath}")
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    # Исправляем SQL запросы
    content, sql_changes = fix_sql_queries(content)
    
    # Исправляем boolean значения
    content, bool_changes = fix_boolean_values(content)
    
    total_changes = sql_changes + bool_changes
    
    if total_changes > 0:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"  ✅ Исправлено: {sql_changes} SQL запросов, {bool_changes} boolean значений")
        return True
    else:
        print(f"  ⏭️  Изменений не требуется")
        return False

def main():
    """Главная функция"""
    
    handlers_dir = 'main/backend/handlers'
    
    if not os.path.exists(handlers_dir):
        print(f"❌ Директория не найдена: {handlers_dir}")
        sys.exit(1)
    
    print(f"🔍 Поиск .go файлов в {handlers_dir}")
    
    go_files = [f for f in os.listdir(handlers_dir) if f.endswith('.go')]
    
    if not go_files:
        print(f"❌ .go файлы не найдены в {handlers_dir}")
        sys.exit(1)
    
    print(f"📝 Найдено {len(go_files)} файлов\n")
    
    modified_files = 0
    
    for filename in sorted(go_files):
        filepath = os.path.join(handlers_dir, filename)
        if process_file(filepath):
            modified_files += 1
    
    print(f"\n✅ Готово! Изменено файлов: {modified_files}/{len(go_files)}")

if __name__ == '__main__':
    main()
