import Header from '../components/layout/Header';
import Sidebar from '../components/layout/Sidebar';

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <Header />
      <div className="min-h-screen">
        <div className="max-w-[1150px] mx-auto flex gap-2.5 pt-2.5 px-2 sm:px-4">
          {/* Left Sidebar */}
          <aside className="hidden lg:block w-[180px] flex-shrink-0">
            <Sidebar />
          </aside>

          {/* Main Content */}
          <main className="flex-1 min-w-0">
            {children}
          </main>
        </div>
      </div>
    </>
  );
}
