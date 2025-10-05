# Dockerfile for React frontend with Bun
FROM oven/bun:1

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Expose port 3000
EXPOSE 3000

# Run the dev server
CMD ["bun", "run", "dev", "--host", "0.0.0.0"]
