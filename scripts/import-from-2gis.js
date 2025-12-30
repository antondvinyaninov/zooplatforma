// Импорт организаций из 2GIS API
const fetch = require('node-fetch');
const sqlite3 = require('sqlite3').verbose();

const TWOGIS_API_KEY = '7bbfde23-157e-4494-ad6a-4da44c75283c';
const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';
const DB_PATH = '../database/data.db';

// Координаты Сарапула
const SARAPUL_COORDS = {
  lon: 53.803333,
  lat: 56.466667,
};

// Поисковые запросы для 2GIS (по Удмуртской республике)
const SEARCH_QUERIES = [
  'ветеринарная клиника Удмуртия',
  'ветклиника Ижевск',
  'зоомагазин Ижевск',
  'товары для животных Удмуртия',
  'ветеринарная клиника Сарапул',
  'зоомагазин Сарапул',
];

// Поиск организаций через 2GIS API
async function search2GIS(query) {
  try {
    const url = `https://catalog.api.2gis.com/3.0/items?key=${TWOGIS_API_KEY}&q=${encodeURIComponent(query)}&page_size=10&fields=items.contact_groups,items.address`;
    
    console.log(`   Запрос к 2GIS...`);
    
    const response = await fetch(url);
    const data = await response.json();
    
    // Отладка
    if (data.meta?.code !== 200) {
      console.error(`   Ошибка API:`, JSON.stringify(data.meta, null, 2));
      return [];
    }
    
    const items = data.result?.items || [];
    
    console.log(`   Получено ${items.length} результатов (total: ${data.result?.total})`);
    
    return items;
  } catch (error) {
    console.error(`   Error: ${error.message}`);
    return [];
  }
}

// Поиск ИНН через DaData по названию
async function findINNByName(name) {
  try {
    const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party', {
      method: 'POST',
      headers: {
        'Authorization': `Token ${DADATA_TOKEN}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: name,
        count: 3,
        locations: [{
          city: 'Сарапул'
        }]
      }),
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    
    if (data.suggestions && data.suggestions.length > 0) {
      // Ищем наиболее подходящую организацию
      for (const suggestion of data.suggestions) {
        const suggestionName = suggestion.data.name?.short_with_opf?.toLowerCase() || '';
        const queryName = name.toLowerCase();
        
        // Проверяем совпадение названия
        if (suggestionName.includes(queryName) || queryName.includes(suggestionName)) {
          return suggestion;
        }
      }
      
      // Если точного совпадения нет, берем первую
      return data.suggestions[0];
    }
    
    return null;
  } catch (error) {
    return null;
  }
}

// Получение города по building_id
async function getCityByBuildingId(buildingId) {
  try {
    const url = `https://catalog.api.2gis.com/3.0/items/byid?id=${buildingId}&key=${TWOGIS_API_KEY}`;
    const response = await fetch(url);
    const data = await response.json();
    
    if (data.meta?.code === 200 && data.result?.items?.[0]?.full_name) {
      const fullName = data.result.items[0].full_name;
      // full_name имеет формат "Город, адрес"
      const parts = fullName.split(',');
      if (parts.length > 0) {
        return parts[0].trim();
      }
    }
    
    return null;
  } catch (error) {
    return null;
  }
}
function determineOrgType(name, rubrics) {
  const nameLower = name.toLowerCase();
  const rubricsStr = rubrics.map(r => r.name).join(' ').toLowerCase();
  
  if (nameLower.includes('ано') || nameLower.includes('фонд') || nameLower.includes('помощ') || nameLower.includes('защит')) {
    return 'foundation';
  }
  
  if (nameLower.includes('ветеринар') || rubricsStr.includes('ветеринар')) {
    return 'vet_clinic';
  }
  
  if (nameLower.includes('зоомагазин') || nameLower.includes('зоотовар') || rubricsStr.includes('зоомагазин') || rubricsStr.includes('товары для животных')) {
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

// Извлечение телефона из contact_groups
function extractPhone(contactGroups) {
  if (!contactGroups || contactGroups.length === 0) return null;
  
  for (const group of contactGroups) {
    if (group.contacts) {
      for (const contact of group.contacts) {
        if (contact.type === 'phone' && contact.text) {
          return contact.text;
        }
      }
    }
  }
  
  return null;
}

// Извлечение сайта из contact_groups
function extractWebsite(contactGroups) {
  if (!contactGroups || contactGroups.length === 0) return null;
  
  for (const group of contactGroups) {
    if (group.contacts) {
      for (const contact of group.contacts) {
        if (contact.type === 'website' && contact.url) {
          return contact.url;
        }
      }
    }
  }
  
  return null;
}

// Вставка организации в базу данных
async function insertOrganization(db, org, dadataOrg = null, userId = 1) {
  return new Promise(async (resolve, reject) => {
    const name = org.name || 'Без названия';
    const rubrics = org.rubrics || [];
    const type = determineOrgType(name, rubrics);
    
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
    
    // Контакты из 2GIS
    const phone = extractPhone(org.contact_groups);
    const website = extractWebsite(org.contact_groups);
    
    // Адрес из 2GIS
    const address = org.address_name || org.address?.name || null;
    
    // Извлекаем город из full_name или запрашиваем по building_id
    let city = null;
    if (org.full_name) {
      // full_name имеет формат "Город, адрес"
      const parts = org.full_name.split(',');
      if (parts.length > 0) {
        city = parts[0].trim();
      }
    } else if (org.address?.building_id) {
      // Если full_name нет, пробуем получить город по building_id
      city = await getCityByBuildingId(org.address.building_id);
    }
    
    // Координаты из 2GIS
    const coords = org.point?.lon && org.point?.lat ? org.point : null;
    
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
      website,
      address,
      city,
      coords?.lat || null,
      coords?.lon || null,
      rubrics.map(r => r.name).join(', ') || null, // description
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
        if (city) console.log(`     Город: ${city}`);
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
  console.log('🗺️  Поиск организаций в Сарапуле через 2GIS API...\n');

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
    const results = await search2GIS(query);
    console.log(`   Найдено: ${results.length} организаций`);
    
    if (results.length > 0) {
      console.log(`   Примеры: ${results.slice(0, 3).map(r => r.name).join(', ')}`);
    }
    
    results.forEach(org => {
      const id = org.id;
      if (id && !allOrganizations.has(id)) {
        allOrganizations.set(id, org);
      }
    });
    
    // Задержка между запросами
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  console.log(`\n\n📊 Всего уникальных организаций: ${allOrganizations.size}`);
  console.log('\n💾 Импорт в базу данных...\n');

  let added = 0;
  let skipped = 0;

  for (const [id, org] of allOrganizations.entries()) {
    try {
      // Пытаемся найти ИНН через DaData
      const dadataOrg = await findINNByName(org.name);
      
      const result = await insertOrganization(db, org, dadataOrg);
      if (result.added) added++;
      if (result.skipped) skipped++;
      
      // Задержка между запросами к DaData
      await new Promise(resolve => setTimeout(resolve, 300));
    } catch (error) {
      console.error(`  ❌ Ошибка при добавлении "${org.name}":`, error.message);
    }
  }

  console.log(`\n\n✨ Импорт завершен!`);
  console.log(`   ✅ Добавлено: ${added}`);
  console.log(`   ⏭️  Пропущено: ${skipped}`);
  console.log(`   📊 Всего: ${allOrganizations.size}`);

  db.close();
}

main().catch(console.error);
