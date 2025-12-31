'use client';

import { ReactNode } from 'react';

interface TableWidgetProps {
  title: string;
  actions?: ReactNode;
  children: ReactNode;
}

export default function TableWidget({ title, actions, children }: TableWidgetProps) {
  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
        <h3 className="text-xl font-bold text-gray-900">{title}</h3>
        {actions && <div>{actions}</div>}
      </div>
      {children}
    </div>
  );
}
