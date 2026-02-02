#!/usr/bin/env python3
import re
import sys

def fix_sql_queries(content):
    """Оборачивает SQL запросы с ? в ConvertPlaceholders()"""
    
    # Паттерн для поиска DB.Query/Exec/QueryRow с ? в строке
    pattern = r'(database\.DB\.(Query|Exec|QueryRow)\()("(?:[^"\\]|\\.)*\?(?:[^"\\]|\\.)*")'
    
    def replace_query(match):
        prefix = match.group(1)  # database.DB.Query(
        query_str = match.group(3)  # "SELECT ... WHERE id = ?"
        
        # Проверяем, не обёрнут ли уже в ConvertPlaceholders
        if 'ConvertPlaceholders' in match.group(0):
            return match.group(0)
        
        # Оборачиваем запрос
        return f'{prefix}ConvertPlaceholders({query_str})'
    
    # Заменяем все вхождения
    fixed = re.sub(pattern, replace_query, content)
    
    # Также исправляем is_read = 1/0 на TRUE/FALSE
    fixed = re.sub(r'is_read\s*=\s*1\b', 'is_read = TRUE', fixed)
    fixed = re.sub(r'is_read\s*=\s*0\b', 'is_read = FALSE', fixed)
    fixed = re.sub(r'is_deleted\s*=\s*1\b', 'is_deleted = TRUE', fixed)
    fixed = re.sub(r'is_deleted\s*=\s*0\b', 'is_deleted = FALSE', fixed)
    
    return fixed

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: fix_postgres.py <file>")
        sys.exit(1)
    
    filename = sys.argv[1]
    
    with open(filename, 'r') as f:
        content = f.read()
    
    fixed_content = fix_sql_queries(content)
    
    with open(filename, 'w') as f:
        f.write(fixed_content)
    
    print(f"Fixed: {filename}")
