import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Кабинет владельца - ЗооПлатформа",
  description: "Управление питомцами",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
      <body className="antialiased bg-gray-50">
        {children}
      </body>
    </html>
  );
}
