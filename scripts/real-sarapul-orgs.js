// Реальные организации Сарапула (данные из открытых источников)
const fetch = require('node-fetch');
const sqlite3 = require('sqlite3').verbose();

const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';
const DB_PATH = '../database/data.db';

// Реальные ИНН организаций Сарапула, связанных с животными
// Источники: 2GIS, Яндекс.Карты, Google Maps, сайты организаций
const ORGANIZATIONS_DATA = [
  {
    inn: '1838026383',
    name: 'АНО "Зоопомощь"',
    type: 'foundation',
  },
  // Добавьте сюда ИНН других реальных организаций Сарапула
  // Можно найти через:
  // - 2GIS: https://2gis.ru/sarapul
  // - Яндекс.Карты: https://yandex.ru/maps/11118/sarapul/
  // - Google Maps: поиск "ветеринарная клиника Сарапул"
  // - Сайты организаций
];

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
      return null;
    }

    const data = await response.json();
    return data.suggestions && data.suggestions.length > 0 ? data.suggestions[0] : null;
  } catch (error) {
    return null;
  }
}

// Определение типа организации
function determineOrgType(org, providedType) {
  if (providedType) return providedType;
  
  const name = (org.data.name?.full_with_opf || '').toLowerCase();
  const okved = org.data.okved || '';
  
  if (name.includes('ано') || name.includes('фонд') || name.includes('помощ') || name.includes('защит')) {
    return 'foundation';
  }
  
  if (name.includes('ветеринар') || okved.includes('75.00')) {
    return 'vet_clinic';
  }
  
  if (name.includes('зоомагазин') || name.includes('зоотовар') || okved.includes('47.76')) {
    return 'pet_shop';
  }
  
  if (name.includes('приют')) {
    return 'shelter';
  }
  
  if (name.includes('кинолог') || name.includes('дрессировка')) {
    return 'kennel';
  }
  
  return 'other';
}

// Вставка организации в базу данных
function insertOrganization(db, org, type, userId = 1) {
  return new Promise((resolve, reject) => {
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
      1, // is_verified (реальные организации с ИНН)
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
        console.log(`  ✅ ${data.name?.short_with_opf || org.value} (${type})`);
        console.log(`     ИНН: ${data.inn}`);
        console.log(`     Адрес: ${data.address?.value || 'не указан'}`);
        if (phone) console.log(`     Телефон: ${phone}`);
        
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
  console.log('🔍 Импорт реальных организаций Сарапула...\n');
  console.log('📝 Для добавления новых организаций:');
  console.log('   1. Найдите организацию на 2GIS, Яндекс.Картах или Google Maps');
  console.log('   2. Найдите ИНН организации (на сайте, в справочниках)');
  console.log('   3. Добавьте ИНН в массив ORGANIZATIONS_DATA в этом скрипте');
  console.log('   4. Запустите скрипт снова\n');

  const db = new sqlite3.Database(DB_PATH, (err) => {
    if (err) {
      console.error('Ошибка подключения к БД:', err.message);
      process.exit(1);
    }
  });

  let added = 0;
  let skipped = 0;
  let notFound = 0;

  for (const orgData of ORGANIZATIONS_DATA) {
    console.log(`\n📋 Поиск: ${orgData.name || orgData.inn}`);
    
    const org = await searchByINN(orgData.inn);
    
    if (!org) {
      console.log(`  ❌ Организация не найдена в DaData`);
      notFound++;
      continue;
    }

    try {
      const type = determineOrgType(org, orgData.type);
      const result = await insertOrganization(db, org, type);
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
  console.log(`   📊 Всего: ${ORGANIZATIONS_DATA.length}`);

  db.close();
}

main().catch(console.error);
