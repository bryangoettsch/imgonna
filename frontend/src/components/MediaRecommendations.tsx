import React from 'react';
import { 
  MicrophoneIcon, 
  PlayCircleIcon, 
  BookOpenIcon, 
  GlobeAltIcon 
} from '@heroicons/react/24/outline';
import { MediaCategory } from './MediaCategory';
import type { MediaRecommendations as MediaRecommendationsType } from '../api/goals';

interface MediaRecommendationsProps {
  recommendations: MediaRecommendationsType;
}

export const MediaRecommendations: React.FC<MediaRecommendationsProps> = ({ recommendations }) => {
  if (!recommendations) {
    return null;
  }

  const hasAnyRecommendations = 
    (recommendations.podcasts?.length > 0) ||
    (recommendations.streaming?.length > 0) ||
    (recommendations.books?.length > 0) ||
    (recommendations.websites?.length > 0);

  if (!hasAnyRecommendations) {
    return null;
  }

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-semibold text-gray-900 mb-6">
        Media Recommendations
      </h2>
      <div className="space-y-8">
        <MediaCategory
          title="Podcasts"
          items={recommendations.podcasts || []}
          icon={<MicrophoneIcon className="h-5 w-5" />}
        />
        <MediaCategory
          title="Streaming Media"
          items={recommendations.streaming || []}
          icon={<PlayCircleIcon className="h-5 w-5" />}
        />
        <MediaCategory
          title="Books"
          items={recommendations.books || []}
          icon={<BookOpenIcon className="h-5 w-5" />}
        />
        <MediaCategory
          title="Websites"
          items={recommendations.websites || []}
          icon={<GlobeAltIcon className="h-5 w-5" />}
        />
      </div>
    </div>
  );
};