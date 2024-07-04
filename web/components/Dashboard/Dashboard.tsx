/* eslint-disable react/no-unescaped-entities */

"use client";

import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";

import React from "react";
import { Sidebar } from "lucide-react";
import { Header } from "./Header";
import { DatabaseCard } from "../Cards/DatabaseCard";
import { LogCard } from "../Cards/LogCard";

export function Dashboard() {
  return (
    <div className="flex min-h-screen w-full bg-gradient-to-br from-background via-primary/5 to-secondary/10">
      <div className="absolute inset-0 background-pattern opacity-5"></div>
      <Sidebar />
      <div className="flex flex-1 flex-col gap-4 p-4 sm:gap-8 sm:p-6 ml-14">
        <Header />
        <main className="grid flex-1 items-start gap-4 sm:px-6 sm:py-0 md:gap-8 lg:grid-cols-2 xl:grid-cols-3">
          <Card className="col-span-2 lg:col-span-1 xl:col-span-2 glassmorphic">
            <CardHeader>
              <CardTitle>Databases</CardTitle>
              <CardDescription>
                Manage and monitor your database connections and performance.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
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
          <Card className="col-span-2 lg:col-span-1 xl:col-span-1 glassmorphic">
            <CardHeader>
              <CardTitle>Logs</CardTitle>
              <CardDescription>
                Monitor and analyze your application and infrastructure logs.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-1">
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
          <Card className="col-span-2 lg:col-span-1 xl:col-span-1 glassmorphic">
            <CardHeader>
              <CardTitle>APIs</CardTitle>
              <CardDescription>
                Monitor and manage your API integrations and usage.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-1" />
            </CardContent>
          </Card>
        </main>
      </div>
    </div>
  );
}

export default Dashboard;
