/* eslint-disable react/no-unescaped-entities */
"use client";

import React, { useState, useEffect, ReactNode } from "react";
import Link from "next/link";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
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
  Cpu,
  Cloud,
  Share2,
  Code,
  GitBranch,
} from "lucide-react";
import Dancing3DHeading from "./3d-dancing-header";

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
    icon: <LineChart className="w-6 h-6" />,
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

interface GlowingCardProps {
  children: ReactNode;
  className?: string;
}

const GlowingCard: React.FC<GlowingCardProps> = ({
  children,
  className = "",
}) => (
  <motion.div
    className={`glassmorphic p-6 rounded-lg floating-card glow group ${className}`}
    whileHover={{ scale: 1.05, transition: { duration: 0.2 } }}
  >
    {children}
  </motion.div>
);

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
          <div className="grid gap-4 md:grid-cols-2 md:gap-16">
            <Dancing3DHeading />
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
                  href="/dashboard"
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
              <GlowingCard key={index}>
                <div className="flex items-center mb-4">
                  <div className="mr-4 text-primary group-hover:text-accent transition-colors duration-300">
                    {feature.icon}
                  </div>
                  <h3 className="text-xl font-semibold group-hover:text-accent transition-colors duration-300">
                    {feature.title}
                  </h3>
                </div>
                <p className="text-muted-foreground group-hover:text-foreground transition-colors duration-300">
                  {feature.description}
                </p>
              </GlowingCard>
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
            {["aggregate", "visualize", "analyze"].map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-4 py-2 rounded-md transition-colors ${
                  activeTab === tab
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary/30 text-secondary-foreground"
                }`}
              >
                {tab.charAt(0).toUpperCase() + tab.slice(1)}
              </button>
            ))}
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
                    {[
                      { icon: Database, text: "Databases" },
                      { icon: Terminal, text: "APIs" },
                      { icon: Chrome, text: "Web Scraping" },
                      { icon: Terminal, text: "Logs" },
                    ].map((item, index) => (
                      <div key={index} className="flex items-center gap-2">
                        <item.icon className="w-6 h-6 text-primary" />
                        <span>{item.text}</span>
                      </div>
                    ))}
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
                    {[
                      "Anomaly detection",
                      "Predictive analytics",
                      "Natural language processing",
                      "Pattern recognition",
                    ].map((item, index) => (
                      <li key={index}>{item}</li>
                    ))}
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
          <GlowingCard className="p-8">
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
                  {[
                    { icon: Code, text: "Custom data connectors" },
                    { icon: GitBranch, text: "Version control integration" },
                    { icon: Cloud, text: "Cloud deployment options" },
                    { icon: Share2, text: "Collaboration features" },
                  ].map((item, index) => (
                    <li key={index} className="flex items-center">
                      <item.icon className="w-5 h-5 mr-2 text-primary" />
                      {item.text}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </GlowingCard>
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
            {[
              {
                title: "Open Source",
                content:
                  "Contribute to DataVinci's core and help shape the future of data analysis.",
              },
              {
                title: "Developer Forum",
                content:
                  "Connect with other developers, share insights, and get help from the community.",
              },
              {
                title: "Resources",
                content:
                  "Access tutorials, documentation, and best practices to make the most of DataVinci.",
              },
            ].map((item, index) => (
              <GlowingCard key={index}>
                <CardHeader>
                  <CardTitle>{item.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground group-hover:text-foreground transition-colors duration-300">
                    {item.content}
                  </p>
                </CardContent>
              </GlowingCard>
            ))}
          </div>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 1 }}
        >
          <GlowingCard className="p-8">
            <h2 className="text-3xl font-bold mb-8 text-center text-foreground">
              What Developers Are Saying
            </h2>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
              {[
                {
                  content:
                    "DataVinci has revolutionized our data pipeline. It's intuitive, powerful, and saves us countless hours every week.",
                  author: "Sarah Chen, Senior Data Engineer",
                },
                {
                  content:
                    "The AI-powered insights have given us a competitive edge. It's like having a data scientist on call 24/7.",
                  author: "Alex Rodriguez, CTO",
                },
                {
                  content:
                    "DataVinci's extensibility is a game-changer. We've integrated it seamlessly with our existing tools and workflows.",
                  author: "Jamie Taylor, Lead Developer",
                },
              ].map((testimonial, index) => (
                <GlowingCard key={index} className="p-6 fallback-bg">
                  <p className="italic mb-4 text-foreground text-shadow group-hover:text-accent transition-colors duration-300">
                    "{testimonial.content}"
                  </p>
                  <p className="font-semibold text-foreground group-hover:text-accent transition-colors duration-300">
                    - {testimonial.author}
                  </p>
                </GlowingCard>
              ))}
            </div>
          </GlowingCard>
        </motion.section>

        <motion.section
          className="py-16 md:py-24"
          initial="hidden"
          animate={isVisible ? "visible" : "hidden"}
          variants={fadeIn}
          transition={{ duration: 0.5, delay: 1.2 }}
        >
          <GlowingCard className="p-8 text-center">
            <h2 className="text-3xl font-bold mb-6">
              Ready to Transform Your Data?
            </h2>
            <p className="text-xl text-muted-foreground mb-8 group-hover:text-foreground transition-colors duration-300">
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
          </GlowingCard>
        </motion.section>
      </main>

      <footer className="w-full py-6 bg-secondary/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {[
              {
                title: "Product",
                links: ["Features", "Pricing", "Case Studies", "API"],
              },
              {
                title: "Company",
                links: ["About Us", "Careers", "Blog", "Contact"],
              },
              {
                title: "Resources",
                links: ["Documentation", "Tutorials", "Community", "GitHub"],
              },
              {
                title: "Legal",
                links: [
                  "Privacy Policy",
                  "Terms of Service",
                  "Cookie Policy",
                  "GDPR",
                ],
              },
            ].map((section, index) => (
              <div key={index}>
                <h3 className="text-lg font-semibold mb-4">{section.title}</h3>
                <ul className="space-y-2">
                  {section.links.map((link, linkIndex) => (
                    <li key={linkIndex}>
                      <Link
                        href="https://github.com/doziestar/datavinci"
                        className="text-muted-foreground hover:text-primary"
                      >
                        {link}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
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

export default LandingPage;
