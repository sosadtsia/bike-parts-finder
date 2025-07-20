import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate, Link } from 'react-router-dom';
import { searchParts, getCategories, getBrands } from '../services/api';

const Search = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);

  const [searchResults, setSearchResults] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [categories, setCategories] = useState([]);
  const [brands, setBrands] = useState([]);

  // Filter states
  const [filters, setFilters] = useState({
    q: queryParams.get('q') || '',
    category: queryParams.get('category') || '',
    brand: queryParams.get('brand') || '',
    minPrice: queryParams.get('minPrice') || '',
    maxPrice: queryParams.get('maxPrice') || '',
    sortBy: queryParams.get('sortBy') || 'relevance',
  });

  // Pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const pageSize = 10;

  // Load filters data
  useEffect(() => {
    const loadFilterData = async () => {
      try {
        const [categoriesData, brandsData] = await Promise.all([
          getCategories(),
          getBrands(),
        ]);
        setCategories(categoriesData);
        setBrands(brandsData);
      } catch (err) {
        console.error('Failed to load filter data:', err);
      }
    };

    loadFilterData();
  }, []);

  // Search for parts when filters change or on initial load
  useEffect(() => {
    const fetchSearchResults = async () => {
      setLoading(true);
      setError(null);

      try {
        const searchParams = {
          ...filters,
          page: currentPage,
          pageSize,
        };

        // Remove empty filters
        Object.keys(searchParams).forEach(
          key => searchParams[key] === '' && delete searchParams[key]
        );

        const results = await searchParts(searchParams);
        setSearchResults(results.items || []);
        setTotalPages(Math.ceil(results.total / pageSize) || 1);
      } catch (err) {
        setError('Failed to fetch search results. Please try again later.');
        console.error('Search error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchSearchResults();
  }, [filters, currentPage]);

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams();

    Object.entries(filters).forEach(([key, value]) => {
      if (value) {
        params.set(key, value);
      }
    });

    if (currentPage > 1) {
      params.set('page', currentPage);
    }

    navigate({ search: params.toString() }, { replace: true });
  }, [filters, currentPage, navigate]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
    setCurrentPage(1); // Reset to first page when filters change
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
    window.scrollTo(0, 0);
  };

  return (
    <div className="flex flex-col md:flex-row gap-6">
      {/* Filters Sidebar */}
      <aside className="md:w-64 bg-white p-6 rounded-lg shadow-md self-start">
        <h2 className="text-xl font-bold mb-4">Filters</h2>

        <form className="space-y-4">
          <div>
            <label htmlFor="q" className="block text-gray-700 mb-1">Search</label>
            <input
              type="text"
              id="q"
              name="q"
              value={filters.q}
              onChange={handleFilterChange}
              className="input w-full"
              placeholder="Search terms..."
            />
          </div>

          <div>
            <label htmlFor="category" className="block text-gray-700 mb-1">Category</label>
            <select
              id="category"
              name="category"
              value={filters.category}
              onChange={handleFilterChange}
              className="input w-full"
            >
              <option value="">All Categories</option>
              {categories.map(category => (
                <option key={category.id} value={category.id}>{category.name}</option>
              ))}
            </select>
          </div>

          <div>
            <label htmlFor="brand" className="block text-gray-700 mb-1">Brand</label>
            <select
              id="brand"
              name="brand"
              value={filters.brand}
              onChange={handleFilterChange}
              className="input w-full"
            >
              <option value="">All Brands</option>
              {brands.map(brand => (
                <option key={brand.id} value={brand.id}>{brand.name}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-gray-700 mb-1">Price Range</label>
            <div className="flex items-center gap-2">
              <input
                type="number"
                name="minPrice"
                value={filters.minPrice}
                onChange={handleFilterChange}
                className="input w-full"
                placeholder="Min"
                min="0"
              />
              <span>-</span>
              <input
                type="number"
                name="maxPrice"
                value={filters.maxPrice}
                onChange={handleFilterChange}
                className="input w-full"
                placeholder="Max"
                min="0"
              />
            </div>
          </div>

          <div>
            <label htmlFor="sortBy" className="block text-gray-700 mb-1">Sort By</label>
            <select
              id="sortBy"
              name="sortBy"
              value={filters.sortBy}
              onChange={handleFilterChange}
              className="input w-full"
            >
              <option value="relevance">Relevance</option>
              <option value="price_low">Price: Low to High</option>
              <option value="price_high">Price: High to Low</option>
              <option value="newest">Newest First</option>
            </select>
          </div>

          <button
            type="button"
            onClick={() => {
              setFilters({
                q: '',
                category: '',
                brand: '',
                minPrice: '',
                maxPrice: '',
                sortBy: 'relevance',
              });
              setCurrentPage(1);
            }}
            className="btn btn-outline w-full"
          >
            Reset Filters
          </button>
        </form>
      </aside>

      {/* Search Results */}
      <div className="flex-1">
        <h1 className="text-2xl font-bold mb-4">
          {filters.q
            ? `Search Results for "${filters.q}"`
            : filters.category
              ? `${filters.category} Parts`
              : 'All Bike Parts'}
        </h1>

        {loading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
            <p className="mt-2 text-gray-600">Loading results...</p>
          </div>
        ) : error ? (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            <p>{error}</p>
          </div>
        ) : searchResults.length === 0 ? (
          <div className="bg-gray-100 p-8 rounded-lg text-center">
            <h2 className="text-xl font-semibold mb-2">No results found</h2>
            <p className="text-gray-600 mb-4">Try adjusting your search or filter criteria</p>
            <button
              onClick={() => setFilters({
                q: '',
                category: '',
                brand: '',
                minPrice: '',
                maxPrice: '',
                sortBy: 'relevance',
              })}
              className="btn btn-primary"
            >
              Clear Filters
            </button>
          </div>
        ) : (
          <>
            <p className="text-gray-600 mb-4">
              Showing {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, searchResults.length * currentPage)} of {searchResults.length * totalPages} results
            </p>

            <div className="space-y-6">
              {searchResults.map(part => (
                <div key={part.id} className="card hover:shadow-lg transition-shadow">
                  <div className="flex flex-col md:flex-row gap-4">
                    <div className="w-full md:w-48 h-48 bg-gray-200 rounded-md flex items-center justify-center">
                      {part.imageUrl ? (
                        <img
                          src={part.imageUrl}
                          alt={part.name}
                          className="max-w-full max-h-full object-contain"
                        />
                      ) : (
                        <span className="text-gray-400">No image</span>
                      )}
                    </div>

                    <div className="flex-1">
                      <h2 className="text-xl font-semibold mb-2">
                        <Link to={`/parts/${part.id}`} className="text-primary-600 hover:underline">
                          {part.name}
                        </Link>
                      </h2>

                      <div className="mb-2 flex items-center">
                        <span className="bg-gray-200 text-gray-700 px-2 py-1 rounded text-sm mr-2">
                          {part.brand}
                        </span>
                        <span className="text-gray-500">{part.category}</span>
                      </div>

                      <p className="text-gray-700 mb-4 line-clamp-2">
                        {part.description || 'No description available'}
                      </p>

                      <div className="flex flex-wrap justify-between items-center">
                        <div>
                          <p className="text-2xl font-bold text-primary-600">
                            ${part.price.toFixed(2)}
                          </p>
                          {part.originalPrice && part.originalPrice > part.price && (
                            <p className="text-sm text-gray-500 line-through">
                              ${part.originalPrice.toFixed(2)}
                            </p>
                          )}
                        </div>

                        <Link to={`/parts/${part.id}`} className="btn btn-primary">
                          View Details
                        </Link>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center mt-8">
                <nav className="flex items-center">
                  <button
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="btn btn-outline px-3 py-1 mr-2 disabled:opacity-50"
                  >
                    Previous
                  </button>

                  {Array.from({ length: totalPages }, (_, i) => i + 1)
                    .filter(page =>
                      page === 1 ||
                      page === totalPages ||
                      Math.abs(page - currentPage) < 2
                    )
                    .map((page, index, array) => (
                      <React.Fragment key={page}>
                        {index > 0 && array[index - 1] !== page - 1 && (
                          <span className="mx-1">...</span>
                        )}
                        <button
                          onClick={() => handlePageChange(page)}
                          className={`w-8 h-8 flex items-center justify-center rounded-full mx-1 ${
                            currentPage === page
                              ? 'bg-primary-500 text-white'
                              : 'bg-gray-200 hover:bg-gray-300 text-gray-800'
                          }`}
                        >
                          {page}
                        </button>
                      </React.Fragment>
                    ))}

                  <button
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                    className="btn btn-outline px-3 py-1 ml-2 disabled:opacity-50"
                  >
                    Next
                  </button>
                </nav>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default Search;
