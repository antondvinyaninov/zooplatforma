import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  
  // Пустой turbopack config для совместимости с Next.js 16
  turbopack: {},
  
  // Отключаем некоторые проверки для ускорения
  typescript: {
    // Не блокируем сборку при ошибках TypeScript в dev режиме
    ignoreBuildErrors: process.env.NODE_ENV === 'development',
  },
};

export default nextConfig;
