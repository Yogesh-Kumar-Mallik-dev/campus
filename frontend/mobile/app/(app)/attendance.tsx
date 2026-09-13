/**
 * BLOCK_MOBILE_ATTENDANCE_SCREEN_001
 * Subsystem: Rank 4 - Attendance Management System (attendance)
 * Purpose:   Mobile student attendance percentage tracker, biometric status, shortage alerts (<75%), and medical leave submission.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, TextInput, Alert } from 'react-native';

export default function MobileAttendanceScreen() {
  const [activeTab, setActiveTab] = useState<'OVERVIEW' | 'LEAVE'>('OVERVIEW');
  const [fromDate, setFromDate] = useState('2026-09-14');
  const [toDate, setToDate] = useState('2026-09-15');
  const [reason, setReason] = useState('');
  const [isLeaveSubmitted, setIsLeaveSubmitted] = useState(false);

  const subjects = [
    { code: 'CS101', name: 'Intro to CS', attended: 28, total: 30, pct: 93.3, shortage: false },
    { code: 'MATH201', name: 'Discrete Maths', attended: 22, total: 32, pct: 68.7, shortage: true },
    { code: 'ENG102', name: 'Technical Writing', attended: 18, total: 20, pct: 90.0, shortage: false },
    { code: 'PHY105', name: 'Applied Physics', attended: 15, total: 24, pct: 62.5, shortage: true },
  ];

  const overallAttended = subjects.reduce((acc, s) => acc + s.attended, 0);
  const overallTotal = subjects.reduce((acc, s) => acc + s.total, 0);
  const overallPct = Math.round((overallAttended / overallTotal) * 100);

  const handleApplyLeave = () => {
    if (!reason.trim()) {
      Alert.alert('Validation Error', 'Please enter a medical reason.');
      return;
    }
    setIsLeaveSubmitted(true);
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Attendance Tracker</Text>
        <Text className="text-xs text-slate-300">
          Rank 4 Attendance · Multi-Mode Sync & Shortage Monitor (&lt;75%)
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('OVERVIEW')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'OVERVIEW' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'OVERVIEW' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Attendance
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('LEAVE')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'LEAVE' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'LEAVE' ? 'text-slate-900' : 'text-slate-600'}`}>
            Medical Leave
          </Text>
        </TouchableOpacity>
      </View>

      {activeTab === 'OVERVIEW' ? (
        <View className="space-y-4">
          {/* Overall Attendance Card */}
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Overall Aggregate</Text>
              <View className={`px-2 py-0.5 rounded-full ${overallPct >= 75 ? 'bg-emerald-100' : 'bg-rose-100'}`}>
                <Text className={`text-[10px] font-bold ${overallPct >= 75 ? 'text-emerald-800' : 'text-rose-800'}`}>
                  {overallPct >= 75 ? 'COMPLIANT' : 'SHORTAGE ALERT'}
                </Text>
              </View>
            </View>

            <View className="flex-row items-baseline gap-2">
              <Text className={`text-4xl font-extrabold ${overallPct >= 75 ? 'text-slate-900' : 'text-rose-600'}`}>
                {overallPct}%
              </Text>
              <Text className="text-xs text-slate-500">
                ({overallAttended} / {overallTotal} classes attended)
              </Text>
            </View>

            <View className="w-full bg-slate-100 h-2.5 rounded-full mt-3 overflow-hidden">
              <View
                className={`h-full rounded-full ${overallPct >= 75 ? 'bg-emerald-500' : 'bg-rose-500'}`}
                style={{ width: `${overallPct}%` }}
              />
            </View>
            <Text className="text-[11px] text-slate-400 mt-2">
              Institutions require min 75.0% for end-semester exam admit card generation.
            </Text>
          </View>

          {/* Subject-Wise Breakdown */}
          <View className="space-y-2">
            <Text className="text-sm font-bold text-slate-900 px-1">Subject Wise Breakdown</Text>
            {subjects.map((sub) => (
              <View key={sub.code} className="p-4 bg-white rounded-xl border border-slate-200 flex-row justify-between items-center">
                <View>
                  <View className="flex-row items-center gap-1.5">
                    <Text className="text-xs font-bold text-slate-900">{sub.code}</Text>
                    <Text className="text-xs text-slate-500">· {sub.name}</Text>
                  </View>
                  <Text className="text-[11px] text-slate-400 mt-0.5">
                    {sub.attended} of {sub.total} classes ({sub.total - sub.attended} absent)
                  </Text>
                </View>

                <View className="items-end">
                  <Text className={`text-base font-bold ${sub.shortage ? 'text-rose-600' : 'text-emerald-600'}`}>
                    {sub.pct.toFixed(1)}%
                  </Text>
                  {sub.shortage && (
                    <Text className="text-[10px] font-bold text-rose-500">Shortage (&lt;75%)</Text>
                  )}
                </View>
              </View>
            ))}
          </View>
        </View>
      ) : (
        /* Medical Leave Form */
        isLeaveSubmitted ? (
          <View className="p-5 bg-emerald-50 border border-emerald-200 rounded-2xl space-y-3">
            <Text className="text-base font-bold text-emerald-900">✓ Medical Application Filed!</Text>
            <Text className="text-xs text-emerald-700 leading-relaxed">
              Your leave certificate for dates {fromDate} to {toDate} has been transmitted to the Academic Dean for condonation review.
            </Text>
            <TouchableOpacity
              onPress={() => { setIsLeaveSubmitted(false); setReason(''); }}
              className="mt-2 py-2 px-4 bg-emerald-600 rounded-xl items-center"
            >
              <Text className="text-xs font-bold text-white">Submit Another Request</Text>
            </TouchableOpacity>
          </View>
        ) : (
          <View className="p-5 bg-white rounded-2xl border border-slate-200 space-y-4">
            <Text className="text-sm font-bold text-slate-900">Apply for Medical Condonation</Text>
            <Text className="text-xs text-slate-500">
              Approved medical leaves automatically convert absent records to excused attendance.
            </Text>

            <View className="space-y-1">
              <Text className="text-xs font-medium text-slate-600">From Date (YYYY-MM-DD)</Text>
              <TextInput
                value={fromDate}
                onChangeText={setFromDate}
                className="p-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
              />
            </View>

            <View className="space-y-1">
              <Text className="text-xs font-medium text-slate-600">To Date (YYYY-MM-DD)</Text>
              <TextInput
                value={toDate}
                onChangeText={setToDate}
                className="p-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
              />
            </View>

            <View className="space-y-1">
              <Text className="text-xs font-medium text-slate-600">Medical Reason / Hospitalization</Text>
              <TextInput
                value={reason}
                onChangeText={setReason}
                placeholder="e.g. Hospitalized due to viral infection"
                placeholderTextColor="#94a3b8"
                multiline
                numberOfLines={3}
                className="p-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900 min-h-[80px]"
              />
            </View>

            <TouchableOpacity
              onPress={handleApplyLeave}
              className="w-full py-3 bg-slate-900 rounded-xl items-center shadow-sm active:bg-slate-800"
            >
              <Text className="text-xs font-bold text-white">Submit Medical Certificate</Text>
            </TouchableOpacity>
          </View>
        )
      )}
    </ScrollView>
  );
}
