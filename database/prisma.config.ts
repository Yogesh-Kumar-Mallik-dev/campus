import { defineConfig, env } from 'prisma/config';

export default defineConfig({
  schema: 'schema.prisma',
  datasource: {
    url: process.env.DATABASE_URL || 'postgresql://campus:campus_secret@localhost:5432/campus_db?schema=public',
  },
});
