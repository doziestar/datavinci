"use client";

import React, { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ChromeIcon, GithubIcon } from "./Icons/icons";
import { Canvas, useFrame } from "@react-three/fiber";
import { Sphere, MeshDistortMaterial } from "@react-three/drei";
import * as THREE from "three";

const AnimatedSphere = () => {
  const sphereRef = useRef<THREE.Mesh>(null);
  const [distortionAmount, setDistortionAmount] = useState(0.3);

  useFrame(({ clock }) => {
    const elapsedTime = clock.getElapsedTime();

    if (sphereRef.current) {
      sphereRef.current.rotation.x = elapsedTime * 0.1;
      sphereRef.current.rotation.y = elapsedTime * 0.15;
    }

    setDistortionAmount(0.3 + Math.sin(elapsedTime) * 0.1);
  });

  return (
    <Sphere args={[1, 100, 200]} scale={2.5} ref={sphereRef}>
      <MeshDistortMaterial
        color="#8a2be2"
        attach="material"
        distort={distortionAmount}
        speed={1.5}
        roughness={0}
      />
    </Sphere>
  );
};

const Background3D = () => {
  return (
    <div className="absolute inset-0 -z-10">
      <Canvas camera={{ position: [0, 0, 5] }}>
        <ambientLight intensity={0.5} />
        <directionalLight position={[10, 10, 5]} intensity={1} />
        <AnimatedSphere />
      </Canvas>
    </div>
  );
};

export function Auth() {
  const [isSignUp, setIsSignUp] = useState(false);
  const toggleForm = () => setIsSignUp(!isSignUp);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), 500);
    return () => clearTimeout(timer);
  }, []);

  const fadeIn = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <div className="flex min-h-[100dvh] flex-col items-center justify-center bg-gradient-to-br from-background via-primary/5 to-secondary/10 relative overflow-hidden">
      <Background3D />
      <div className="absolute inset-0 background-pattern opacity-5"></div>
      <motion.div
        className="w-full max-w-md space-y-8 z-10"
        initial="hidden"
        animate={isVisible ? "visible" : "hidden"}
        variants={fadeIn}
        transition={{ duration: 0.5 }}
      >
        <Card className="glassmorphic glow group">
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
                  className="mt-1 block w-full bg-secondary/30"
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
                  className="mt-1 block w-full bg-secondary/30"
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
                      className="mt-1 block w-full bg-secondary/30"
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
                <Button
                  variant="outline"
                  className="w-full justify-center bg-secondary/30"
                >
                  <ChromeIcon className="mr-2 h-5 w-5" />
                  Google
                </Button>
                <Button
                  variant="outline"
                  className="w-full justify-center bg-secondary/30"
                >
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
