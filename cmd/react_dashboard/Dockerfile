FROM node:20.10.0
WORKDIR /app
COPY cmd/react_dashboard/package*.json ./
RUN npm ci
COPY cmd/react_dashboard ./
EXPOSE 3000
CMD ["npm", "run", "start", "--", "--host", "0.0.0.0"]
