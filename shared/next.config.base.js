// Базовая конфигурация Next.js для всех frontend сервисов
// Совместимость с Next.js 16 (Turbopack по умолчанию)

module.exports = {
  // Пустой turbopack config для совместимости с Next.js 16
  turbopack: {},
  
  // Не блокируем dev сервер при ошибках TypeScript
  typescript: {
    ignoreBuildErrors: process.env.NODE_ENV === 'development',
  },
};
