// Прямой поиск организаций в Сарапуле
const fetch = require('node-fetch');

const DADATA_TOKEN = '300ba9e25ef32f0d6ea7c41826b2255b138e19e2';

async function searchByCity() {
  try {
    const response = await fetch('https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/party', {
      method: 'POST',
      headers: {
        'Authorization': `Token ${DADATA_TOKEN}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: 'Сарапул',
        count: 20,
        status: ['ACTIVE'],
        locations: [{
          kladr_id: '1800000700000' // КЛАДР код Сарапула
        }]
      }),
    });

    const data = await response.json();
    console.log('Найдено организаций:', data.suggestions?.length || 0);
    
    if (data.suggestions && data.suggestions.length > 0) {
      data.suggestions.forEach((org, index) => {
        console.log(`\n${index + 1}. ${org.value}`);
        console.log(`   ИНН: ${org.data.inn}`);
        console.log(`   Адрес: ${org.data.address?.value || 'не указан'}`);
        console.log(`   Город: ${org.data.address?.data?.city || 'не указан'}`);
      });
    }
  } catch (error) {
    console.error('Ошибка:', error.message);
  }
}

searchByCity();
