'use client';

import { useEffect, useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import ChatList from './components/ChatList';
import ChatHeader from './components/ChatHeader';
import MessageList from './components/MessageList';
import MessageInput from './components/MessageInput';
import { Chat, Message } from './types';

export default function MessengerPage() {
  const router = useRouter();
  const { user, isLoading: authLoading } = useAuth();
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [messageText, setMessageText] = useState('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<number | null>(null);
  const [sending, setSending] = useState(false);
  const [isFetchingChats, setIsFetchingChats] = useState(false);
  const chatsLoaded = useRef(false);

  // Проверка авторизации
  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  // Блокируем скролл страницы
  useEffect(() => {
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, []);

  // Загружаем чаты только один раз при монтировании
  useEffect(() => {
    if (user && !chatsLoaded.current) {
      chatsLoaded.current = true;
      fetchChats();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.id]);

  // Обработка query параметра ?user=ID для открытия чата с конкретным пользователем
  useEffect(() => {
    // Безопасное получение query параметра
    const urlParams = new URLSearchParams(window.location.search);
    const userIdParam = urlParams.get('user');
    
    if (userIdParam && chats.length > 0 && !selectedChatId) {
      const targetUserId = parseInt(userIdParam);
      
      // Ищем существующий чат с этим пользователем
      const existingChat = chats.find(
        chat => chat.other_user?.id === targetUserId
      );
      
      if (existingChat) {
        // Открываем существующий чат
        setSelectedChatId(existingChat.id);
      } else {
        // Создаем временный чат (с отрицательным ID)
        const tempChat: Chat = {
          id: -targetUserId, // Временный ID
          other_user: {
            id: targetUserId,
            name: 'Загрузка...',
            last_name: '',
          },
          unread_count: 0,
        };
        
        setChats(prev => [tempChat, ...prev]);
        setSelectedChatId(tempChat.id);
        
        // Загружаем данные пользователя
        fetchUserData(targetUserId);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chats, selectedChatId, user?.id]);

  const fetchUserData = async (userId: number) => {
    try {
      const token = localStorage.getItem('auth_token');
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };
      
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch(`/api/users/${userId}`, {
        credentials: 'include',
        headers,
      });
      
      if (response.ok) {
        const data = await response.json();
        
        // Обновляем временный чат с реальными данными пользователя
        setChats(prev => prev.map(chat => {
          if (chat.id === -userId) {
            return {
              ...chat,
              other_user: data.success ? data.data : chat.other_user,
            };
          }
          return chat;
        }));
      }
    } catch (error) {
      console.error('Failed to fetch user data:', error);
    }
  };

  // Загружаем сообщения при выборе чата
  useEffect(() => {
    if (selectedChatId) {
      fetchMessages(selectedChatId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedChatId]);

  const fetchMessages = async (chatId: number) => {
    // Если это временный чат (ID < 0), не загружаем сообщения
    if (chatId < 0) {
      setMessages([]);
      return;
    }

    try {
      const token = localStorage.getItem('auth_token');
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };
      
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch(`/api/chats/${chatId}/messages`, {
        credentials: 'include',
        headers,
      });

      if (response.ok) {
        const data = await response.json();
        setMessages(data);
      } else {
        console.error('Failed to fetch messages:', response.status);
        setMessages([]);
      }
    } catch (error) {
      console.error('Error fetching messages:', error);
      setMessages([]);
    }
  };

  const fetchChats = async () => {
    if (isFetchingChats) return;
    
    setIsFetchingChats(true);
    try {
      const token = localStorage.getItem('auth_token');
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };
      
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch('/api/chats', {
        credentials: 'include',
        headers,
      });

      if (response.ok) {
        const data = await response.json();
        setChats(data);
      } else {
        console.error('Failed to fetch chats:', response.status);
      }
    } catch (error) {
      console.error('Ошибка загрузки чатов:', error);
    } finally {
      setLoading(false);
      setIsFetchingChats(false);
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!messageText.trim() || !selectedChatId || sending) return;

    const selectedChat = chats.find(c => c.id === selectedChatId);
    
    if (!selectedChat?.other_user?.id) {
      console.error('Не найден получатель сообщения');
      return;
    }

    setSending(true);
    
    try {
      const token = localStorage.getItem('auth_token');
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };
      
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch(`/api/messages/send`, {
        method: 'POST',
        headers,
        credentials: 'include',
        body: JSON.stringify({
          receiver_id: selectedChat.other_user.id,
          content: messageText.trim(),
        }),
      });

      if (response.ok) {
        setMessageText('');
        
        // Если это был временный чат (ID < 0), нужно обновить список чатов
        if (selectedChatId < 0) {
          const updatedChatsResponse = await fetch('/api/chats', {
            credentials: 'include',
            headers,
          });
          
          if (updatedChatsResponse.ok) {
            const updatedChats = await updatedChatsResponse.json();
            setChats(updatedChats);
            
            const realChat = updatedChats.find((chat: Chat) => chat.other_user?.id === selectedChat.other_user?.id);
            
            if (realChat) {
              setSelectedChatId(realChat.id);
              fetchMessages(realChat.id);
            }
          }
        } else {
          // Обычный чат - просто перезагружаем сообщения
          fetchMessages(selectedChatId);
        }
      }
    } catch (error) {
      console.error('Ошибка отправки сообщения:', error);
    } finally {
      setSending(false);
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0 || !selectedChatId) return;

    const selectedChat = chats.find(c => c.id === selectedChatId);
    if (!selectedChat?.other_user?.id) {
      console.error('Не найден получатель сообщения');
      return;
    }

    setSending(true);

    try {
      const token = localStorage.getItem('auth_token');
      const formData = new FormData();
      formData.append('receiver_id', selectedChat.other_user.id.toString());
      
      if (messageText.trim()) {
        formData.append('content', messageText.trim());
      }

      Array.from(files).forEach(file => {
        formData.append('media', file);
      });

      const headers: HeadersInit = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch('/api/messages/send-media', {
        method: 'POST',
        credentials: 'include',
        headers,
        body: formData,
      });

      if (response.ok) {
        setMessageText('');
        
        if (selectedChatId < 0) {
          await fetchChats();
          const updatedChatsResponse = await fetch('/api/chats', {
            credentials: 'include',
            headers,
          });
          
          if (updatedChatsResponse.ok) {
            const updatedChats = await updatedChatsResponse.json();
            const realChat = updatedChats.find((chat: Chat) => chat.other_user?.id === selectedChat.other_user?.id);
            
            if (realChat) {
              setSelectedChatId(realChat.id);
              fetchMessages(realChat.id);
            }
          }
        } else {
          fetchMessages(selectedChatId);
          fetchChats();
        }
      } else {
        const errorText = await response.text();
        console.error('Ошибка отправки медиа:', errorText);
        alert('Ошибка отправки файла');
      }
    } catch (error) {
      console.error('Ошибка отправки медиа:', error);
      alert('Ошибка отправки файла');
    } finally {
      setSending(false);
    }
  };

  const handleCloseChat = () => {
    setSelectedChatId(null);
    setMessages([]);
  };

  const selectedChat = chats.find(c => c.id === selectedChatId);

  return (
    <div className="h-[calc(100vh-74px)] bg-white rounded-lg shadow-sm border border-gray-200 flex overflow-hidden">
      {/* Левая панель - список чатов */}
      <ChatList
        chats={chats}
        loading={loading}
        isCollapsed={isCollapsed}
        selectedChatId={selectedChatId}
        currentUserId={user?.id}
        onToggleCollapse={() => setIsCollapsed(!isCollapsed)}
        onSelectChat={setSelectedChatId}
      />

      {/* Правая часть - окно чата */}
      <div className="flex-1 flex flex-col bg-gray-50">
        {selectedChatId ? (
          <>
            {/* Шапка чата */}
            <ChatHeader 
              user={selectedChat?.other_user || null}
              onClose={handleCloseChat}
            />

            {/* Область сообщений */}
            <MessageList 
              messages={messages}
              currentUserId={user?.id}
            />

            {/* Поле ввода */}
            <MessageInput
              messageText={messageText}
              sending={sending}
              onMessageChange={setMessageText}
              onSendMessage={handleSendMessage}
              onFileSelect={handleFileSelect}
            />
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center" style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23e3f2fd' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
            backgroundColor: '#e3f2fd'
          }}>
            <div className="text-center text-gray-400">
              <svg className="w-24 h-24 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <p className="text-lg font-medium">Выберите чат</p>
              <p className="text-sm mt-1">Чтобы начать общение</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
