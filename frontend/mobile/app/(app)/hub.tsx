/**
 * BLOCK_MOBILE_HUB_SCREEN_001
 * Subsystem: Rank 17 - The Hub Root Super-App (hub)
 * Purpose:   Mobile Super-App Cockpit with Persona Switcher, Live KPIs, 1-Tap Action Shortcuts, and Subsystem Telemetry.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';

type Persona = 'STUDENT' | 'FACULTY' | 'WARDEN' | 'ADMIN';

interface Shortcut {
  id: string;
  title: string;
  desc: string;
  icon: string;
  route: string;
}

export default function MobileHubScreen() {
  const [persona, setPersona] = useState<Persona>('STUDENT');

  const getKPIs = () => {
    switch (persona) {
      case 'STUDENT':
        return [
          { label: 'Attendance', val: '89.2%', sub: 'Eligible', color: 'text-emerald-400' },
          { label: 'Pending Dues', val: '₹25,000', sub: '1 Invoice', color: 'text-amber-400' },
          { label: 'Gate Pass', val: 'APPROVED', sub: 'Curfew 9:30 PM', color: 'text-indigo-400' },
          { label: 'Library', val: '2 Issued', sub: '0 Overdue', color: 'text-blue-400' },
        ];
      case 'FACULTY':
        return [
          { label: 'Roll Call', val: '94.5%', sub: 'Today Avg', color: 'text-emerald-400' },
          { label: 'Grading', val: '18 Pending', sub: 'Study Hub', color: 'text-amber-400' },
          { label: 'Mentees', val: '8 Active', sub: 'Mentorship', color: 'text-indigo-400' },
          { label: 'Circulars', val: '2 Unread', sub: 'Notices', color: 'text-blue-400' },
        ];
      case 'WARDEN':
        return [
          { label: 'Gate Passes', val: '14 Pending', sub: 'Action Req', color: 'text-rose-400' },
          { label: 'Occupancy', val: '96.4%', sub: 'Hostel Block', color: 'text-indigo-400' },
          { label: 'Curfew Log', val: '0 Breaches', sub: 'Last 24h', color: 'text-emerald-400' },
          { label: 'Helpdesk', val: '4 Tickets', sub: 'Maintenance', color: 'text-amber-400' },
        ];
      default:
        return [
          { label: 'Active SOS', val: '0 Active', sub: 'All Safe', color: 'text-emerald-400' },
          { label: 'Fee Inflow', val: '₹10.5L', sub: 'This Month', color: 'text-emerald-400' },
          { label: 'Onboarding', val: '42 Verified', sub: 'Applicants', color: 'text-blue-400' },
          { label: 'Audit Trail', val: 'SYNCED', sub: 'Merkle Check', color: 'text-indigo-400' },
        ];
    }
  };

  const getShortcuts = (): Shortcut[] => {
    switch (persona) {
      case 'STUDENT':
        return [
          { id: '1', title: 'SOS Emergency', desc: '1-Tap distress trigger', icon: '🚨', route: '/sos' },
          { id: '2', title: 'Mark Attendance', desc: 'Geofence check-in', icon: '📍', route: '/attendance' },
          { id: '3', title: 'Pay Fee Online', desc: 'UPI & Card gateway', icon: '💳', route: '/billing' },
          { id: '4', title: 'Gate Pass Apply', desc: 'Night out permit', icon: '🚪', route: '/hostel' },
          { id: '5', title: 'Mess Dining QR', desc: 'Scan meal coupon', icon: '🍽️', route: '/mess' },
          { id: '6', title: 'Whistleblower', desc: 'Confidential report', icon: '🛡️', route: '/whistleblower' },
        ];
      case 'FACULTY':
        return [
          { id: 'f1', title: 'Class Roll Call', desc: 'Start QR session', icon: '📋', route: '/attendance' },
          { id: 'f2', title: 'Grade Homework', desc: '18 submissions', icon: '📝', route: '/studyhub' },
          { id: 'f3', title: 'Broadcast Notice', desc: 'Publish circular', icon: '📢', route: '/notices' },
          { id: 'f4', title: 'Mentee Reviews', desc: 'Academic scores', icon: '🌱', route: '/mentorship' },
        ];
      case 'WARDEN':
        return [
          { id: 'w1', title: 'Approve Passes', desc: '14 pending requests', icon: '🔑', route: '/hostel' },
          { id: 'w2', title: 'Curfew Infraction', desc: 'Disciplinary log', icon: '⚠️', route: '/hostel' },
          { id: 'w3', title: 'Bed Occupancy', desc: 'Room allocations', icon: '🛏️', route: '/hostel' },
        ];
      default:
        return [
          { id: 'a1', title: 'Applicant Pipeline', desc: 'Review admissions', icon: '🎓', route: '/onboarding' },
          { id: 'a2', title: 'Ledger Accounts', desc: 'Financial audit', icon: '📊', route: '/billing' },
          { id: 'a3', title: 'Audit Merkle Log', desc: 'Compliance proofs', icon: '🔒', route: '/audit' },
          { id: 'a4', title: 'Public Portal CRM', desc: 'Lead triage desk', icon: '🌐', route: '/portal' },
        ];
    }
  };

  const handleShortcutPress = (sc: Shortcut) => {
    Alert.alert(sc.title, `Triggered 1-Tap shortcut for ${sc.title}.\nNavigating to ${sc.route}`);
  };

  return (
    <ScrollView className="flex-1 bg-slate-900 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-gradient-to-r from-slate-900 via-indigo-950 to-blue-950 border border-indigo-800 rounded-3xl mb-4 space-y-2 shadow-lg">
        <View className="flex-row items-center justify-between">
          <View className="bg-indigo-500/20 px-2.5 py-0.5 rounded-full border border-indigo-400">
            <Text className="text-[10px] font-black text-indigo-300">CAMPUS SUPER-APP HUB</Text>
          </View>
          <View className="flex-row items-center space-x-1">
            <View className="w-2 h-2 rounded-full bg-emerald-400" />
            <Text className="text-[10px] text-emerald-400 font-bold">16 Subsystems Online</Text>
          </View>
        </View>

        <Text className="text-xl font-black text-white">Unified Cockpit</Text>
        <Text className="text-xs text-slate-300">
          Switched Persona: <Text className="font-bold text-white uppercase">{persona}</Text>
        </Text>

        {/* Persona Selector Tabs */}
        <View className="flex-row bg-slate-800/80 p-1 rounded-xl mt-2">
          {(['STUDENT', 'FACULTY', 'WARDEN', 'ADMIN'] as Persona[]).map((p) => (
            <TouchableOpacity
              key={p}
              onPress={() => setPersona(p)}
              className={`flex-1 py-1.5 rounded-lg items-center ${persona === p ? 'bg-indigo-600' : ''}`}
            >
              <Text className={`text-[10px] font-bold ${persona === p ? 'text-white' : 'text-slate-400'}`}>
                {p}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Live KPIs Grid */}
      <View className="flex-row flex-wrap gap-2 mb-4">
        {getKPIs().map((kpi, idx) => (
          <View
            key={idx}
            className="flex-1 min-w-[45%] bg-slate-800 p-3.5 rounded-2xl border border-slate-700 space-y-0.5"
          >
            <Text className="text-[10px] font-bold text-slate-400 uppercase">{kpi.label}</Text>
            <Text className={`text-base font-black ${kpi.color}`}>{kpi.val}</Text>
            <Text className="text-[10px] text-slate-400">{kpi.sub}</Text>
          </View>
        ))}
      </View>

      {/* 1-Tap Action Shortcuts */}
      <View className="space-y-3 mb-6">
        <Text className="text-xs font-black text-slate-300 uppercase tracking-wider">
          ⚡ 1-Tap Quick Action Shortcuts ({persona})
        </Text>

        <View className="space-y-2">
          {getShortcuts().map((sc) => (
            <TouchableOpacity
              key={sc.id}
              onPress={() => handleShortcutPress(sc)}
              className="bg-slate-800 p-3.5 rounded-2xl border border-slate-700 flex-row items-center space-x-3 active:bg-slate-700"
            >
              <Text className="text-2xl">{sc.icon}</Text>
              <View className="flex-1">
                <Text className="text-xs font-bold text-white">{sc.title}</Text>
                <Text className="text-[10px] text-slate-400">{sc.desc}</Text>
              </View>
              <Text className="font-mono text-[10px] text-indigo-400 font-bold">{sc.route} →</Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Subsystem Telemetry Map */}
      <View className="bg-slate-800 p-4 rounded-2xl border border-slate-700 space-y-2 mb-8">
        <Text className="text-xs font-black text-slate-300">Topological Telemetry Radar</Text>
        <Text className="text-[10px] text-slate-400">
          Central Auth, Audit Merkle, Onboarding, Attendance, Billing, Notices, Hostel, Mess, Library, Study Hub, Mentorship, Events, Helpdesk, SOS, Whistleblower, Portal.
        </Text>
        <View className="p-2 bg-emerald-950 border border-emerald-700 rounded-xl flex-row justify-between items-center mt-1">
          <Text className="text-[10px] font-bold text-emerald-300">✓ All 16 Modular Engines Healthy</Text>
          <Text className="text-[10px] text-emerald-400 font-mono">0.42ms RTT</Text>
        </View>
      </View>
    </ScrollView>
  );
}
