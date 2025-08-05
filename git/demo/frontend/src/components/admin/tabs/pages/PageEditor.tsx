import React from 'react';

const PageEditor = ({ pageId }: { pageId: string }) => {
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Bearbeite: {pageId}</h2>
      <textarea
        className="w-full h-64 border p-2"
        placeholder={`Inhalt der Seite "${pageId}" bearbeiten...`}
      />
    </div>
  );
};

export default PageEditor;
