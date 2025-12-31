import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Кабинет ветеринарной клиники - ЗооПлатформа",
  description: "Управление ветеринарной клиникой",
  icons: {
    icon: '/favicon.svg',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
      <head>
        <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
      </head>
      <body className="antialiased bg-gray-50">
        {children}
      </body>
    </html>
  );
}
