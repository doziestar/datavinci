/* eslint-disable react/no-unescaped-entities */
"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { motion, useScroll, useTransform } from "framer-motion";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import {
  Database,
  Terminal,
  Chrome,
  TrendingUp,
  Package,
  Info,
  ChevronDown,
  Code,
  Cpu,
  GitBranch,
  Cloud,
  Share2,
} from "lucide-react";

const mockData = [
  { name: "Jan", value: 4000 },
  { name: "Feb", value: 3000 },
  { name: "Mar", value: 5000 },
  { name: "Apr", value: 4500 },
  { name: "May", value: 6000 },
  { name: "Jun", value: 5500 },
];

const features = [
  {
    icon: <Database className="w-6 h-6" />,
    title: "Multi-source Integration",
    description:
      "Connect to various databases, APIs, and data sources effortlessly.",
  },
  {
    icon: <LineChartIcon className="w-6 h-6" />,
    title: "Advanced Visualizations",
    description: "Create stunning, interactive charts and dashboards.",
  },
  {
    icon: <Cpu className="w-6 h-6" />,
    title: "AI-powered Analysis",
    description:
      "Leverage machine learning for predictive analytics and anomaly detection.",
  },
  {
    icon: <Terminal className="w-6 h-6" />,
    title: "Real-time Processing",
    description: "Process and analyze data streams in real-time.",
  },
  {
    icon: <Cloud className="w-6 h-6" />,
    title: "Cloud Integration",
    description:
      "Seamlessly integrate with major cloud providers and services.",
  },
  {
    icon: <Share2 className="w-6 h-6" />,
    title: "Collaboration Tools",
    description:
      "Foster team collaboration with built-in sharing and version control.",
  },
];

