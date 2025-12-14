'use client';

import { useState } from 'react';
import OrganizationSearch from '../../components/shared/OrganizationSearch';
import YandexMap from '../../components/shared/YandexMap';

interface Organization {
  name: string;
  fullName?: string;
  inn: string;
  ogrn?: string;
  kpp?: string;
  address?: {
    full?: string;
    postalCode?: string;
    region?: string;
    city?: string;
    settlement?: string;
    street?: string;
    house?: string;
    flat?: string;
    geoLat?: string;
    geoLon?: string;
  };
  phones?: string[];
  emails?: string[];
  director?: {
    name?: string;
    post?: string;
  };
  registrationDate?: string;
  status?: string;
  opf?: string;
}

export default function OrganizationsPage() {
  const [selectedOrg, setSelectedOrg] = useState<Organization | null>(null);

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Организации</h1>
        <p className="text-gray-600">Поиск организаций по ИНН, ОГРН или названию</p>
      </div>

      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Найти организацию</h2>
        
        <OrganizationSearch
          placeholder="Поиск по ИНН, ОГРН или названию..."
          onSelect={(org) => {
            setSelectedOrg({
              name: org.data.name?.short_with_opf || org.value,
              fullName: org.data.name?.full_with_opf,
              inn: org.data.inn || '',
              ogrn: org.data.ogrn,
              kpp: org.data.kpp,
              address: {
                full: org.data.address?.unrestricted_value,
                postalCode: org.data.address?.data?.postal_code,
                region: org.data.address?.data?.region_with_type,
                city: org.data.address?.data?.city_with_type || org.data.address?.data?.settlement_with_type,
                street: org.data.address?.data?.street_with_type,
                house: org.data.address?.data?.house,
                flat: org.data.address?.data?.flat,
                geoLat: org.data.address?.data?.geo_lat,
                geoLon: org.data.address?.data?.geo_lon,
              },
              phones: org.data.phones?.map(p => p.value || '').filter(Boolean),
              emails: org.data.emails?.map(e => e.value || '').filter(Boolean),
              director: {
                name: org.data.management?.name,
                post: org.data.management?.post,
              },
              registrationDate: org.data.state?.registration_date,
              status: org.data.state?.status,
              opf: org.data.opf?.full,
            });
          }}
        />

        {selectedOrg && (
          <div className="mt-6 space-y-4">
            {/* Основная информация */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                🏢 Основная информация
              </h3>
              <div className="space-y-2 text-sm">
                <div>
                  <span className="text-gray-600">Название:</span>{' '}
                  <span className="font-medium text-gray-900">{selectedOrg.name}</span>
                </div>
                {selectedOrg.fullName && selectedOrg.fullName !== selectedOrg.name && (
                  <div>
                    <span className="text-gray-600">Полное название:</span>{' '}
                    <span className="text-gray-900">{selectedOrg.fullName}</span>
                  </div>
                )}
                {selectedOrg.opf && (
                  <div>
                    <span className="text-gray-600">Организационно-правовая форма:</span>{' '}
                    <span className="text-gray-900">{selectedOrg.opf}</span>
                  </div>
                )}
                {selectedOrg.status && (
                  <div>
                    <span className="text-gray-600">Статус:</span>{' '}
                    <span className={`font-medium ${selectedOrg.status === 'ACTIVE' ? 'text-green-600' : 'text-red-600'}`}>
                      {selectedOrg.status === 'ACTIVE' ? 'Действующая' : 'Ликвидирована'}
                    </span>
                  </div>
                )}
                {selectedOrg.registrationDate && (
                  <div>
                    <span className="text-gray-600">Дата регистрации:</span>{' '}
                    <span className="text-gray-900">
                      {new Date(parseInt(selectedOrg.registrationDate)).toLocaleDateString('ru-RU')}
                    </span>
                  </div>
                )}
              </div>
            </div>

            {/* Реквизиты */}
            <div className="p-4 bg-gray-50 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                📋 Реквизиты
              </h3>
              <div className="space-y-2 text-sm">
                <div>
                  <span className="text-gray-600">ИНН:</span>{' '}
                  <span className="font-mono text-gray-900">{selectedOrg.inn}</span>
                </div>
                {selectedOrg.ogrn && (
                  <div>
                    <span className="text-gray-600">ОГРН:</span>{' '}
                    <span className="font-mono text-gray-900">{selectedOrg.ogrn}</span>
                  </div>
                )}
                {selectedOrg.kpp && (
                  <div>
                    <span className="text-gray-600">КПП:</span>{' '}
                    <span className="font-mono text-gray-900">{selectedOrg.kpp}</span>
                  </div>
                )}
              </div>
            </div>

            {/* Адрес */}
            {selectedOrg.address?.full && (
              <div className="p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                  📍 Адрес
                </h3>
                <div className="space-y-2 text-sm">
                  <div>
                    <span className="text-gray-600">Полный адрес:</span>{' '}
                    <span className="text-gray-900">{selectedOrg.address.full}</span>
                  </div>
                  {selectedOrg.address.postalCode && (
                    <div>
                      <span className="text-gray-600">Индекс:</span>{' '}
                      <span className="font-mono text-gray-900">{selectedOrg.address.postalCode}</span>
                    </div>
                  )}
                  {selectedOrg.address.region && (
                    <div>
                      <span className="text-gray-600">Регион:</span>{' '}
                      <span className="text-gray-900">{selectedOrg.address.region}</span>
                    </div>
                  )}
                  {selectedOrg.address.city && (
                    <div>
                      <span className="text-gray-600">Город:</span>{' '}
                      <span className="text-gray-900">{selectedOrg.address.city}</span>
                    </div>
                  )}
                  {selectedOrg.address.street && (
                    <div>
                      <span className="text-gray-600">Улица:</span>{' '}
                      <span className="text-gray-900">{selectedOrg.address.street}</span>
                    </div>
                  )}
                  {selectedOrg.address.house && (
                    <div>
                      <span className="text-gray-600">Дом:</span>{' '}
                      <span className="text-gray-900">{selectedOrg.address.house}</span>
                    </div>
                  )}
                  {selectedOrg.address.flat && (
                    <div>
                      <span className="text-gray-600">Квартира/Офис:</span>{' '}
                      <span className="text-gray-900">{selectedOrg.address.flat}</span>
                    </div>
                  )}
                  {selectedOrg.address.geoLat && selectedOrg.address.geoLon && (
                    <div>
                      <span className="text-gray-600">Координаты:</span>{' '}
                      <span className="font-mono text-gray-900">
                        {selectedOrg.address.geoLat}, {selectedOrg.address.geoLon}
                      </span>
                    </div>
                  )}
                </div>

                {/* Карта */}
                <div className="mt-4">
                  <YandexMap
                    address={selectedOrg.address.full}
                    organizationName={selectedOrg.name}
                    latitude={selectedOrg.address.geoLat ? parseFloat(selectedOrg.address.geoLat) : undefined}
                    longitude={selectedOrg.address.geoLon ? parseFloat(selectedOrg.address.geoLon) : undefined}
                  />
                </div>
              </div>
            )}

            {/* Контакты */}
            {(selectedOrg.phones?.length || selectedOrg.emails?.length) && (
              <div className="p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                  📞 Контакты
                </h3>
                <div className="space-y-2 text-sm">
                  {selectedOrg.phones && selectedOrg.phones.length > 0 && (
                    <div>
                      <span className="text-gray-600">Телефоны:</span>
                      <div className="mt-1 space-y-1">
                        {selectedOrg.phones.map((phone, idx) => (
                          <div key={idx}>
                            <a href={`tel:${phone}`} className="text-blue-600 hover:text-blue-800 font-mono">
                              {phone}
                            </a>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                  {selectedOrg.emails && selectedOrg.emails.length > 0 && (
                    <div>
                      <span className="text-gray-600">Email:</span>
                      <div className="mt-1 space-y-1">
                        {selectedOrg.emails.map((email, idx) => (
                          <div key={idx}>
                            <a href={`mailto:${email}`} className="text-blue-600 hover:text-blue-800">
                              {email}
                            </a>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Руководство */}
            {selectedOrg.director?.name && (
              <div className="p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                  👤 Руководство
                </h3>
                <div className="space-y-2 text-sm">
                  <div>
                    <span className="text-gray-600">{selectedOrg.director.post || 'Руководитель'}:</span>{' '}
                    <span className="text-gray-900">{selectedOrg.director.name}</span>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h3 className="font-semibold text-blue-900 mb-2">💡 Возможности поиска:</h3>
        <ul className="text-sm text-blue-800 space-y-1">
          <li>• Поиск по ИНН (например: 7707083893)</li>
          <li>• Поиск по ОГРН (например: 1027700132195)</li>
          <li>• Поиск по названию (например: Сбербанк)</li>
          <li>• Полная информация: реквизиты, адрес, контакты</li>
          <li>• Координаты для отображения на карте</li>
          <li>• Телефоны и email для связи</li>
          <li>• Почтовый индекс для отправки корреспонденции</li>
        </ul>
      </div>
    </div>
  );
}
