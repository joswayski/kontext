FROM node:22-alpine AS development-dependencies-env
COPY web/ /app
WORKDIR /app
RUN npm ci

FROM node:22-alpine AS production-dependencies-env
COPY web/package.json web/package-lock.json /app/
WORKDIR /app
RUN npm ci --omit=dev

FROM node:22-alpine AS build-env
COPY web/ /app/
COPY --from=development-dependencies-env /app/node_modules /app/node_modules
WORKDIR /app
RUN npm run build

FROM node:22-alpine
COPY web/package.json web/package-lock.json /app/
COPY --from=production-dependencies-env /app/node_modules /app/node_modules
COPY --from=build-env /app/build /app/build
WORKDIR /app
CMD ["npm", "run", "start"]
