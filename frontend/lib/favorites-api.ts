// API для работы с избранными питомцами

const API_URL = process.env.NEXT_PUBLIC_API_URL || 
  (typeof window !== 'undefined' && process.env.NODE_ENV === 'production'
    ? '' // Пустая строка = относительные пути (/api/...)
    : 'http://localhost:8000');

export interface Favorite {
  id: number;
  user_id: number;
  pet_id: number;
  created_at: string;
}

// Получить список избранных питомцев
export async function getFavorites(): Promise<Favorite[]> {
  const response = await fetch(`${API_URL}/api/favorites`, {
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch favorites');
  }

  const data = await response.json();
  return data.data || [];
}

// Добавить питомца в избранное
export async function addFavorite(petId: number): Promise<Favorite> {
  const response = await fetch(`${API_URL}/api/favorites`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ pet_id: petId }),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to add favorite');
  }

  const data = await response.json();
  return data.data;
}

// Удалить питомца из избранного
export async function removeFavorite(petId: number): Promise<void> {
  const response = await fetch(`${API_URL}/api/favorites/${petId}`, {
    method: 'DELETE',
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Failed to remove favorite');
  }
}
