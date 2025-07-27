import React, { useState, useEffect } from 'react';
import SearchForm from '../components/SearchForm';
import PartsList from '../components/PartsList';
import { searchParts } from '../services/api';

function HomePage() {
  const [parts, setParts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchParams, setSearchParams] = useState({ query: '', category: '' });

  useEffect(() => {
    // Load initial results
    fetchParts();
  }, []);

  const fetchParts = async (query = '', category = '') => {
    setLoading(true);
    setError(null);

    try {
      const data = await searchParts(query, category);
      setParts(data);
    } catch (err) {
      setError(err.message || 'Failed to fetch parts');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (query, category) => {
    setSearchParams({ query, category });
    fetchParts(query, category);
  };

  return (
    <>
      <SearchForm onSearch={handleSearch} />

      <section id="results-section" className="mb-8">
        <h2 className="text-xl font-bold mb-4">Results</h2>
        <PartsList parts={parts} loading={loading} error={error} />
      </section>
    </>
  );
}

export default HomePage;
