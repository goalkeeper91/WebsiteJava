import { useState } from 'react';

const PageEditor = ({ pageId }: { pageId: string }) => {

  const renderForm = () => {
    switch (pageId) {
      default:
        return <p>Kein Editor für diese Seite vorhanden.</p>;
    }
  };

  return (
    <div className="p-4">
      <h2 className="text-3xl font-bold mb-6">Bearbeite Seite: {pageId}</h2>
      {renderForm()}
    </div>
  );
};

export default PageEditor;
