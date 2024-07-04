import React from "react";
import Link from "next/link";
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "@/components/ui/tooltip";
import {
  DatabaseIcon,
  LayoutGridIcon,
  LayersIcon,
  CodeIcon,
  GlobeIcon,
  SettingsIcon,
} from "../Icons/icons";

type IconType = React.ComponentType<React.SVGProps<SVGSVGElement>>;

interface SidebarLinkProps {
  href: string;
  icon: IconType;
  label: string;
}

const SidebarLink: React.FC<SidebarLinkProps> = ({
  href,
  icon: IconComponent,
  label,
}) => (
  <Tooltip>
    <TooltipTrigger asChild>
      <Link
        href={href}
        className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground hover:bg-accent md:h-8 md:w-8"
        prefetch={false}
      >
        <div className="h-5 w-5">
          <IconComponent />
        </div>
        <span className="sr-only">{label}</span>
      </Link>
    </TooltipTrigger>
    <TooltipContent side="right">{label}</TooltipContent>
  </Tooltip>
);

export const Sidebar: React.FC = () => (
  <aside className="fixed inset-y-0 left-0 z-10 hidden w-14 flex-col border-r border-border bg-background/50 backdrop-blur-md sm:flex">
    <nav className="flex flex-col items-center gap-4 px-2 sm:py-5">
      <TooltipProvider>
        <Link
          href="#"
          className="group flex h-9 w-9 shrink-0 items-center justify-center gap-2 rounded-full bg-primary text-lg font-semibold text-primary-foreground md:h-8 md:w-8 md:text-base"
          prefetch={false}
        >
          <div className="h-4 w-4 transition-all group-hover:scale-110">
            <DatabaseIcon />
          </div>
          <span className="sr-only">Data Dashboard</span>
        </Link>
        <SidebarLink href="#" icon={LayoutGridIcon} label="Databases" />
        <SidebarLink href="#" icon={LayersIcon} label="Logs" />
        <SidebarLink href="#" icon={CodeIcon} label="APIs" />
        <SidebarLink href="#" icon={GlobeIcon} label="Scraped Data" />
      </TooltipProvider>
    </nav>
    <nav className="mt-auto flex flex-col items-center gap-4 px-2 sm:py-5">
      <TooltipProvider>
        <SidebarLink href="#" icon={SettingsIcon} label="Settings" />
      </TooltipProvider>
    </nav>
  </aside>
);
