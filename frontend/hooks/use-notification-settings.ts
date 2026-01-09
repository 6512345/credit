"use client"

import { useState, useEffect } from "react"

interface NotificationSettings {
  systemNotifications: boolean
  notificationSound: boolean
  enabledTypes: string[]
}

const DEFAULT_SETTINGS: NotificationSettings = {
  systemNotifications: true,
  notificationSound: false,
  enabledTypes: ['transfer', 'community', 'red_envelope_receive', 'distribute', 'payment', 'receive'],
}

const STORAGE_KEY = 'ldc-notification-settings'

export function useNotificationSettings() {
  const [settings, setSettings] = useState<NotificationSettings>(DEFAULT_SETTINGS)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        setSettings({ ...DEFAULT_SETTINGS, ...JSON.parse(stored) })
      } catch (e) {
        console.error('Failed to parse notification settings', e)
      }
    }
  }, [])

  const updateSettings = (newSettings: Partial<NotificationSettings>) => {
    const updated = { ...settings, ...newSettings }
    setSettings(updated)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated))
  }

  const toggleMute = () => {
    updateSettings({ systemNotifications: !settings.systemNotifications })
  }

  const setSoundEnabled = (enabled: boolean) => {
    // Current requirement: display as disabled/grayed out, but if we were to implement it:
    updateSettings({ notificationSound: enabled })
  }

  const toggleType = (type: string) => {
    const current = settings.enabledTypes || []
    const next = current.includes(type)
      ? current.filter(t => t !== type)
      : [...current, type]
    updateSettings({ enabledTypes: next })
  }

  return {
    isMuted: !settings.systemNotifications,
    toggleMute,
    soundEnabled: settings.notificationSound,
    setSoundEnabled,
    enabledTypes: settings.enabledTypes || DEFAULT_SETTINGS.enabledTypes,
    toggleType,
    mounted
  }
}
