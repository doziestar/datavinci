import React from "react";
import AuthForm from "@/components/forms/AuthForm";

const AuthPage: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/20 to-secondary/20">
      <AuthForm />
    </div>
  );
};

export default AuthPage;
