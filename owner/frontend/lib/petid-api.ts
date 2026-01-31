/**
 * PetID API Client
 * 
 * Клиент для взаимодействия с питомцами через Owner backend
 * Owner backend проксирует запросы к PetID service
 */

import { Pet, MedicalEvent, CreatePetData, ApiResponse } from '../types/owner';

const OWNER_API_URL = process.env.NEXT_PUBLIC_OWNER_API_URL || 'http://localhost:8400';
const AUTH_SERVICE_URL = 'http://localhost:7100';

/**
 * Получить JWT токен через Auth Service
 */
async function getAuthToken(): Promise<string | null> {
  try {
    console.log('🔍 Getting auth token from Auth Service...');
    const response = await fetch(`${AUTH_SERVICE_URL}/api/auth/me`, {
      credentials: 'include', // Отправляет cookie с токеном
    });
    
    if (!response.ok) {
      console.error('❌ Failed to get auth token:', response.status);
      return null;
    }
    
    const data = await response.json();
    const token = data.data?.token || null;
    
    if (token) {
      console.log('✅ Auth token received:', token.substring(0, 20) + '...');
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
 * Базовая функция для выполнения запросов к Owner API
 */
async function ownerFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  // Получаем токен через Main API
  const token = await getAuthToken();
  
  if (!token) {
    console.error('❌ No auth token available for request:', endpoint);
    return { error: 'Unauthorized' };
  }
  
  console.log('🚀 Making request to Owner API:', endpoint);
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`, // Передаем токен в заголовке
    ...(options.headers as Record<string, string>),
  };
  
  try {
    const response = await fetch(`${OWNER_API_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include', // Важно для передачи cookies с JWT
    });
    
    console.log('📥 Response from Owner API:', response.status, endpoint);
    
    // Обработка ошибок аутентификации
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
    console.error('Owner API Error:', error);
    return { 
      error: error instanceof Error ? error.message : 'Network error' 
    };
  }
}

/**
 * Интерфейс для породы
 */
export interface Breed {
  id: number;
  name: string;
  name_en: string;
}

/**
 * PetID API Client
 */
