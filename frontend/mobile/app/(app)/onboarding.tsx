/**
 * BLOCK_MOBILE_ONBOARDING_SCREEN_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Mobile KYC submission and real-time admissions verification status tracking.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, TextInput } from 'react-native';

export default function MobileOnboardingScreen() {
  const [activeTab, setActiveTab] = useState<'REGISTER' | 'STATUS'>('REGISTER');
  const [category, setCategory] = useState<'STUDENT' | 'FACULTY' | 'STAFF'>('STUDENT');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [program, setProgram] = useState('BTECH_CSE');
  const [isSubmitted, setIsSubmitted] = useState(false);

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Admissions & Registration</Text>
        <Text className="text-xs text-slate-300">
          Rank 3 Onboarding · Fast-Track Digital KYC & Roll Number Generation
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('REGISTER')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'REGISTER' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'REGISTER' ? 'text-slate-900' : 'text-slate-600'}`}>
            New KYC Application
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('STATUS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'STATUS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'STATUS' ? 'text-slate-900' : 'text-slate-600'}`}>
            Tracking Status (2)
          </Text>
        </TouchableOpacity>
      </View>

      {activeTab === 'REGISTER' ? (
        isSubmitted ? (
          <View className="p-5 bg-emerald-50 border border-emerald-200 rounded-2xl space-y-3">
            <Text className="text-base font-bold text-emerald-900">✓ Application Submitted!</Text>
            <Text className="text-xs text-emerald-700">
              Your application for {firstName} {lastName} ({category}) is queued for verification.
            </Text>
            <TouchableOpacity
              onPress={() => { setIsSubmitted(false); setFirstName(''); setLastName(''); }}
              className="mt-2 py-2 px-4 bg-emerald-600 rounded-xl items-center"
            >
              <Text className="text-xs font-bold text-white">Submit Another Form</Text>
            </TouchableOpacity>
          </View>
        ) : (
          <View className="p-5 bg-white rounded-2xl border border-slate-200 space-y-4">
            <Text className="text-sm font-bold text-slate-900">Registration Category</Text>
            <View className="flex-row gap-2">
              {(['STUDENT', 'FACULTY', 'STAFF'] as const).map((cat) => (
                <TouchableOpacity
                  key={cat}
                  onPress={() => setCategory(cat)}
                  className={`flex-1 py-2 rounded-xl items-center border ${category === cat ? 'bg-slate-900 border-slate-900' : 'border-slate-200 bg-slate-50'}`}
                >
                  <Text className={`text-xs font-bold ${category === cat ? 'text-white' : 'text-slate-700'}`}>
                    {cat}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text className="text-xs font-bold text-slate-500 uppercase mt-2">Personal Details</Text>
            <TextInput
              placeholder="First Name *"
              value={firstName}
              onChangeText={setFirstName}
              className="h-11 px-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
            />
            <TextInput
              placeholder="Last Name *"
              value={lastName}
              onChangeText={setLastName}
              className="h-11 px-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
            />
            <TextInput
              placeholder="Email Address *"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              className="h-11 px-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
            />
            <TextInput
              placeholder="Mobile Phone Number *"
              value={phone}
              onChangeText={setPhone}
              keyboardType="phone-pad"
              className="h-11 px-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
            />

            <Text className="text-xs font-bold text-slate-500 uppercase mt-2">Target Program / Dept</Text>
            <TextInput
              placeholder="Program Code (e.g. BTECH_CSE)"
              value={program}
              onChangeText={setProgram}
              className="h-11 px-3 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-900"
            />

            <TouchableOpacity
              onPress={() => setIsSubmitted(true)}
              className="mt-4 py-3 bg-slate-900 rounded-xl items-center"
            >
              <Text className="text-xs font-bold text-white uppercase tracking-wider">Submit Application →</Text>
            </TouchableOpacity>
          </View>
        )
      ) : (
        <View className="space-y-3">
          <View className="p-4 bg-white rounded-xl border border-slate-200 space-y-2">
            <View className="flex-row justify-between items-center">
              <Text className="text-sm font-bold text-slate-900">Aarav Sharma</Text>
              <View className="px-2 py-0.5 rounded-full bg-blue-100">
                <Text className="text-[10px] font-bold text-blue-800">UNDER REVIEW</Text>
              </View>
            </View>
            <Text className="text-xs text-slate-500">Student · B.Tech CSE (2026-2027)</Text>
            <Text className="text-[11px] text-slate-400">Docs Attached: Aadhaar, 12th Marksheet, Photo</Text>
          </View>

          <View className="p-4 bg-white rounded-xl border border-slate-200 space-y-2">
            <View className="flex-row justify-between items-center">
              <Text className="text-sm font-bold text-slate-900">Dr. Ananya Iyer</Text>
              <View className="px-2 py-0.5 rounded-full bg-emerald-100">
                <Text className="text-[10px] font-bold text-emerald-800">ENROLLED</Text>
              </View>
            </View>
            <Text className="text-xs text-slate-500">Faculty · Employee ID: EMP-CSE-0001</Text>
            <Text className="text-[11px] text-emerald-600 font-semibold">Verified & Active Profile</Text>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
