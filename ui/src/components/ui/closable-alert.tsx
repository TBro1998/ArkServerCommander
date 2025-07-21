"use client";

import * as React from "react"
import { type VariantProps } from "class-variance-authority"
import { X } from "lucide-react"

import { cn } from "@/lib/utils"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"

interface ClosableAlertProps extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof Alert> {
  title?: string;
  onClose: () => void;
}

const ClosableAlert = React.forwardRef<HTMLDivElement, ClosableAlertProps>(
  ({ className, variant, title, children, onClose, ...props }, ref) => {
    const [show, setShow] = React.useState(true);

    const handleClose = () => {
      setShow(false);
      onClose();
    };

    if (!show) return null;

    return (
      <Alert ref={ref} className={cn(className)} variant={variant} {...props}>
        <div className="flex justify-between items-start">
          <div>
            {title && <AlertTitle>{title}</AlertTitle>}
            <AlertDescription>{children}</AlertDescription>
          </div>
          <button onClick={handleClose} className="p-1 -m-1">
            <X className="h-4 w-4" />
          </button>
        </div>
      </Alert>
    )
  }
)
ClosableAlert.displayName = "ClosableAlert"

export { ClosableAlert }