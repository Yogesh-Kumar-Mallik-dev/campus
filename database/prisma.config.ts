import { defineConfig } from 'prisma/config';

export default defineConfig({
  schema: 'schema.prisma',
  datasource: {
    url: process.env.DATABASE_URL || 'postgresql://campus_admin:campus_secret_dev_pass_2026@localhost:5432/campus_core?schema=public',
  },
});
