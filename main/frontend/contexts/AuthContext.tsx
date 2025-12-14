'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiClient, authApi, User } from '../lib/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (name: string, email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Проверяем авторизацию при загрузке
    const checkAuth = async () => {
      const response = await authApi.me();
      
      if (response.success && response.data) {
        setUser(response.data);
        setToken('authenticated');
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

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        register,
        logout,
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
