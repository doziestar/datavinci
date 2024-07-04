import React, { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import {
  MenuIcon,
  SearchIcon,
  BellIcon,
  SunIcon,
  MoonIcon,
} from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { AnimatedButton } from "../Buttons/AnimatedButton";

interface HeaderProps {
  onToggleSidebar: () => void;
  showToggle: boolean;
  isSidebarExpanded: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  onToggleSidebar,
  showToggle,
  isSidebarExpanded,
}) => {
  const [isDarkMode, setIsDarkMode] = useState(false);

  const toggleDarkMode = () => {
    setIsDarkMode(!isDarkMode);
    // Implement actual dark mode toggle logic here
  };

  return (
    <header className="sticky top-0 z-40 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center px-4 lg:px-6">
        {showToggle && (
          <AnimatedButton
            variant="ghost"
            size="icon"
            onClick={onToggleSidebar}
            className={`mr-4 z-50 ${
              isSidebarExpanded ? "fixed left-[200px]" : ""
            } transition-all duration-300 ease-in-out`}
          >
            <MenuIcon className="h-5 w-5" />
            <span className="sr-only">Toggle Sidebar</span>
          </AnimatedButton>
        )}
        <div className="flex-1" /> {/* Spacer */}
        <div className="flex items-center space-x-4">
          <AnimatedButton size="icon" variant="ghost">
            <SearchIcon className="h-5 w-5" />
            <span className="sr-only">Search</span>
          </AnimatedButton>
          <AnimatedButton size="icon" variant="ghost">
            <BellIcon className="h-5 w-5" />
            <span className="sr-only">Notifications</span>
          </AnimatedButton>
          <AnimatedButton size="icon" variant="ghost" onClick={toggleDarkMode}>
            {isDarkMode ? (
              <SunIcon className="h-5 w-5" />
            ) : (
              <MoonIcon className="h-5 w-5" />
            )}
            <span className="sr-only">Toggle Dark Mode</span>
          </AnimatedButton>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <AnimatedButton
                variant="ghost"
                size="icon"
                className="relative h-8 w-8 rounded-full"
              >
                <Avatar className="h-8 w-8">
                  <AvatarImage
                    src="https://avatars.githubusercontent.com/u/39489124?v=4"
                    alt="@doziestar"
                  />
                  <AvatarFallback>SC</AvatarFallback>
                </Avatar>
              </AnimatedButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Settings</DropdownMenuItem>
              <DropdownMenuItem>Logout</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
};
