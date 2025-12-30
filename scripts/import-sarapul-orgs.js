// Скрипт для импорта организаций Сарапула из DaData
const fetch = require('node-fetch');
const sqlite3 = require('sqlite3').verbose();

const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';
const DB_PATH = '../database/data.db';

// Поисковые запросы для разных типов организаций
const SEARCH_QUERIES = [
  { query: 'ветеринар', okved: '75.00' },  // Ветеринарная деятельность
  { query: 'зоо', okved: '47.76' },        // Торговля зоотоварами
  { query: 'зоомагазин', okved: null },    // Зоомагазины
  { query: 'зоотовары', okved: null },     // Зоотовары
  { query: 'корм', okved: null },          // Корма для животных
];

// Определение типа организации по названию и ОКВЭД
function determineOrgType(org) {
  const name = (org.data.name?.full_with_opf || '').toLowerCase();
  const okved = org.data.okved || '';
  
  if (name.includes('ветеринар') || okved.includes('75.00')) {
    return 'vet_clinic';
  }
  if (name.includes('зоомагазин') || name.includes('зоотовар') || okved.includes('47.76')) {
    return 'pet_shop';
  }
  if (name.includes('приют') || name.includes('защита животных')) {
    return 'shelter';
  }
  if (name.includes('фонд') || name.includes('ано')) {
    return 'foundation';
  }
  if (name.includes('кинолог') || name.includes('дрессировка')) {
    return 'kennel';
  }
  
  return 'other';
}

// Поиск организаций через DaData
async function searchOrganizations(searchConfig) {
  try {
    const requestBody = {
      query: searchConfig.query,
      count: 20,
      status: ['ACTIVE'], // Только действующие
      locations: [{
        city: 'Сарапул'
      }]
    };

    // Добавляем фильтр по ОКВЭД если указан
    if (searchConfig.okved) {
      requestBody.okved = [searchConfig.okved];
    }

    const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party', {
      method: 'POST',
      headers: {
        'Authorization': `Token ${DADATA_TOKEN}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestBody),
    });

    if (!response.ok) {
      console.error(`DaData error for "${searchConfig.query}":`, response.status);
      return [];
    }

    const data = await response.json();
    return data.suggestions || [];
  } catch (error) {
    console.error(`Error searching "${searchConfig.query}":`, error.message);
    return [];
  }
}

// Вставка организации в базу данных
function insertOrganization(db, org, userId = 1) {
  return new Promise((resolve, reject) => {
    const type = determineOrgType(org);
    const data = org.data;
    
    // Фильтруем только организации из Сарапула
    const city = data.address?.data?.city || '';
    if (city !== 'Сарапул') {
      console.log(`  ⏭️  Пропущено (не Сарапул): ${data.name?.short_with_opf || org.value} (${city})`);
      resolve({ skipped: true });
      return;
    }
    
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
      data.address?.data?.city || 'Сарапул',
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

// Основная функция
async function main() {
  console.log('🔍 Поиск организаций в Сарапуле...\n');

  const db = new sqlite3.Database(DB_PATH, (err) => {
    if (err) {
      console.error('Ошибка подключения к БД:', err.message);
      process.exit(1);
    }
  });

  const allOrganizations = new Map(); // Используем Map для удаления дубликатов по ИНН

  // Поиск по всем запросам
  for (const searchConfig of SEARCH_QUERIES) {
    console.log(`\n📋 Поиск: "${searchConfig.query}"${searchConfig.okved ? ` (ОКВЭД: ${searchConfig.okved})` : ''}`);
    const results = await searchOrganizations(searchConfig);
    console.log(`   Найдено: ${results.length} организаций`);
    
    results.forEach(org => {
      const inn = org.data.inn;
      if (inn && !allOrganizations.has(inn)) {
        allOrganizations.set(inn, org);
      }
    });
    
    // Задержка между запросами
    await new Promise(resolve => setTimeout(resolve, 500));
  }

  console.log(`\n\n📊 Всего уникальных организаций: ${allOrganizations.size}`);
  console.log('\n💾 Импорт в базу данных...\n');

  let added = 0;
  let skipped = 0;

  for (const org of allOrganizations.values()) {
    try {
      const result = await insertOrganization(db, org);
      if (result.added) added++;
      if (result.skipped) skipped++;
    } catch (error) {
      console.error(`  ❌ Ошибка: ${error.message}`);
    }
  }

  console.log(`\n\n✨ Импорт завершен!`);
  console.log(`   ✅ Добавлено: ${added}`);
  console.log(`   ⏭️  Пропущено: ${skipped}`);
  console.log(`   📊 Всего: ${allOrganizations.size}`);

  db.close();
}

main().catch(console.error);
