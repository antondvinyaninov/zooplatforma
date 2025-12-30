const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getHeaders(): HeadersInit {
    return {
      'Content-Type': 'application/json',
    };
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    try {
      const result = await response.json();

      if (!response.ok) {
        return {
          success: false,
          error: result.error || `API Error: ${response.statusText}`,
        };
      }

      // Бэкенд возвращает {success: true, data: ...}, извлекаем data
      return {
        success: true,
        data: result.data !== undefined ? result.data : result,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method: 'GET',
        headers: this.getHeaders(),
        credentials: 'include', // Include cookies
      });

      return this.handleResponse<T>(response);
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  async post<T>(endpoint: string, body: unknown): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method: 'POST',
        headers: this.getHeaders(),
        credentials: 'include', // Include cookies
        body: JSON.stringify(body),
      });

      return this.handleResponse<T>(response);
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  async put<T>(endpoint: string, body: unknown): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method: 'PUT',
        headers: this.getHeaders(),
        credentials: 'include', // Include cookies
        body: JSON.stringify(body),
      });

      return this.handleResponse<T>(response);
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method: 'DELETE',
        headers: this.getHeaders(),
        credentials: 'include', // Include cookies
      });

      return this.handleResponse<T>(response);
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }
}

export const apiClient = new ApiClient(API_URL);

// API методы
export const authApi = {
  register: (name: string, email: string, password: string) =>
    apiClient.post<{ user: User }>('/api/auth/register', { name, email, password }),
  
  login: (email: string, password: string) =>
    apiClient.post<{ user: User }>('/api/auth/login', { email, password }),
  
  logout: () =>
    apiClient.post<{ message: string }>('/api/auth/logout', {}),
  
  me: () =>
    apiClient.get<User>('/api/auth/me'),
};

