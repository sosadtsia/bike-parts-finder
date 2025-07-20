import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const searchParts = async (params) => {
  try {
    const response = await api.get('/api/parts/search', { params });
    return response.data;
  } catch (error) {
    console.error('Error searching parts:', error);
    throw error;
  }
};

export const getPartById = async (id) => {
  try {
    const response = await api.get(`/api/parts/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching part ${id}:`, error);
    throw error;
  }
};

export const requestScrapingJob = async (url) => {
  try {
    const response = await api.post('/api/scrape', { url });
    return response.data;
  } catch (error) {
    console.error('Error requesting scraping job:', error);
    throw error;
  }
};

export const getCategories = async () => {
  try {
    const response = await api.get('/api/categories');
    return response.data;
  } catch (error) {
    console.error('Error fetching categories:', error);
    throw error;
  }
};

export const getBrands = async () => {
  try {
    const response = await api.get('/api/brands');
    return response.data;
  } catch (error) {
    console.error('Error fetching brands:', error);
    throw error;
  }
};

export default api;
