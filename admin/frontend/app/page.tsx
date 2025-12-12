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
} from '@heroicons/react/24/outline';
import styles from './page.module.css';

export default function AdminDashboard() {
  const router = useRouter();
  const [stats, setStats] = useState<Stats | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [posts, setPosts] = useState<Post[]>([]);
  const [logs, setLogs] = useState<Log[]>([]);
  const [activeTab, setActiveTab] = useState<'users' | 'posts' | 'stats' | 'logs'>('stats');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [adminUser, setAdminUser] = useState<{ email: string; role: string } | null>(null);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  useEffect(() => {
    loadData();
  }, [activeTab]);

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
          <h1 className={styles.logo}>
            <span className={styles.logoIcon}>🐾</span>
            ЗооПлатформа
          </h1>
          <div className={styles.headerRight}>
            {adminUser && (
              <>
                <span className={styles.balance}>0,00 ₽</span>
                <button className={styles.iconBtn}>
                  <ClipboardDocumentListIcon className="w-4 h-4" />
                </button>
                <button className={styles.iconBtn}>
                  <BellIcon className="w-4 h-4" />
                </button>
                <button className={styles.iconBtn}>
                  <Cog6ToothIcon className="w-4 h-4" />
                </button>
                <div className={styles.userMenu}>
                  <span className={styles.userName}>{adminUser.email}</span>
                  <button className={styles.logoutBtn} onClick={handleLogout}>
                    ▼
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </header>

      <div className={styles.layout}>
        <aside className={`${styles.sidebar} ${sidebarCollapsed ? styles.sidebarCollapsed : ''}`}>
          <nav className={styles.nav}>
            <button
              className={`${styles.navItem} ${activeTab === 'users' ? styles.navItemActive : ''}`}
              onClick={() => setActiveTab('users')}
              title="Пользователи"
            >
              <UsersIcon />
              {!sidebarCollapsed && <span className={styles.navLabel}>Пользователи</span>}
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'posts' ? styles.navItemActive : ''}`}
              onClick={() => setActiveTab('posts')}
              title="Посты"
            >
              <DocumentTextIcon />
              {!sidebarCollapsed && <span className={styles.navLabel}>Посты</span>}
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'stats' ? styles.navItemActive : ''}`}
              onClick={() => setActiveTab('stats')}
              title="Статистика"
            >
              <ChartBarIcon />
              {!sidebarCollapsed && <span className={styles.navLabel}>Статистика</span>}
            </button>
            <button
              className={`${styles.navItem} ${activeTab === 'logs' ? styles.navItemActive : ''}`}
              onClick={() => setActiveTab('logs')}
              title="Логирование"
            >
              <DocumentDuplicateIcon />
              {!sidebarCollapsed && <span className={styles.navLabel}>Логирование</span>}
            </button>
          </nav>

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
          
          {loading ? (
            <div className={styles.loading}>Загрузка...</div>
          ) : (
            <>
              {activeTab === 'users' && (
                <>
                  <div className={styles.toolbar}>
                    <button className={styles.createBtn}>Создать пользователя</button>
                    <div className={styles.filters}>
                      <div className={styles.searchWrapper}>
                        <MagnifyingGlassIcon className="w-3.5 h-3.5 text-gray-400" />
                        <input 
                          type="text" 
                          placeholder="Поиск" 
                          className={styles.searchInput}
                        />
                      </div>
                    </div>
                  </div>
                  
                  <div className={styles.tableWrapper}>
                    <table className={styles.dataTable}>
                      <thead>
                        <tr>
                          <th>Пользователи ▼</th>
                          <th>Email</th>
                          <th>Посты</th>
                          <th>Питомцы</th>
                          <th>Дата регистрации</th>
                          <th>Действия</th>
                        </tr>
                      </thead>
                      <tbody>
                        {users.map(user => (
                          <tr key={user.id}>
                            <td>
                              <div className={styles.userCell}>
                                <span className={styles.userName}>{user.name}</span>
                              </div>
                            </td>
                            <td>{user.email}</td>
                            <td>0</td>
                            <td>0</td>
                            <td>{new Date(user.created_at).toLocaleDateString('ru-RU')}</td>
                            <td>
                              <button
                                className={styles.deleteBtn}
                                onClick={() => handleDeleteUser(user.id)}
                              >
                                <TrashIcon className="w-3.5 h-3.5" />
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    <div className={styles.tableFooter}>
                      Итого: {users.length} пользователей
                    </div>
                  </div>
                </>
              )}

              {activeTab === 'posts' && (
                <>
                  <div className={styles.toolbar}>
                    <div className={styles.filters}>
                      <select className={styles.filterSelect}>
                        <option>Все статусы</option>
                        <option>На модерации</option>
                        <option>Одобрено</option>
                        <option>Отклонено</option>
                      </select>
                      <div className={styles.searchWrapper}>
                        <MagnifyingGlassIcon className="w-3.5 h-3.5 text-gray-400" />
                        <input 
                          type="text" 
                          placeholder="Поиск" 
                          className={styles.searchInput}
                        />
                      </div>
                    </div>
                  </div>
                  
                  <div className={styles.tableWrapper}>
                    <table className={styles.dataTable}>
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>Содержание</th>
                          <th>Автор</th>
                          <th>Статус</th>
                          <th>Дата</th>
                          <th>Действия</th>
                        </tr>
                      </thead>
                      <tbody>
                        {posts.map(post => (
                          <tr key={post.id}>
                            <td>{post.id}</td>
                            <td>{post.content.substring(0, 60)}...</td>
                            <td>User #{post.user_id}</td>
                            <td>
                              <span className={`${styles.statusBadge} ${styles[`status${post.status}`]}`}>
                                {post.status === 'pending' ? 'На модерации' : 
                                 post.status === 'approved' ? 'Одобрен' : 'Отклонен'}
                              </span>
                            </td>
                            <td>{new Date(post.created_at).toLocaleDateString('ru-RU')}</td>
                            <td>
                              <div className={styles.actionBtns}>
                                {post.status === 'pending' && (
                                  <>
                                    <button
                                      className={styles.approveBtn}
                                      onClick={() => handleModeratePost(post.id, 'approved')}
                                      title="Одобрить"
                                    >
                                      <CheckIcon className="w-3.5 h-3.5" />
                                    </button>
                                    <button
                                      className={styles.rejectBtn}
                                      onClick={() => handleModeratePost(post.id, 'rejected')}
                                      title="Отклонить"
                                    >
                                      <XMarkIcon className="w-3.5 h-3.5" />
                                    </button>
                                  </>
                                )}
                                <button
                                  className={styles.deleteBtn}
                                  onClick={() => handleDeletePost(post.id)}
                                  title="Удалить"
                                >
                                  <TrashIcon className="w-3.5 h-3.5" />
                                </button>
                              </div>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    <div className={styles.tableFooter}>
                      Итого: {posts.length} постов
                    </div>
                  </div>
                </>
              )}

              {activeTab === 'stats' && stats && (
                <div className={styles.statsContainer}>
                  <h2 className={styles.statsTitle}>Статистика платформы</h2>
                  <div className={styles.statsGrid}>
                    <div className={styles.statCard}>
                      <div className={styles.statValue}>{stats.total_users}</div>
                      <div className={styles.statLabel}>Всего пользователей</div>
                    </div>
                    <div className={styles.statCard}>
                      <div className={styles.statValue}>{stats.total_posts}</div>
                      <div className={styles.statLabel}>Всего постов</div>
                    </div>
                    <div className={styles.statCard}>
                      <div className={styles.statValue}>{stats.pending_posts}</div>
                      <div className={styles.statLabel}>На модерации</div>
                    </div>
                    <div className={styles.statCard}>
                      <div className={styles.statValue}>{stats.active_users_today}</div>
                      <div className={styles.statLabel}>Активных сегодня</div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'logs' && (
                <>
                  <div className={styles.toolbar}>
                    <h2 className={styles.pageTitle}>Логи системы</h2>
                    <div className={styles.filters}>
                      <select className={styles.filterSelect}>
                        <option>Все действия</option>
                        <option>Создание</option>
                        <option>Обновление</option>
                        <option>Удаление</option>
                        <option>Ошибки</option>
                      </select>
                      <div className={styles.searchWrapper}>
                        <MagnifyingGlassIcon className="w-3.5 h-3.5 text-gray-400" />
                        <input 
                          type="text" 
                          placeholder="Поиск" 
                          className={styles.searchInput}
                        />
                      </div>
                    </div>
                  </div>
                  
                  <div className={styles.tableWrapper}>
                    <table className={styles.dataTable}>
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>Администратор</th>
                          <th>Действие</th>
                          <th>Детали</th>
                          <th>Дата и время</th>
                        </tr>
                      </thead>
                      <tbody>
                        {logs.length === 0 ? (
                          <tr>
                            <td colSpan={5} style={{ textAlign: 'center', padding: '40px', color: '#99a2ad' }}>
                              Логи отсутствуют
                            </td>
                          </tr>
                        ) : (
                          logs.map(log => (
                            <tr key={log.id}>
                              <td>{log.id}</td>
                              <td>{log.admin_name}</td>
                              <td>
                                <span className={`${styles.logAction} ${styles[`action${log.action.toLowerCase().replace(/\s/g, '')}`]}`}>
                                  {log.action}
                                </span>
                              </td>
                              <td>{log.details}</td>
                              <td>{new Date(log.created_at).toLocaleString('ru-RU')}</td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                    <div className={styles.tableFooter}>
                      Итого: {logs.length} записей
                    </div>
                  </div>
                </>
              )}
            </>
          )}
        </main>
      </div>
    </div>
  );
}
