FROM oven/bun:latest AS builder

WORKDIR /app
COPY bun.lockb package.json ./
RUN bun install
COPY . .

ARG API_URL
ENV VITE_API_URL=${API_URL}
RUN bun run build

FROM nginx:alpine AS runner

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
