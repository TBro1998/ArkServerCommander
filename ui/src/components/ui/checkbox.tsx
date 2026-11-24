import * as React from "react"
import { Check } from "lucide-react"

import { cn } from "@/lib/utils"

// 原生 input 实现，只借外观：选中时填充主色，勾号由 lucide 图标叠加
function Checkbox({ className, ...props }: React.ComponentProps<"input">) {
  return (
    <span className="relative inline-flex size-4 shrink-0 items-center justify-center">
      <input
        type="checkbox"
        data-slot="checkbox"
        className={cn(
          "peer border-input dark:bg-input/30 size-4 shrink-0 appearance-none rounded-[3px] border shadow-xs outline-none transition-colors",
          "checked:border-primary checked:bg-primary",
          "focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
          "disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        {...props}
      />
      <Check className="text-primary-foreground pointer-events-none absolute size-3 opacity-0 peer-checked:opacity-100" />
    </span>
  )
}

export { Checkbox }
