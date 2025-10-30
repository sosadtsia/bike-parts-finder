# Build stage
FROM node:25.1.0-alpine AS build

WORKDIR /app

# Copy only package files first to leverage Docker cache
COPY package.json .npmrc ./

# Install dependencies
RUN npm install

# Copy the rest of the app and build
COPY . .
ENV NODE_OPTIONS=--no-deprecation
RUN npm run build

# Production stage
FROM node:25.1.0-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

# Copy necessary files from build stage
COPY --from=build /app/next.config.js ./
COPY --from=build /app/public ./public
COPY --from=build /app/package.json ./package.json

# Copy the standalone output and static files
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static

# Expose port
EXPOSE 3000

# Start the app
CMD ["node", "server.js"]
