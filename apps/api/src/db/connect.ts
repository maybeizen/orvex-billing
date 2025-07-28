import mongoose, { type ConnectOptions } from "mongoose";

export default class Database {
  private uri: string;
  private options: ConnectOptions;

  public connection?: mongoose.Connection;

  constructor(uri: string, options?: Partial<ConnectOptions>) {
    this.uri = uri;
    this.options = options ?? {};
  }

  public async connect() {
    if (this.connection?.readyState === 1) {
      console.log("Already connected to MongoDB");
      return;
    }

    try {
      const startTime = Date.now();

      console.log("Connecting to MongoDB...");
      await mongoose
        .connect(this.uri, this.options)
        .then(() =>
          console.log(`Connected to MongoDB in ${Date.now() - startTime}ms`)
        );
      this.connection = mongoose.connection;

      this.connection.on("error", (error) => {
        console.error(`MongoDB connection error: ${error}`);
      });
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error(`Database connection failed: ${error.message}`);
      } else {
        console.error(`Database connection failed: ${error}`);
      }
    }
  }

  public async destroy() {
    if (!this.connection || this.connection.readyState !== 1) {
      console.warn("No active database connection to destroy.");
      return;
    }

    try {
      await this.connection.close();
      console.log("Database connection closed.");
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error(`Failed to destroy database instance: ${error.message}`);
      } else {
        console.error(`Failed to destroy database instance: ${error}`);
      }
    }
  }
}
