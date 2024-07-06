import React from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  DatabaseIcon,
  LayoutGridIcon,
  LayersIcon,
  CodeIcon,
  GlobeIcon,
  SettingsIcon,
} from "lucide-react";
import { LucideIcon } from "lucide-react";
import { useAnimations } from "@/hooks/animation/useAnimation";

interface SidebarLinkProps {
  href: string;
  icon: LucideIcon;
  label: string;
}

const SidebarLink: React.FC<SidebarLinkProps> = ({
  href,
  icon: Icon,
  label,
}) => {
  const { jelly, wobble } = useAnimations();

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <motion.div
            whileHover="hover"
            variants={jelly}
            onHoverStart={() => wobble.set(10)}
            onHoverEnd={() => wobble.set(0)}
          >
            <Link
              href={href}
              className="flex h-10 w-10 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground hover:bg-accent"
            >
              <motion.div style={{ rotate: wobble }}>
                <Icon className="h-5 w-5" />
              </motion.div>
              <span className="sr-only">{label}</span>
            </Link>
          </motion.div>
        </TooltipTrigger>
        <TooltipContent side="right">{label}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

interface SidebarProps {
  isExpanded: boolean;
  setIsExpanded: (isExpanded: boolean) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  isExpanded,
  setIsExpanded,
}) => {
  return (
    <motion.aside
      initial={false}
      animate={{ width: isExpanded ? 240 : 64 }}
      transition={{ duration: 0.3 }}
      className="fixed inset-y-0 left-0 z-50 bg-background/80 backdrop-blur-md border-r border-border flex flex-col"
    >
      <div className="flex h-16 items-center justify-between px-4">
        {isExpanded ? (
          <Link
            href="#"
            className="flex items-center gap-2 text-lg font-semibold"
          >
            <DatabaseIcon className="h-6 w-6" />
            <span>Data Dashboard</span>
          </Link>
        ) : (
          <DatabaseIcon className="h-6 w-6 mx-auto" />
        )}
      </div>
      <nav className="flex flex-col items-center gap-4 py-4">
        <SidebarLink href="#" icon={LayoutGridIcon} label="Databases" />
        <SidebarLink href="#" icon={LayersIcon} label="Logs" />
        <SidebarLink href="#" icon={CodeIcon} label="APIs" />
        <SidebarLink href="#" icon={GlobeIcon} label="Scraped Data" />
      </nav>
      <div className="mt-auto pb-4">
        <SidebarLink href="#" icon={SettingsIcon} label="Settings" />
      </div>
    </motion.aside>
  );
};
