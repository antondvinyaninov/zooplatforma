import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Кабинет приюта - ЗооПлатформа",
  description: "Управление приютом для животных",
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