export function LandingPage() {
  const [isVisible, setIsVisible] = useState(false);
  const [activeTab, setActiveTab] = useState("aggregate");
  const { scrollYProgress } = useScroll();
  const opacity = useTransform(scrollYProgress, [0, 0.5], [1, 0]);
  const scale = useTransform(scrollYProgress, [0, 0.5], [1, 0.8]);

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), 500);
    return () => clearTimeout(timer);
  }, []);

  const fadeIn = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <div className="flex flex-col min-h-[100dvh] items-center bg-gradient-to-br from-background via-primary/5 to-secondary/10">
      <div className="absolute inset-0 background-pattern opacity-5"></div>
      <main className="relative flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 overflow-hidden">
        <motion.section
          className="w-full pt-12 md:pt-24 lg:pt-32"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5 }}
          style={{ opacity, scale }}
        >
          <div className="space-y-10 xl:space-y-16">
            <div className="grid gap-4 md:grid-cols-2 md:gap-16">
              <div>
                <motion.h1
                  className="lg:leading-tighter text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl xl:text-7xl bg-clip-text text-transparent bg-gradient-to-r from-primary to-secondary"
                  variants={fadeIn}
                >
                  Empower Your Data with DataVinci
                </motion.h1>
              </div>
              <div className="flex flex-col items-start space-y-4">
                <motion.p
                  className="text-muted-foreground md:text-xl"
                  variants={fadeIn}
                >
                  DataVinci is the ultimate data management and visualization
                  platform for developers. Harness the power of your data with
                  advanced analytics, AI-driven insights, and seamless
                  integration.
                </motion.p>
                <motion.div variants={fadeIn} className="flex space-x-4">
                  <Link
                    href="#"
                    className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-6 py-2 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
                    prefetch={false}
                  >
                    Get Started
                  </Link>
                  <Link
                    href="#"
                    className="inline-flex h-10 items-center justify-center rounded-md border border-primary px-6 py-2 text-sm font-medium text-primary shadow transition-colors hover:bg-primary/10 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
                    prefetch={false}
                  >
                    View Demo
                  </Link>
                </motion.div>
              </div>
            </div>
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <h2 className="text-3xl font-bold mb-8 text-center">
            Powerful Features for Developers
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                className="glassmorphic p-6 rounded-lg"
                variants={fadeIn}
                transition={{ duration: 0.5, delay: 0.1 * index }}
              >
                <div className="flex items-center mb-4">
                  <div className="mr-4 text-primary">{feature.icon}</div>
                  <h3 className="text-xl font-semibold">{feature.title}</h3>
                </div>
                <p className="text-muted-foreground">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 0.4 }}
        >
          <h2 className="text-3xl font-bold mb-8 text-center">
            Explore DataVinci
          </h2>
          <div className="flex justify-center space-x-4 mb-8">
            <button
              onClick={() => setActiveTab("aggregate")}
              className={`px-4 py-2 rounded-md transition-colors ${
                activeTab === "aggregate"
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary/30 text-secondary-foreground"
              }`}
            >
              Aggregate
            </button>
            <button
              onClick={() => setActiveTab("visualize")}
              className={`px-4 py-2 rounded-md transition-colors ${
                activeTab === "visualize"
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary/30 text-secondary-foreground"
              }`}
            >
              Visualize
            </button>
            <button
              onClick={() => setActiveTab("analyze")}
              className={`px-4 py-2 rounded-md transition-colors ${
                activeTab === "analyze"
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary/30 text-secondary-foreground"
              }`}
            >
              Analyze
            </button>
          </div>

          <Card className="w-full glassmorphic">
            <CardContent className="p-6">
              {activeTab === "aggregate" && (
                <div className="space-y-4">
                  <h3 className="text-2xl font-bold">
                    Aggregate Data from Multiple Sources
                  </h3>
                  <p className="text-muted-foreground">
                    Connect seamlessly to various data sources:
                  </p>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="flex items-center gap-2">
                      <Database className="w-6 h-6 text-primary" />
                      <span>Databases</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Terminal className="w-6 h-6 text-primary" />
                      <span>APIs</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Chrome className="w-6 h-6 text-primary" />
                      <span>Web Scraping</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Terminal className="w-6 h-6 text-primary" />
                      <span>Logs</span>
                    </div>
                  </div>
                </div>
              )}
              {activeTab === "visualize" && (
                <div className="space-y-4">
                  <h3 className="text-2xl font-bold">Visualize Your Data</h3>
                  <p className="text-muted-foreground">
                    Create interactive and insightful visualizations:
                  </p>
                  <div className="h-64 w-full">
                    <ResponsiveContainer width="100%" height="100%">
                      <LineChart data={mockData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" />
                        <YAxis />
                        <Tooltip />
                        <Line
                          type="monotone"
                          dataKey="value"
                          stroke="#8884d8"
                        />
                      </LineChart>
                    </ResponsiveContainer>
                  </div>
                </div>
              )}
              {activeTab === "analyze" && (
                <div className="space-y-4">
                  <h3 className="text-2xl font-bold">AI-Powered Analysis</h3>
                  <p className="text-muted-foreground">
                    Leverage advanced AI models for data analysis:
                  </p>
                  <ul className="list-disc list-inside space-y-2">
                    <li>Anomaly detection</li>
                    <li>Predictive analytics</li>
                    <li>Natural language processing</li>
                    <li>Pattern recognition</li>
                  </ul>
                </div>
              )}
            </CardContent>
          </Card>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 0.6 }}
        >
          <div className="glassmorphic p-8 rounded-lg">
            <h2 className="text-3xl font-bold mb-6">
              Built for Developers, by Developers
            </h2>
            <div className="grid md:grid-cols-2 gap-8">
              <div>
                <h3 className="text-xl font-semibold mb-4">Powerful API</h3>
                <p className="text-muted-foreground mb-4">
                  Integrate DataVinci seamlessly into your existing workflows
                  with our comprehensive API.
                </p>
                <pre className="bg-secondary/20 p-4 rounded-md overflow-x-auto">
                  <code className="text-sm">
                    {`import datavinci

# Initialize the client
client = datavinci.Client(api_key="your_api_key")

# Fetch and analyze data
data = client.fetch_data("sales_2023")
insights = client.analyze(data, model="predictive")

# Visualize results
chart = client.visualize(insights, type="line_chart")
chart.save("sales_forecast.png")`}
                  </code>
                </pre>
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-4">
                  Extensible Architecture
                </h3>
                <p className="text-muted-foreground mb-4">
                  DataVinci's plugin system allows you to extend functionality
                  and integrate with your favorite tools.
                </p>
                <ul className="space-y-2">
                  <li className="flex items-center">
                    <Code className="w-5 h-5 mr-2 text-primary" />
                    Custom data connectors
                  </li>
                  <li className="flex items-center">
                    <GitBranch className="w-5 h-5 mr-2 text-primary" />
                    Version control integration
                  </li>
                  <li className="flex items-center">
                    <Cloud className="w-5 h-5 mr-2 text-primary" />
                    Cloud deployment options
                  </li>
                  <li className="flex items-center">
                    <Share2 className="w-5 h-5 mr-2 text-primary" />
                    Collaboration features
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 0.8 }}
        >
          <h2 className="text-3xl font-bold mb-8 text-center">
            Join the DataVinci Community
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            <Card className="glassmorphic">
              <CardHeader>
                <CardTitle>Open Source</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Contribute to DataVinci's core and help shape the future of
                  data analysis.
                </p>
              </CardContent>
            </Card>
            <Card className="glassmorphic">
              <CardHeader>
                <CardTitle>Developer Forum</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Connect with other developers, share insights, and get help
                  from the community.
                </p>
              </CardContent>
            </Card>
            <Card className="glassmorphic">
              <CardHeader>
                <CardTitle>Resources</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Access tutorials, documentation, and best practices to make
                  the most of DataVinci.
                </p>
              </CardContent>
            </Card>
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 1 }}
        >
          <div className="glassmorphic p-8 rounded-lg">
            <h2 className="text-3xl font-bold mb-8 text-center">
              What Developers Are Saying
            </h2>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
              <Card className="glassmorphic">
                <CardContent className="pt-6">
                  <p className="italic mb-4">
                    "DataVinci has revolutionized our data pipeline. It's
                    intuitive, powerful, and saves us countless hours every
                    week."
                  </p>
                  <p className="font-semibold">
                    - Sarah Chen, Senior Data Engineer
                  </p>
                </CardContent>
              </Card>
              <Card className="glassmorphic">
                <CardContent className="pt-6">
                  <p className="italic mb-4">
                    "The AI-powered insights have given us a competitive edge.
                    It's like having a data scientist on call 24/7."
                  </p>
                  <p className="font-semibold">- Alex Rodriguez, CTO</p>
                </CardContent>
              </Card>
              <Card className="glassmorphic">
                <CardContent className="pt-6">
                  <p className="italic mb-4">
                    "DataVinci's extensibility is a game-changer. We've
                    integrated it seamlessly with our existing tools and
                    workflows."
                  </p>
                  <p className="font-semibold">
                    - Jamie Taylor, Lead Developer
                  </p>
                </CardContent>
              </Card>
            </div>
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 1.2 }}
        >
          <div className="glassmorphic p-8 rounded-lg text-center">
            <h2 className="text-3xl font-bold mb-6">
              Ready to Transform Your Data?
            </h2>
            <p className="text-xl text-muted-foreground mb-8">
              Join thousands of developers who are already leveraging DataVinci
              to unlock the full potential of their data.
            </p>
            <div className="flex justify-center space-x-4">
              <Link
                href="#"
                className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-8 py-2 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
                prefetch={false}
              >
                Start Free Trial
              </Link>
              <Link
                href="#"
                className="inline-flex h-10 items-center justify-center rounded-md border border-primary px-8 py-2 text-sm font-medium text-primary shadow transition-colors hover:bg-primary/10 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
                prefetch={false}
              >
                Schedule a Demo
              </Link>
            </div>
          </div>
        </motion.section>
      </main>

      <footer className="w-full py-6 bg-secondary/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-lg font-semibold mb-4">Product</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Features
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Pricing
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Case Studies
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    API
                  </Link>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="text-lg font-semibold mb-4">Company</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    About Us
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Careers
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Blog
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Contact
                  </Link>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="text-lg font-semibold mb-4">Resources</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Documentation
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Tutorials
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Community
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    GitHub
                  </Link>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="text-lg font-semibold mb-4">Legal</h3>
              <ul className="space-y-2">
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Privacy Policy
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Terms of Service
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    Cookie Policy
                  </Link>
                </li>
                <li>
                  <Link
                    href="#"
                    className="text-muted-foreground hover:text-primary"
                  >
                    GDPR
                  </Link>
                </li>
              </ul>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-secondary/30 text-center">
            <p className="text-muted-foreground">
              &copy; 2024 DataVinci. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}

function DatabaseIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <ellipse cx="12" cy="5" rx="9" ry="3" />
      <path d="M3 5V19A9 3 0 0 0 21 19V5" />
      <path d="M3 12A9 3 0 0 0 21 12" />
    </svg>
  );
}

function InfoIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="10" />
      <path d="M12 16v-4" />
      <path d="M12 8h.01" />
    </svg>
  );
}

function LineChartIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M3 3v18h18" />
      <path d="m19 9-5 5-4-4-3 3" />
    </svg>
  );
}

function PackageIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="m7.5 4.27 9 5.15" />
      <path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z" />
      <path d="m3.3 7 8.7 5 8.7-5" />
      <path d="M12 22V12" />
    </svg>
  );
}

function PiIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <line x1="9" x2="9" y1="4" y2="20" />
      <path d="M4 7c0-1.7 1.3-3 3-3h13" />
      <path d="M18 20c-1.7 0-3-1.3-3-3V4" />
    </svg>
  );
}

function PieChartIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M21.21 15.89A10 10 0 1 1 8 2.83" />
      <path d="M22 12A10 10 0 0 0 12 2v10z" />
    </svg>
  );
}

function ScatterChartIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="7.5" cy="7.5" r=".5" fill="currentColor" />
      <circle cx="18.5" cy="5.5" r=".5" fill="currentColor" />
      <circle cx="11.5" cy="11.5" r=".5" fill="currentColor" />
      <circle cx="7.5" cy="16.5" r=".5" fill="currentColor" />
      <circle cx="17.5" cy="14.5" r=".5" fill="currentColor" />
      <path d="M3 3v18h18" />
    </svg>
  );
}

function TerminalIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <polyline points="4 17 10 11 4 5" />
      <line x1="12" x2="20" y1="19" y2="19" />
    </svg>
  );
}

function TrendingUpIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <polyline points="22 7 13.5 15.5 8.5 10.5 2 17" />
      <polyline points="16 7 22 7 22 13" />
    </svg>
  );
}
