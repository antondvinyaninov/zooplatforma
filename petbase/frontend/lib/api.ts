const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8100';
const MAIN_API_URL = 'http://localhost:8000';
const ADMIN_API_URL = 'http://localhost:9000';

export const petbaseApi = {
  // Auth через SSO
  async checkAuth() {
    try {
      // Проверяем авторизацию через главный backend
      const response = await fetch(`${MAIN_API_URL}/api/auth/me`, {
        credentials: 'include',
      });
      const data = await response.json();
      
      if (!data.success) {
        return { success: false, error: 'Не авторизован' };
      }

      // Проверяем права суперадмина
      const adminResponse = await fetch(`${ADMIN_API_URL}/api/admin/auth/me`, {
        credentials: 'include',
      });
      const adminData = await adminResponse.json();

      if (!adminData.success || adminData.data?.role !== 'superadmin') {
        return { success: false, error: 'Доступ только для суперадминистраторов' };
      }

      return { success: true, data: adminData.data };
    } catch (error) {
      return { success: false, error: 'Ошибка подключения' };
    }
  },

  // Species API
  async getSpecies() {
    try {
      const response = await fetch(`${API_URL}/api/species`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки видов' };
    }
  },

  async createSpecies(data: { name: string; name_en: string; description: string; icon: string }) {
    try {
      const response = await fetch(`${API_URL}/api/species`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data),
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка создания вида' };
    }
  },

  async deleteSpecies(id: number) {
    try {
      const response = await fetch(`${API_URL}/api/species/${id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка удаления вида' };
    }
  },

  // Breeds API
  async getBreeds() {
    try {
      const response = await fetch(`${API_URL}/api/breeds`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки пород' };
    }
  },

  async getBreedsBySpecies(speciesId: number) {
    try {
      const response = await fetch(`${API_URL}/api/breeds/species/${speciesId}`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки пород' };
    }
  },

  async createBreed(data: any) {
    try {
      const response = await fetch(`${API_URL}/api/breeds`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data),
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка создания породы' };
    }
  },

  async deleteBreed(id: number) {
    try {
      const response = await fetch(`${API_URL}/api/breeds/${id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка удаления породы' };
    }
  },

  // Cards API
  async getCards() {
    try {
      const response = await fetch(`${API_URL}/api/cards`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки карточек' };
    }
  },

  async getCardsByBreed(breedId: number) {
    try {
      const response = await fetch(`${API_URL}/api/cards/breed/${breedId}`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки карточек' };
    }
  },

  async createCard(data: any) {
    try {
      const response = await fetch(`${API_URL}/api/cards`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data),
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка создания карточки' };
    }
  },

  async deleteCard(id: number) {
    try {
      const response = await fetch(`${API_URL}/api/cards/${id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка удаления карточки' };
    }
  },
};
