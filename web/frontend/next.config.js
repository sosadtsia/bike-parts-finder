/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Configure output directory to match CRA's build folder if needed
  distDir: 'build',
  // Ensure images from external sources can be loaded
  images: {
    domains: [],
  },
  // Add output configuration for standalone mode (for Docker)
  output: 'standalone',
};

module.exports = nextConfig;
