'use client';

import PetCard from '@/app/components/pets/PetCard';

export default function PetCardsDemoPage() {
  // Тестовые данные для разных типов карточек
  const pets = [
    {
      id: 1,
      name: 'Барсик',
      species: 'Кошка',
      breed: 'Британская короткошерстная',
      gender: 'male' as const,
      birth_date: '2024-09-01',
      color: 'Ласковый, спокойный, не агрессивный',
      photo: null,
      card_type: 'looking_for_home' as const,
      card_description: 'Милый котенок 3 месяца, приучен к лотку. Очень ласковый и игривый. Ищет добрую семью, где его будут любить.',
      card_location_city: 'Москва',
      is_vaccinated: true,
      is_sterilized: false,
      chip_number: null,
    },
    {
      id: 2,
      name: 'Рекс',
      species: 'Собака',
      breed: 'Лабрадор',
      gender: 'male' as const,
      birth_date: '2019-05-15',
      color: 'Золотистый',
      photo: null,
      card_type: 'lost' as const,
      card_description: 'Потерялся лабрадор Рекс, 5 лет. Очень дружелюбный, откликается на имя. Последний раз видели в парке Сокольники.',
      card_location_city: 'Москва',
      card_lost_date: '2024-12-28',
      card_lost_location: 'Парк Сокольники, около главного входа',
      card_reward_amount: 10000,
      is_vaccinated: true,
      is_sterilized: true,
      chip_number: '123456789',
    },
    {
      id: 3,
      name: 'Мурка',
      species: 'Кошка',
      breed: 'Дворовая',
      gender: 'female' as const,
      birth_date: '2020-03-10',
      color: 'Трехцветная',
      photo: null,
      card_type: 'fundraising' as const,
      card_description: 'Мурке нужна срочная операция на лапке после травмы. Помогите спасти жизнь кошечке!',
      card_location_city: 'Москва',
      card_fundraising_goal: 50000,
      card_fundraising_current: 28500,
      card_fundraising_purpose: 'Операция на лапке после травмы',
      is_vaccinated: true,
      is_sterilized: true,
      chip_number: null,
    },
    {
      id: 4,
      name: 'Найденыш',
      species: 'Собака',
      breed: 'Неизвестна',
      gender: 'male' as const,
      birth_date: null,
      color: 'Черный с белыми пятнами',
      photo: null,
      card_type: 'found' as const,
      card_description: 'Найден щенок около метро Сокольники. Очень ласковый, но напуган. Ищем владельца или новый дом.',
      card_location_city: 'Москва',
      card_location_address: 'ул. Сокольническая, д. 15',
      is_vaccinated: false,
      is_sterilized: false,
      chip_number: null,
    },
  ];

  const handleAction = (petName: string, action: string) => {
    console.log(`Action "${action}" for pet "${petName}"`);
    alert(`Действие "${action}" для питомца "${petName}"`);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Демо: Карточки питомцев
        </h1>
        <p className="text-gray-600">
          Примеры карточек питомцев для разных типов объявлений
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {pets.map((pet) => (
          <div key={pet.id}>
            <PetCard 
              pet={pet} 
              onAction={(action) => handleAction(pet.name, action)}
            />
          </div>
        ))}
      </div>

      {/* Пример в посте */}
      <div className="mt-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-4">
          Пример карточки в посте
        </h2>
        
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          {/* Автор поста */}
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-bold">
              А
            </div>
            <div>
              <div className="font-semibold text-gray-900">Антон Двинянинов</div>
              <div className="text-sm text-gray-500">2 часа назад</div>
            </div>
          </div>

          {/* Текст поста */}
          <div className="mb-4">
            <p className="text-gray-800">
              Друзья, помогите найти дом для этого замечательного котенка! 
              Барсик очень ласковый и игривый, приучен к лотку. 
              Ищет любящую семью 🏠❤️
            </p>
          </div>

          {/* Карточка питомца */}
          <div className="max-w-md">
            <PetCard 
              pet={pets[0]} 
              onAction={(action) => handleAction(pets[0].name, action)}
            />
          </div>

          {/* Действия поста */}
          <div className="mt-4 pt-4 border-t border-gray-200 flex items-center gap-6 text-gray-500">
            <button className="flex items-center gap-2 hover:text-blue-500">
              <span>❤️</span>
              <span>12</span>
            </button>
            <button className="flex items-center gap-2 hover:text-blue-500">
              <span>💬</span>
              <span>5</span>
            </button>
            <button className="flex items-center gap-2 hover:text-blue-500">
              <span>🔄</span>
              <span>3</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
