import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Home = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const navigate = useNavigate();

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchTerm.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchTerm)}`);
    }
  };

  return (
    <div className="flex flex-col items-center">
      <section className="w-full py-16 bg-primary-500 text-white text-center rounded-lg mb-12">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl md:text-5xl font-bold mb-6">
            Find the Perfect Bike Parts
          </h1>
          <p className="text-xl mb-8 max-w-3xl mx-auto">
            Compare prices and availability across multiple retailers to get the best deals on bicycle components.
          </p>
          <form onSubmit={handleSearch} className="max-w-xl mx-auto flex flex-col md:flex-row gap-4">
            <input
              type="text"
              placeholder="Search for bike parts (e.g., 'Shimano XT brakes')"
              className="input flex-grow py-3 px-4 text-gray-800 text-lg"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <button type="submit" className="btn btn-secondary py-3 px-8 text-lg">
              Search
            </button>
          </form>
        </div>
      </section>

      <section className="w-full mb-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8 text-center">Popular Categories</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {['Drivetrain', 'Wheels & Tires', 'Brakes', 'Suspension', 'Handlebars', 'Saddles', 'Pedals', 'Accessories'].map((category) => (
              <div
                key={category}
                className="card cursor-pointer hover:shadow-lg transition-shadow"
                onClick={() => navigate(`/search?category=${encodeURIComponent(category)}`)}
              >
                <h3 className="text-xl font-semibold mb-2">{category}</h3>
                <p className="text-gray-600">Browse {category.toLowerCase()} components</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="w-full mb-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8 text-center">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="card text-center">
              <div className="text-primary-500 text-5xl mb-4">1</div>
              <h3 className="text-xl font-semibold mb-2">Search</h3>
              <p className="text-gray-600">
                Enter the bike parts you're looking for or browse by category.
              </p>
            </div>
            <div className="card text-center">
              <div className="text-primary-500 text-5xl mb-4">2</div>
              <h3 className="text-xl font-semibold mb-2">Compare</h3>
              <p className="text-gray-600">
                View prices, availability, and specifications from multiple sources.
              </p>
            </div>
            <div className="card text-center">
              <div className="text-primary-500 text-5xl mb-4">3</div>
              <h3 className="text-xl font-semibold mb-2">Save</h3>
              <p className="text-gray-600">
                Get the best deal by finding the lowest price or closest retailer.
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Home;
