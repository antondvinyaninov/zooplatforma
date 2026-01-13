'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { adminApi, type Stats, type User, type Post, type Log } from '@/lib/api';
import {
  UsersIcon,
  DocumentTextIcon,
  ChartBarIcon,
  ClipboardDocumentListIcon,
  BellIcon,
  Cog6ToothIcon,
  MagnifyingGlassIcon,
  TrashIcon,
  CheckIcon,
  XMarkIcon,
  ChevronLeftIcon,
  DocumentDuplicateIcon,
  BuildingOfficeIcon,
} from '@heroicons/react/24/outline';
import styles from './page.module.css';
import StatsWidget from './components/widgets/StatsWidget';
import TableWidget from './components/widgets/TableWidget';
import ChartWidget from './components/widgets/ChartWidget';

export default function AdminDashboard() {
  const router = useRouter();
  const [stats, setStats] = useState<Stats | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [logs, setLogs] = useState<Log[]>([]);
  const [activeTab, setActiveTab] = useState<'users' | 'posts' | 'stats' | 'logs' | 'organizations'>('stats');
  const [isClient, setIsClient] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [adminUser, setAdminUser] = useState<{ email: string; role: string } | null>(null);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  // Инициализация из localStorage после монтирования
  useEffect(() => {
    setIsClient(true);
    const savedTab = localStorage.getItem('adminActiveTab') as 'users' | 'posts' | 'stats' | 'logs' | 'organizations' | null;
    if (savedTab) {
      setActiveTab(savedTab);
    }
  }, []);

  useEffect(() => {
    if (isClient) {
      loadData();
    }
  }, [activeTab, isClient]);

  // Закрытие меню при клике вне его
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      const target = event.target as HTMLElement;
      if (userMenuOpen && !target.closest(`.${styles.userMenuWrapper}`)) {
        setUserMenuOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [userMenuOpen]);

  const loadData = async () => {
    setLoading(true);
    setError('');

    const meResult = await adminApi.me();
    if (!meResult.success) {
      router.push('/auth');
      return;
    }

    if (meResult.data) {
      setAdminUser(meResult.data);
    }

    if (activeTab === 'stats') {
      const result = await adminApi.getStats();
      if (result.success && result.data) {
        setStats(result.data);
      } else {
        setError(result.error || 'Ошибка загрузки статистики');
      }
    } else if (activeTab === 'users') {
      const result = await adminApi.getUsers();
      if (result.success && result.data) {
        setUsers(result.data);
      } else {
        setError(result.error || 'Ошибка загрузки пользователей');
      }
    } else if (activeTab === 'posts') {
      const result = await adminApi.getPosts();
      if (result.success && result.data) {
        setPosts(result.data);
      } else {
        setError(result.error || 'Ошибка загрузки постов');
      }
    } else if (activeTab === 'logs') {
      const result = await adminApi.getLogs();
      console.log('Logs result:', result);
      if (result.success) {
        console.log('Logs data:', result.data);
        setLogs(result.data || []);
      } else {
        // Если логов нет - это нормально, просто показываем пустой список
        setLogs([]);
        console.log('Logs error:', result.error);
      }
    }

    setLoading(false);
  };

  const handleDeleteUser = async (id: number) => {
    if (!confirm('Удалить пользователя?')) return;
    const result = await adminApi.deleteUser(id);
    if (result.success) {
      loadData();
    } else {
      alert(result.error);
    }
  };

  const handleDeletePost = async (id: number) => {
    if (!confirm('Удалить пост?')) return;
    const result = await adminApi.deletePost(id);
    if (result.success) {
      loadData();
    } else {
      alert(result.error);
    }
  };

  const handleModeratePost = async (id: number, status: string) => {
    const result = await adminApi.updatePost(id, { status });
    if (result.success) {
      loadData();
    } else {
      alert(result.error);
    }
  };

  const handleLogout = () => {
    // Удаляем cookie
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; domain=localhost';
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    // Редирект на страницу входа
    router.push('/auth');
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className="flex items-center gap-4">
            <button
              className="lg:hidden p-2 hover:bg-gray-100 rounded-lg transition-colors"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <h1 className={styles.logo}>
              <img src="/favicon.svg" alt="ЗооАдминка" className={styles.logoImage} />
              <span className={styles.logoText}>ЗооАдминка</span>
            </h1>
          </div>
          <div className={styles.headerRight}>
            {adminUser && (
              <>
                <button className={styles.iconBtn}>
                  <ClipboardDocumentListIcon className="w-4 h-4" />
                </button>
                <button className={styles.iconBtn}>
                  <BellIcon className="w-4 h-4" />
                </button>
                <button className={styles.iconBtn}>
                  <Cog6ToothIcon className="w-4 h-4" />
                </button>
                <div className={styles.userMenuWrapper}>
                  <button 
                    className={styles.userMenuButton}
                    onClick={() => setUserMenuOpen(!userMenuOpen)}
                  >
                    <div className={styles.userAvatar}>
                      <UsersIcon className="w-4 h-4" />
                    </div>
                    <span className={styles.userMenuArrow}>▼</span>
                  </button>
                  
                  {userMenuOpen && (
                    <div className={styles.userMenuDropdown}>
                      <div className={styles.userMenuHeader}>
                        <div className={styles.userAvatar}>
                          <UsersIcon className="w-4 h-4" />
                        </div>
                        <div className={styles.userInfo}>
                          <div className={styles.userEmail}>{adminUser.email}</div>
                          <div className={styles.userRole}>{adminUser.role}</div>
                        </div>
                      </div>
                      
                      <div className={styles.userMenuItems}>
                        <a 
                          href="http://localhost:3000" 
                          target="_blank"
                          rel="noopener noreferrer"
                          className={styles.userMenuItem}
                        >
                          Главная страница
                        </a>
                        <button 
                          className={styles.userMenuItem}
                          onClick={handleLogout}
                        >
                          Выйти
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      </header>

      <div className={styles.layout}>
        {/* Mobile overlay */}
        {mobileMenuOpen && (
          <div 
            className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
            onClick={() => setMobileMenuOpen(false)}
          />
        )}
        
        <aside className={`${styles.sidebar} ${sidebarCollapsed ? styles.sidebarCollapsed : ''} ${mobileMenuOpen ? styles.sidebarMobile : styles.sidebarHidden}`}>
          {isClient ? (
            <nav className={styles.nav}>
              <button
                className={`${styles.navItem} ${activeTab === 'users' ? styles.navItemActive : ''}`}
                onClick={() => {
                  setActiveTab('users');
                  localStorage.setItem('adminActiveTab', 'users');
                }}
                title="Пользователи"
              >
                <UsersIcon />
                {!sidebarCollapsed && <span className={styles.navLabel}>Пользователи</span>}
              </button>
              <button
                className={`${styles.navItem} ${activeTab === 'posts' ? styles.navItemActive : ''}`}
                onClick={() => {
                  setActiveTab('posts');
                  localStorage.setItem('adminActiveTab', 'posts');
                }}
                title="Посты"
              >
                <DocumentTextIcon />
                {!sidebarCollapsed && <span className={styles.navLabel}>Посты</span>}
              </button>
              <button
                className={`${styles.navItem} ${activeTab === 'stats' ? styles.navItemActive : ''}`}
                onClick={() => {
                  setActiveTab('stats');
                  localStorage.setItem('adminActiveTab', 'stats');
                }}
                title="Статистика"
              >
                <ChartBarIcon />
                {!sidebarCollapsed && <span className={styles.navLabel}>Статистика</span>}
              </button>
              <button
                className={`${styles.navItem} ${activeTab === 'logs' ? styles.navItemActive : ''}`}
                onClick={() => {
                  setActiveTab('logs');
                  localStorage.setItem('adminActiveTab', 'logs');
                }}
                title="Логирование"
              >
                <DocumentDuplicateIcon />
                {!sidebarCollapsed && <span className={styles.navLabel}>Логирование</span>}
              </button>
              <button
                className={`${styles.navItem} ${activeTab === 'organizations' ? styles.navItemActive : ''}`}
                onClick={() => {
                  setActiveTab('organizations');
                  localStorage.setItem('adminActiveTab', 'organizations');
                  // Редирект на отдельную страницу
                  window.location.href = '/organizations';
                }}
                title="Организации"
              >
                <BuildingOfficeIcon />
                {!sidebarCollapsed && <span className={styles.navLabel}>Организации</span>}
              </button>
            </nav>
          ) : (
            <nav className={styles.nav}>
              <div style={{ height: '200px' }}></div>
            </nav>
          )}

          <button 
            className={styles.collapseBtn}
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            title={sidebarCollapsed ? 'Развернуть меню' : 'Свернуть меню'}
          >
            <ChevronLeftIcon className={`${styles.collapseIcon} ${sidebarCollapsed ? styles.collapseIconRotated : ''}`} />
            {!sidebarCollapsed && <span className={styles.collapseLabel}>Свернуть</span>}
          </button>
        </aside>

        <main className={styles.main}>
          {error && <div className={styles.error}>{error}</div>}
          
          {!isClient || loading ? (
            <div className={styles.loading}>Загрузка...</div>
          ) : (
            <>
              {activeTab === 'users' && (
                <div className="space-y-8">
                  <div>
                    <h2 className="text-3xl font-bold text-gray-900 mb-2">Управление пользователями</h2>
                    <p className="text-gray-600">Просмотр и управление учетными записями</p>
                  </div>

                  <TableWidget
                    title={`Пользователи (${users.length})`}
                    actions={
                      <>
                        <div className="relative">
                          <MagnifyingGlassIcon className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                          <input
                            type="text"
                            placeholder="Поиск..."
                            className="pl-9 pr-4 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                        <button className="px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors">
                          Создать пользователя
                        </button>
                      </>
                    }
                  >
                    <div className="overflow-x-auto">
                      <table className="w-full">
                        <thead>
                          <tr className="border-b border-gray-100">
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Пользователь</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Email</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Посты</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Питомцы</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата регистрации</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Действия</th>
                          </tr>
                        </thead>
                        <tbody>
                          {users.map(user => (
                            <tr key={user.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                              <td className="py-3 px-4">
                                <div className="flex items-center gap-3">
                                  <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
                                    <span className="text-sm font-medium text-blue-600">
                                      {user.name.charAt(0).toUpperCase()}
                                    </span>
                                  </div>
                                  <span className="text-sm font-medium text-gray-900">{user.name}</span>
                                </div>
                              </td>
                              <td className="py-3 px-4 text-sm text-gray-600">{user.email}</td>
                              <td className="py-3 px-4 text-sm text-gray-600">0</td>
                              <td className="py-3 px-4 text-sm text-gray-600">0</td>
                              <td className="py-3 px-4 text-sm text-gray-600">
                                {new Date(user.created_at).toLocaleDateString('ru-RU')}
                              </td>
                              <td className="py-3 px-4">
                                <button
                                  onClick={() => handleDeleteUser(user.id)}
                                  className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                  title="Удалить"
                                >
                                  <TrashIcon className="w-4 h-4" />
                                </button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </TableWidget>
                </div>
              )}

              {activeTab === 'posts' && (
                <div className="space-y-8">
                  <div>
                    <h2 className="text-3xl font-bold text-gray-900 mb-2">Модерация постов</h2>
                    <p className="text-gray-600">Проверка и управление публикациями</p>
                  </div>

                  <TableWidget
                    title={`Посты (${posts.length})`}
                    actions={
                      <>
                        <select className="px-4 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                          <option>Все статусы</option>
                          <option>На модерации</option>
                          <option>Одобрено</option>
                          <option>Отклонено</option>
                        </select>
                        <div className="relative">
                          <MagnifyingGlassIcon className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                          <input
                            type="text"
                            placeholder="Поиск..."
                            className="pl-9 pr-4 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                      </>
                    }
                  >
                    <div className="overflow-x-auto">
                      <table className="w-full">
                        <thead>
                          <tr className="border-b border-gray-100">
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">ID</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Содержание</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Автор</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Статус</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата</th>
                            <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Действия</th>
                          </tr>
                        </thead>
                        <tbody>
                          {posts.map(post => (
                            <tr key={post.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                              <td className="py-3 px-4 text-sm text-gray-600">#{post.id}</td>
                              <td className="py-3 px-4 text-sm text-gray-900 max-w-md truncate">
                                {post.content}
                              </td>
                              <td className="py-3 px-4 text-sm text-gray-600">User #{post.user_id}</td>
                              <td className="py-3 px-4">
                                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                                  post.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                  post.status === 'approved' ? 'bg-green-100 text-green-800' :
                                  'bg-red-100 text-red-800'
                                }`}>
                                  {post.status === 'pending' ? 'На модерации' :
                                   post.status === 'approved' ? 'Одобрен' : 'Отклонен'}
                                </span>
                              </td>
                              <td className="py-3 px-4 text-sm text-gray-600">
                                {new Date(post.created_at).toLocaleDateString('ru-RU')}
                              </td>
                              <td className="py-3 px-4">
                                <div className="flex items-center gap-2">
                                  {post.status === 'pending' && (
                                    <>
                                      <button
                                        onClick={() => handleModeratePost(post.id, 'approved')}
                                        className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors"
                                        title="Одобрить"
                                      >
                                        <CheckIcon className="w-4 h-4" />
                                      </button>
                                      <button
                                        onClick={() => handleModeratePost(post.id, 'rejected')}
                                        className="p-2 text-orange-600 hover:bg-orange-50 rounded-lg transition-colors"
                                        title="Отклонить"
                                      >
                                        <XMarkIcon className="w-4 h-4" />
                                      </button>
                                    </>
                                  )}
                                  <button
                                    onClick={() => handleDeletePost(post.id)}
                                    className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                    title="Удалить"
                                  >
                                    <TrashIcon className="w-4 h-4" />
                                  </button>
                                </div>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </TableWidget>
                </div>
              )}

              {activeTab === 'stats' && stats && (
                <div className="space-y-8">
                  <div>
                    <h2 className="text-3xl font-bold text-gray-900 mb-2">Обзор платформы</h2>
                    <p className="text-gray-600">Основные метрики и статистика</p>
                  </div>
                  <div>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                      <StatsWidget
                        title="Всего пользователей"
                        value={stats.total_users}
                        icon={<UsersIcon className="w-6 h-6" />}
                        color="blue"
                        trend={{ value: '+12% за месяц', isPositive: true }}
                      />
                      <StatsWidget
                        title="Всего постов"
                        value={stats.total_posts}
                        icon={<DocumentTextIcon className="w-6 h-6" />}
                        color="green"
                        trend={{ value: '+8% за неделю', isPositive: true }}
                      />
                      <StatsWidget
                        title="На модерации"
                        value={stats.pending_posts}
                        icon={<ClipboardDocumentListIcon className="w-6 h-6" />}
                        color="orange"
                      />
                      <StatsWidget
                        title="Активных сегодня"
                        value={stats.active_users_today}
                        icon={<ChartBarIcon className="w-6 h-6" />}
                        color="purple"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <ChartWidget title="Активность пользователей" subtitle="За последние 7 дней">
                      <div className="h-64 flex items-center justify-center text-gray-400">
                        График активности (в разработке)
                      </div>
                    </ChartWidget>
                    
                    <ChartWidget title="Новые регистрации" subtitle="За последний месяц">
                      <div className="h-64 flex items-center justify-center text-gray-400">
                        График регистраций (в разработке)
                      </div>
                    </ChartWidget>
                  </div>

                  <TableWidget 
                    title="Последние действия"
                    actions={
                      <button className="text-sm text-blue-600 hover:text-blue-700">
                        Посмотреть все
                      </button>
                    }
                  >
                    <div className="text-sm text-gray-600">
                      <div className="space-y-3">
                        <div className="flex items-center justify-between py-2 border-b border-gray-100">
                          <span>Новый пользователь зарегистрирован</span>
                          <span className="text-gray-400">2 мин назад</span>
                        </div>
                        <div className="flex items-center justify-between py-2 border-b border-gray-100">
                          <span>Пост отправлен на модерацию</span>
                          <span className="text-gray-400">15 мин назад</span>
                        </div>
                        <div className="flex items-center justify-between py-2">
                          <span>Пост одобрен модератором</span>
                          <span className="text-gray-400">1 час назад</span>
                        </div>
                      </div>
                    </div>
                  </TableWidget>
                </div>
              )}

              {activeTab === 'logs' && (
                <div className="space-y-8">
                  <div>
                    <h2 className="text-3xl font-bold text-gray-900 mb-2">Логи системы</h2>
                    <p className="text-gray-600">История действий администраторов</p>
                  </div>

                  <TableWidget
                    title={`Логи (${logs.length})`}
                    actions={
                      <>
                        <select className="px-4 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                          <option>Все действия</option>
                          <option>Создание</option>
                          <option>Обновление</option>
                          <option>Удаление</option>
                          <option>Ошибки</option>
                        </select>
                        <div className="relative">
                          <MagnifyingGlassIcon className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                          <input
                            type="text"
                            placeholder="Поиск..."
                            className="pl-9 pr-4 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                      </>
                    }
                  >
                    {logs.length === 0 ? (
                      <div className="text-center py-12 text-gray-400">
                        <DocumentDuplicateIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
                        <p>Логи отсутствуют</p>
                      </div>
                    ) : (
                      <div className="overflow-x-auto">
                        <table className="w-full">
                          <thead>
                            <tr className="border-b border-gray-100">
                              <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">ID</th>
                              <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Администратор</th>
                              <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Действие</th>
                              <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Детали</th>
                              <th className="text-left py-3 px-4 text-sm font-semibold text-gray-700">Дата и время</th>
                            </tr>
                          </thead>
                          <tbody>
                            {logs.map(log => (
                              <tr key={log.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                                <td className="py-3 px-4 text-sm text-gray-600">#{log.id}</td>
                                <td className="py-3 px-4 text-sm text-gray-900">{log.admin_name}</td>
                                <td className="py-3 px-4">
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                    {log.action}
                                  </span>
                                </td>
                                <td className="py-3 px-4 text-sm text-gray-600">{log.details}</td>
                                <td className="py-3 px-4 text-sm text-gray-600">
                                  {new Date(log.created_at).toLocaleString('ru-RU')}
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    )}
                  </TableWidget>
                </div>
              )}
            </>
          )}
        </main>
      </div>
    </div>
  );
}
