/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Ensure images from external sources can be loaded
  images: {
    domains: [],
  },
  // Add output configuration for standalone mode (for Docker)
  output: 'standalone',
};

module.exports = nextConfig;
