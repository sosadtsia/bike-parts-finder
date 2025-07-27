import React from 'react';
import { Link } from 'react-router-dom';

function Header() {
  return (
    <header className="bg-blue-600 shadow-md">
      <div className="container mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/" className="text-white text-2xl font-bold">Bike Parts Finder</Link>
        <nav>
          <ul className="flex space-x-6 text-white">
            <li><Link to="/" className="hover:text-blue-200">Home</Link></li>
            <li><Link to="/" className="hover:text-blue-200">Search</Link></li>
            <li><Link to="/" className="hover:text-blue-200">Categories</Link></li>
            <li><Link to="/" className="hover:text-blue-200">About</Link></li>
          </ul>
        </nav>
      </div>
    </header>
  );
}

export default Header;