export const petidApi = {
  /**
   * Получить список пород для указанного вида животного
   */
  async getBreeds(species: 'dog' | 'cat'): Promise<ApiResponse<Breed[]>> {
    const response = await ownerFetch<{ breeds: Breed[] }>(`/api/breeds?species=${species}`);
    
    if (response.error) {
      return { error: response.error };
    }
    
    return { data: response.data?.breeds || [] };
  },

  /**
   * Получить список питомцев текущего пользователя
   */
  async getPets(): Promise<ApiResponse<Pet[]>> {
    const response = await ownerFetch<{ pets: Pet[] }>('/api/pets');
    
    if (response.error) {
      return { error: response.error };
    }
    
    return { data: response.data?.pets || [] };
  },
  
  /**
   * Получить информацию о конкретном питомце
   */
  async getPet(id: number): Promise<ApiResponse<Pet>> {
    const response = await ownerFetch<{ pet: Pet }>(`/api/pets/${id}`);
    
    if (response.error) {
      return { error: response.error };
    }
    
    if (!response.data?.pet) {
      return { error: 'Pet not found' };
    }
    
    return { data: response.data.pet };
  },
  
  /**
   * Создать нового питомца
   */
  async createPet(data: CreatePetData): Promise<ApiResponse<Pet>> {
    const response = await ownerFetch<{ pet: Pet }>('/api/pets', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    
    if (response.error) {
      return { error: response.error };
    }
    
    if (!response.data?.pet) {
      return { error: 'Failed to create pet' };
    }
    
    return { 
      data: response.data.pet,
      message: 'PetID создан, ожидает подтверждения от ветклиники'
    };
  },
  
  /**
   * Обновить информацию о питомце
   */
  async updatePet(id: number, data: Partial<Pet>): Promise<ApiResponse<Pet>> {
    const response = await ownerFetch<{ pet: Pet }>(`/api/pets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    
    if (response.error) {
      return { error: response.error };
    }
    
    if (!response.data?.pet) {
      return { error: 'Failed to update pet' };
    }
    
    return { data: response.data.pet };
  },
  
  /**
   * Загрузить фото питомца
   */
  async uploadPhoto(id: number, file: File): Promise<ApiResponse<string>> {
    // Получаем токен через Main API
    const token = await getAuthToken();
    
    if (!token) {
      return { error: 'Unauthorized' };
    }
    
    const formData = new FormData();
    formData.append('photo', file);
    
    try {
      const response = await fetch(`${OWNER_API_URL}/api/pets/${id}/photo`, {
        method: 'POST',
        body: formData,
        headers: {
          'Authorization': `Bearer ${token}`, // Передаем токен
        },
        credentials: 'include',
      });
      
      if (response.status === 401) {
        return { error: 'Unauthorized' };
      }
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return { 
          error: errorData.message || errorData.error || 'Upload failed' 
        };
      }
      
      const data = await response.json();
      return { data: data.photo_url || data.photoUrl };
      
    } catch (error) {
      console.error('Photo upload error:', error);
      return { 
        error: error instanceof Error ? error.message : 'Upload failed' 
      };
    }
  },
  
  /**
   * Получить медицинскую историю питомца (события от клиник)
   */
  async getMedicalHistory(petId: number): Promise<ApiResponse<MedicalEvent[]>> {
    const response = await ownerFetch<{ events: MedicalEvent[] }>(
      `/api/petid/${petId}/medical`
    );
    
    if (response.error) {
      return { error: response.error };
    }
    
    return { data: response.data?.events || [] };
  },
  
  /**
   * Создать событие в истории питомца
   */
  async createEvent(
    petId: number,
    event: {
      event_type: string;
      event_date: string;
      description: string;
      veterinarian?: string;
      next_visit_date?: string;
    }
  ): Promise<ApiResponse<MedicalEvent>> {
    const response = await ownerFetch<{ event: MedicalEvent }>(
      `/api/petid/${petId}/events`,
      {
        method: 'POST',
        body: JSON.stringify(event),
      }
    );
    
    if (response.error) {
      return { error: response.error };
    }
    
    if (!response.data?.event) {
      return { error: 'Failed to create event' };
    }
    
    return { data: response.data.event };
  },
  
  /**
   * Получить все события питомца
   */
  async getEvents(petId: number): Promise<ApiResponse<MedicalEvent[]>> {
    const response = await ownerFetch<{ events: MedicalEvent[] }>(
      `/api/petid/${petId}/events`
    );
    
    if (response.error) {
      return { error: response.error };
    }
    
    return { data: response.data?.events || [] };
  },
};

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
  async getEvents(petId: number): Promise<ApiResponse<{ events: MedicalEvent[] }>> {
    try {
      const response = await ownerFetch<{ events: MedicalEvent[] }>(
        `/api/pets/${petId}/events`
      );
      return response;
    } catch (error) {
      return { error: 'Ошибка загрузки событий' };
    }
  },

  async createEvent(petId: number, data: any): Promise<ApiResponse<any>> {
    try {
      const response = await ownerFetch(
        `/api/pets/${petId}/events`,
        {
          method: 'POST',
          body: JSON.stringify(data),
        }
      );
      return response;
    } catch (error) {
      return { error: 'Ошибка создания события' };
    }
  },

  async deleteEvent(petId: number, eventId: number): Promise<ApiResponse<any>> {
    try {
      const response = await ownerFetch(
        `/api/pets/${petId}/events/${eventId}`,
        {
          method: 'DELETE',
        }
      );
      return response;
    } catch (error) {
      return { error: 'Ошибка удаления события' };
    }
  },
};

export default petidApi;

