'use client';

import { useEffect, useState, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { getMediaUrl } from '@/lib/utils';

interface User {
  id: number;
  name: string;
  last_name: string;
  avatar?: string;
  is_online?: boolean; // Добавлено поле онлайн статуса
  last_seen?: string; // Добавлено поле последней активности
}

interface Message {
  id: number;
  sender_id: number;
  content: string;
  created_at: string;
  is_read?: boolean;
  attachments?: MessageAttachment[];
}

interface MessageAttachment {
  id: number;
  message_id: number;
  file_path: string;
  file_type: string;
  file_size: number;
  created_at: string;
  file_name?: string;
}

interface Chat {
  id: number;
  other_user?: User;
  last_message?: Message;
  last_message_at?: string;
  unread_count: number;
}

export default function MessengerPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, isLoading: authLoading } = useAuth();
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [messageText, setMessageText] = useState('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<number | null>(null);
  const [sending, setSending] = useState(false);
  const [showAttachMenu, setShowAttachMenu] = useState(false);
  const [showPetsModal, setShowPetsModal] = useState(false);
  const [selectedPets, setSelectedPets] = useState<number[]>([]);
  const [userPets, setUserPets] = useState<any[]>([]);
  const [curatorPets, setCuratorPets] = useState<any[]>([]);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Функция прокрутки к последнему сообщению
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

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

  // Загружаем чаты
  useEffect(() => {
    if (user) {
      fetchChats();
      fetchUserPets();
    }
  }, [user]);

  // Обрабатываем параметр user из URL (открытие диалога с конкретным пользователем)
  useEffect(() => {
    const userIdParam = searchParams.get('user');
    if (userIdParam && chats.length > 0 && !loading) {
      const targetUserId = parseInt(userIdParam);
      
      // Ищем существующий чат с этим пользователем
      const existingChat = chats.find(chat => chat.other_user?.id === targetUserId);
      
      if (existingChat) {
        // Если чат существует - выбираем его
        setSelectedChatId(existingChat.id);
      } else {
        // Если чата нет - загружаем информацию о пользователе и создаем временный чат
        loadUserAndCreateTempChat(targetUserId);
      }
      
      // Убираем параметр из URL
      router.replace('/messenger', { scroll: false });
    }
  }, [searchParams, chats, loading, router]);

  // Загружаем сообщения при выборе чата
  useEffect(() => {
    if (selectedChatId) {
      fetchMessages(selectedChatId);
    }
  }, [selectedChatId]);

  // Прокручиваем к последнему сообщению при изменении списка сообщений
  useEffect(() => {
    // Небольшая задержка чтобы дождаться загрузки изображений
    const timer = setTimeout(() => {
      scrollToBottom();
    }, 100);
    
    return () => clearTimeout(timer);
  }, [messages]);

  const fetchMessages = async (chatId: number) => {
    // Если это временный чат (ID < 0), не загружаем сообщения
    if (chatId < 0) {
      setMessages([]);
      return;
    }
    
    try {
      const response = await fetch(`http://localhost:8000/api/chats/${chatId}/messages`, {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setMessages(data);
      }
    } catch (error) {
      console.error('Ошибка загрузки сообщений:', error);
    }
  };

  const fetchChats = async () => {
    try {
      const response = await fetch('http://localhost:8000/api/chats', {
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setChats(data);
      }
    } catch (error) {
      console.error('Ошибка загрузки чатов:', error);
    } finally {
      setLoading(false);
    }
  };

  // Загружаем питомцев пользователя
  const fetchUserPets = async () => {
    try {
      // Загружаем животных, которыми владеет пользователь
      const ownedResponse = await fetch(`http://localhost:8100/api/pets/user/${user?.id}`, {
        credentials: 'include',
      });

      if (ownedResponse.ok) {
        const data = await ownedResponse.json();
        const pets = data.data || data || [];
        setUserPets(Array.isArray(pets) ? pets : []);
        console.log('✅ Загружены мои питомцы:', pets);
      } else {
        setUserPets([]);
      }

      // Пока нет рабочего endpoint для животных на попечении
      setCuratorPets([]);
    } catch (error) {
      console.error('Ошибка загрузки питомцев:', error);
      setUserPets([]);
      setCuratorPets([]);
    }
  };

  // Переключаем выбор питомца
  const togglePetSelection = (petId: number) => {
    setSelectedPets(prev => 
      prev.includes(petId) 
        ? prev.filter(id => id !== petId)
        : [...prev, petId]
    );
  };

  // Отправляем выбранных животных
  const handleSendPets = async () => {
    if (selectedPets.length === 0 || !selectedChatId) return;

    const selectedChat = chats.find(c => c.id === selectedChatId);
    if (!selectedChat?.other_user?.id) {
      console.error('❌ Не найден получатель сообщения');
      return;
    }

    setSending(true);
    setShowPetsModal(false);

    try {
      // Отправляем каждого животного отдельным сообщением
      for (const petId of selectedPets) {
        // Находим информацию о животном
        const pet = userPets.find(p => p.id === petId);
        if (!pet) continue;

        // Формируем JSON с информацией о животном
        const petData = {
          type: 'pet',
          id: pet.id,
          name: pet.name,
          species: pet.species || 'Животное',
          photo: pet.photo || null,
        };

        // Отправляем как JSON в контенте сообщения
        const response = await fetch(`http://localhost:8000/api/messages/send`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            receiver_id: selectedChat.other_user.id,
            content: JSON.stringify(petData),
          }),
        });

        if (!response.ok) {
          console.error(`❌ Ошибка отправки питомца ${petId}`);
        }
      }

      // Очищаем выбор и перезагружаем сообщения
      setSelectedPets([]);
      
      // Если это был временный чат, обновляем список
      if (selectedChatId < 0) {
        await fetchChats();
        const updatedChatsResponse = await fetch('http://localhost:8000/api/chats', {
          credentials: 'include',
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
    } catch (error) {
      console.error('❌ Ошибка отправки животных:', error);
      alert('Ошибка отправки животных');
    } finally {
      setSending(false);
    }
  };

  const loadUserAndCreateTempChat = async (targetUserId: number) => {
    try {
      // Загружаем информацию о пользователе
      const response = await fetch(`http://localhost:8000/api/users/${targetUserId}`, {
        credentials: 'include',
      });

      if (response.ok) {
        const userData = await response.json();
        
        // Создаем временный чат с этим пользователем
        const tempChat: Chat = {
          id: -targetUserId, // Временный отрицательный ID
          other_user: {
            id: userData.id,
            name: userData.name,
            last_name: userData.last_name || '',
            avatar: userData.avatar,
            is_online: userData.is_online,
            last_seen: userData.last_seen,
          },
          unread_count: 0,
        };
        
        // Добавляем временный чат в начало списка
        setChats([tempChat, ...chats]);
        setSelectedChatId(tempChat.id);
        setMessages([]); // Пустой список сообщений
      }
    } catch (error) {
      console.error('Ошибка загрузки пользователя:', error);
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!messageText.trim() || !selectedChatId || sending) return;

    // Находим выбранный чат чтобы получить receiver_id
    const selectedChat = chats.find(c => c.id === selectedChatId);
    if (!selectedChat?.other_user?.id) {
      console.error('❌ Не найден получатель сообщения');
      return;
    }

    setSending(true);
    
    try {
      console.log('📤 Отправка сообщения:', {
        receiver_id: selectedChat.other_user.id,
        content: messageText.trim()
      });

      const response = await fetch(`http://localhost:8000/api/messages/send`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          receiver_id: selectedChat.other_user.id,
          content: messageText.trim(),
        }),
      });

      console.log('📥 Ответ сервера:', {
        status: response.status,
        ok: response.ok
      });

      if (response.ok) {
        const responseData = await response.json();
        console.log('✅ Полученные данные:', responseData);
        
        // Очищаем поле ввода
        setMessageText('');
        
        // Если это был временный чат (ID < 0), нужно обновить список чатов и найти реальный чат
        if (selectedChatId < 0) {
          await fetchChats();
          
          // Находим реальный чат с этим пользователем
          const updatedChatsResponse = await fetch('http://localhost:8000/api/chats', {
            credentials: 'include',
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
          // Обычный чат - просто перезагружаем сообщения
          fetchMessages(selectedChatId);
          fetchChats();
        }
      } else {
        const errorText = await response.text();
        console.error('❌ Ошибка от сервера:', errorText);
      }
    } catch (error) {
      console.error('❌ Ошибка отправки сообщения:', error);
    } finally {
      setSending(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage(e as any);
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0 || !selectedChatId) return;

    const selectedChat = chats.find(c => c.id === selectedChatId);
    if (!selectedChat?.other_user?.id) {
      console.error('❌ Не найден получатель сообщения');
      return;
    }

    setSending(true);
    setShowAttachMenu(false);

    try {
      const formData = new FormData();
      formData.append('receiver_id', selectedChat.other_user.id.toString());
      
      // Добавляем текст сообщения если есть
      if (messageText.trim()) {
        formData.append('content', messageText.trim());
      }

      // Добавляем файлы
      Array.from(files).forEach(file => {
        formData.append('media', file);
      });

      const response = await fetch('http://localhost:8000/api/messages/send-media', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      if (response.ok) {
        // Очищаем поле ввода
        setMessageText('');
        
        // Если это был временный чат, обновляем список
        if (selectedChatId < 0) {
          await fetchChats();
          const updatedChatsResponse = await fetch('http://localhost:8000/api/chats', {
            credentials: 'include',
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
        console.error('❌ Ошибка отправки медиа:', errorText);
        alert('Ошибка отправки файла');
      }
    } catch (error) {
      console.error('❌ Ошибка отправки медиа:', error);
      alert('Ошибка отправки файла');
    } finally {
      setSending(false);
      // Очищаем input
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const formatTime = (dateString?: string) => {
    if (!dateString) return '';
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

  return (
    <div className="h-[calc(100vh-74px)] bg-white rounded-lg shadow-sm border border-gray-200 flex overflow-hidden">
      {/* Левая панель - список чатов */}
      <div className={`${isCollapsed ? 'w-[70px]' : 'w-[340px]'} border-r border-gray-200 flex flex-col transition-all duration-300`}>
        {/* Шапка */}
        <div className="p-4 border-b border-gray-200">
          <div className={`flex items-center ${isCollapsed ? 'justify-center' : 'justify-between'} mb-3`}>
            {!isCollapsed && <h2 className="text-xl font-semibold text-gray-900">Чаты</h2>}
            <button
              onClick={() => setIsCollapsed(!isCollapsed)}
              className="p-2 hover:bg-gray-100 rounded-full transition-colors"
              title={isCollapsed ? 'Развернуть' : 'Свернуть'}
            >
              <svg 
                className={`w-5 h-5 text-gray-600 transition-transform ${isCollapsed ? 'rotate-180' : ''}`} 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
              </svg>
            </button>
          </div>

          {/* Поиск - только в развернутом виде */}
          {!isCollapsed && (
            <div className="relative">
              <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input
                type="text"
                placeholder="Поиск"
                className="w-full pl-10 pr-4 py-2 bg-gray-100 border-0 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              />
            </div>
          )}
        </div>

        {/* Список чатов */}
        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            </div>
          ) : chats.length === 0 ? (
            <div className="flex items-center justify-center h-full text-gray-400">
              <div className="text-center">
                <svg className="w-16 h-16 mx-auto mb-2 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                {!isCollapsed && <p className="text-sm">Нет чатов</p>}
              </div>
            </div>
          ) : (
            chats.map((chat) => (
              <button
                key={chat.id}
                onClick={() => setSelectedChatId(chat.id)}
                className={`w-full ${isCollapsed ? 'p-2' : 'p-3'} flex ${isCollapsed ? 'justify-center' : 'items-start gap-3'} hover:bg-gray-50 transition-colors border-b border-gray-100 ${selectedChatId === chat.id ? 'bg-blue-50' : ''}`}
                title={isCollapsed ? `${chat.other_user?.name} ${chat.other_user?.last_name}` : ''}
              >
                {/* Аватар */}
                <div className="relative flex-shrink-0">
                  {chat.other_user?.avatar ? (
                    <img
                      src={`http://localhost:8000${chat.other_user.avatar}`}
                      alt={chat.other_user.name}
                      className={`${isCollapsed ? 'w-10 h-10' : 'w-12 h-12'} rounded-full object-cover`}
                    />
                  ) : (
                    <div className={`${isCollapsed ? 'w-10 h-10' : 'w-12 h-12'} rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold`}>
                      {chat.other_user?.name?.[0]?.toUpperCase() || '?'}
                    </div>
                  )}
                  {chat.unread_count > 0 && (
                    <div className="absolute -top-1 -right-1 min-w-[20px] h-5 px-1.5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center font-bold">
                      {chat.unread_count > 9 ? '9+' : chat.unread_count}
                    </div>
                  )}
                </div>

                {/* Информация - только в развернутом виде */}
                {!isCollapsed && (
                  <div className="flex-1 min-w-0 text-left">
                    <div className="flex items-center justify-between mb-1">
                      <span className={`font-medium truncate ${chat.unread_count > 0 ? 'font-semibold' : ''}`}>
                        {chat.other_user?.name} {chat.other_user?.last_name}
                      </span>
                      {chat.last_message_at && (
                        <span className="text-xs text-gray-500 ml-2 flex-shrink-0">
                          {formatTime(chat.last_message_at)}
                        </span>
                      )}
                    </div>
                    
                    {chat.last_message && (
                      <p className={`text-sm truncate ${chat.unread_count > 0 ? 'text-gray-900 font-medium' : 'text-gray-600'}`}>
                        {chat.last_message.sender_id === user?.id && 'Вы: '}
                        {(() => {
                          // Проверяем, является ли это сообщением с животным
                          try {
                            if (chat.last_message.content && chat.last_message.content.startsWith('{')) {
                              const parsed = JSON.parse(chat.last_message.content);
                              if (parsed.type === 'pet') {
                                return `🐾 ${parsed.name}`;
                              }
                            }
                          } catch (e) {
                            // Не JSON, обычное сообщение
                          }
                          return chat.last_message.content || '📎 Вложение';
                        })()}
                      </p>
                    )}
                  </div>
                )}
              </button>
            ))
          )}
        </div>
      </div>

      {/* Правая часть - окно чата */}
      <div className="flex-1 flex flex-col bg-gray-50">
        {selectedChatId ? (
          <>
            {/* Шапка чата */}
            <div className="bg-white border-b border-gray-200 p-2 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <button 
                  onClick={() => {
                    setSelectedChatId(null);
                    setMessages([]);
                  }}
                  className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                  title="Закрыть чат"
                >
                  <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
                
                {(() => {
                  const selectedChat = chats.find(c => c.id === selectedChatId);
                  const otherUser = selectedChat?.other_user;
                  
                  if (!otherUser) return null;
                  
                  return (
                    <div className="flex items-center gap-3">
                      {otherUser.avatar ? (
                        <img
                          src={`http://localhost:8000${otherUser.avatar}`}
                          alt={otherUser.name}
                          className="w-10 h-10 rounded-full object-cover"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                          {otherUser.name?.[0]?.toUpperCase() || '?'}
                        </div>
                      )}
                      
                      <div>
                        <h3 className="font-semibold text-gray-900">
                          {otherUser.name} {otherUser.last_name}
                        </h3>
                        {otherUser.is_online ? (
                          <p className="text-sm text-green-600">онлайн</p>
                        ) : otherUser.last_seen ? (
                          <p className="text-sm text-gray-500">
                            {(() => {
                              const lastSeen = new Date(otherUser.last_seen);
                              const now = new Date();
                              const diffMs = now.getTime() - lastSeen.getTime();
                              const diffMins = Math.floor(diffMs / 60000);
                              const diffHours = Math.floor(diffMs / 3600000);
                              const diffDays = Math.floor(diffMs / 86400000);

                              if (diffMins < 1) return 'был(а) только что';
                              if (diffMins < 60) return `был(а) ${diffMins} мин. назад`;
                              if (diffHours < 24) return `был(а) ${diffHours} ч. назад`;
                              if (diffDays === 1) return 'был(а) вчера';
                              if (diffDays < 7) return `был(а) ${diffDays} дн. назад`;
                              return `был(а) ${lastSeen.toLocaleDateString('ru-RU')}`;
                            })()}
                          </p>
                        ) : (
                          <p className="text-sm text-gray-400">не в сети</p>
                        )}
                      </div>
                    </div>
                  );
                })()}
              </div>
            </div>

            {/* Область сообщений */}
            <div className="flex-1 overflow-y-auto p-4" style={{
              backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23e3f2fd' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
              backgroundColor: '#e3f2fd'
            }}>
              <div className="max-w-4xl mx-auto space-y-3">
                {messages.map((message, index) => {
                  const isMyMessage = message.sender_id === user?.id;
                  
                  let messageTime = '';
                  try {
                    if (message.created_at) {
                      const date = new Date(message.created_at);
                      if (!isNaN(date.getTime())) {
                        messageTime = date.toLocaleTimeString('ru-RU', { 
                          hour: '2-digit', 
                          minute: '2-digit' 
                        });
                      }
                    }
                  } catch (e) {
                    console.error('Ошибка парсинга даты:', message.created_at);
                  }

                  // Проверяем, является ли это сообщением с животным
                  let petData = null;
                  try {
                    if (message.content && message.content.startsWith('{')) {
                      const parsed = JSON.parse(message.content);
                      if (parsed.type === 'pet') {
                        petData = parsed;
                      }
                    }
                  } catch (e) {
                    // Не JSON, обычное сообщение
                  }

                  // Пропускаем сообщения без контента и без вложений
                  if (!message.content && (!message.attachments || message.attachments.length === 0)) {
                    return null;
                  }

                  return (
                    <div key={`${message.id}-${index}`} className={`flex items-end gap-2 ${isMyMessage ? 'justify-end' : ''}`}>
                      {/* Сообщение с животным */}
                      {petData ? (
                        <button
                          onClick={() => router.push(`/pets/${petData.id}`)}
                          className={`${isMyMessage ? 'bg-blue-500 text-white rounded-2xl rounded-br-md hover:bg-blue-600' : 'bg-white text-gray-900 rounded-2xl rounded-bl-md hover:bg-gray-50'} p-3 shadow-sm max-w-[70%] transition-colors text-left`}
                        >
                          <div className="flex items-center gap-3">
                            {/* Фото животного */}
                            <div className="w-16 h-16 rounded-lg bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center text-white font-semibold overflow-hidden flex-shrink-0">
                              {petData.photo ? (
                                <img 
                                  src={getMediaUrl(petData.photo)} 
                                  alt={petData.name} 
                                  className="w-full h-full object-cover"
                                  onError={(e) => {
                                    const img = e.target as HTMLImageElement;
                                    img.style.display = 'none';
                                  }}
                                />
                              ) : (
                                <span className="text-2xl">🐾</span>
                              )}
                            </div>

                            {/* Информация о животном */}
                            <div className="flex-1">
                              <div className={`font-semibold text-sm ${isMyMessage ? 'text-white' : 'text-gray-900'}`}>
                                {petData.name}
                              </div>
                              <div className={`text-xs ${isMyMessage ? 'text-blue-100' : 'text-gray-600'}`}>
                                {petData.species}
                              </div>
                              {messageTime && (
                                <div className={`text-xs mt-1 ${isMyMessage ? 'text-blue-100' : 'text-gray-500'}`}>
                                  {messageTime}
                                </div>
                              )}
                            </div>
                          </div>
                        </button>
                      ) : (
                        <div className={`${isMyMessage ? 'bg-blue-500 text-white rounded-2xl rounded-br-md' : 'bg-white text-gray-900 rounded-2xl rounded-bl-md'} px-4 py-2 shadow-sm max-w-[70%]`}>
                        {/* Вложения */}
                        {message.attachments && message.attachments.length > 0 && (
                          <div className="mb-2 space-y-2">
                            {message.attachments.map((attachment) => (
                              <div key={attachment.id}>
                                {attachment.file_type === 'image' && (
                                  <img
                                    src={`http://localhost:8000${attachment.file_path}`}
                                    alt="Изображение"
                                    className="rounded-lg max-w-full h-auto cursor-pointer hover:opacity-90 transition-opacity"
                                    onClick={() => window.open(`http://localhost:8000${attachment.file_path}`, '_blank')}
                                    onLoad={() => scrollToBottom()}
                                  />
                                )}
                                {attachment.file_type === 'video' && (
                                  <video
                                    src={`http://localhost:8000${attachment.file_path}`}
                                    controls
                                    className="rounded-lg max-w-full h-auto"
                                    onLoadedData={() => scrollToBottom()}
                                  />
                                )}
                                {attachment.file_type !== 'image' && attachment.file_type !== 'video' && (() => {
                                  const fileName = attachment.file_path.split('/').pop() || 'Файл';
                                  const ext = fileName.split('.').pop()?.toLowerCase() || '';
                                  
                                  // Определяем иконку и цвет по расширению
                                  let iconColor = 'text-gray-600';
                                  let bgColor = 'bg-gray-100';
                                  let label = ext.toUpperCase() || 'FILE';
                                  let icon = null;
                                  
                                  if (['pdf'].includes(ext)) {
                                    iconColor = 'text-red-600';
                                    bgColor = 'bg-red-100';
                                    label = 'PDF';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
                                        <path d="M14 2v6h6M9 13h6M9 17h6M9 9h1"/>
                                      </svg>
                                    );
                                  } else if (['doc', 'docx'].includes(ext)) {
                                    iconColor = 'text-blue-600';
                                    bgColor = 'bg-blue-100';
                                    label = 'DOC';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
                                        <path d="M14 2v6h6M9 13h6M9 17h6M9 9h1"/>
                                      </svg>
                                    );
                                  } else if (['xls', 'xlsx'].includes(ext)) {
                                    iconColor = 'text-green-600';
                                    bgColor = 'bg-green-100';
                                    label = 'XLS';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
                                        <path d="M14 2v6h6M8 13h8M8 17h8M8 9h2"/>
                                      </svg>
                                    );
                                  } else if (['ppt', 'pptx'].includes(ext)) {
                                    iconColor = 'text-orange-600';
                                    bgColor = 'bg-orange-100';
                                    label = 'PPT';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
                                        <path d="M14 2v6h6M10 13h4M10 17h4"/>
                                      </svg>
                                    );
                                  } else if (['zip', 'rar', '7z', 'tar', 'gz'].includes(ext)) {
                                    iconColor = 'text-purple-600';
                                    bgColor = 'bg-purple-100';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                                      </svg>
                                    );
                                  } else if (['txt', 'md', 'json', 'xml', 'csv'].includes(ext)) {
                                    iconColor = 'text-gray-600';
                                    bgColor = 'bg-gray-100';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                      </svg>
                                    );
                                  } else if (['js', 'ts', 'jsx', 'tsx', 'py', 'java', 'cpp', 'c', 'go', 'rs', 'php', 'rb', 'swift'].includes(ext)) {
                                    iconColor = 'text-indigo-600';
                                    bgColor = 'bg-indigo-100';
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                                      </svg>
                                    );
                                  } else {
                                    // По умолчанию
                                    icon = (
                                      <svg className={`w-6 h-6 ${isMyMessage ? 'text-white' : iconColor}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                      </svg>
                                    );
                                  }
                                  
                                  return (
                                    <a
                                      href={`http://localhost:8000${attachment.file_path}`}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      className={`flex items-center gap-3 p-3 rounded-lg ${isMyMessage ? 'bg-blue-600' : 'bg-white border border-gray-200'} hover:opacity-80 transition-opacity`}
                                    >
                                      <div className={`w-10 h-10 rounded-lg ${isMyMessage ? 'bg-blue-700' : bgColor} flex items-center justify-center flex-shrink-0`}>
                                        {icon}
                                      </div>
                                      <div className="flex-1 min-w-0">
                                        <p className={`text-sm font-medium truncate ${isMyMessage ? 'text-white' : 'text-gray-900'}`}>
                                          {fileName}
                                        </p>
                                        <div className="flex items-center gap-2">
                                          <span className={`text-xs font-medium ${isMyMessage ? 'text-blue-100' : iconColor}`}>
                                            {label}
                                          </span>
                                          <span className={`text-xs ${isMyMessage ? 'text-blue-100' : 'text-gray-500'}`}>
                                            • {(attachment.file_size / 1024).toFixed(1)} KB
                                          </span>
                                        </div>
                                      </div>
                                      <svg className={`w-5 h-5 ${isMyMessage ? 'text-white' : 'text-gray-400'} flex-shrink-0`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                      </svg>
                                    </a>
                                  );
                                })()}
                              </div>
                            ))}
                          </div>
                        )}
                        
                        {/* Текст сообщения */}
                        {message.content && <p>{message.content}</p>}
                        
                        {messageTime && (
                          <div className={`flex items-center ${isMyMessage ? 'justify-end' : ''} gap-1 mt-1`}>
                            <span className={`text-xs ${isMyMessage ? 'text-blue-100' : 'text-gray-500'}`}>
                              {messageTime}
                            </span>
                            {isMyMessage && (
                              <>
                                {message.is_read ? (
                                  // Две галочки - прочитано
                                  <>
                                    <svg className="w-4 h-4 text-blue-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                    <svg className="w-4 h-4 text-blue-200 -ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                  </>
                                ) : (
                                  // Одна галочка - отправлено, но не прочитано
                                  <svg className="w-4 h-4 text-blue-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                  </svg>
                                )}
                              </>
                            )}
                          </div>
                        )}
                        </div>
                      )}
                    </div>
                  );
                })}
                {/* Элемент для прокрутки к последнему сообщению */}
                <div ref={messagesEndRef} />
              </div>
            </div>

            {/* Поле ввода */}
            <div className="bg-white border-t border-gray-200 p-4 relative">
              {/* Скрытый input для выбора файлов */}
              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept="image/*,video/*,application/*"
                onChange={handleFileSelect}
                className="hidden"
              />
              
              {/* Меню вложений */}
              {showAttachMenu && (
                <div className="absolute bottom-16 left-4 bg-white rounded-xl shadow-lg p-1 w-48 border border-gray-200">
                  <button 
                    onClick={() => {
                      setShowAttachMenu(false);
                      if (fileInputRef.current) {
                        fileInputRef.current.accept = 'image/*';
                        fileInputRef.current.click();
                      }
                    }}
                    className="w-full flex items-center gap-3 px-3 py-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                    </div>
                    <span className="text-gray-900 text-sm">Фото</span>
                  </button>

                  <button 
                    onClick={() => {
                      setShowAttachMenu(false);
                      if (fileInputRef.current) {
                        fileInputRef.current.accept = 'video/*';
                        fileInputRef.current.click();
                      }
                    }}
                    className="w-full flex items-center gap-3 px-3 py-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                    </div>
                    <span className="text-gray-900 text-sm">Видео</span>
                  </button>

                  <button 
                    onClick={() => {
                      setShowAttachMenu(false);
                      if (fileInputRef.current) {
                        fileInputRef.current.accept = '*/*';
                        fileInputRef.current.click();
                      }
                    }}
                    className="w-full flex items-center gap-3 px-3 py-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <svg className="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                    </div>
                    <span className="text-gray-900 text-sm">Файл</span>
                  </button>

                  <div className="border-t border-gray-200 my-1"></div>

                  <button 
                    onClick={() => {
                      setShowAttachMenu(false);
                      setShowPetsModal(true);
                    }}
                    className="w-full flex items-center gap-3 px-3 py-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <div className="w-8 h-8 bg-amber-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <span className="text-lg">🐾</span>
                    </div>
                    <span className="text-gray-900 text-sm">Животное</span>
                  </button>
                </div>
              )}

              <div className="flex items-center gap-3">
                <button 
                  onClick={() => setShowAttachMenu(!showAttachMenu)}
                  className="p-2 hover:bg-gray-100 rounded-full transition-colors flex-shrink-0"
                >
                  <svg className="w-6 h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                </button>

                <div className="flex-1 relative">
                  <input
                    type="text"
                    placeholder="Сообщение"
                    value={messageText}
                    onChange={(e) => setMessageText(e.target.value)}
                    onKeyPress={handleKeyPress}
                    disabled={sending}
                    className="w-full px-4 py-2.5 bg-gray-100 border-0 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm pr-12 disabled:opacity-50"
                  />
                  <button className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 hover:bg-gray-200 rounded-full transition-colors">
                    <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                </div>

                <button 
                  onClick={handleSendMessage}
                  disabled={!messageText.trim() || sending}
                  className="w-10 h-10 bg-blue-500 hover:bg-blue-600 rounded-full transition-all flex items-center justify-center shadow-md hover:shadow-lg active:scale-95 flex-shrink-0 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-blue-500"
                >
                  {sending ? (
                    <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                  ) : (
                    <svg className="w-5 h-5 text-white translate-x-0.5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                    </svg>
                  )}
                </button>
              </div>
            </div>
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

      {/* Модальное окно выбора питомцев */}
      {showPetsModal && (
        <div className="fixed inset-0 bg-black/50 flex items-end z-50">
          <div className="bg-white w-full rounded-t-2xl flex flex-col max-h-[80vh]">
            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-gray-200">
              <h3 className="text-lg font-semibold text-gray-900">Выберите животное</h3>
              <button 
                onClick={() => {
                  setShowPetsModal(false);
                  setSelectedPets([]);
                }}
                className="p-2 hover:bg-gray-100 rounded-full transition-colors"
              >
                <svg className="w-6 h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto p-4">
              {userPets.length === 0 && curatorPets.length === 0 ? (
                <div className="text-center py-12">
                  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-100 flex items-center justify-center">
                    <span className="text-4xl">🐾</span>
                  </div>
                  <h4 className="text-base font-semibold text-gray-900 mb-2">
                    У вас пока нет питомцев
                  </h4>
                  <p className="text-sm text-gray-500">
                    Добавьте своего первого питомца, чтобы делиться им в сообщениях
                  </p>
                </div>
              ) : (
                <div className="space-y-4">
                  {/* Мои животные */}
                  {userPets.length > 0 && (
                    <div>
                      <h4 className="text-sm font-semibold text-gray-900 mb-2 px-1">Мои животные</h4>
                      <div className="space-y-2">
                        {userPets.map((pet) => {
                          const petPhotoUrl = getMediaUrl(pet.photo);
                          
                          return (
                            <button
                              key={pet.id}
                              onClick={() => togglePetSelection(pet.id)}
                              className={`w-full flex items-center gap-3 p-3 rounded-lg border-2 transition-all ${
                                selectedPets.includes(pet.id)
                                  ? 'border-blue-500 bg-blue-50'
                                  : 'border-gray-200 hover:border-gray-300 bg-white'
                              }`}
                            >
                              {/* Pet Photo */}
                              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center text-white font-semibold overflow-hidden flex-shrink-0">
                                {petPhotoUrl ? (
                                  <img 
                                    src={petPhotoUrl} 
                                    alt={pet.name} 
                                    className="w-full h-full object-cover"
                                    onError={(e) => {
                                      const img = e.target as HTMLImageElement;
                                      img.style.display = 'none';
                                    }}
                                  />
                                ) : (
                                  pet.name[0]?.toUpperCase()
                                )}
                              </div>

                              {/* Pet Info */}
                              <div className="flex-1 text-left">
                                <div className="font-semibold text-sm text-gray-900">{pet.name}</div>
                                <div className="text-xs text-gray-600">{pet.species || 'Животное'}</div>
                              </div>

                              {/* Checkmark */}
                              {selectedPets.includes(pet.id) && (
                                <div className="w-6 h-6 rounded-full bg-blue-500 flex items-center justify-center flex-shrink-0">
                                  <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                  </svg>
                                </div>
                              )}
                            </button>
                          );
                        })}
                      </div>
                    </div>
                  )}

                  {/* Животные на попечении */}
                  {curatorPets.length > 0 && (
                    <div>
                      <h4 className="text-sm font-semibold text-gray-900 mb-2 px-1">На попечении</h4>
                      <div className="space-y-2">
                        {curatorPets.map((pet) => {
                          const petPhotoUrl = getMediaUrl(pet.photo);
                          
                          return (
                            <button
                              key={pet.id}
                              onClick={() => togglePetSelection(pet.id)}
                              className={`w-full flex items-center gap-3 p-3 rounded-lg border-2 transition-all ${
                                selectedPets.includes(pet.id)
                                  ? 'border-blue-500 bg-blue-50'
                                  : 'border-gray-200 hover:border-gray-300 bg-white'
                              }`}
                            >
                              {/* Pet Photo */}
                              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-amber-400 to-orange-400 flex items-center justify-center text-white font-semibold overflow-hidden flex-shrink-0">
                                {petPhotoUrl ? (
                                  <img 
                                    src={petPhotoUrl} 
                                    alt={pet.name} 
                                    className="w-full h-full object-cover"
                                    onError={(e) => {
                                      const img = e.target as HTMLImageElement;
                                      img.style.display = 'none';
                                    }}
                                  />
                                ) : (
                                  pet.name[0]?.toUpperCase()
                                )}
                              </div>

                              {/* Pet Info */}
                              <div className="flex-1 text-left">
                                <div className="font-semibold text-sm text-gray-900">{pet.name}</div>
                                <div className="text-xs text-gray-600">{pet.species || 'Животное'}</div>
                              </div>

                              {/* Checkmark */}
                              {selectedPets.includes(pet.id) && (
                                <div className="w-6 h-6 rounded-full bg-blue-500 flex items-center justify-center flex-shrink-0">
                                  <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                  </svg>
                                </div>
                              )}
                            </button>
                          );
                        })}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="border-t border-gray-200 p-4 bg-gray-50 flex gap-3">
              <button 
                onClick={() => {
                  setShowPetsModal(false);
                  setSelectedPets([]);
                }}
                className="flex-1 px-4 py-2.5 border border-gray-300 text-gray-900 rounded-full text-sm font-semibold hover:bg-gray-100 transition-colors"
              >
                Отмена
              </button>
              <button 
                onClick={handleSendPets}
                disabled={selectedPets.length === 0 || sending}
                className="flex-1 px-4 py-2.5 bg-blue-500 text-white rounded-full text-sm font-semibold hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {sending ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
                    Отправка...
                  </>
                ) : (
                  <>
                    Отправить ({selectedPets.length})
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
