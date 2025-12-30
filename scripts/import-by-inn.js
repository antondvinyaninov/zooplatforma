// Импорт организаций по ИНН через DaData
const fetch = require('node-fetch');
const sqlite3 = require('sqlite3').verbose();

const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';
const DB_PATH = '../database/data.db';

// ИНН организаций Сарапула, связанных с животными
const INN_LIST = [
  '1838026383', // Зоопомощь Сарапул
  // Добавьте сюда ИНН других организаций Сарапула
  // Примеры (нужно проверить реальные ИНН):
  // '1838XXXXXX', // Ветклиника
  // '1838YYYYYY', // Зоомагазин
];

// Определение типа организации по названию и ОКВЭД
function determineOrgType(org) {
  const name = (org.data.name?.full_with_opf || '').toLowerCase();
  const okved = org.data.okved || '';
  
  // Фонды и АНО
  if (name.includes('ано') || name.includes('фонд') || name.includes('помощ') || name.includes('защит')) {
    return 'foundation';
  }
  
  // Ветеринарные клиники
  if (name.includes('ветеринар') || okved.includes('75.00')) {
    return 'vet_clinic';
  }
  
  // Зоомагазины
  if (name.includes('зоомагазин') || name.includes('зоотовар') || okved.includes('47.76')) {
    return 'pet_shop';
  }
  
  // Приюты
  if (name.includes('приют')) {
    return 'shelter';
  }
  
  // Кинологические центры
  if (name.includes('кинолог') || name.includes('дрессировка')) {
    return 'kennel';
  }
  
  return 'other';
}

// Поиск организации по ИНН через DaData
async function searchByINN(inn) {
  try {
    const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party', {
      method: 'POST',
      headers: {
        'Authorization': `Token ${DADATA_TOKEN}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: inn,
        count: 1,
      }),
    });

    if (!response.ok) {
      console.error(`DaData error for INN ${inn}:`, response.status);
      return null;
    }

    const data = await response.json();
    return data.suggestions && data.suggestions.length > 0 ? data.suggestions[0] : null;
  } catch (error) {
    console.error(`Error searching INN ${inn}:`, error.message);
    return null;
  }
}

// Вставка организации в базу данных
function insertOrganization(db, org, userId = 1) {
  return new Promise((resolve, reject) => {
    const type = determineOrgType(org);
    const data = org.data;
    
    const query = `
      INSERT OR IGNORE INTO organizations (
        name, short_name, legal_form, type,
        inn, ogrn, kpp, registration_date,
        email, phone, website,
        address_full, address_postal_code, address_region, address_city,
        address_street, address_house, address_office,
        geo_lat, geo_lon,
        description, bio,
        director_name, director_position,
        owner_user_id,
        is_verified, status,
        created_at, updated_at
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
    `;

    const phone = data.phones && data.phones[0] ? data.phones[0].value : null;
    const email = data.emails && data.emails[0] ? data.emails[0].value : null;
    
    const params = [
      data.name?.full_with_opf || org.value,
      data.name?.short_with_opf || null,
      data.opf?.full || null,
      type,
      data.inn || null,
      data.ogrn || null,
      data.kpp || null,
      data.state?.registration_date || null,
      email,
      phone,
      null, // website
      data.address?.unrestricted_value || data.address?.value || null,
      data.address?.data?.postal_code || null,
      data.address?.data?.region_with_type || null,
      data.address?.data?.city || null,
      data.address?.data?.street_with_type || null,
      data.address?.data?.house || null,
      data.address?.data?.flat || null,
      data.address?.data?.geo_lat ? parseFloat(data.address.data.geo_lat) : null,
      data.address?.data?.geo_lon ? parseFloat(data.address.data.geo_lon) : null,
      null, // description
      null, // bio
      data.management?.name || null,
      data.management?.post || null,
      userId,
      0, // is_verified
      'active',
    ];

    db.run(query, params, function(err) {
      if (err) {
        if (err.message.includes('UNIQUE constraint')) {
          console.log(`  ⏭️  Пропущено (уже есть): ${data.name?.short_with_opf || org.value}`);
          resolve({ skipped: true });
        } else {
          reject(err);
        }
      } else if (this.changes > 0) {
        console.log(`  ✅ Добавлено: ${data.name?.short_with_opf || org.value} (${type})`);
        console.log(`     ИНН: ${data.inn}`);
        console.log(`     Адрес: ${data.address?.value || 'не указан'}`);
        console.log(`     Телефон: ${phone || 'не указан'}`);
        
        // Добавляем владельца в organization_members
        db.run(`
          INSERT INTO organization_members (organization_id, user_id, role, can_post, can_edit, can_manage_members)
          VALUES (?, ?, 'owner', 1, 1, 1)
        `, [this.lastID, userId], (err) => {
          if (err) console.error('Error adding member:', err.message);
        });
        
        resolve({ added: true, id: this.lastID });
      } else {
        resolve({ skipped: true });
      }
    });
  });
}

async function main() {
  console.log('🔍 Поиск организаций по ИНН через DaData...\n');

  const db = new sqlite3.Database(DB_PATH, (err) => {
    if (err) {
      console.error('Ошибка подключения к БД:', err.message);
      process.exit(1);
    }
  });

  let added = 0;
  let skipped = 0;
  let notFound = 0;

  for (const inn of INN_LIST) {
    console.log(`\n📋 Поиск организации с ИНН: ${inn}`);
    
    const org = await searchByINN(inn);
    
    if (!org) {
      console.log(`  ❌ Организация не найдена`);
      notFound++;
      continue;
    }

    try {
      const result = await insertOrganization(db, org);
      if (result.added) added++;
      if (result.skipped) skipped++;
    } catch (error) {
      console.error(`  ❌ Ошибка: ${error.message}`);
    }

    // Задержка между запросами
    await new Promise(resolve => setTimeout(resolve, 500));
  }

  console.log(`\n\n✨ Импорт завершен!`);
  console.log(`   ✅ Добавлено: ${added}`);
  console.log(`   ⏭️  Пропущено: ${skipped}`);
  console.log(`   ❌ Не найдено: ${notFound}`);
  console.log(`   📊 Всего: ${INN_LIST.length}`);

  db.close();
}

main().catch(console.error);
