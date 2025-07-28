import app from "./app";
import mongoose from "mongoose";
import { env } from "./utils/env";

const PORT = Number(env.PORT) || 3001;

async function startServer() {
  try {
    app.listen(PORT, () => {
      console.log(`🚀 Server is running at http://localhost:${PORT}`);
    });
  } catch (error) {
    console.error("❌ Failed to start server:", error);
    process.exit(1);
  }
}

startServer();
