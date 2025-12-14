'use client';

export default function RightPanel() {
  return (
    <div className="space-y-2.5">
      {/* Online Friends */}
      <div className="bg-white rounded-xl shadow-sm p-4 border border-gray-100">
        <h3 className="text-sm font-semibold text-gray-800 mb-3 flex items-center gap-2">
          <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
          Друзья онлайн
        </h3>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-3 hover:bg-gray-50 p-2 rounded-lg cursor-pointer transition-all group">
              <div className="relative">
                <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center shadow-sm group-hover:scale-105 transition-transform">
                  <span className="text-sm">👤</span>
                </div>
                <div className="absolute bottom-0 right-0 w-2.5 h-2.5 bg-green-500 rounded-full border-2 border-white"></div>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">Пользователь {i}</p>
                <p className="text-xs text-gray-500">Онлайн</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Suggestions */}
      <div className="bg-white rounded-xl shadow-sm p-4 border border-gray-100">
        <h3 className="text-sm font-semibold text-gray-800 mb-3">Рекомендации</h3>
        <div className="space-y-2">
          {[1, 2].map((i) => (
            <div key={i} className="flex items-center gap-3 hover:bg-gray-50 p-2 rounded-lg cursor-pointer transition-all group">
              <div className="w-10 h-10 bg-gradient-to-br from-pink-400 to-orange-400 rounded-full flex items-center justify-center shadow-sm group-hover:scale-105 transition-transform">
                <span className="text-sm">👥</span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">Сообщество {i}</p>
                <p className="text-xs text-gray-500">1.2K подписчиков</p>
              </div>
              <button 
                className="px-3 py-1 text-white text-xs font-medium rounded-lg transition-colors opacity-0 group-hover:opacity-100"
                style={{ backgroundColor: '#1B76FF' }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = '#0D5FE0';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = '#1B76FF';
                }}
              >
                +
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
