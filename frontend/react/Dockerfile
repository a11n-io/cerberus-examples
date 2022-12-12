# frontend build stage
FROM node:lts-alpine as frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3001
CMD ["npm", "start"]

