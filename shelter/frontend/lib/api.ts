const API_URL = 'http://localhost:8200/api';

interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
}

export const shelterApi = {
  async checkHealth(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/health`);
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка подключения к серверу' };
    }
  },

  async getAnimals(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/animals`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки животных' };
    }
  },

  async getStats(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/stats`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки статистики' };
    }
  },

  async getOrganization(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/organization`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки информации о приюте' };
    }
  },

  async getVolunteers(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/volunteers`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки волонтеров' };
    }
  },

  async getAdoptionRequests(): Promise<ApiResponse> {
    try {
      const response = await fetch(`${API_URL}/adoption-requests`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки заявок' };
    }
  },
};
