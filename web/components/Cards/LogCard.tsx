import React from "react";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { LayersIcon } from "lucide-react";

type LogCardProps = {
  title: string;
  description: string;
  entries: number;
};

export const LogCard = ({ title, description, entries }: LogCardProps) => (
  <Card className="glassmorphic">
    <CardHeader>
      <CardTitle>{title}</CardTitle>
      <CardDescription>{description}</CardDescription>
    </CardHeader>
    <CardContent>
      <div className="flex items-center justify-between">
        <div className="text-4xl font-bold">{entries}</div>
        <LayersIcon className="h-8 w-8 text-muted-foreground" />
      </div>
      <div className="text-sm text-muted-foreground">Entries</div>
    </CardContent>
    <CardFooter>
      <Button variant="outline" size="sm">
        Analyze
      </Button>
    </CardFooter>
  </Card>
);
