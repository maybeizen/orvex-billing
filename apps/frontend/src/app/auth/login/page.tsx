"use client";

import axios from "axios";
import React, { useState } from "react";
import { Button } from "@/components/ui/button";

export default function Login() {
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    try {
      e.preventDefault();
      setError(null);

      setLoading(true);
      const res = await axios.post(
        "http://localhost:3001/api/v1/auth/login",
        { email, password },
        { withCredentials: true }
      );

      setLoading(false);
      console.log(`Signed in! ${res.data}`);
    } catch (error: any) {
      console.error(error);
      setError(error.response?.data?.message || "Something went wrong");
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-950 w-full">
      <form
        onSubmit={handleSubmit}
        className="max-w-md w-full bg-neutral-900 border border-neutral-800 p-8 rounded-lg flex flex-col gap-4"
      >
        <h1 className="text-xl font-semibold text-white mb-4 text-center">
          Login
        </h1>
        {error && <p className="text-red-500 mb-2 text-center">{error}</p>}
        <div className="flex flex-col gap-1">
          <label htmlFor="email" className="text-white text-sm">
            Email address
          </label>
          <input
            id="email"
            name="email"
            type="email"
            placeholder="example@email.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="p-2 rounded bg-neutral-800 text-white border border-neutral-700 focus:outline-none focus:ring-2 focus:ring-violet-500"
            autoComplete="email"
          />
        </div>
        <div className="flex flex-col gap-1">
          <div className="flex justify-between">
            <label htmlFor="password" className="text-white text-sm">
              Password
            </label>
            <a
              href="/auth/forgot-password"
              className="text-indigo-500 text-sm hover:underline hover:cursor-pointer"
            >
              Forgot Password
            </a>
          </div>
          <div className="flex items-center gap-2 bg-neutral-800 border border-neutral-700 rounded">
            <input
              id="password"
              name="password"
              placeholder="Enter a password"
              type={showPassword ? "text" : "password"}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="p-2 w-full rounded bg-neutral-800 text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
              autoComplete="current-password"
            />
            <button
              type="button"
              onClick={() => {
                setShowPassword(showPassword ? false : true);
              }}
              className="w-10 h-10 flex items-center justify-center rounded bg-neutral-800 hover:bg-neutral-700 transition"
            >
              <i
                className={`fas ${
                  showPassword ? "fa-eye-slash" : "fa-eye"
                } text-neutral-300 text-lg`}
              />
            </button>
          </div>
        </div>

        <Button
          variant="primary"
          type="submit"
          loading={loading}
          fullWidth
          icon="fas fa-user"
        >
          {loading ? "Logging in..." : "Submit"}
        </Button>

        <div className="flex flex-col items-center justify-center py-6 w-full">
          <div className="flex items-center w-full max-w-md mb-6">
            <span className="flex-1 h-px bg-neutral-700" />
            <p className="mx-4 text-xs text-neutral-500 uppercase tracking-wider whitespace-nowrap">
              or continue with
            </p>
            <span className="flex-1 h-px bg-neutral-700" />
          </div>
          <div className="flex flex-row justify-center items-center gap-4 w-full max-w-md">
            <Button
              variant="glass"
              icon="fab fa-google"
              iconPosition="right"
              fullWidth
            >
              Google
            </Button>
            <Button
              variant="glass"
              icon="fab fa-discord"
              iconPosition="right"
              fullWidth
            >
              Discord
            </Button>
          </div>
        </div>
      </form>
    </div>
  );
}
