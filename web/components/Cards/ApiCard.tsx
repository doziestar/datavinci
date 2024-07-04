import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ActivityIcon } from "lucide-react";

interface APICardProps {
  name: string;
  endpoint: string;
  method: "GET" | "POST" | "PUT" | "DELETE";
  status: "Active" | "Inactive" | "Deprecated";
}

export const APICard: React.FC<APICardProps> = ({
  name,
  endpoint,
  method,
  status,
}) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "Active":
        return "bg-green-500";
      case "Inactive":
        return "bg-yellow-500";
      case "Deprecated":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{name}</CardTitle>
        <Badge variant="outline" className={getStatusColor(status)}>
          {status}
        </Badge>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{endpoint}</div>
        <p className="text-xs text-muted-foreground">Method: {method}</p>
        <div className="mt-4 flex items-center">
          <ActivityIcon className="h-4 w-4 text-muted-foreground" />
          <span className="ml-2 text-xs text-muted-foreground">
            Last called 2 hours ago
          </span>
        </div>
      </CardContent>
    </Card>
  );
};
