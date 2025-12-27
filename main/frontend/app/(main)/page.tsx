'use client';

import { useState } from 'react';
import PostsFeed from '../components/posts/PostsFeed';
import FeedFilters from '../components/posts/FeedFilters';
import RightPanel from '../components/layout/RightPanel';

type FilterType = 'for-you' | 'following' | 'city' | 'lost' | 'found' | 'looking-for-home';

export default function Home() {
  const [activeFilter, setActiveFilter] = useState<FilterType>('for-you');

  return (
    <div className="flex gap-4">
      {/* Main Feed - оптимальная ширина */}
      <div className="w-full xl:w-[600px] xl:flex-shrink-0">
        <PostsFeed activeFilter={activeFilter} />
      </div>

      {/* Right Panel - уменьшенная ширина */}
      <aside className="hidden xl:block w-[320px]">
        <div className="sticky top-[48px] space-y-2.5">
          <FeedFilters activeFilter={activeFilter} onFilterChange={setActiveFilter} />
          <RightPanel />
        </div>
      </aside>
    </div>
  );
}
