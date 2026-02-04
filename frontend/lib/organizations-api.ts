// API клиент для работы с организациями

const API_URL = process.env.NEXT_PUBLIC_API_URL || 
  (typeof window !== 'undefined' && process.env.NODE_ENV === 'production'
    ? '' // Пустая строка = относительные пути (/api/...)
    : 'http://localhost:8000');

export interface Organization {
  id: number;
  name: string;
  short_name?: string;
  legal_form?: string;
  type: string; // shelter, vet_clinic, pet_shop, foundation, kennel, other
  
  // Юридическая информация
  inn?: string;
  ogrn?: string;
  kpp?: string;
  registration_date?: string;
  
  // Контактная информация
  email?: string;
  phone?: string;
  website?: string;
  
  // Адрес
  address_full?: string;
  address_postal_code?: string;
  address_region?: string;
  address_city?: string;
  address_street?: string;
  address_house?: string;
  address_office?: string;
  
  // Координаты
  geo_lat?: number;
  geo_lon?: number;
  
  // Описание
  description?: string;
  bio?: string;
  
  // Медиа
  logo?: string;
  cover_photo?: string;
  
  // Руководство
  director_name?: string;
  director_position?: string;
  
  // Владелец
  owner_user_id: number;
  
  // Настройки приватности
  profile_visibility: string;
  show_phone: string;
  show_email: string;
  allow_messages: string;
  
  // Статус
  is_verified: boolean;
  is_active: boolean;
  status: string;
  
  // Метаданные
  created_at: string;
  updated_at: string;
}

export interface OrganizationMember {
  id: number;
  organization_id: number;
  user_id: number;
  role: string; // owner, admin, moderator, member
  position?: string;
  can_post: boolean;
  can_edit: boolean;
  can_manage_members: boolean;
  joined_at: string;
  user_name?: string;
  user_avatar?: string;
}

export interface CreateOrganizationData {
  name: string;
  short_name?: string;
  legal_form?: string;
  type: string;
  inn?: string;
  ogrn?: string;
  kpp?: string;
  registration_date?: string;
  email?: string;
  phone?: string;
  website?: string;
  address_full?: string;
  address_postal_code?: string;
  address_region?: string;
  address_city?: string;
  address_street?: string;
  address_house?: string;
  address_office?: string;
  geo_lat?: number;
  geo_lon?: number;
  description?: string;
  bio?: string;
  director_name?: string;
  director_position?: string;
}

export interface UpdateOrganizationData {
  name?: string;
  short_name?: string;
  legal_form?: string;
  type?: string;
  email?: string;
  phone?: string;
  website?: string;
  address_full?: string;
  address_city?: string;
  geo_lat?: number;
  geo_lon?: number;
  description?: string;
  bio?: string;
  director_name?: string;
  director_position?: string;
  profile_visibility?: string;
  show_phone?: string;
  show_email?: string;
  allow_messages?: string;
}

export const organizationsApi = {
  // Создать организацию
  async create(data: CreateOrganizationData) {
    const response = await fetch(`${API_URL}/api/organizations`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Получить организацию по ID
  async getById(id: number) {
    const response = await fetch(`${API_URL}/api/organizations/${id}`, {
      credentials: 'include',
    });
    return response.json();
  },

  // Обновить организацию
  async update(id: number, data: UpdateOrganizationData) {
    const response = await fetch(`${API_URL}/api/organizations/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Удалить организацию
  async delete(id: number) {
    const response = await fetch(`${API_URL}/api/organizations/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    return response.json();
  },

  // Получить организации пользователя
  async getUserOrganizations(userId: number) {
    const response = await fetch(`${API_URL}/api/organizations/user/${userId}`, {
      credentials: 'include',
    });
    return response.json();
  },

  // Получить участников организации
  async getMembers(orgId: number) {
    const response = await fetch(`${API_URL}/api/organizations/members/${orgId}`, {
      credentials: 'include',
    });
    return response.json();
  },

  // Загрузить логотип
  async uploadLogo(orgId: number, file: File) {
    const formData = new FormData();
    formData.append('logo', file);

    const response = await fetch(`${API_URL}/api/organizations/${orgId}/logo`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    });
    return response.json();
  },

  // Загрузить обложку
  async uploadCover(orgId: number, file: File) {
    const formData = new FormData();
    formData.append('cover', file);

    const response = await fetch(`${API_URL}/api/organizations/${orgId}/cover`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    });
    return response.json();
  },

  // Подать заявку на вступление
  async requestJoin(orgId: number) {
    const response = await fetch(`${API_URL}/api/organizations/${orgId}/join-request`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });
    return response.json();
  },

  // Добавить участника
  async addMember(orgId: number, userId: number, role: string, position?: string) {
    const response = await fetch(`${API_URL}/api/organizations/${orgId}/members`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ user_id: userId, role, position }),
    });
    return response.json();
  },

  // Обновить участника
  async updateMember(memberId: number, role: string, position?: string) {
    const response = await fetch(`${API_URL}/api/organizations/members/${memberId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ role, position }),
    });
    return response.json();
  },

  // Удалить участника
  async removeMember(memberId: number) {
    const response = await fetch(`${API_URL}/api/organizations/members/${memberId}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    return response.json();
  },
};

// Типы организаций
export const ORGANIZATION_TYPES = {
  shelter: 'Приют для животных',
  vet_clinic: 'Ветеринарная клиника',
  pet_shop: 'Зоомагазин',
  foundation: 'Фонд помощи животным',
  kennel: 'Кинологический центр',
  other: 'Другое',
};

// Получить название типа организации
export function getOrganizationTypeName(type: string): string {
  return ORGANIZATION_TYPES[type as keyof typeof ORGANIZATION_TYPES] || 'Другое';
}
