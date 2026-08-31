# Stage 1: Build the Nuxt application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install all dependencies
RUN npm install

# Copy the rest of the application files
COPY . .

# Build the Nuxt application
RUN npm run build

# Stage 2: Production runtime container
FROM node:20-alpine AS runner

# Create a non-root system group and user for container security hardening (SonarQube)
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy built application from build stage
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package.json /app/package.json

# Set ownership to non-root user
RUN chown -R appuser:appgroup /app

# Run container under non-root context
USER appuser

# Expose default port
EXPOSE 3000

ENV PORT=3000
ENV HOST=0.0.0.0
ENV NODE_ENV=production

# Run the Nuxt application using the built server
CMD ["node", ".output/server/index.mjs"]
