// Утилиты для работы с медиа

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

/**
 * Получить полный URL для медиа файла
 * @param url - относительный URL из базы данных (например, /uploads/users/1/avatars/file.png)
 * @returns полный URL (например, http://localhost:8000/uploads/users/1/avatars/file.png)
 */
export function getMediaUrl(url: string | null | undefined): string | null {
  if (!url) return null;
  
  // Если URL уже полный, возвращаем как есть
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  
  // Добавляем базовый URL
  return `${API_URL}${url}`;
}
