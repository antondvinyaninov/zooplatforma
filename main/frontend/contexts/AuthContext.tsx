'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiClient, authApi, User } from '../lib/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (name: string, email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  refreshUser: () => Promise<void>;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Проверяем авторизацию при загрузке (только на клиенте)
    if (typeof window === 'undefined') {
      setIsLoading(false);
      return;
    }

    const checkAuth = async () => {
      const response = await authApi.me();
      
      if (response.success && response.data) {
        // Backend возвращает { user: User, token: string }
        const userData = response.data as any;
        const userId = userData.user?.id || userData.id;
        
        if (userId) {
          // Получаем полные данные профиля из Main Backend
          const profileResponse = await apiClient.get<User>(`/api/users/${userId}`);
          
          if (profileResponse.success && profileResponse.data) {
            setUser(profileResponse.data);
            setToken('authenticated');
            setIsLoading(false);
            return;
          }
        }
        
        // Fallback: используем данные из Auth Service
        if (userData.user) {
          setUser(userData.user);
          setToken('authenticated');
        } else {
          setUser(response.data as User);
          setToken('authenticated');
        }
      }
      
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const login = async (email: string, password: string) => {
    const response = await authApi.login(email, password);
    
    if (response.success && response.data) {
      const { user } = response.data;
      setToken('authenticated'); // Токен в cookie, просто флаг
      setUser(user);
      return { success: true };
    }

    return { success: false, error: response.error };
  };

  const register = async (name: string, email: string, password: string) => {
    const response = await authApi.register(name, email, password);
    
    if (response.success && response.data) {
      const { user } = response.data;
      setToken('authenticated'); // Токен в cookie, просто флаг
      setUser(user);
      return { success: true };
    }

    return { success: false, error: response.error };
  };

  const logout = async () => {
    await authApi.logout();
    setToken(null);
    setUser(null);
  };

  const refreshUser = async () => {
    // Сначала получаем базовые данные из Auth Service (для проверки авторизации)
    const authResponse = await authApi.me();
    
    if (authResponse.success && authResponse.data) {
      const userData = authResponse.data as any;
      const userId = userData.user?.id || userData.id;
      
      if (userId) {
        // Затем получаем полные данные профиля из Main Backend
        const profileResponse = await apiClient.get<User>(`/api/users/${userId}`);
        
        if (profileResponse.success && profileResponse.data) {
          setUser(profileResponse.data);
          return;
        }
      }
      
      // Fallback: используем данные из Auth Service если не удалось получить из Main Backend
      if (userData.user) {
        setUser(userData.user);
      } else {
        setUser(authResponse.data as User);
      }
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        register,
        logout,
        refreshUser,
        isAuthenticated: !!token,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
