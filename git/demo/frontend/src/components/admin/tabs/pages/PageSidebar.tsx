const pages = [
  { id: 'home', name: 'Startseite' },
  { id: 'about', name: 'Über uns' },
  { id: 'all videos', name: 'Alle Videos' },
];

const PageSidebar = ({ selectedPageId, onSelectPage }: {
  selectedPageId: string | null;
  onSelectPage: (id: string) => void;
}) => {
  return (
    <ul className="p-4 space-y-2">
      {pages.map((page) => (
        <li
          key={page.id}
          className={`cursor-pointer p-2 rounded ${
            selectedPageId === page.id ? 'bg-blue-500 text-white' : 'hover:bg-blue-700'
          }`}
          onClick={() => onSelectPage(page.id)}
        >
          {page.name}
        </li>
      ))}
    </ul>
  );
};

export default PageSidebar;
