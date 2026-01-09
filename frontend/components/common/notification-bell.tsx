"use client"

import { Bell, BellOff } from "lucide-react"
import { motion, AnimatePresence } from "motion/react"
import { Button } from "@/components/ui/button"
import { useNotificationSettings } from "@/hooks/use-notification-settings"
import { cn } from "@/lib/utils"

interface NotificationBellProps {
  className?: string
}

export function NotificationBell({ className }: NotificationBellProps) {
  const { isMuted, toggleMute, mounted } = useNotificationSettings()

  // Prevent hydration mismatch by rendering a static icon initially or nothing
  if (!mounted) {
    return (
        <Button variant="ghost" size="icon" className={cn("size-9 text-muted-foreground hover:text-foreground", className)}>
          <Bell className="size-[18px]" />
          <span className="sr-only">通知</span>
        </Button>
    )
  }

  return (
    <Button 
      variant="ghost" 
      size="icon" 
      className={cn("size-9 text-muted-foreground hover:text-foreground", className)} 
      onClick={toggleMute}
    >
      <AnimatePresence mode="wait" initial={false}>
        <motion.div
          key={isMuted ? "muted" : "active"}
          initial={{ rotate: -90, opacity: 0 }}
          animate={{ rotate: 0, opacity: 1 }}
          exit={{ rotate: 90, opacity: 0 }}
          transition={{ duration: 0.2 }}
        >
          {isMuted ? (
             <BellOff className="size-[18px]" />
          ) : (
             <Bell className="size-[18px]" />
          )}
        </motion.div>
      </AnimatePresence>
      <span className="sr-only">{isMuted ? "开启通知" : "关闭通知"}</span>
    </Button>
  )
}
