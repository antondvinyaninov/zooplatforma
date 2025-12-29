# Ветеринарный паспорт питомца

## Обзор

Система ЗооБаза поддерживает полноценный цифровой ветеринарный паспорт для каждого питомца.

## Структура данных

### Основная информация (таблица `pets`)

**Идентификация:**
- `name` - Кличка
- `species` - Вид животного
- `breed` - Порода
- `gender` - Пол
- `birth_date` - Дата рождения
- `color` - Окрас
- `distinctive_marks` - Особые приметы (шрамы, отметины)
- `size` - Размер (small/medium/large)
- `weight` - Вес (кг)

**Документы:**
- `chip_number` - Номер микрочипа
- `tattoo_number` - Номер клейма (татуировки)
- `ear_tag_number` - Номер ушной бирки (ОСВВ - Отлов-Стерилизация-Вакцинация-Возврат, ставится в приютах)
- `passport_number` - Номер ветеринарного паспорта
- `pedigree_number` - Номер родословной
- `registration_org` - Организация регистрации (РКФ, FCI и т.д.)

**Владелец:**
- `user_id` - ID пользователя в системе
- `owner_name` - ФИО владельца
- `owner_address` - Адрес владельца
- `owner_phone` - Телефон владельца
- `owner_email` - Email владельца

**Медицинская информация:**
- `is_sterilized` - Стерилизован/кастрирован
- `is_vaccinated` - Вакцинирован
- `blood_type` - Группа крови
- `allergies` - Аллергии
- `chronic_diseases` - Хронические заболевания
- `current_medications` - Текущие лекарства
- `health_notes` - Заметки о здоровье

**Характер:**
- `character_traits` - Черты характера
- `special_needs` - Особые потребности

**Статус:**
- `status` - Статус (home/looking_for_home/lost/found/adopted)
- `status_updated_at` - Дата изменения статуса

**История:**
- `story` - История питомца
- `created_at` - Дата создания карточки
- `updated_at` - Дата последнего обновления

### Прививки (таблица `pet_vaccinations`)

- `vaccination_date` - Дата прививки
- `vaccine_name` - Название вакцины (Нобивак, Эурикан и т.д.)
- `vaccine_series` - Серия вакцины
- `disease` - От какой болезни (бешенство, чума, энтерит и т.д.)
- `next_vaccination_date` - Дата следующей прививки
- `veterinarian_name` - ФИО ветеринара
- `clinic_name` - Название клиники
- `clinic_stamp` - Печать клиники (путь к изображению)
- `notes` - Примечания

### Обработки от паразитов (таблица `pet_parasite_treatments`)

- `treatment_date` - Дата обработки
- `treatment_type` - Тип: external (блохи, клещи), internal (глисты), both (комплексная)
- `drug_name` - Название препарата
- `drug_series` - Серия препарата
- `dosage` - Дозировка
- `next_treatment_date` - Дата следующей обработки
- `veterinarian_name` - ФИО ветеринара
- `notes` - Примечания

### Операции и процедуры (таблица `pet_surgeries`)

- `surgery_date` - Дата операции
- `surgery_type` - Тип операции/процедуры
- `description` - Описание
- `veterinarian_name` - ФИО ветеринара
- `clinic_name` - Название клиники
- `anesthesia_type` - Тип анестезии
- `complications` - Осложнения
- `outcome` - Исход
- `notes` - Примечания

### Болезни и лечение (таблица `pet_medical_history`)

- `diagnosis_date` - Дата диагноза
- `diagnosis` - Диагноз
- `symptoms` - Симптомы
- `treatment` - Лечение
- `medications` - Назначенные лекарства
- `recovery_date` - Дата выздоровления
- `veterinarian_name` - ФИО ветеринара
- `clinic_name` - Название клиники
- `notes` - Примечания

## API Endpoints

### Основная информация питомца
- `GET /api/pets/:id` - Получить полную информацию о питомце
- `PUT /api/pets/:id` - Обновить информацию о питомце

### Прививки
- `GET /api/pets/:id/vaccinations` - Получить список прививок
- `POST /api/pets/:id/vaccinations` - Добавить прививку
- `PUT /api/pets/:id/vaccinations/:vaccination_id` - Обновить прививку
- `DELETE /api/pets/:id/vaccinations/:vaccination_id` - Удалить прививку

### Обработки от паразитов
- `GET /api/pets/:id/parasite-treatments` - Получить список обработок
- `POST /api/pets/:id/parasite-treatments` - Добавить обработку
- `PUT /api/pets/:id/parasite-treatments/:treatment_id` - Обновить обработку
- `DELETE /api/pets/:id/parasite-treatments/:treatment_id` - Удалить обработку

### Операции
- `GET /api/pets/:id/surgeries` - Получить список операций
- `POST /api/pets/:id/surgeries` - Добавить операцию
- `PUT /api/pets/:id/surgeries/:surgery_id` - Обновить операцию
- `DELETE /api/pets/:id/surgeries/:surgery_id` - Удалить операцию

### Медицинская история
- `GET /api/pets/:id/medical-history` - Получить медицинскую историю
- `POST /api/pets/:id/medical-history` - Добавить запись
- `PUT /api/pets/:id/medical-history/:record_id` - Обновить запись
- `DELETE /api/pets/:id/medical-history/:record_id` - Удалить запись

## Использование

### Создание полной карточки питомца

```json
{
  "user_id": 1,
  "name": "Барсик",
  "species": "Кошка",
  "breed": "Британская короткошёрстная",
  "gender": "male",
  "birth_date": "2020-05-15",
  "color": "Серый",
  "distinctive_marks": "Белое пятно на груди",
  "size": "medium",
  "weight": 5.2,
  "chip_number": "123456789012345",
  "passport_number": "RU-12345",
  "owner_name": "Иванов Иван Иванович",
  "owner_address": "г. Москва, ул. Ленина, д. 1, кв. 10",
  "owner_phone": "+7 (999) 123-45-67",
  "owner_email": "ivanov@example.com",
  "is_sterilized": true,
  "is_vaccinated": true,
  "blood_type": "A",
  "character_traits": "Спокойный, ласковый",
  "status": "home"
}
```

### Добавление прививки

```json
{
  "vaccination_date": "2024-01-15",
  "vaccine_name": "Нобивак Tricat Trio",
  "vaccine_series": "A123456",
  "disease": "Панлейкопения, калицивироз, ринотрахеит",
  "next_vaccination_date": "2025-01-15",
  "veterinarian_name": "Петров П.П.",
  "clinic_name": "Ветклиника Айболит",
  "notes": "Перенёс хорошо"
}
```

## Миграции

- `012_add_pet_passport_fields.sql` - Добавление полей паспорта
- `013_create_medical_records_tables.sql` - Создание таблиц медицинских записей

## Следующие шаги

1. Создать API endpoints для медицинских записей
2. Добавить формы для ввода данных в frontend
3. Реализовать загрузку фотографий печатей клиник
4. Добавить напоминания о предстоящих прививках/обработках
5. Экспорт паспорта в PDF

---

**Последнее обновление:** 29 декабря 2024
