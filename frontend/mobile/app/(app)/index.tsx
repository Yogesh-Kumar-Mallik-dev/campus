/**
 * BLOCK_MOBILE_INDEX_SCREEN_001
 * Subsystem: Mobile Hub Navigation Gateway
 * Purpose:   Mobile Super-App Directory providing quick access to all 17 topological subsystems.
 */

import React from 'react';
import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { useRouter } from 'expo-router';

interface SubsystemNav {
  id: string;
  title: string;
  desc: string;
  icon: string;
  route: string;
  rank: number;
}

const subsystems: SubsystemNav[] = [
  { id: 'hub', title: 'The Hub Super-App', desc: 'Unified multi-persona cockpit & live KPIs', icon: '⚡', route: '/(app)/hub', rank: 17 },
  { id: 'sos', title: 'SOS Emergency Response', desc: '1-Tap distress trigger & dispatcher radar', icon: '🚨', route: '/(app)/sos', rank: 14 },
  { id: 'whistleblower', title: 'Anonymous Whistleblower', desc: 'Zero-knowledge anti-ragging & grievance vault', icon: '🛡️', route: '/(app)/whistleblower', rank: 15 },
  { id: 'portal', title: 'Public Web Portal', desc: 'Academic catalog & admissions leads', icon: '🌐', route: '/(app)/portal', rank: 16 },
  { id: 'attendance', title: 'Attendance Management', desc: 'Geofenced biometric & roll call sessions', icon: '📍', route: '/(app)/attendance', rank: 4 },
  { id: 'billing', title: 'Billing & Ledger', desc: 'Online fee payments & invoice receipts', icon: '💳', route: '/(app)/billing', rank: 5 },
  { id: 'notices', title: 'Campus Notices', desc: 'Targeted announcements & push alerts', icon: '📢', route: '/(app)/notices', rank: 6 },
  { id: 'hostel', title: 'Hostel Management', desc: 'Room beds, gate passes & curfew tracking', icon: '🏢', route: '/(app)/hostel', rank: 7 },
  { id: 'mess', title: 'Mess Management', desc: 'Meal tokens, QR dining & rebate claims', icon: '🍽️', route: '/(app)/mess', rank: 8 },
  { id: 'library', title: 'E-Library System', desc: 'Catalog search, book borrow & e-reader', icon: '📚', route: '/(app)/library', rank: 9 },
  { id: 'studyhub', title: 'Study Hub & LMS', desc: 'Courses, assignments & peer reviews', icon: '📝', route: '/(app)/studyhub', rank: 10 },
  { id: 'mentorship', title: 'Mentorship Tracker', desc: 'Faculty mentor reviews & progress KPIs', icon: '🌱', route: '/(app)/mentorship', rank: 11 },
  { id: 'events', title: 'Events Organisation', desc: 'Festival schedules & QR gate check-ins', icon: '🎉', route: '/(app)/events', rank: 12 },
  { id: 'helpdesk', title: 'Helpdesk & Query', desc: 'SLA ticket triage & 1-tap urgent escalation', icon: '🎫', route: '/(app)/helpdesk', rank: 13 },
  { id: 'onboarding', title: 'Student Onboarding', desc: 'Digital admissions KYC & sequence generator', icon: '🎓', route: '/(app)/onboarding', rank: 3 },
  { id: 'audit', title: 'Audit Ledger', desc: 'SHA-256 Merkle compliance verification', icon: '🔒', route: '/(app)/audit', rank: 2 },
];

export default function HubHomeScreen() {
  const router = useRouter();

  return (
    <ScrollView className="flex-1 bg-slate-900 p-4">
      {/* Super-App Header Banner */}
      <View className="p-5 bg-gradient-to-r from-slate-900 via-indigo-950 to-blue-950 border border-indigo-800 rounded-3xl mb-4 space-y-2 shadow-lg">
        <View className="flex-row items-center justify-between">
          <View className="bg-emerald-500/20 px-2.5 py-0.5 rounded-full border border-emerald-500/40">
            <Text className="text-[10px] font-black text-emerald-300">17 MODULAR SUBSYSTEMS ACTIVE</Text>
          </View>
          <Text className="text-[10px] text-slate-300 font-mono">v2.0.0-PROD</Text>
        </View>

        <Text className="text-xl font-black text-white">Campus Management System</Text>
        <Text className="text-xs text-slate-300 leading-relaxed">
          Enterprise strict modular monolith with topological persona cockpits and tamper-evident audit ledger.
        </Text>

        <TouchableOpacity
          onPress={() => router.push('/(app)/hub')}
          className="bg-indigo-600 py-2.5 px-4 rounded-xl items-center mt-2 self-start shadow-md flex-row space-x-2"
        >
          <Text className="text-xs font-black text-white">Open The Hub Cockpit →</Text>
        </TouchableOpacity>
      </View>

      <Text className="text-xs font-black uppercase text-slate-400 mb-3 px-1 tracking-wider">
        Topological Subsystems Directory
      </Text>

      {/* Subsystems List */}
      <View className="space-y-2.5 mb-8">
        {subsystems.map((sub) => (
          <TouchableOpacity
            key={sub.id}
            onPress={() => router.push(sub.route as any)}
            activeOpacity={0.7}
            className="p-3.5 bg-slate-800 rounded-2xl border border-slate-700 flex-row items-center justify-between shadow-sm active:bg-slate-700"
          >
            <View className="flex-row items-center space-x-3 flex-1 mr-2">
              <Text className="text-2xl">{sub.icon}</Text>
              <View className="flex-1">
                <View className="flex-row items-center space-x-2">
                  <Text className="text-xs font-bold text-white">{sub.title}</Text>
                  <View className="bg-slate-700 px-1.5 py-0.2 rounded">
                    <Text className="text-[9px] font-mono text-slate-300">R{sub.rank}</Text>
                  </View>
                </View>
                <Text className="text-[10px] text-slate-400" numberOfLines={1}>
                  {sub.desc}
                </Text>
              </View>
            </View>
            <Text className="text-xs font-bold text-indigo-400">Launch →</Text>
          </TouchableOpacity>
        ))}
      </View>
    </ScrollView>
  );
}
