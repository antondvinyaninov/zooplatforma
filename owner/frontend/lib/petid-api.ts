/**
 * PetID API Client
 * 
 * Клиент для взаимодействия с питомцами через Owner backend
 * Owner backend проксирует запросы к PetID service
 */

import { Pet, MedicalEvent, CreatePetData, ApiResponse } from '../types/owner';

const OWNER_API_URL = process.env.NEXT_PUBLIC_OWNER_API_URL || 'http://localhost:8400';

/**
 * Базовая функция для выполнения запросов к Owner API
 */
async function ownerFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };
  
  try {
    const response = await fetch(`${OWNER_API_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include', // Важно для передачи cookies с JWT
    });
    
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
    const formData = new FormData();
    formData.append('photo', file);
    
    try {
      const response = await fetch(`${OWNER_API_URL}/api/pets/${id}/photo`, {
        method: 'POST',
        body: formData,
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

export default petidApi;
