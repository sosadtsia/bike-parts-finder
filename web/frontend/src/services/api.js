import axios from 'axios';

// Use environment variable with fallback
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const searchParts = async (query = '', category = '') => {
  try {
    const params = new URLSearchParams();
    if (query) params.append('q', query);
    if (category) params.append('category', category);

    const response = await api.get(`/parts/search?${params.toString()}`);
    return response.data;
  } catch (error) {
    throw handleError(error);
  }
};

export const getPartById = async (partId) => {
  try {
    const response = await api.get(`/parts/${partId}`);
    return response.data;
  } catch (error) {
    throw handleError(error);
  }
};

const handleError = (error) => {
  if (error.response) {
    // The request was made and the server responded with a status code
    // that falls out of the range of 2xx
    const { data, status } = error.response;
    return {
      message: data.message || 'An error occurred with the API response',
      code: status,
      error: data.error || 'API Error',
    };
  } else if (error.request) {
    // The request was made but no response was received
    return {
      message: 'No response received from server. Please check your connection.',
      code: 0,
      error: 'Network Error',
    };
  } else {
    // Something happened in setting up the request that triggered an Error
    return {
      message: error.message || 'An unexpected error occurred',
      code: 0,
      error: 'Request Error',
    };
  }
};

const apiService = {
  searchParts,
  getPartById,
};

export default apiService;
