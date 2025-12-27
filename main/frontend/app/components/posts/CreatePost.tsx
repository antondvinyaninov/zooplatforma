'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient } from '@/lib/api';
import {
  PhotoIcon,
  MapPinIcon,
  XMarkIcon,
  DocumentDuplicateIcon,
  ClockIcon,
  EllipsisHorizontalIcon,
  UserIcon,
  ListBulletIcon,
  VideoCameraIcon,
} from '@heroicons/react/24/outline';
import { GiPawPrint } from 'react-icons/gi';
import { MdPets } from 'react-icons/md';
import PollCreator, { PollData } from '../polls/PollCreator';

interface CreatePostProps {
  onPostCreated?: () => void;
}

type ReplySettingType = 'anyone' | 'followers' | 'following' | 'mentions';

export default function CreatePost({ onPostCreated }: CreatePostProps) {
  const { user } = useAuth();
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [showPollCreator, setShowPollCreator] = useState(false);
  const [pollData, setPollData] = useState<PollData | null>(null);
  const [showReplySettings, setShowReplySettings] = useState(false);
  const [replySetting, setReplySetting] = useState<ReplySettingType>('anyone');
  const [verifyReplies, setVerifyReplies] = useState(false);
  const [showDrafts, setShowDrafts] = useState(false);
  const [showScheduleModal, setShowScheduleModal] = useState(false);
  const [showSaveDraftDialog, setShowSaveDraftDialog] = useState(false);
  const [scheduledDate, setScheduledDate] = useState<Date | null>(null);
  const [scheduledTime, setScheduledTime] = useState('19:00');
  const [currentMonth, setCurrentMonth] = useState(new Date());
  const [drafts, setDrafts] = useState<any[]>([]);
  const [draftMenuOpen, setDraftMenuOpen] = useState<number | null>(null);

  const handleSubmit = async () => {
    if (!content.trim() && !pollData) return;

    setIsSubmitting(true);
    try {
      const postData: any = {
        content,
        attached_pets: [],
        attachments: [],
        tags: [],
      };

      // Add poll if exists
      if (pollData) {
        postData.poll = pollData;
      }

      // If scheduled, add scheduled_at field
      if (scheduledDate && scheduledTime) {
        const [hours, minutes] = scheduledTime.split(':');
        const scheduledDateTime = new Date(scheduledDate);
        scheduledDateTime.setHours(parseInt(hours), parseInt(minutes), 0, 0);
        postData.scheduled_at = scheduledDateTime.toISOString();
        postData.status = 'scheduled'; // Mark as scheduled
      } else {
        postData.status = 'published'; // Publish immediately
      }

      await apiClient.post('/api/posts', postData);

      setContent('');
      setPollData(null);
      setShowPollCreator(false);
      setScheduledDate(null);
      setScheduledTime('19:00');
      setShowModal(false);
      onPostCreated?.();
    } catch (error) {
      console.error('Ошибка создания поста:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Close modal on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && showModal) {
        setShowModal(false);
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [showModal]);

  // Load drafts when drafts modal opens
  useEffect(() => {
    if (showDrafts) {
      loadDrafts();
      setDraftMenuOpen(null); // Reset menu state
    }
  }, [showDrafts]);

  const loadDrafts = async () => {
    try {
      const response = await apiClient.get('/api/posts/drafts');
      setDrafts(response.data || []);
    } catch (error) {
      console.error('Ошибка загрузки черновиков:', error);
    }
  };

  const saveDraft = async () => {
    if (!content.trim()) return;

    try {
      await apiClient.post('/api/posts', {
        content,
        attached_pets: [],
        attachments: [],
        tags: [],
        status: 'draft',
      });
      setContent('');
      setShowSaveDraftDialog(false);
      setShowModal(false);
      setShowDrafts(true); // Open drafts after saving
    } catch (error) {
      console.error('Ошибка сохранения черновика:', error);
    }
  };

  const loadDraft = (draft: any) => {
    setContent(draft.content);
    setShowDrafts(false);
    // TODO: Delete draft from server after loading
  };

  const deleteDraft = async (draftId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await apiClient.delete(`/api/posts/${draftId}`);
      setDrafts(drafts.filter(d => d.id !== draftId));
      setDraftMenuOpen(null);
    } catch (error) {
      console.error('Ошибка удаления черновика:', error);
    }
  };

  // Close modal on backdrop click
  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget) {
      setShowModal(false);
    }
  };

  const replySettingLabels: Record<ReplySettingType, string> = {
    anyone: 'Кто угодно',
    followers: 'Ваши подписчики',
    following: 'Профили, на которые вы подписаны',
    mentions: 'Упомянутые профили',
  };

  // Calendar helpers
  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay();
    
    return { daysInMonth, startingDayOfWeek, year, month };
  };

  const monthNames = [
    'январь', 'февраль', 'март', 'апрель', 'май', 'июнь',
    'июль', 'август', 'сентябрь', 'октябрь', 'ноябрь', 'декабрь'
  ];

  const weekDays = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'];

  const goToPreviousMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1));
  };

  const goToNextMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1));
  };

  const isToday = (day: number, month: number, year: number) => {
    const today = new Date();
    return day === today.getDate() && 
           month === today.getMonth() && 
           year === today.getFullYear();
  };

  const isSelectedDate = (day: number, month: number, year: number) => {
    if (!scheduledDate) return false;
    return day === scheduledDate.getDate() && 
           month === scheduledDate.getMonth() && 
           year === scheduledDate.getFullYear();
  };

  return (
    <>
      {/* Trigger Button */}
      <div className="flex items-start gap-3 p-4">
        {/* Avatar */}
        <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-white font-semibold flex-shrink-0 overflow-hidden">
          {user?.avatar ? (
            <img src={user.avatar} alt={user.name} className="w-full h-full object-cover" />
          ) : (
            <UserIcon className="w-6 h-6 text-gray-500" />
          )}
        </div>

        {/* Input Trigger */}
        <button
          onClick={() => setShowModal(true)}
          className="flex-1 text-left text-gray-500 text-[15px] py-2 hover:text-gray-700 transition-colors"
        >
          Что нового?
        </button>
      </div>

      {/* Modal */}
      {showModal && (
        <div 
          className="fixed inset-0 z-50 flex items-center justify-center bg-white/40 p-4 overflow-y-auto backdrop-blur-sm"
          onClick={handleBackdropClick}
        >
          <div className="bg-white rounded-2xl w-full max-w-[600px] shadow-2xl my-8">
            {/* Modal Header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200">
              <button
                onClick={() => {
                  if (content.trim()) {
                    setShowSaveDraftDialog(true);
                  } else {
                    setShowModal(false);
                  }
                }}
                className="text-[15px] text-gray-700 hover:text-gray-900 font-medium"
              >
                Отмена
              </button>
              <h3 className="font-bold text-[16px]">Новая метка</h3>
              <div className="flex items-center gap-1">
                {/* Drafts Button */}
                <button
                  onClick={() => {
                    if (content.trim()) {
                      setShowSaveDraftDialog(true);
                    } else {
                      setShowDrafts(true);
                    }
                  }}
                  className="p-1.5 hover:bg-gray-100 rounded-full transition-colors relative"
                  title="Черновики"
                >
                  <DocumentDuplicateIcon className="w-6 h-6 text-gray-600" strokeWidth={2} />
                </button>
                
                {/* Schedule Button */}
                <button
                  onClick={() => setShowScheduleModal(true)}
                  className="p-1.5 hover:bg-gray-100 rounded-full transition-colors"
                  title="Запланировать публикацию"
                >
                  <ClockIcon className="w-6 h-6 text-gray-600" strokeWidth={2} />
                </button>
              </div>
            </div>

            {/* Modal Content */}
            <div className="px-4 py-3">
              {/* User Info */}
              <div className="flex items-start gap-3 mb-2">
                {/* Avatar */}
                <div className="w-9 h-9 rounded-full bg-gray-300 flex items-center justify-center text-white font-semibold flex-shrink-0 text-[14px] overflow-hidden">
                  {user?.avatar ? (
                    <img src={user.avatar} alt={user.name} className="w-full h-full object-cover" />
                  ) : (
                    <UserIcon className="w-5 h-5 text-gray-500" />
                  )}
                </div>
                <div className="flex-1">
                  {/* Username */}
                  <div className="font-semibold text-[15px] mb-1">{user?.name || 'Пользователь'}</div>
                  
                  {/* Text Input - no border, just placeholder */}
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Что нового у вашего питомца?"
                    rows={1}
                    className="w-full text-[15px] border-none focus:outline-none resize-none text-gray-900 placeholder-gray-400 overflow-hidden"
                    autoFocus
                    style={{
                      minHeight: '24px',
                      maxHeight: '400px',
                      height: 'auto',
                    }}
                    onInput={(e) => {
                      const target = e.target as HTMLTextAreaElement;
                      target.style.height = 'auto';
                      target.style.height = target.scrollHeight + 'px';
                    }}
                  />
                </div>
              </div>

              {/* Attachment Icons */}
              <div className="flex items-center gap-0.5 mb-2 ml-12">
                <button
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  title="Добавить фото"
                >
                  <PhotoIcon className="w-5 h-5 text-gray-400" strokeWidth={2} />
                </button>
                <button
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  title="Добавить видео"
                >
                  <VideoCameraIcon className="w-5 h-5 text-gray-400" strokeWidth={2} />
                </button>
                <button
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  title="Прикрепить карточку питомца"
                >
                  <MdPets className="w-5 h-5 text-gray-400" />
                </button>
                <button
                  onClick={() => setShowPollCreator(!showPollCreator)}
                  className={`p-2 hover:bg-gray-100 rounded-lg transition-colors ${showPollCreator || pollData ? 'bg-blue-50' : ''}`}
                  title="Голосование"
                >
                  <ListBulletIcon className={`w-5 h-5 ${showPollCreator || pollData ? 'text-blue-600' : 'text-gray-400'}`} strokeWidth={2} />
                </button>
                <button
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  title="Местоположение"
                >
                  <MapPinIcon className="w-5 h-5 text-gray-400" strokeWidth={2} />
                </button>
              </div>

              {/* Poll Creator */}
              {showPollCreator && !pollData && (
                <PollCreator
                  onPollChange={(poll) => {
                    setPollData(poll);
                    if (poll) {
                      setShowPollCreator(false);
                    }
                  }}
                  onClose={() => setShowPollCreator(false)}
                />
              )}

              {/* Poll Preview */}
              {pollData && !showPollCreator && (
                <div className="mb-3 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                  <div className="flex items-start justify-between mb-2">
                    <div className="flex-1">
                      <div className="font-semibold text-[14px] text-gray-900 mb-2">{pollData.question}</div>
                      <div className="space-y-1">
                        {pollData.options.map((option, index) => (
                          <div key={index} className="text-[13px] text-gray-700 px-2 py-1 bg-white rounded">
                            {option}
                          </div>
                        ))}
                      </div>
                      <div className="text-[12px] text-gray-500 mt-2">
                        {pollData.multiple_choice ? 'Множественный выбор' : 'Один вариант'} • 
                        {pollData.allow_vote_changes ? ' Можно изменить' : ' Нельзя изменить'}
                        {pollData.anonymous_voting && ' • Анонимно'}
                        {pollData.expires_at && ' • Ограничен по времени'}
                      </div>
                    </div>
                    <button
                      onClick={() => {
                        setPollData(null);
                        setShowPollCreator(false);
                      }}
                      className="ml-2 p-1 hover:bg-blue-100 rounded-full transition-colors"
                    >
                      <XMarkIcon className="w-4 h-4 text-blue-700" strokeWidth={2} />
                    </button>
                  </div>
                </div>
              )}

              {/* Bottom Section */}
              <div className="pt-3 border-t border-gray-200">
                {/* Scheduled Info */}
                {scheduledDate && (
                  <div className="mb-3 p-2.5 bg-blue-50 border border-blue-200 rounded-lg flex items-center justify-between">
                    <div className="flex items-center gap-2 text-[13px] text-blue-700">
                      <ClockIcon className="w-4 h-4" strokeWidth={2} />
                      <span>
                        Запланировано: {scheduledDate.toLocaleDateString('ru-RU')} в {scheduledTime}
                      </span>
                    </div>
                    <button
                      onClick={() => {
                        setScheduledDate(null);
                        setScheduledTime('19:00');
                      }}
                      className="text-blue-700 hover:text-blue-900"
                    >
                      <XMarkIcon className="w-4 h-4" strokeWidth={2} />
                    </button>
                  </div>
                )}

                {/* Reply Settings Dropdown */}
                {showReplySettings && (
                  <div className="mb-4 bg-white border border-gray-200 rounded-xl shadow-lg overflow-hidden">
                    <div className="p-4">
                      <h4 className="text-[13px] text-gray-500 mb-3">Кто может отвечать и цитировать</h4>
                      
                      {/* Options */}
                      <div className="space-y-0">
                        {(['anyone', 'followers', 'following', 'mentions'] as ReplySettingType[]).map((option) => (
                          <button
                            key={option}
                            onClick={() => setReplySetting(option)}
                            className="w-full flex items-center justify-between py-3 px-2 hover:bg-gray-50 rounded-lg transition-colors"
                          >
                            <span className="text-[15px] text-gray-900">{replySettingLabels[option]}</span>
                            {replySetting === option && (
                              <svg className="w-5 h-5 text-black" fill="currentColor" viewBox="0 0 20 20">
                                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                              </svg>
                            )}
                          </button>
                        ))}
                      </div>
                    </div>

                    {/* Verify Toggle */}
                    <div className="border-t border-gray-200 p-4">
                      <button
                        onClick={() => setVerifyReplies(!verifyReplies)}
                        className="w-full flex items-center justify-between"
                      >
                        <span className="text-[15px] text-gray-900">Проверять и одобрять ответы</span>
                        <div className={`relative w-11 h-6 rounded-full transition-colors ${verifyReplies ? 'bg-black' : 'bg-gray-300'}`}>
                          <div className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform ${verifyReplies ? 'translate-x-5' : 'translate-x-0'}`} />
                        </div>
                      </button>
                    </div>
                  </div>
                )}

                {/* Action Buttons */}
                <div className="flex items-center justify-between">
                  <button 
                    onClick={() => setShowReplySettings(!showReplySettings)}
                    className="text-[13px] text-gray-500 hover:text-gray-900 flex items-center gap-1.5"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                    </svg>
                    Настройки ответов
                  </button>

                  <button
                    onClick={handleSubmit}
                    disabled={(!content.trim() && !pollData) || isSubmitting}
                    className="px-5 py-1.5 bg-black text-white rounded-full text-[14px] font-semibold hover:bg-gray-800 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                  >
                    {isSubmitting ? 'Публикация...' : scheduledDate ? 'Запланировать' : 'Поставить метку'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Drafts Modal */}
      {showDrafts && (
        <div 
          className="fixed inset-0 z-50 flex items-center justify-center bg-white/40 p-4 overflow-y-auto backdrop-blur-sm"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setShowDrafts(false);
              setDraftMenuOpen(null);
            }
          }}
        >
          <div className="bg-white rounded-2xl w-full max-w-[600px] shadow-2xl my-8">
            <div className="flex items-center justify-between p-4 border-b border-gray-200">
              <button
                onClick={() => {
                  setShowDrafts(false);
                  setDraftMenuOpen(null);
                }}
                className="p-1 hover:bg-gray-100 rounded-full transition-colors"
              >
                <XMarkIcon className="w-6 h-6 text-gray-600" strokeWidth={2} />
              </button>
              <h3 className="font-bold text-[16px]">Черновики</h3>
              <div className="w-8"></div>
            </div>
            <div className="max-h-[500px] overflow-y-auto">
              {drafts.length === 0 ? (
                <div className="p-8 text-center text-gray-500">
                  <p>Черновиков пока нет</p>
                </div>
              ) : (
                <div>
                  {drafts.map((draft) => (
                    <div
                      key={draft.id}
                      className="relative border-b border-gray-200 hover:bg-gray-50 transition-colors"
                    >
                      <button
                        onClick={() => loadDraft(draft)}
                        className="w-full p-4 text-left"
                      >
                        <div className="flex items-start gap-3">
                          <div className="w-9 h-9 rounded-full bg-gradient-to-br from-purple-400 to-pink-400 flex items-center justify-center text-white font-semibold flex-shrink-0 text-[14px]">
                            {user?.name?.[0]?.toUpperCase() || 'U'}
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="text-[15px] text-gray-900 line-clamp-3">{draft.content}</div>
                            <div className="text-[13px] text-gray-500 mt-1">
                              {new Date(draft.created_at).toLocaleDateString('ru-RU')}
                            </div>
                          </div>
                        </div>
                      </button>
                      
                      {/* Three dots menu */}
                      <div className="absolute top-4 right-4">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setDraftMenuOpen(draftMenuOpen === draft.id ? null : draft.id);
                          }}
                          className="p-1 hover:bg-gray-200 rounded-full transition-colors"
                        >
                          <EllipsisHorizontalIcon className="w-5 h-5 text-gray-600" strokeWidth={2} />
                        </button>
                        
                        {/* Dropdown menu */}
                        {draftMenuOpen === draft.id && (
                          <div className="absolute right-0 top-full mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg overflow-hidden z-10">
                            <button
                              onClick={(e) => deleteDraft(draft.id, e)}
                              className="w-full px-4 py-3 text-left text-[15px] text-red-600 hover:bg-gray-50 transition-colors"
                            >
                              Удалить черновик
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Save Draft Dialog */}
      {showSaveDraftDialog && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-white/40 p-4 backdrop-blur-sm">
          <div className="bg-white rounded-2xl w-full max-w-[400px] shadow-2xl overflow-hidden">
            <div className="p-6 text-center">
              <h3 className="font-bold text-[18px] mb-2">Сохранить как черновик?</h3>
              <p className="text-[14px] text-gray-500">
                Сохранить как черновик, чтобы отредактировать и отправить позже.
              </p>
            </div>
            <div className="border-t border-gray-200">
              <button
                onClick={saveDraft}
                className="w-full py-3 text-[15px] font-semibold hover:bg-gray-50 transition-colors"
              >
                Сохранить
              </button>
            </div>
            <div className="border-t border-gray-200">
              <button
                onClick={() => {
                  setShowSaveDraftDialog(false);
                  setShowModal(false);
                  setContent('');
                }}
                className="w-full py-3 text-[15px] text-red-500 font-semibold hover:bg-gray-50 transition-colors"
              >
                Не сохранять
              </button>
            </div>
            <div className="border-t border-gray-200">
              <button
                onClick={() => setShowSaveDraftDialog(false)}
                className="w-full py-3 text-[15px] hover:bg-gray-50 transition-colors"
              >
                Отмена
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Schedule Modal */}
      {showScheduleModal && (
        <div 
          className="fixed inset-0 z-[60] flex items-center justify-center bg-white/40 p-4 backdrop-blur-sm"
          onClick={(e) => e.target === e.currentTarget && setShowScheduleModal(false)}
        >
          <div className="bg-white rounded-2xl w-full max-w-[500px] shadow-2xl">
            <div className="flex items-center justify-between p-4 border-b border-gray-200">
              <button
                onClick={() => setShowScheduleModal(false)}
                className="p-1 hover:bg-gray-100 rounded-full transition-colors"
              >
                <XMarkIcon className="w-6 h-6 text-gray-600" strokeWidth={2} />
              </button>
              <h3 className="font-bold text-[16px]">Запланировать публикацию</h3>
              <div className="w-8"></div>
            </div>
            
            <div className="p-6">
              {/* Calendar Header */}
              <div className="flex items-center justify-between mb-4">
                <h4 className="text-[18px] font-bold">
                  {monthNames[currentMonth.getMonth()]} {currentMonth.getFullYear()} г.
                </h4>
                <div className="flex gap-2">
                  <button
                    onClick={goToPreviousMonth}
                    className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                  </button>
                  <button
                    onClick={goToNextMonth}
                    className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                </div>
              </div>

              {/* Week Days */}
              <div className="grid grid-cols-7 gap-2 mb-2">
                {weekDays.map((day) => (
                  <div key={day} className="text-center text-[13px] text-gray-400 font-medium py-2">
                    {day}
                  </div>
                ))}
              </div>

              {/* Calendar Days */}
              <div className="grid grid-cols-7 gap-2">
                {(() => {
                  const { daysInMonth, startingDayOfWeek, year, month } = getDaysInMonth(currentMonth);
                  const days = [];
                  
                  // Empty cells before first day
                  for (let i = 0; i < startingDayOfWeek; i++) {
                    days.push(<div key={`empty-${i}`} className="aspect-square" />);
                  }
                  
                  // Days of month
                  for (let day = 1; day <= daysInMonth; day++) {
                    const today = isToday(day, month, year);
                    const selected = isSelectedDate(day, month, year);
                    
                    days.push(
                      <button
                        key={day}
                        onClick={() => setScheduledDate(new Date(year, month, day))}
                        className={`aspect-square rounded-full flex items-center justify-center text-[15px] transition-colors
                          ${selected ? 'bg-black text-white font-bold' : 
                            today ? 'bg-gray-200 font-semibold' : 
                            'hover:bg-gray-100'}`}
                      >
                        {day}
                      </button>
                    );
                  }
                  
                  return days;
                })()}
              </div>

              {/* Time Picker */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <div className="flex items-center gap-3">
                  <ClockIcon className="w-5 h-5 text-gray-400" strokeWidth={2} />
                  <input
                    type="time"
                    value={scheduledTime}
                    onChange={(e) => setScheduledTime(e.target.value)}
                    className="flex-1 px-4 py-2 text-[15px] bg-gray-50 rounded-lg border border-gray-200 focus:outline-none focus:border-gray-300"
                  />
                </div>
              </div>

              {/* Action Button */}
              <button
                onClick={() => {
                  if (scheduledDate) {
                    setShowScheduleModal(false);
                    // Date and time are now set, user can publish
                  }
                }}
                disabled={!scheduledDate}
                className="w-full mt-6 py-3 bg-black text-white rounded-full text-[15px] font-semibold hover:bg-gray-800 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
              >
                Готово
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
