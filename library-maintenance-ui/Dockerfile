# ---------- Build Stage ----------
FROM node:20-alpine AS builder

WORKDIR /app

COPY library-maintenance-api/package*.json ./
RUN npm install

COPY library-maintenance-ui/. .
RUN npm run build


# ---------- Runtime Stage ----------
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
