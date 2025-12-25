import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../backend/docs/swagger.json', // Path to Go docs
  output: 'src/client', // Where types will be generated
  client: '@hey-api/client-axios',
});
