FROM node:latest

# Expose default port
EXPOSE 3000

# Workdir
WORKDIR /opt/idlemon

# Copy project
COPY . .

# Npm install
RUN npm i

# TypeScript compiler
RUN npm run tsc

# Start app
ENTRYPOINT npm start
