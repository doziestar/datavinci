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
import { DatabaseIcon } from "lucide-react";

type DatabaseCardProps = {
  title: string;
  description: string;
  connections: number;
};

export const DatabaseCard = ({
  title,
  description,
  connections,
}: DatabaseCardProps) => (
  <Card className="glassmorphic">
    <CardHeader>
      <CardTitle>{title}</CardTitle>
      <CardDescription>{description}</CardDescription>
    </CardHeader>
    <CardContent>
      <div className="flex items-center justify-between">
        <div className="text-4xl font-bold">{connections}</div>
        <DatabaseIcon className="h-8 w-8 text-muted-foreground" />
      </div>
      <div className="text-sm text-muted-foreground">Connections</div>
    </CardContent>
    <CardFooter>
      <Button variant="outline" size="sm">
        Manage
      </Button>
    </CardFooter>
  </Card>
);
