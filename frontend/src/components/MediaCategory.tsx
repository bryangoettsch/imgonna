import React from 'react';
import type { MediaItem } from '../api/goals';

interface MediaCategoryProps {
  title: string;
  items: MediaItem[];
  icon: React.ReactNode;
}

export const MediaCategory: React.FC<MediaCategoryProps> = ({ title, items, icon }) => {
  if (!items || items.length === 0) {
    return null;
  }

  return (
    <div>
      <div className="flex items-center space-x-3 mb-4">
        <span className="text-gray-600">{icon}</span>
        <h3 className="font-semibold text-gray-900 text-lg">{title}</h3>
        <span className="text-sm text-gray-500">({items.length})</span>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {items.map((item, index) => (
          <div key={index} className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-lg transition-all hover:border-gray-300 min-h-[120px] flex flex-col">
            <h4 className="font-medium text-gray-900 mb-2">
              {item.link ? (
                <a
                  href={item.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:text-blue-600 transition-colors"
                >
                  {item.title}
                </a>
              ) : (
                item.title
              )}
            </h4>
            {item.description && (
              <p className="text-sm text-gray-600 flex-grow">{item.description}</p>
            )}
            {item.platform && (
              <span className="inline-block text-xs text-gray-500 mt-2">
                {item.platform}
              </span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};