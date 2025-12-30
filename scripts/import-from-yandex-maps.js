// Импорт организаций из Яндекс.Карт
const fetch = require('node-fetch');
const sqlite3 = require('sqlite3').verbose();

const YANDEX_API_KEY = 'ece8ef8e-8782-426f-951d-79e965468547';
const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';
const DB_PATH = '../database/data.db';

// Поисковые запросы для Яндекс.Карт
const SEARCH_QUERIES = [
  'ветеринарная клиника',
  'зоомагазин',
  'приют для животных',
  'товары для животных',
  'ветеринарная аптека',
];

// Координаты Сарапула
const SARAPUL_COORDS = {
  lon: 53.803333,
  lat: 56.466667,
};

// Поиск организаций через Яндекс.Карты API
async function searchYandexMaps(query) {
  try {
    const url = `https://search-maps.yandex.ru/v1/?apikey=${YANDEX_API_KEY}&text=${encodeURIComponent(query + ' Сарапул')}&lang=ru_RU&ll=${SARAPUL_COORDS.lon},${SARAPUL_COORDS.lat}&spn=0.1,0.1&results=50`;
    
    const response = await fetch(url);
    
    if (!response.ok) {
      console.error(`Yandex Maps error for "${query}":`, response.status);
      return [];
    }

    const data = await response.json();
    return data.features || [];
  } catch (error) {
    console.error(`Error searching "${query}":`, error.message);
    return [];
  }
}

// Поиск ИНН через DaData по названию и адресу
async function findINNByNameAndAddress(name, address) {
  try {
    const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party', {
      method: 'POST',
      headers: {
        'Authorization': `Token ${DADATA_TOKEN}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: name,
        count: 5,
        locations: [{
          city: 'Сарапул'
        }]
      }),
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    
    // Ищем наиболее подходящую организацию
    if (data.suggestions && data.suggestions.length > 0) {
      // Берем первую, так как DaData сортирует по релевантности
      return data.suggestions[0];
    }
    
    return null;
  } catch (error) {
    return null;
  }
}

// Определение типа организации
function determineOrgType(name, categories) {
  const nameLower = name.toLowerCase();
  const categoriesStr = categories.join(' ').toLowerCase();
  
  if (nameLower.includes('ано') || nameLower.includes('фонд') || nameLower.includes('помощ') || nameLower.includes('защит')) {
    return 'foundation';
  }
  
  if (nameLower.includes('ветеринар') || categoriesStr.includes('ветеринар')) {
    return 'vet_clinic';
  }
  
  if (nameLower.includes('зоомагазин') || nameLower.includes('зоотовар') || categoriesStr.includes('зоомагазин') || categoriesStr.includes('товары для животных')) {
    return 'pet_shop';
  }
  
  if (nameLower.includes('приют')) {
    return 'shelter';
  }
  
  if (nameLower.includes('кинолог') || nameLower.includes('дрессировка')) {
    return 'kennel';
  }
  
  return 'other';
}

// Вставка организации в базу данных
function insertOrganization(db, org, dadataOrg = null, userId = 1) {
  return new Promise((resolve, reject) => {
    const props = org.properties;
    const coords = org.geometry.coordinates;
    
    const name = props.name || props.CompanyMetaData?.name || 'Без названия';
    const description = props.description || null;
    const categories = props.CompanyMetaData?.Categories || [];
    const type = determineOrgType(name, categories);
    
    // Данные из DaData (если найдены)
    const inn = dadataOrg?.data?.inn || null;
    const ogrn = dadataOrg?.data?.ogrn || null;
    const kpp = dadataOrg?.data?.kpp || null;
    const legalForm = dadataOrg?.data?.opf?.full || null;
    const fullName = dadataOrg?.data?.name?.full_with_opf || name;
    const shortName = dadataOrg?.data?.name?.short_with_opf || name;
    const registrationDate = dadataOrg?.data?.state?.registration_date || null;
    const directorName = dadataOrg?.data?.management?.name || null;
    const directorPosition = dadataOrg?.data?.management?.post || null;
    
    // Контакты
    const phone = props.CompanyMetaData?.Phones?.[0]?.formatted || null;
    const url = props.CompanyMetaData?.url || null;
    
    // Адрес
    const address = props.description || props.CompanyMetaData?.address || null;
    
    const query = `
      INSERT OR IGNORE INTO organizations (
        name, short_name, legal_form, type,
        inn, ogrn, kpp, registration_date,
        phone, website,
        address_full, address_city,
        geo_lat, geo_lon,
        description, bio,
        director_name, director_position,
        owner_user_id,
        is_verified, status,
        created_at, updated_at
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
    `;

    const params = [
      fullName,
      shortName,
      legalForm,
      type,
      inn,
      ogrn,
      kpp,
      registrationDate,
      phone,
      url,
      address,
      'Сарапул',
      coords[1], // lat
      coords[0], // lon
      description,
      null, // bio
      directorName,
      directorPosition,
      userId,
      inn ? 1 : 0, // is_verified (если есть ИНН, считаем верифицированной)
      'active',
    ];

    db.run(query, params, function(err) {
      if (err) {
        if (err.message.includes('UNIQUE constraint')) {
          resolve({ skipped: true });
        } else {
          reject(err);
        }
      } else if (this.changes > 0) {
        console.log(`  ✅ ${name} (${type})`);
        if (inn) console.log(`     ИНН: ${inn}`);
        if (phone) console.log(`     Телефон: ${phone}`);
        if (address) console.log(`     Адрес: ${address}`);
        
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
  console.log('🗺️  Поиск организаций в Сарапуле через Яндекс.Карты...\n');

  const db = new sqlite3.Database(DB_PATH, (err) => {
    if (err) {
      console.error('Ошибка подключения к БД:', err.message);
      process.exit(1);
    }
  });

  const allOrganizations = new Map();

  // Поиск по всем запросам
  for (const query of SEARCH_QUERIES) {
    console.log(`\n📋 Поиск: "${query}"`);
    const results = await searchYandexMaps(query);
    console.log(`   Найдено: ${results.length} организаций`);
    
    results.forEach(org => {
      const name = org.properties?.name || org.properties?.CompanyMetaData?.name;
      if (name && !allOrganizations.has(name)) {
        allOrganizations.set(name, org);
      }
    });
    
    // Задержка между запросами
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  console.log(`\n\n📊 Всего уникальных организаций: ${allOrganizations.size}`);
  console.log('\n💾 Импорт в базу данных...\n');

  let added = 0;
  let skipped = 0;

  for (const [name, org] of allOrganizations.entries()) {
    try {
      // Пытаемся найти ИНН через DaData
      const address = org.properties?.description || org.properties?.CompanyMetaData?.address;
      const dadataOrg = await findINNByNameAndAddress(name, address);
      
      const result = await insertOrganization(db, org, dadataOrg);
      if (result.added) added++;
      if (result.skipped) skipped++;
      
      // Задержка между запросами к DaData
      await new Promise(resolve => setTimeout(resolve, 300));
    } catch (error) {
      console.error(`  ❌ Ошибка при добавлении "${name}":`, error.message);
    }
  }

  console.log(`\n\n✨ Импорт завершен!`);
  console.log(`   ✅ Добавлено: ${added}`);
  console.log(`   ⏭️  Пропущено: ${skipped}`);
  console.log(`   📊 Всего: ${allOrganizations.size}`);

  db.close();
}

main().catch(console.error);
