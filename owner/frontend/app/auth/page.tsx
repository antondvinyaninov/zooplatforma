'use client';

export default function AuthPage() {
  const handleLogin = () => {
    // Редирект на основной сайт для авторизации
    window.location.href = 'http://localhost:3000/auth?redirect=http://localhost:7000';
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <div className="text-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Кабинет владельца
          </h1>
          <p className="text-gray-600">
            Управление питомцами
          </p>
        </div>

        <div className="space-y-4">
          <button
            onClick={handleLogin}
            className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
          >
            Войти через ЗооПлатформу
          </button>

          <p className="text-sm text-gray-500 text-center">
            Используется единая авторизация (SSO)
          </p>
        </div>
      </div>
    </div>
  );
}
