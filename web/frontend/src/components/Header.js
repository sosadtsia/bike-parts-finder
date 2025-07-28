import React from 'react';
import { Link } from 'react-router-dom';

function Header() {
  return (
    <header className="bg-green-600 shadow-md">
      <div className="container mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/" className="text-white text-2xl font-bold">Bike Parts Finder</Link>
        <nav>
          <ul className="flex space-x-6 text-white">
            <li><Link to="/" className="hover:text-green-100">Home</Link></li>
            <li><Link to="/" className="hover:text-green-100">Search</Link></li>
            <li><Link to="/" className="hover:text-green-100">Categories</Link></li>
            <li><Link to="/" className="hover:text-green-100">About</Link></li>
          </ul>
        </nav>
      </div>
    </header>
  );
}

export default Header;
