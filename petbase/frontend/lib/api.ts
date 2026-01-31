const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8100';
const AUTH_SERVICE_URL = 'http://localhost:7100';

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
  async getEvents(petId: number) {
    try {
      const response = await fetch(`${API_URL}/api/pets/${petId}/events`, {
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка загрузки событий' };
    }
  },

  async createEvent(petId: number, data: any) {
    try {
      const response = await fetch(`${API_URL}/api/pets/${petId}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data),
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка создания события' };
    }
  },

  async deleteEvent(petId: number, eventId: number) {
    try {
      const response = await fetch(`${API_URL}/api/pets/${petId}/events/${eventId}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      return await response.json();
    } catch (error) {
      return { success: false, error: 'Ошибка удаления события' };
    }
  },
};


export const petbaseApi = {
  // Auth через Auth Service (7100)
  async checkAuth() {
    try {
      // Проверяем авторизацию через Auth Service
      const response = await fetch(`${AUTH_SERVICE_URL}/api/auth/me`, {
        credentials: 'include',
      });
      const data = await response.json();
      
      if (!data.success || !data.data?.user) {
        return { success: false, error: 'Не авторизован' };
      }

      // Проверяем права суперадмина
      const user = data.data.user;
      if (user.role !== 'superadmin') {
        return { success: false, error: 'Доступ только для суперадминистраторов' };
      }

      return { success: true, data: { email: user.email, role: user.role, avatar: user.avatar } };
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
