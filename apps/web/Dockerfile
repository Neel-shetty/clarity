# Build stage
FROM node:24-alpine3.21 AS builder
WORKDIR /app

# Install dependencies
COPY package.json pnpm-lock.yaml tsconfig.json ./
RUN npm install -g pnpm && \
    pnpm install

# Copy source and build
COPY . .
RUN pnpm run build

# Production stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
