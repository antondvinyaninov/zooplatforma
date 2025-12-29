-- Миграция 015: Добавление куратора и локации питомца
-- Дата: 2024-12-29

-- Добавляем поле curator_id (ID куратора в системе)
ALTER TABLE pets ADD COLUMN curator_id INTEGER DEFAULT NULL;

-- Добавляем поле curator_name (ФИО куратора, если не в системе)
ALTER TABLE pets ADD COLUMN curator_name TEXT DEFAULT '';

-- Добавляем поле curator_phone (телефон куратора)
ALTER TABLE pets ADD COLUMN curator_phone TEXT DEFAULT '';

-- Добавляем поле location (текущая локация: home/street/foster/shelter)
ALTER TABLE pets ADD COLUMN location TEXT DEFAULT 'home';

-- Добавляем поле foster_address (адрес передержки)
ALTER TABLE pets ADD COLUMN foster_address TEXT DEFAULT '';

-- Добавляем поле shelter_name (название приюта)
ALTER TABLE pets ADD COLUMN shelter_name TEXT DEFAULT '';

-- Комментарии к полям:
-- curator_id: ID пользователя-куратора в системе (может быть NULL)
-- curator_name: ФИО куратора (если куратор не зарегистрирован в системе)
-- curator_phone: Телефон куратора для связи
-- location: Текущая локация питомца
--   - 'home' - дома у владельца
--   - 'street' - на улице (бездомный)
--   - 'foster' - на передержке
--   - 'shelter' - в приюте
-- foster_address: Адрес передержки (если location = 'foster')
-- shelter_name: Название приюта (если location = 'shelter')

-- Примечание: 
-- - Если есть owner (user_id), то location обычно 'home'
-- - Если есть curator, но нет owner, то location может быть 'street', 'foster', 'shelter'
-- - Если нет ни owner, ни curator, то животное прошло через приют и пристроено
