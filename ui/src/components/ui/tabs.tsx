"use client"

import * as React from "react"
import * as TabsPrimitive from "@radix-ui/react-tabs"

import { cn } from "@/lib/utils"

function Tabs({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.Root>) {
  return (
    <TabsPrimitive.Root
      data-slot="tabs"
      className={cn("flex flex-col gap-2", className)}
      {...props}
    />
  )
}

// line：页面级标签（详情页、表单），底部一条线，选中项下方高亮；默认是分段按钮样式
function TabsList({
  className,
  variant = "pill",
  ...props
}: React.ComponentProps<typeof TabsPrimitive.List> & { variant?: "pill" | "line" }) {
  return (
    <TabsPrimitive.List
      data-slot="tabs-list"
      data-variant={variant}
      className={cn(
        "group/tabs text-muted-foreground inline-flex w-fit items-center",
        variant === "line"
          ? "h-10 min-w-full justify-start gap-1 border-b"
          : "bg-muted h-9 justify-center rounded-lg p-[3px]",
        className
      )}
      {...props}
    />
  )
}

function TabsTrigger({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.Trigger>) {
  return (
    <TabsPrimitive.Trigger
      data-slot="tabs-trigger"
      className={cn(
        "group-data-[variant=line]/tabs:relative group-data-[variant=line]/tabs:h-full group-data-[variant=line]/tabs:flex-none group-data-[variant=line]/tabs:rounded-none group-data-[variant=line]/tabs:px-3 group-data-[variant=line]/tabs:text-muted-foreground group-data-[variant=line]/tabs:hover:text-foreground group-data-[variant=line]/tabs:data-[state=active]:text-foreground group-data-[variant=line]/tabs:data-[state=active]:shadow-none group-data-[variant=line]/tabs:data-[state=active]:bg-transparent group-data-[variant=line]/tabs:dark:data-[state=active]:bg-transparent group-data-[variant=line]/tabs:dark:data-[state=active]:border-transparent group-data-[variant=line]/tabs:after:absolute group-data-[variant=line]/tabs:after:inset-x-2 group-data-[variant=line]/tabs:after:-bottom-px group-data-[variant=line]/tabs:after:h-0.5 group-data-[variant=line]/tabs:after:rounded-full group-data-[variant=line]/tabs:data-[state=active]:after:bg-primary group-data-[variant=line]/tabs:data-[state=active]:after:shadow-[0_0_8px] group-data-[variant=line]/tabs:data-[state=active]:after:shadow-primary group-data-[variant=line]/tabs:data-[state=active]:text-primary data-[state=active]:bg-background dark:data-[state=active]:text-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:outline-ring dark:data-[state=active]:border-input dark:data-[state=active]:bg-input/30 text-foreground dark:text-muted-foreground inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
        className
      )}
      {...props}
    />
  )
}

function TabsContent({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.Content>) {
  return (
    <TabsPrimitive.Content
      data-slot="tabs-content"
      className={cn("flex-1 outline-none", className)}
      {...props}
    />
  )
}

export { Tabs, TabsList, TabsTrigger, TabsContent }
