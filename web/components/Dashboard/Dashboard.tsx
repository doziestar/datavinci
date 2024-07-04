/* eslint-disable react/no-unescaped-entities */

"use client";

import React, { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Header } from "./Header";
import { DatabaseCard } from "../Cards/DatabaseCard";
import { LogCard } from "../Cards/LogCard";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { Sidebar } from "./SideBar";
import { APICard } from "../Cards/ApiCard";

interface DashboardProps {}

export const Dashboard: React.FC<DashboardProps> = () => {
  const [isSidebarExpanded, setIsSidebarExpanded] = useState(false);
  const [isDesktop, setIsDesktop] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      setIsDesktop(window.innerWidth >= 1024);
    };
    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return (
    <div className="flex h-screen overflow-hidden bg-gradient-to-br from-background via-primary/5 to-secondary/10">
      <Sidebar
        isExpanded={isSidebarExpanded}
        setIsExpanded={setIsSidebarExpanded}
      />
      <div className="flex-1 flex flex-col overflow-hidden ml-16 lg:ml-16">
        <Header
          onToggleSidebar={() => setIsSidebarExpanded(!isSidebarExpanded)}
          showToggle={true}
        />
        <main className="flex-1 overflow-auto p-4 lg:p-6">
          <div className="max-w-[2000px] mx-auto space-y-6">
            <Card className="glassmorphic">
              <CardHeader>
                <CardTitle>Databases</CardTitle>
                <CardDescription>
                  Manage and monitor your database connections and performance.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                  <DatabaseCard
                    title="MySQL"
                    description="View connection details, query logs, and performance metrics."
                    connections={12}
                  />
                  <DatabaseCard
                    title="PostgreSQL"
                    description="View connection details, query logs, and performance metrics."
                    connections={8}
                  />
                  <DatabaseCard
                    title="MongoDB"
                    description="View connection details, query logs, and performance metrics."
                    connections={6}
                  />
                  <DatabaseCard
                    title="Redis"
                    description="View connection details, query logs, and performance metrics."
                    connections={4}
                  />
                </div>
              </CardContent>
            </Card>

            <div className="grid gap-6 lg:grid-cols-2">
              <Card className="glassmorphic">
                <CardHeader>
                  <CardTitle>Logs</CardTitle>
                  <CardDescription>
                    Monitor and analyze your application and infrastructure
                    logs.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <LogCard
                      title="Application Logs"
                      description="View and search your application logs for errors and warnings."
                      entries={1234}
                    />
                    <LogCard
                      title="Infrastructure Logs"
                      description="View and search your infrastructure logs for system events."
                      entries={3456}
                    />
                  </div>
                </CardContent>
              </Card>

              <Card className="glassmorphic">
                <CardHeader>
                  <CardTitle>APIs</CardTitle>
                  <CardDescription>
                    Monitor and manage your API integrations and usage.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <APICard
                      name="User Authentication"
                      endpoint="/api/auth"
                      method="POST"
                      status="Active"
                    />
                    <APICard
                      name="Product Catalog"
                      endpoint="/api/products"
                      method="GET"
                      status="Active"
                    />
                    <APICard
                      name="Order Processing"
                      endpoint="/api/orders"
                      method="PUT"
                      status="Inactive"
                    />
                    <APICard
                      name="Legacy Payment Gateway"
                      endpoint="/api/v1/payments"
                      method="POST"
                      status="Deprecated"
                    />
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default Dashboard;
