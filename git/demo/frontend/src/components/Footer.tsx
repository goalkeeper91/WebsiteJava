import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-800 text-gray-300 p-4 mt-auto">
      <div className="container mx-auto text-center">
        &copy; {new Date().getFullYear()} Streamer Website. All rights reserved.
      </div>
    </footer>
  );
};

export default Footer;
