"use client"

import * as React from "react"
import Link from "next/link"
import { Bell, Volume2, Check } from "lucide-react"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { useNotificationSettings } from "@/hooks/use-notification-settings"
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

// 定义系统的通知类型
const NOTIFICATION_TYPES = [
  { id: 'transfer', label: '积分转移', description: '账户之间的积分互转' },
  { id: 'community', label: '社区划转', description: '来自社区操作的积分变动' },
  { id: 'red_envelope_receive', label: '红包收入', description: '领取他人发放的红包' },
  { id: 'distribute', label: '系统分发', description: '系统自动分发的积分奖励' },
  { id: 'receive', label: '他人转入', description: '收到他人的直接转账' },
  { id: 'payment', label: '消费支出', description: '购买商品或服务的支出' },
]

function GeneralNotificationSection() {
  const { isMuted, toggleMute, soundEnabled, setSoundEnabled, mounted } = useNotificationSettings()

  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-medium text-sm text-foreground">通用设置</h2>
        <p className="text-xs text-muted-foreground">
          控制全局通知开关和音效
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* System Notifications */}
        <div className={cn(
          "relative flex flex-col justify-between rounded-xl border p-4 transition-all hover:bg-muted/30 hover:border-border/80",
          !isMuted ? "border-primary/50 bg-primary/5" : "border-border"
        )}>
          <div className="space-y-2">
            <div className="flex items-start justify-between">
              <div className={cn("p-2 rounded-lg transition-colors", !isMuted ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground")}>
                <Bell className="size-5" />
              </div>
              <Switch
                id="system-notifications"
                checked={mounted ? !isMuted : true}
                onCheckedChange={toggleMute}
              />
            </div>
            <div>
              <Label htmlFor="system-notifications" className="font-medium">系统通知</Label>
              <p className="text-xs text-muted-foreground mt-1">
                {isMuted ? "使用网页弹窗即时通知消息" : "使用系统弹窗即时通知消息"}
              </p>
            </div>
          </div>
        </div>

        {/* Sound Notifications */}
        <div className="relative flex flex-col justify-between rounded-xl border border-border p-4 transition-all opacity-60">
          <div className="space-y-2">
            <div className="flex items-start justify-between">
              <div className="p-2 rounded-lg bg-muted text-muted-foreground">
                <Volume2 className="size-5" />
              </div>
              <Switch
                id="notification-sound"
                checked={soundEnabled}
                onCheckedChange={setSoundEnabled}
                disabled
              />
            </div>
            <div>
              <Label htmlFor="notification-sound" className="font-medium">通知音效</Label>
              <p className="text-xs text-muted-foreground mt-1">
                接收通知时播放提示音
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function FilterNotificationSection() {
  const { enabledTypes, toggleType, mounted } = useNotificationSettings()

  if (!mounted) return null

  // 计算已选中的数量
  const startCount = enabledTypes.length
  
  return (
    <div className="space-y-6">
      <div>
        <div className="flex items-center gap-2">
            <h2 className="font-medium text-sm text-foreground">过滤类别</h2>
            <Badge variant="secondary" className="text-[10px] h-5 px-1.5 font-normal text-muted-foreground">
                已选 {startCount}/{NOTIFICATION_TYPES.length}
            </Badge>
        </div>
        <p className="text-xs text-muted-foreground">
          选择您希望接收通知的交易类型
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {NOTIFICATION_TYPES.map((type) => {
          const isChecked = enabledTypes.includes(type.id)
          return (
            <div
              key={type.id}
              className={cn(
                "flex items-start space-x-3 rounded-lg border p-3 transition-colors cursor-pointer",
                isChecked ? "bg-accent/40 border-primary/20" : "hover:bg-muted/50 border-transparent hover:border-border"
              )}
              onClick={() => toggleType(type.id)}
            >
              <Checkbox
                id={`type-${type.id}`}
                checked={isChecked}
                onCheckedChange={() => toggleType(type.id)}
                className="mt-0.5"
              />
              <div className="space-y-1">
                <Label
                  htmlFor={`type-${type.id}`}
                  className="font-medium text-sm leading-none cursor-pointer"
                >
                  {type.label}
                </Label>
                <p className="text-xs text-muted-foreground line-clamp-2">
                  {type.description}
                </p>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default function NotificationSettingsPage() {
  return (
    <div className="py-6 space-y-8">
      {/* Breadcrumb Header */}
      <div className="font-semibold">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href="/settings" className="text-base text-primary">设置</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage className="text-base font-semibold">通知设置</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </div>

      <GeneralNotificationSection />
      
      <div className="w-full h-px bg-border/40" />
      
      <FilterNotificationSection />
    </div>
  )
}