export const usersApi = {
  getAll: () => apiClient.get<User[]>('/api/users'),
  
  getById: (id: number) => apiClient.get<User>(`/api/users/${id}`),
  
  create: (user: Partial<User>) => apiClient.post<User>('/api/users', user),
  
  update: (id: number, user: Partial<User>) => apiClient.put<User>(`/api/users/${id}`, user),
  
  delete: (id: number) => apiClient.delete<{ message: string }>(`/api/users/${id}`),
  
  // Обновление своего профиля
  updateProfile: (data: { 
    name?: string; 
    last_name?: string;
    bio?: string; 
    phone?: string; 
    location?: string;
    profile_visibility?: string;
    show_phone?: string;
    show_email?: string;
    allow_messages?: string;
    show_online?: string;
  }) =>
    apiClient.put<User>('/api/profile', data),
  
  // Загрузка аватара
  uploadAvatar: async (file: File): Promise<ApiResponse<{ avatar_url: string; message: string }>> => {
    const formData = new FormData();
    formData.append('avatar', file);
    
    try {
      const response = await fetch(`${API_URL}/api/profile/avatar`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      
      const result = await response.json();
      
      if (!response.ok) {
        return {
          success: false,
          error: result.error || 'Ошибка загрузки аватара',
        };
      }
      
      return {
        success: true,
        data: result.data,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  },
  
  // Загрузка обложки
  uploadCover: async (file: File): Promise<ApiResponse<{ cover_url: string; message: string }>> => {
    const formData = new FormData();
    formData.append('cover', file);
    
    try {
      const response = await fetch(`${API_URL}/api/profile/cover`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      
      const result = await response.json();
      
      if (!response.ok) {
        return {
          success: false,
          error: result.error || 'Ошибка загрузки обложки',
        };
      }
      
      return {
        success: true,
        data: result.data,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  },
  
  // Удаление аватара
  deleteAvatar: async (): Promise<ApiResponse<{ message: string }>> => {
    try {
      const response = await fetch(`${API_URL}/api/profile/avatar/delete`, {
        method: 'DELETE',
        credentials: 'include',
      });
      
      const result = await response.json();
      
      if (!response.ok) {
        return {
          success: false,
          error: result.error || 'Ошибка удаления аватара',
        };
      }
      
      return {
        success: true,
        data: result.data,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  },
  
  // Удаление обложки
  deleteCover: async (): Promise<ApiResponse<{ message: string }>> => {
    try {
      const response = await fetch(`${API_URL}/api/profile/cover/delete`, {
        method: 'DELETE',
        credentials: 'include',
      });
      
      const result = await response.json();
      
      if (!response.ok) {
        return {
          success: false,
          error: result.error || 'Ошибка удаления обложки',
        };
      }
      
      return {
        success: true,
        data: result.data,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  },
};

// API методы для профиля
export const profileApi = {
  update: (data: {
    bio?: string;
    phone?: string;
    location?: string;
    avatar?: string;
    cover_photo?: string;
  }) => apiClient.put<User>('/api/profile', data),
};

// API методы для постов
export const postsApi = {
  getAll: () => apiClient.get<Post[]>('/api/posts'),
  
  getUserPosts: (userId: number) => apiClient.get<Post[]>(`/api/posts/user/${userId}`),
  
  getPetPosts: (petId: number) => apiClient.get<Post[]>(`/api/posts/pet/${petId}`),
  
  create: (data: { content: string; post_type?: string }) =>
    apiClient.post<Post>('/api/posts', data),
  
  update: (id: number, data: { content?: string; post_type?: string }) =>
    apiClient.put<Post>(`/api/posts/${id}`, data),
  
  delete: (id: number) => apiClient.delete<{ message: string }>(`/api/posts/${id}`),
};

// API методы для комментариев
export const commentsApi = {
  getPostComments: (postId: number) => apiClient.get<Comment[]>(`/api/comments/post/${postId}`),
  
  create: (postId: number, data: { content: string; parent_id?: number; reply_to_user_id?: number }) =>
    apiClient.post<Comment>(`/api/comments/post/${postId}`, data),
  
  delete: (commentId: number) => apiClient.delete<{ message: string }>(`/api/comments/${commentId}`),
};

// API методы для питомцев
export const petsApi = {
  getUserPets: (userId: number) => apiClient.get<Pet[]>(`/api/pets/user/${userId}`),
  
  create: (data: { name: string; species?: string; photo?: string }) =>
    apiClient.post<Pet>('/api/pets', data),
  
  delete: (id: number) => apiClient.delete<{ message: string }>(`/api/pets/${id}`),
};

// Типы
export interface User {
  id: number;
  name: string;
  last_name?: string;
  email: string;
  bio?: string;
  phone?: string;
  location?: string;
  avatar?: string;
  cover_photo?: string;
  profile_visibility?: string;
  show_phone?: string;
  show_email?: string;
  allow_messages?: string;
  show_online?: string;
  created_at?: string;
}

export interface Post {
  id: number;
  author_id: number;
  author_type: string;
  content: string;
  attached_pets: number[];
  attachments: Attachment[];
  tags: string[];
  status: string;
  scheduled_at?: string;
  created_at: string;
  updated_at: string;
  is_deleted?: boolean;
  user?: User;
  pets?: Pet[];
  poll?: Poll;
  comments_count?: number;
}

export interface Attachment {
  url: string;
  type: string;
  file_name?: string;
  size?: number;
}

export interface Poll {
  id: number;
  post_id: number;
  question: string;
  options: PollOption[];
  multiple_choice: boolean;
  allow_vote_changes: boolean;
  anonymous_voting: boolean;
  end_date?: string;
  expires_at?: string;
  created_at: string;
  total_votes: number;
  total_voters: number;
  user_voted: boolean;
  user_votes?: number[];
  is_expired: boolean;
  voters?: PollVoter[];
}

export interface PollOption {
  id: number;
  poll_id: number;
  option_text: string;
  votes_count: number;
  percentage: number;
  user_voted: boolean;
}

export interface PollVoter {
  user_id: number;
  user_name: string;
  user_avatar?: string;
  voted_at: string;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  created_at: string;
  parent_id?: number;
  reply_to_user_id?: number;
  user?: User;
  reply_to_user?: User;
  replies?: Comment[];
}

export interface Pet {
  id: number;
  user_id: number;
  name: string;
  species: string;
  breed?: string;
  gender?: string;
  birth_date?: string;
  color?: string;
  photo?: string;
  status?: string;
  is_sterilized?: boolean;
  is_vaccinated?: boolean;
  chip_number?: string;
  region?: string;
  city?: string;
  urgent?: boolean;
  contact_name?: string;
  contact_phone?: string;
  organization_id?: number;
  created_at: string;
}
