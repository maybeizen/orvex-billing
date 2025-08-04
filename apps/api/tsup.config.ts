import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/index.ts"],
  outDir: "dist",
  target: "es2020",
  format: ["esm"],
  splitting: false,
  sourcemap: true,
  clean: true,
  dts: true,
  minify: false,
  skipNodeModulesBundle: true,
  tsconfig: "./tsconfig.json",
  onSuccess: "node dist/index.js",
});
