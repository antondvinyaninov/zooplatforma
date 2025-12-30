'use client';

import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { useUnreadMessages } from '@/app/hooks/useUnreadMessages';
import {
  HomeIcon,
  UserIcon,
  ChatBubbleLeftIcon,
  PhoneIcon,
  UserGroupIcon,
  VideoCameraIcon,
  Cog6ToothIcon,
  ShoppingBagIcon,
  HeartIcon,
  NewspaperIcon,
  CalendarDaysIcon,
  BuildingOfficeIcon,
  AcademicCapIcon,
  DocumentTextIcon,
  ClipboardDocumentListIcon,
  InformationCircleIcon,
  ChartBarIcon,
  WrenchScrewdriverIcon,
  RectangleStackIcon,
} from '@heroicons/react/24/outline';

const additionalLinks = [
  { name: 'О платформе', href: '/about', icon: InformationCircleIcon },
  { name: 'Статистика', href: '/statistics', icon: ChartBarIcon },
  { name: 'Техподдержка', href: '/support', icon: WrenchScrewdriverIcon },
  { name: 'Команда', href: '/team', icon: UserGroupIcon },
];

export default function Sidebar() {
  const { user } = useAuth();
  const unreadCount = useUnreadMessages();

  const mainNavigation = [
    { name: 'Метки', href: '/', icon: DocumentTextIcon },
    { name: 'Профиль', href: user ? `/id${user.id}` : '/profile', icon: UserIcon },
    { name: 'Мессенджер', href: '/messenger', icon: ChatBubbleLeftIcon, badge: unreadCount > 0 ? unreadCount.toString() : undefined },
    { name: 'Организации', href: '/org', icon: BuildingOfficeIcon },
    { name: 'Каталог', href: '/catalog', icon: RectangleStackIcon },
    { name: 'Сервисы', href: '/services', icon: Cog6ToothIcon },
  ];
  return (
    <div className="sticky top-[48px]">
      <nav className="space-y-0">
        {mainNavigation.map((item) => {
          return (
            <div key={item.name}>
              {/* Разделитель перед Сервисами */}
              {item.name === 'Сервисы' && (
                <div className="border-t border-gray-300 my-3 mx-2"></div>
              )}

              <Link
                href={item.href}
                className="flex items-center space-x-2 px-2 py-1.5 rounded-lg transition-colors duration-200 group hover:bg-gray-200 ml-2"
              >
                <item.icon className="w-5 h-5 flex-shrink-0 text-gray-600" strokeWidth={2} />
                <span className="text-[13px] text-gray-700 font-medium flex-1 whitespace-nowrap overflow-hidden text-ellipsis">{item.name}</span>
                {item.badge && (
                  <span className="text-white text-[11px] px-2 py-0.5 rounded-full flex-shrink-0 font-semibold" style={{ backgroundColor: '#FC2B2B' }}>
                    {item.badge}
                  </span>
                )}
              </Link>
            </div>
          );
        })}
      </nav>

      {/* Разделитель перед дополнительными ссылками */}
      <div className="border-t border-gray-300 my-3 mx-2"></div>

      <nav className="space-y-0">
        {additionalLinks.map((link) => (
          <Link
            key={link.name}
            href={link.href}
            className="flex items-center space-x-2 px-2 py-1 rounded-lg transition-colors duration-200 group hover:bg-gray-200 ml-2"
          >
            <link.icon className="w-4 h-4 flex-shrink-0 text-gray-500" strokeWidth={2} />
            <span className="text-[12px] text-gray-600 font-medium whitespace-nowrap overflow-hidden text-ellipsis">{link.name}</span>
          </Link>
        ))}
      </nav>
    </div>
  );
}
