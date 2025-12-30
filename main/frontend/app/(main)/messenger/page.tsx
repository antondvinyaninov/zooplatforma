'use client';

import { useState, useEffect, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

interface User {
  id: number;
  name: string;
  last_name: string;
  avatar?: string;
  email?: string;
}

interface Message {
  id: number;
  chat_id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  is_read: boolean;
  read_at?: string;
  created_at: string;
  sender?: User;
}

interface Chat {
  id: number;
  user1_id: number;
  user2_id: number;
  last_message_id?: number;
  last_message_at?: string;
  created_at: string;
  other_user?: User;
  last_message?: Message;
  unread_count: number;
}

export default function MessengerPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const searchParams = new URLSearchParams(typeof window !== 'undefined' ? window.location.search : '');
  const userIdParam = searchParams.get('user');
  
  const [chats, setChats] = useState<Chat[]>([]);
  const [selectedChat, setSelectedChat] = useState<Chat | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Проверяем только на клиенте
    if (typeof window === 'undefined') return;
    
    if (!isLoading && !isAuthenticated) {
      router.push('/auth');
      return;
    }
    
    if (user) {
      fetchChats();
    }
  }, [user, isLoading, isAuthenticated, router]);

  useEffect(() => {
    // Если есть параметр user в URL, пытаемся открыть диалог с этим пользователем
    if (userIdParam && chats.length > 0 && !selectedChat && user) {
      const targetUserId = parseInt(userIdParam);
      const existingChat = chats.find(chat => 
        chat.other_user?.id === targetUserId
      );
      
      if (existingChat) {
        setSelectedChat(existingChat);
      } else {
        // Создаем временный "чат" для нового диалога
        const newChat: Chat = {
          id: 0, // Временный ID
          user1_id: user.id,
          user2_id: targetUserId,
          created_at: new Date().toISOString(),
          unread_count: 0,
          other_user: {
            id: targetUserId,
            name: 'Загрузка...',
            last_name: '',
            email: '',
          }
        };
        setSelectedChat(newChat);
        
        // Загружаем информацию о пользователе
        fetch(`http://localhost:8000/api/users/${targetUserId}`, {
          credentials: 'include',
        })
          .then(res => res.json())
          .then(data => {
            // API возвращает { success: true, data: {...} }
            const userData = data.data || data;
            setSelectedChat(prev => prev ? {
              ...prev,
              other_user: userData
            } : null);
          })
          .catch(err => console.error('Error loading user:', err));
      }
    }
  }, [userIdParam, chats, selectedChat, user]);

  useEffect(() => {
    // Загружаем сообщения только для реальных чатов (не временных)
    if (selectedChat && selectedChat.id > 0) {
      fetchMessages(selectedChat.id);
    } else if (selectedChat && selectedChat.id === 0) {
      // Для нового диалога просто очищаем сообщения
      setMessages([]);
    }
  }, [selectedChat]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const fetchChats = async () => {
    try {
      const response = await fetch('http://localhost:8000/api/chats', {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setChats(data);
      } else if (response.status === 401) {
        // Пользователь не авторизован - редирект на /auth
        router.push('/auth');
      }
    } catch (error) {
      // Тихо игнорируем ошибки сети при загрузке
      console.debug('Could not fetch chats:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMessages = async (chatId: number) => {
    try {
      const response = await fetch(`http://localhost:8000/api/chats/${chatId}/messages`, {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setMessages(data);
        
        // Обновляем счетчик непрочитанных в списке чатов
        setChats(prevChats =>
          prevChats.map(chat =>
            chat.id === chatId ? { ...chat, unread_count: 0 } : chat
          )
        );
      }
    } catch (error) {
      console.error('Error fetching messages:', error);
    }
  };

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newMessage.trim() || !selectedChat || sending) return;

    const isNewChat = selectedChat.id === 0;
    const receiverId = selectedChat.other_user?.id;

    if (!receiverId) {
      console.error('Receiver ID is missing');
      return;
    }

    setSending(true);
    try {
      const response = await fetch('http://localhost:8000/api/messages/send', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          receiver_id: receiverId,
          content: newMessage.trim(),
        }),
      });

      if (response.ok) {
        const message = await response.json();
        setMessages(prev => [...prev, message]);
        setNewMessage('');
        
        // Если это был новый чат, обновляем список чатов и выбираем созданный
        if (isNewChat) {
          const chatsResponse = await fetch('http://localhost:8000/api/chats', {
            credentials: 'include',
          });
          
          if (chatsResponse.ok) {
            const updatedChats = await chatsResponse.json();
            setChats(updatedChats);
            
            // Находим и выбираем только что созданный чат
            const newChat = updatedChats.find((chat: Chat) => 
              chat.other_user?.id === receiverId
            );
            
            if (newChat) {
              setSelectedChat(newChat);
            }
          }
        } else {
          // Просто обновляем список чатов для обновления last_message
          fetchChats();
        }
      } else {
        const errorData = await response.json();
        console.error('Error sending message:', errorData);
      }
    } catch (error) {
      console.error('Error sending message:', error);
    } finally {
      setSending(false);
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) {
      return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    } else if (days === 1) {
      return 'Вчера';
    } else if (days < 7) {
      return date.toLocaleDateString('ru-RU', { weekday: 'short' });
    } else {
      return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-500">Загрузка...</div>
      </div>
    );
  }

  return (
    <div className="flex h-[calc(100vh-64px)] bg-gray-50">
      {/* Список диалогов */}
      <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
        <div className="p-4 border-b border-gray-200">
          <h1 className="text-xl font-semibold">Мессенджер</h1>
        </div>

        <div className="flex-1 overflow-y-auto">
          {chats.length === 0 ? (
            <div className="p-4 text-center text-gray-500">
              Нет диалогов
            </div>
          ) : (
            chats.map(chat => (
              <button
                key={chat.id}
                onClick={() => setSelectedChat(chat)}
                className={`w-full p-4 flex items-start gap-3 hover:bg-gray-50 transition-colors border-b border-gray-100 ${
                  selectedChat?.id === chat.id ? 'bg-blue-50' : ''
                }`}
              >
                {/* Аватар */}
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-semibold flex-shrink-0">
                  {chat.other_user?.name?.[0]?.toUpperCase() || '?'}
                </div>

                {/* Информация */}
                <div className="flex-1 min-w-0 text-left">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-medium text-gray-900 truncate">
                      {chat.other_user?.name} {chat.other_user?.last_name}
                    </span>
                    {chat.last_message_at && (
                      <span className="text-xs text-gray-500 ml-2 flex-shrink-0">
                        {formatTime(chat.last_message_at)}
                      </span>
                    )}
                  </div>
                  
                  {chat.last_message && (
                    <div className="flex items-center justify-between">
                      <p className="text-sm text-gray-600 truncate">
                        {chat.last_message.sender_id === user?.id && 'Вы: '}
                        {chat.last_message.content}
                      </p>
                      {chat.unread_count > 0 && (
                        <span className="ml-2 bg-blue-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center flex-shrink-0">
                          {chat.unread_count}
                        </span>
                      )}
                    </div>
                  )}
                </div>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Окно чата */}
      <div className="flex-1 flex flex-col bg-white">
        {selectedChat ? (
          <>
            {/* Шапка чата */}
            <div className="p-4 border-b border-gray-200 flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-semibold">
                {selectedChat.other_user?.name?.[0]?.toUpperCase() || '?'}
              </div>
              <div>
                <h2 className="font-semibold text-gray-900">
                  {selectedChat.other_user?.name} {selectedChat.other_user?.last_name}
                </h2>
                {selectedChat.id === 0 && (
                  <p className="text-xs text-gray-500">Новый диалог</p>
                )}
              </div>
            </div>

            {/* Сообщения */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messages.length === 0 && selectedChat.id === 0 ? (
                <div className="flex items-center justify-center h-full text-gray-500">
                  <div className="text-center">
                    <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-2xl font-semibold">
                      {selectedChat.other_user?.name?.[0]?.toUpperCase() || '?'}
                    </div>
                    <p className="text-lg font-medium text-gray-900 mb-1">
                      {selectedChat.other_user?.name} {selectedChat.other_user?.last_name}
                    </p>
                    <p className="text-sm text-gray-500">
                      Начните диалог, отправив первое сообщение
                    </p>
                  </div>
                </div>
              ) : messages.length === 0 ? (
                <div className="flex items-center justify-center h-full text-gray-500">
                  Нет сообщений
                </div>
              ) : (
                messages.map(message => {
                  const isMyMessage = message.sender_id === user?.id;
                  
                  return (
                    <div
                      key={message.id}
                      className={`flex ${isMyMessage ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                          isMyMessage
                            ? 'bg-blue-500 text-white'
                            : 'bg-gray-100 text-gray-900'
                        }`}
                      >
                        <p className="whitespace-pre-wrap break-words">{message.content}</p>
                        <div
                          className={`text-xs mt-1 ${
                            isMyMessage ? 'text-blue-100' : 'text-gray-500'
                          }`}
                        >
                          {formatTime(message.created_at)}
                        </div>
                      </div>
                    </div>
                  );
                })
              )}
              <div ref={messagesEndRef} />
            </div>

            {/* Форма отправки */}
            <form onSubmit={sendMessage} className="p-4 border-t border-gray-200">
              <div className="flex gap-2">
                <input
                  type="text"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  placeholder="Написать сообщение..."
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={sending}
                />
                <button
                  type="submit"
                  disabled={!newMessage.trim() || sending}
                  className="px-6 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {sending ? 'Отправка...' : 'Отправить'}
                </button>
              </div>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Выберите диалог
          </div>
        )}
      </div>
    </div>
  );
}
