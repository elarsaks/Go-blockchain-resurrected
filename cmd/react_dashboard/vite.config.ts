import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          three: ["three"],
          vendor: ["axios", "react", "react-dom", "styled-components"],
        },
      },
    },
  },
  resolve: {
    alias: {
      api: path.resolve(__dirname, "src/api"),
      assets: path.resolve(__dirname, "src/assets"),
      components: path.resolve(__dirname, "src/components"),
      store: path.resolve(__dirname, "src/store"),
      utils: path.resolve(__dirname, "src/utils"),
    },
  },
  test: {
    environment: "node",
    include: ["src/tests/**/*.test.ts"],
  },
});
