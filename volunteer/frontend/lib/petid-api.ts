/**
 * PetID API Client
 * 
 * Клиент для взаимодействия с питомцами через Volunteer backend
 */

const VOLUNTEER_API_URL = process.env.NEXT_PUBLIC_VOLUNTEER_API_URL || 'http://localhost:8500';
const AUTH_SERVICE_URL = 'http://localhost:7100';

/**
 * Получить JWT токен через Auth Service
 */
async function getAuthToken(): Promise<string | null> {
  try {
    console.log('🔍 Getting auth token from Auth Service...');
    const response = await fetch(`${AUTH_SERVICE_URL}/api/auth/me`, {
      credentials: 'include',
    });
    
    if (!response.ok) {
      console.error('❌ Failed to get auth token:', response.status);
      return null;
    }
    
    const data = await response.json();
    const token = data.data?.token || null;
    
    if (token) {
      console.log('✅ Auth token received');
    } else {
      console.error('❌ No token in response:', data);
    }
    
    return token;
  } catch (error) {
    console.error('❌ Error getting auth token:', error);
    return null;
  }
}

/**
 * Базовая функция для выполнения запросов к Volunteer API
 */
async function volunteerFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<{ data?: T; error?: string }> {
  const token = await getAuthToken();
  
  if (!token) {
    console.error('❌ No auth token available for request:', endpoint);
    return { error: 'Unauthorized' };
  }
  
  console.log('🚀 Making request to Volunteer API:', endpoint);
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...(options.headers as Record<string, string>),
  };
  
  try {
    const response = await fetch(`${VOLUNTEER_API_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include',
    });
    
    console.log('📥 Response from Volunteer API:', response.status, endpoint);
    
    if (response.status === 401) {
      return { error: 'Unauthorized' };
    }
    
    if (response.status === 403) {
      return { error: 'Access denied' };
    }
    
    if (response.status === 404) {
      return { error: 'Not found' };
    }
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return { 
        error: errorData.message || errorData.error || 'Request failed' 
      };
    }
    
    const data = await response.json();
    return { data };
    
  } catch (error) {
    console.error('Volunteer API Error:', error);
    return { 
      error: error instanceof Error ? error.message : 'Network error' 
    };
  }
}

// Pet Events Types
export type PetEventType = 
  | 'general'
  | 'vaccination'
  | 'treatment'
  | 'ownership_change'
  | 'lost'
  | 'found'
  | 'death'
  | 'shelter_intake'
  | 'adoption';

export const eventTypeLabels: Record<PetEventType, string> = {
  general: 'Общее событие',
  vaccination: 'Вакцинация',
  treatment: 'Лечение',
  ownership_change: 'Смена владельца',
  lost: 'Потеря',
  found: 'Найден',
  death: 'Смерть',
  shelter_intake: 'Поступление в приют',
  adoption: 'Усыновление',
};

export const eventTypeIcons: Record<PetEventType, string> = {
  general: '📝',
  vaccination: '💉',
  treatment: '💊',
  ownership_change: '👥',
  lost: '🔍',
  found: '✅',
  death: '🕊️',
  shelter_intake: '🏠',
  adoption: '❤️',
};

export const deathReasonLabels: Record<string, string> = {
  natural: 'Естественная смерть',
  illness: 'Болезнь',
  accident: 'Несчастный случай',
  euthanasia: 'Эвтаназия',
  unknown: 'Неизвестно',
};

// Pet Events API
export const petEventsApi = {
  async getEvents(petId: number): Promise<{ data?: any; error?: string; success?: boolean }> {
    try {
      const response = await volunteerFetch<any>(
        `/api/pets/${petId}/events`
      );
      if (response.error) {
        return { error: response.error };
      }
      return { success: true, data: response.data };
    } catch (error) {
      return { error: 'Ошибка загрузки событий' };
    }
  },

  async createEvent(petId: number, data: any): Promise<{ data?: any; error?: string; success?: boolean }> {
    try {
      const response = await volunteerFetch(
        `/api/pets/${petId}/events`,
        {
          method: 'POST',
          body: JSON.stringify(data),
        }
      );
      if (response.error) {
        return { error: response.error };
      }
      return { success: true, data: response.data };
    } catch (error) {
      return { error: 'Ошибка создания события' };
    }
  },

  async deleteEvent(petId: number, eventId: number): Promise<{ data?: any; error?: string; success?: boolean }> {
    try {
      const response = await volunteerFetch(
        `/api/pets/${petId}/events/${eventId}`,
        {
          method: 'DELETE',
        }
      );
      if (response.error) {
        return { error: response.error };
      }
      return { success: true, data: response.data };
    } catch (error) {
      return { error: 'Ошибка удаления события' };
    }
  },
};
