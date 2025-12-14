'use client';

import { useState } from 'react';
import Feed from '../components/layout/Feed';
import RightPanel from '../components/layout/RightPanel';
import FeedFilters from '../components/shared/FeedFilters';

type PostType = 'all' | 'post' | 'sale' | 'lost' | 'found' | 'help';

export default function Home() {
  const [activeFilter, setActiveFilter] = useState<PostType>('all');

  return (
    <div className="flex gap-2.5">
      {/* Main Feed */}
      <div className="w-full xl:w-[520px] xl:flex-shrink-0">
        <Feed activeFilter={activeFilter} />
      </div>

      {/* Right Panel */}
      <aside className="hidden xl:block flex-1">
        <div className="sticky top-[48px] space-y-2.5">
          <FeedFilters activeFilter={activeFilter} onFilterChange={setActiveFilter} />
          <RightPanel />
        </div>
      </aside>
    </div>
  );
}
