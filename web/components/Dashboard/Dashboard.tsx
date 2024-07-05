"use client";
import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ChromeIcon, GithubIcon } from "../Icons/icons";

export function Auth() {
  const [isSignUp, setIsSignUp] = useState(false);
  const toggleForm = () => setIsSignUp(!isSignUp);

  return (
    <div className="flex min-h-[100dvh] flex-col items-center justify-center bg-gradient-to-br from-background via-primary/5 to-secondary/10">
      <motion.div
        className="w-full max-w-md space-y-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="glassmorphic">
          <CardHeader>
            <CardTitle className="text-center text-3xl font-bold">
              {isSignUp ? "Create your account" : "Sign in to your account"}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-6">
              <div>
                <Label htmlFor="email" className="block text-sm font-medium">
                  Email address
                </Label>
                <Input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="mt-1 block w-full"
                />
              </div>
              <div>
                <Label htmlFor="password" className="block text-sm font-medium">
                  Password
                </Label>
                <Input
                  id="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  className="mt-1 block w-full"
                />
              </div>
              <AnimatePresence>
                {isSignUp && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.3 }}
                  >
                    <Label
                      htmlFor="confirm-password"
                      className="block text-sm font-medium"
                    >
                      Confirm Password
                    </Label>
                    <Input
                      id="confirm-password"
                      type="password"
                      required
                      className="mt-1 block w-full"
                    />
                  </motion.div>
                )}
              </AnimatePresence>
              {!isSignUp && (
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Checkbox id="remember-me" className="h-4 w-4 rounded" />
                    <Label htmlFor="remember-me" className="ml-2 block text-sm">
                      Remember me
                    </Label>
                  </div>
                  <div className="text-sm">
                    <Link
                      href="#"
                      className="font-medium text-primary hover:text-primary/80"
                      prefetch={false}
                    >
                      Forgot your password?
                    </Link>
                  </div>
                </div>
              )}
              <Button type="submit" className="w-full justify-center">
                {isSignUp ? "Sign up" : "Sign in"}
              </Button>
            </form>
            <div className="mt-6">
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-muted" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="bg-background px-2 text-muted-foreground">
                    Or continue with
                  </span>
                </div>
              </div>
              <div className="mt-6 grid grid-cols-2 gap-3">
                <Button variant="outline" className="w-full justify-center">
                  <ChromeIcon className="mr-2 h-5 w-5" />
                  Google
                </Button>
                <Button variant="outline" className="w-full justify-center">
                  <GithubIcon className="mr-2 h-5 w-5" />
                  Github
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
        <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
          <span>
            {isSignUp ? "Already have an account?" : "Don't have an account?"}
          </span>
          <Button
            variant="link"
            onClick={toggleForm}
            className="font-medium text-primary hover:text-primary/80"
          >
            {isSignUp ? "Sign in" : "Sign up"}
          </Button>
        </div>
      </motion.div>
    </div>
  );
}

export default Auth;
