/**
 * BLOCK_MOBILE_HUB_SCREEN_001
 * Purpose: Mobile Hub Dashboard providing navigation cards to role workspaces and features.
 */

import React from 'react';
import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { useRouter } from 'expo-router';

const workspaces = [
  { id: 'student', title: 'Student Hub', desc: 'Timetable, attendance, mess card', icon: '🎓' },
  { id: 'teacher', title: 'Teacher Hub', desc: 'Roll call, gradebook, syllabus', icon: '👨‍🏫' },
  { id: 'admin', title: 'Admin Hub', desc: 'User accounts, fee structures', icon: '⚙️' },
  { id: 'warden', title: 'Warden Hub', desc: 'Hostel rooms, gate pass approvals', icon: '🏢' },
  { id: 'library', title: 'Library Hub', desc: 'Book checkouts, digital library', icon: '📚' },
  { id: 'parent', title: 'Parent Hub', desc: 'Attendance, fee payment, counseling', icon: '👨‍👩‍👧' },
];

export default function HubHomeScreen() {
  const router = useRouter();

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Campus Management System</Text>
        <Text className="text-xs text-slate-300">
          Rank 1 IAM & Rank 2 Audit Active · 30,000 Total Capacity
        </Text>

        <View className="flex-row gap-2 mt-4 flex-wrap">
          <TouchableOpacity
            onPress={() => router.push('/(app)/onboarding')}
            className="px-3 py-2 bg-white rounded-xl"
          >
            <Text className="text-xs font-bold text-slate-900">Onboarding →</Text>
          </TouchableOpacity>
          <TouchableOpacity
            onPress={() => router.push('/(app)/audit')}
            className="px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl"
          >
            <Text className="text-xs font-bold text-white">Audit Logs →</Text>
          </TouchableOpacity>
          <TouchableOpacity
            onPress={() => router.push('/(app)/compliance')}
            className="px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl"
          >
            <Text className="text-xs font-bold text-white">Compliance →</Text>
          </TouchableOpacity>
        </View>
      </View>

      <Text className="text-sm font-bold uppercase text-slate-500 mb-3 px-1">Institutional Workspaces</Text>

      {/* Grid */}
      <View className="space-y-3">
        {workspaces.map((ws) => (
          <TouchableOpacity
            key={ws.id}
            activeOpacity={0.7}
            className="p-4 bg-white rounded-xl border border-slate-200 flex-row items-center justify-between"
          >
            <View className="flex-row items-center space-x-3 flex-1 mr-2">
              <Text className="text-2xl mr-3">{ws.icon}</Text>
              <View className="flex-1">
                <Text className="text-sm font-bold text-slate-900">{ws.title}</Text>
                <Text className="text-xs text-slate-500" numberOfLines={1}>{ws.desc}</Text>
              </View>
            </View>
            <Text className="text-xs font-bold text-blue-600">Open →</Text>
          </TouchableOpacity>
        ))}
      </View>
    </ScrollView>
  );
}
