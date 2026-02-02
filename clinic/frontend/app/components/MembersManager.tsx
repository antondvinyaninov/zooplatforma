'use client';

import { useState, useEffect } from 'react';
import {
  UserGroupIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';

interface Member {
  id: number;
  user_id: number;
  name?: string;
  last_name?: string;
  email?: string;
  avatar?: string;
  role: string;
  position?: string;
  joined_at: string;
}

interface User {
  id: number;
  name?: string;
  last_name?: string;
  email?: string;
  avatar?: string;
}

const roleColors: Record<string, string> = {
  owner: 'bg-purple-100 text-purple-800 border-purple-200',
  admin: 'bg-blue-100 text-blue-800 border-blue-200',
  moderator: 'bg-green-100 text-green-800 border-green-200',
  member: 'bg-gray-100 text-gray-800 border-gray-200',
};

const roleLabels: Record<string, string> = {
  owner: 'Владелец',
  admin: 'Администратор',
  moderator: 'Модератор',
  member: 'Участник',
};

export default function MembersManager() {
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedMember, setSelectedMember] = useState<Member | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<User[]>([]);
  const [searching, setSearching] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newMemberRole, setNewMemberRole] = useState('member');
  const [newMemberPosition, setNewMemberPosition] = useState('');
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  // Helper функция для получения токена из cookie
  const getAuthToken = () => {
    // Сначала пробуем получить из cookie localhost:6300
    let token = document.cookie
      .split('; ')
      .find(row => row.startsWith('auth_token='))
      ?.split('=')[1];
    
    // Если нет, пробуем получить через Main API
    if (!token) {
      console.warn('No auth_token cookie found on localhost:6300');
    }
    
    return token;
  };

  // Helper функция для создания headers с авторизацией
  const getAuthHeaders = (clinicId: string) => {
    const token = getAuthToken();
    return {
      'Content-Type': 'application/json',
      'X-Clinic-ID': clinicId,
      ...(token && { 'Authorization': `Bearer ${token}` }),
    };
  };

  useEffect(() => {
    loadMembers();
  }, []);

  const loadMembers = async () => {
    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      console.log('🔍 localStorage selectedClinicId:', clinicId, 'type:', typeof clinicId);
      if (!clinicId) return;

      // Получаем токен: сначала из cookie, если нет - через Main API
      let token = getAuthToken();
      
      if (!token) {
        // Получаем токен через Main API
        try {
          const meResponse = await fetch('http://localhost:7100/api/auth/me', {
            method: 'GET',
            credentials: 'include', // Отправляет cookie от localhost:3000
          });
          
          if (meResponse.ok) {
            const meData = await meResponse.json();
            if (meData.success && meData.data && meData.data.token) {
              token = meData.data.token;
              console.log('✅ Token received from Main API');
            }
          } else {
            setMessage({ type: 'error', text: 'Пожалуйста, войдите через главный сайт (localhost:3000)' });
            setLoading(false);
            return;
          }
        } catch (error) {
          console.error('Failed to get token from Main API:', error);
          setMessage({ type: 'error', text: 'Ошибка получения токена. Пожалуйста, войдите через главный сайт' });
          setLoading(false);
          return;
        }
      }

      if (!token) {
        setMessage({ type: 'error', text: 'Токен авторизации не найден' });
        setLoading(false);
        return;
      }

      console.log('🔍 Sending request with clinic_id:', clinicId, 'token length:', token?.length);
      
      const response = await fetch('http://localhost:8600/api/members', {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'X-Clinic-ID': clinicId,
          'Authorization': `Bearer ${token}`,
        },
      });

      console.log('Members API response status:', response.status);

      if (response.ok) {
        const result = await response.json();
        console.log('Members API result:', result);
        if (result.success && result.data) {
          setMembers(result.data);
          console.log('✅ Members loaded:', result.data.length);
        } else {
          console.warn('No members data in response');
        }
      } else {
        const errorText = await response.text();
        console.error('Failed to load members:', response.status, errorText);
        setMessage({ type: 'error', text: `Ошибка загрузки участников: ${response.status}` });
      }
    } catch (error) {
      console.error('Failed to load members:', error);
      setMessage({ type: 'error', text: 'Ошибка загрузки участников' });
    }

    setLoading(false);
  };

  const searchUsers = async (query: string) => {
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }

    setSearching(true);

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) return;

      const response = await fetch(`http://localhost:8600/api/users/search?q=${encodeURIComponent(query)}`, {
        method: 'GET',
        credentials: 'include',
        headers: getAuthHeaders(clinicId),
      });

      if (response.ok) {
        const result = await response.json();
        if (result.success && result.data) {
          setSearchResults(result.data);
        }
      }
    } catch (error) {
      console.error('Failed to search users:', error);
    }

    setSearching(false);
  };

  const handleAddMember = async () => {
    if (!selectedUser) return;

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) return;

      const response = await fetch('http://localhost:8600/api/members', {
        method: 'POST',
        credentials: 'include',
        headers: getAuthHeaders(clinicId),
        body: JSON.stringify({
          user_id: selectedUser.id,
          role: newMemberRole,
          position: newMemberPosition,
        }),
      });

      if (response.ok) {
        setMessage({ type: 'success', text: 'Участник успешно добавлен' });
        setShowAddModal(false);
        setSelectedUser(null);
        setSearchQuery('');
        setSearchResults([]);
        setNewMemberRole('member');
        setNewMemberPosition('');
        loadMembers();
      } else {
        const error = await response.text();
        setMessage({ type: 'error', text: error || 'Ошибка при добавлении участника' });
      }
    } catch (error) {
      console.error('Failed to add member:', error);
      setMessage({ type: 'error', text: 'Ошибка при добавлении участника' });
    }
  };

  const handleUpdateMember = async () => {
    if (!selectedMember) return;

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) return;

      const response = await fetch('http://localhost:8600/api/members', {
        method: 'PUT',
        credentials: 'include',
        headers: getAuthHeaders(clinicId),
        body: JSON.stringify({
          member_id: selectedMember.id,
          role: selectedMember.role,
          position: selectedMember.position || '',
        }),
      });

      if (response.ok) {
        setMessage({ type: 'success', text: 'Участник успешно обновлён' });
        setShowEditModal(false);
        setSelectedMember(null);
        loadMembers();
      } else {
        const error = await response.text();
        setMessage({ type: 'error', text: error || 'Ошибка при обновлении участника' });
      }
    } catch (error) {
      console.error('Failed to update member:', error);
      setMessage({ type: 'error', text: 'Ошибка при обновлении участника' });
    }
  };

  const handleRemoveMember = async (member: Member) => {
    if (!confirm(`Удалить ${member.name || member.email} из команды?`)) return;

    try {
      const clinicId = localStorage.getItem('selectedClinicId');
      if (!clinicId) return;

      const response = await fetch('http://localhost:8600/api/members', {
        method: 'DELETE',
        credentials: 'include',
        headers: getAuthHeaders(clinicId),
        body: JSON.stringify({
          member_id: member.id,
        }),
      });

      if (response.ok) {
        setMessage({ type: 'success', text: 'Участник успешно удалён' });
        loadMembers();
      } else {
        const error = await response.text();
        setMessage({ type: 'error', text: error || 'Ошибка при удалении участника' });
      }
    } catch (error) {
      console.error('Failed to remove member:', error);
      setMessage({ type: 'error', text: 'Ошибка при удалении участника' });
    }
  };

  if (loading) {
    return <div className="text-center py-12 text-gray-400">Загрузка...</div>;
  }

  return (
    <div className="space-y-6">
      {message && (
        <div
          className={`p-4 rounded-lg ${
            message.type === 'success'
              ? 'bg-green-50 text-green-800 border border-green-200'
              : 'bg-red-50 text-red-800 border border-red-200'
          }`}
        >
          {message.text}
        </div>
      )}

      <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <UserGroupIcon className="w-6 h-6 text-blue-600" />
            <h3 className="text-xl font-bold text-gray-900">Команда клиники ({members.length})</h3>
          </div>
          <button
            onClick={() => setShowAddModal(true)}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <PlusIcon className="w-5 h-5" />
            Добавить участника
          </button>
        </div>

        <div className="space-y-3">
          {members.map((member) => (
            <div
              key={member.id}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-center gap-4">
                {member.avatar ? (
                  <img
                    src={member.avatar}
                    alt={member.name || member.email}
                    className="w-12 h-12 rounded-full object-cover"
                  />
                ) : (
                  <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold">
                    {(member.name || member.email || '?')[0].toUpperCase()}
                  </div>
                )}
                <div>
                  <div className="font-semibold text-gray-900">
                    {member.name && member.last_name
                      ? `${member.name} ${member.last_name}`
                      : member.name || member.email}
                  </div>
                  {member.position && (
                    <div className="text-sm text-gray-600">{member.position}</div>
                  )}
                  <div className="flex items-center gap-2 mt-1">
                    <span className={`text-xs px-2 py-1 rounded border ${roleColors[member.role]}`}>
                      {roleLabels[member.role]}
                    </span>
                  </div>
                </div>
              </div>

              {member.role !== 'owner' && (
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => {
                      setSelectedMember(member);
                      setShowEditModal(true);
                    }}
                    className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                    title="Редактировать"
                  >
                    <PencilIcon className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleRemoveMember(member)}
                    className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="Удалить"
                  >
                    <TrashIcon className="w-5 h-5" />
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Модальное окно добавления участника */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-8 max-w-md w-full max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-2xl font-bold text-gray-900">Добавить участника</h3>
              <button
                onClick={() => {
                  setShowAddModal(false);
                  setSelectedUser(null);
                  setSearchQuery('');
                  setSearchResults([]);
                }}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>

            <div className="space-y-4">
              {!selectedUser ? (
                <>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Поиск пользователя
                    </label>
                    <div className="relative">
                      <MagnifyingGlassIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                      <input
                        type="text"
                        value={searchQuery}
                        onChange={(e) => {
                          setSearchQuery(e.target.value);
                          searchUsers(e.target.value);
                        }}
                        className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="Имя, фамилия или email"
                      />
                    </div>
                  </div>

                  {searching && (
                    <div className="text-center py-4 text-gray-500">Поиск...</div>
                  )}

                  {searchResults.length > 0 && (
                    <div className="space-y-2 max-h-60 overflow-y-auto">
                      {searchResults.map((user) => (
                        <button
                          key={user.id}
                          onClick={() => setSelectedUser(user)}
                          className="w-full flex items-center gap-3 p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors text-left"
                        >
                          {user.avatar ? (
                            <img
                              src={user.avatar}
                              alt={user.name || user.email}
                              className="w-10 h-10 rounded-full object-cover"
                            />
                          ) : (
                            <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold">
                              {(user.name || user.email || '?')[0].toUpperCase()}
                            </div>
                          )}
                          <div>
                            <div className="font-medium text-gray-900">
                              {user.name && user.last_name
                                ? `${user.name} ${user.last_name}`
                                : user.name || user.email}
                            </div>
                            {user.email && (
                              <div className="text-sm text-gray-600">{user.email}</div>
                            )}
                          </div>
                        </button>
                      ))}
                    </div>
                  )}

                  {searchQuery.length >= 2 && !searching && searchResults.length === 0 && (
                    <div className="text-center py-4 text-gray-500">
                      Пользователи не найдены
                    </div>
                  )}
                </>
              ) : (
                <>
                  <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-lg">
                    {selectedUser.avatar ? (
                      <img
                        src={selectedUser.avatar}
                        alt={selectedUser.name || selectedUser.email}
                        className="w-12 h-12 rounded-full object-cover"
                      />
                    ) : (
                      <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold">
                        {(selectedUser.name || selectedUser.email || '?')[0].toUpperCase()}
                      </div>
                    )}
                    <div className="flex-1">
                      <div className="font-semibold text-gray-900">
                        {selectedUser.name && selectedUser.last_name
                          ? `${selectedUser.name} ${selectedUser.last_name}`
                          : selectedUser.name || selectedUser.email}
                      </div>
                      {selectedUser.email && (
                        <div className="text-sm text-gray-600">{selectedUser.email}</div>
                      )}
                    </div>
                    <button
                      onClick={() => setSelectedUser(null)}
                      className="p-1 hover:bg-gray-200 rounded"
                    >
                      <XMarkIcon className="w-5 h-5" />
                    </button>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Роль *
                    </label>
                    <select
                      value={newMemberRole}
                      onChange={(e) => setNewMemberRole(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                      <option value="member">Участник</option>
                      <option value="moderator">Модератор</option>
                      <option value="admin">Администратор</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Должность
                    </label>
                    <input
                      type="text"
                      value={newMemberPosition}
                      onChange={(e) => setNewMemberPosition(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="Ветеринар, администратор..."
                    />
                  </div>

                  <div className="flex gap-3 pt-4">
                    <button
                      onClick={handleAddMember}
                      className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
                    >
                      Добавить
                    </button>
                    <button
                      onClick={() => setSelectedUser(null)}
                      className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
                    >
                      Назад
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Модальное окно редактирования участника */}
      {showEditModal && selectedMember && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-8 max-w-md w-full">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-2xl font-bold text-gray-900">Редактировать участника</h3>
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setSelectedMember(null);
                }}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-lg">
                {selectedMember.avatar ? (
                  <img
                    src={selectedMember.avatar}
                    alt={selectedMember.name || selectedMember.email}
                    className="w-12 h-12 rounded-full object-cover"
                  />
                ) : (
                  <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold">
                    {(selectedMember.name || selectedMember.email || '?')[0].toUpperCase()}
                  </div>
                )}
                <div>
                  <div className="font-semibold text-gray-900">
                    {selectedMember.name && selectedMember.last_name
                      ? `${selectedMember.name} ${selectedMember.last_name}`
                      : selectedMember.name || selectedMember.email}
                  </div>
                  {selectedMember.email && (
                    <div className="text-sm text-gray-600">{selectedMember.email}</div>
                  )}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Роль *
                </label>
                <select
                  value={selectedMember.role}
                  onChange={(e) =>
                    setSelectedMember({ ...selectedMember, role: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="member">Участник</option>
                  <option value="moderator">Модератор</option>
                  <option value="admin">Администратор</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Должность
                </label>
                <input
                  type="text"
                  value={selectedMember.position || ''}
                  onChange={(e) =>
                    setSelectedMember({ ...selectedMember, position: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Ветеринар, администратор..."
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  onClick={handleUpdateMember}
                  className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
                >
                  Сохранить
                </button>
                <button
                  onClick={() => {
                    setShowEditModal(false);
                    setSelectedMember(null);
                  }}
                  className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
                >
                  Отмена
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
