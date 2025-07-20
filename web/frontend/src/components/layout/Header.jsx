import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const navigate = useNavigate();

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchTerm.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchTerm)}`);
      setSearchTerm('');
    }
  };

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4">
        <div className="flex flex-wrap items-center justify-between py-4">
          <div className="flex items-center">
            <Link to="/" className="flex items-center">
              <span className="text-2xl font-bold text-primary-500">Bike Parts Finder</span>
            </Link>
          </div>

          <div className="lg:hidden">
            <button
              onClick={toggleMenu}
              className="text-gray-600 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary-300 rounded p-2"
            >
              <svg className="h-6 w-6" fill="none" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" viewBox="0 0 24 24" stroke="currentColor">
                {isMenuOpen ? (
                  <path d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>
          </div>

          <div className={`w-full lg:flex lg:w-auto ${isMenuOpen ? 'block' : 'hidden'} mt-4 lg:mt-0`}>
            <form onSubmit={handleSearch} className="flex mb-4 lg:mb-0 lg:mr-4">
              <input
                type="text"
                placeholder="Search parts..."
                className="input flex-grow"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
              <button type="submit" className="btn btn-primary ml-2">
                Search
              </button>
            </form>

            <nav className="flex flex-col lg:flex-row">
              <Link to="/" className="py-2 lg:px-4 text-gray-700 hover:text-primary-500">Home</Link>
              <Link to="/search" className="py-2 lg:px-4 text-gray-700 hover:text-primary-500">Browse</Link>
              <Link to="/about" className="py-2 lg:px-4 text-gray-700 hover:text-primary-500">About</Link>
            </nav>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
