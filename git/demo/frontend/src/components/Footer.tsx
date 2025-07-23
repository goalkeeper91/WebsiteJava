import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-800 text-gray-300 w-full">
      <div className="w-full text-center py-4 px-6">
        &copy; {new Date().getFullYear()} Streamer Website. All rights reserved.
      </div>
    </footer>
  );
};

export default Footer;
