# Build stage
FROM node:24.4.1-alpine AS build

WORKDIR /app

# Copy only package files first to leverage Docker cache
COPY package.json .npmrc ./

# Install dependencies
RUN npm install

# Copy the rest of the app and build
COPY . .
ENV NODE_OPTIONS=--no-deprecation
ENV DISABLE_ESLINT_PLUGIN=true
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built assets from the build stage
COPY --from=build /app/build /usr/share/nginx/html

# Add nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
