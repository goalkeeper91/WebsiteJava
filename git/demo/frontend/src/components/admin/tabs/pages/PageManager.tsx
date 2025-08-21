import { useState } from 'react';
import PageSidebar from './PageSidebar';
import PageEditor from './PageEditor';

const PageManager = () => {
  const [selectedPageId, setSelectedPageId] = useState<string | null>(null);

  return (
    <div className="flex h-full">
      <div className="w-1/4 border-r bg-gray-800">
        <PageSidebar onSelectPage={setSelectedPageId} selectedPageId={selectedPageId} />
      </div>
      <div className="flex-1 p-4">
        {selectedPageId ? (
          <PageEditor pageId={selectedPageId} />
        ) : (
          <div className="text-gray-500">Bitte eine Seite auswählen.</div>
        )}
      </div>
    </div>
  );
};

export default PageManager;