"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { AlertCircle, Mail, Lock, User, ArrowRight } from "lucide-react";

const SleekAuthForm: React.FC = () => {
  const [isSignUp, setIsSignUp] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    setError("Invalid email or password");
  };

  const toggleAuthMode = () => {
    setIsSignUp(!isSignUp);
    setError("");
  };

  const formVariants = {
    hidden: { opacity: 0, y: 50 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.6, ease: "easeOut" },
    },
    exit: { opacity: 0, y: -50, transition: { duration: 0.4, ease: "easeIn" } },
  };

  const inputVariants = {
    hidden: { opacity: 0, x: -50 },
    visible: { opacity: 1, x: 0, transition: { duration: 0.5 } },
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
      <motion.div
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md relative"
      >
        <AnimatePresence mode="wait">
          <motion.form
            key={isSignUp ? "signup" : "signin"}
            variants={formVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            onSubmit={handleSubmit}
            className="bg-gray-800/30 backdrop-blur-xl p-8 rounded-lg shadow-xl border border-gray-700 overflow-hidden"
          >
            <h2 className="text-3xl font-bold mb-2 text-white">
              {isSignUp ? "Create Account" : "Welcome Back"}
            </h2>
            <p className="text-gray-300 mb-6">
              {isSignUp ? "Sign up to get started" : "Sign in to your account"}
            </p>

            <div className="space-y-4">
              <AnimatePresence>
                {isSignUp && (
                  <motion.div
                    variants={inputVariants}
                    initial="hidden"
                    animate="visible"
                    exit="hidden"
                  >
                    <Label htmlFor="name" className="text-gray-300 block mb-1">
                      Name
                    </Label>
                    <div className="relative">
                      <Input
                        id="name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        className="w-full pl-10 pr-3 py-2 bg-gray-700/50 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition duration-200"
                        placeholder="John Doe"
                      />
                      <User
                        className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                        size={18}
                      />
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>

              <motion.div variants={inputVariants}>
                <Label htmlFor="email" className="text-gray-300 block mb-1">
                  Email
                </Label>
                <div className="relative">
                  <Input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full pl-10 pr-3 py-2 bg-gray-700/50 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition duration-200"
                    placeholder="john@example.com"
                  />
                  <Mail
                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                    size={18}
                  />
                </div>
              </motion.div>

              <motion.div variants={inputVariants}>
                <Label htmlFor="password" className="text-gray-300 block mb-1">
                  Password
                </Label>
                <div className="relative">
                  <Input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full pl-10 pr-3 py-2 bg-gray-700/50 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition duration-200"
                    placeholder="••••••••"
                  />
                  <Lock
                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                    size={18}
                  />
                </div>
              </motion.div>
            </div>

            <AnimatePresence>
              {error && (
                <motion.div
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  className="mt-4 p-2 bg-red-500/10 border border-red-500/50 rounded-md flex items-center text-red-400"
                >
                  <AlertCircle size={18} className="mr-2" />
                  {error}
                </motion.div>
              )}
            </AnimatePresence>

            <motion.div whileHover={{ scale: 1.03 }} whileTap={{ scale: 0.98 }}>
              <Button
                type="submit"
                className="w-full mt-6 bg-gradient-to-r from-purple-600 to-blue-600 text-white py-2 rounded-md hover:from-purple-700 hover:to-blue-700 transition duration-300 flex items-center justify-center group"
              >
                {isSignUp ? "Sign Up" : "Sign In"}
                <ArrowRight
                  className="ml-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                  size={18}
                />
              </Button>
            </motion.div>

            <p className="mt-4 text-center text-gray-400">
              {isSignUp ? "Already have an account?" : "Don't have an account?"}{" "}
              <Button
                variant="link"
                onClick={toggleAuthMode}
                className="text-purple-400 hover:text-purple-300 transition duration-200"
              >
                {isSignUp ? "Sign In" : "Sign Up"}
              </Button>
            </p>
          </motion.form>
        </AnimatePresence>
      </motion.div>
    </div>
  );
};

export default SleekAuthForm;
