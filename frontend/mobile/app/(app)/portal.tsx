/**
 * BLOCK_MOBILE_PORTAL_SCREEN_001
 * Subsystem: Rank 16 - Public Web Portal (portal)
 * Purpose:   Mobile Institutional Prospectus, Academic Program Catalog, and 1-Tap Admission Inquiry Lead Capture.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface ProgramItem {
  code: string;
  name: string;
  degree: 'UG' | 'PG' | 'PHD';
  dept: string;
  duration: string;
  fee: number;
}

export default function MobilePortalScreen() {
  const [selectedDegree, setSelectedDegree] = useState<string>('ALL');
  const [searchQuery, setSearchQuery] = useState<string>('');

  // Inquiry modal state
  const [inquiryProgram, setInquiryProgram] = useState<string | null>(null);
  const [name, setName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [phone, setPhone] = useState<string>('');
  const [message, setMessage] = useState<string>('');

  const programs: ProgramItem[] = [
    {
      code: 'BTECH_CSE_AI',
      name: 'B.Tech in Computer Science & AI',
      degree: 'UG',
      dept: 'Computer Science & Eng',
      duration: '4 Years (8 Sems)',
      fee: 250000,
    },
    {
      code: 'BTECH_ECE_VLSI',
      name: 'B.Tech in Electronics & VLSI Design',
      degree: 'UG',
      dept: 'Electronics & Comm',
      duration: '4 Years (8 Sems)',
      fee: 220000,
    },
    {
      code: 'MTECH_CYBER_SEC',
      name: 'M.Tech in Cybersecurity & Cloud Resilience',
      degree: 'PG',
      dept: 'Computer Science & Eng',
      duration: '2 Years (4 Sems)',
      fee: 180000,
    },
    {
      code: 'PHD_QUANTUM_COMP',
      name: 'Ph.D. in Quantum Information Systems',
      degree: 'PHD',
      dept: 'Interdisciplinary Quantum Lab',
      duration: '3 Years (6 Sems)',
      fee: 60000,
    },
  ];

  const filteredPrograms = programs.filter((p) => {
    if (selectedDegree !== 'ALL' && p.degree !== selectedDegree) return false;
    if (searchQuery.trim() !== '') {
      const q = searchQuery.toLowerCase();
      return p.name.toLowerCase().includes(q) || p.code.toLowerCase().includes(q) || p.dept.toLowerCase().includes(q);
    }
    return true;
  });

  const handleSendInquiry = () => {
    if (!name.trim() || !email.trim() || !phone.trim()) {
      Alert.alert('Required Fields', 'Please provide your name, email, and phone number.');
      return;
    }
    Alert.alert(
      'Inquiry Submitted!',
      `Thank you ${name}! Our admissions desk has received your inquiry for ${inquiryProgram} and will contact you within 24 hours.`,
      [{ text: 'OK', onPress: () => setInquiryProgram(null) }]
    );
    setName('');
    setEmail('');
    setPhone('');
    setMessage('');
  };

  return (
    <ScrollView className="flex-1 bg-slate-900 p-4">
      {/* Hero Banner */}
      <View className="p-5 bg-gradient-to-r from-blue-950 to-indigo-950 border border-indigo-800 rounded-2xl mb-4 space-y-2">
        <View className="bg-emerald-500/20 border border-emerald-500/40 px-2.5 py-0.5 rounded-full self-start">
          <Text className="text-[10px] font-bold text-emerald-300">ADMISSIONS OPEN: FALL 2026</Text>
        </View>
        <Text className="text-xl font-black text-white">Empowering Future Leaders & Engineers</Text>
        <Text className="text-xs text-slate-300 leading-relaxed">
          Explore 28+ world-class accredited undergraduate, postgraduate, and doctoral degree tracks.
        </Text>
        <TouchableOpacity
          onPress={() => setInquiryProgram('General Admissions')}
          className="bg-white py-2 px-4 rounded-xl items-center mt-2 self-start shadow-md"
        >
          <Text className="text-xs font-black text-slate-900">Request Prospectus & Info →</Text>
        </TouchableOpacity>
      </View>

      {/* Search & Filter */}
      <View className="space-y-3 mb-4">
        <TextInput
          placeholder="Search programs, departments..."
          placeholderTextColor="#94a3b8"
          value={searchQuery}
          onChangeText={setSearchQuery}
          className="border border-slate-700 rounded-xl p-3 text-xs bg-slate-800 text-white"
        />

        <View className="flex-row gap-2">
          {['ALL', 'UG', 'PG', 'PHD'].map((d) => (
            <TouchableOpacity
              key={d}
              onPress={() => setSelectedDegree(d)}
              className={`px-3 py-1.5 rounded-lg border ${selectedDegree === d ? 'bg-indigo-600 border-indigo-600' : 'bg-slate-800 border-slate-700'}`}
            >
              <Text className={`text-xs font-bold ${selectedDegree === d ? 'text-white' : 'text-slate-400'}`}>
                {d === 'ALL' ? 'All Degrees' : d}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Programs List */}
      <View className="space-y-3 mb-6">
        {filteredPrograms.map((prog) => (
          <View
            key={prog.code}
            className="bg-slate-800 p-4 rounded-2xl border border-slate-700 space-y-2"
          >
            <View className="flex-row justify-between items-center">
              <View className="bg-indigo-950 px-2 py-0.5 rounded border border-indigo-400">
                <Text className="font-mono text-[10px] font-bold text-indigo-300">{prog.code}</Text>
              </View>
              <Text className="text-[10px] font-bold text-emerald-400">
                ₹{(prog.fee / 100000).toFixed(2)} Lakh / yr
              </Text>
            </View>

            <Text className="text-sm font-bold text-white leading-snug">{prog.name}</Text>

            <View className="space-y-0.5">
              <Text className="text-[11px] text-slate-400">🏛️ {prog.dept}</Text>
              <Text className="text-[11px] text-slate-400">⏱️ {prog.duration}</Text>
            </View>

            <TouchableOpacity
              onPress={() => setInquiryProgram(prog.name)}
              className="bg-indigo-600 py-2 rounded-xl items-center mt-1"
            >
              <Text className="text-white text-xs font-bold">Apply / Inquire</Text>
            </TouchableOpacity>
          </View>
        ))}
      </View>

      {/* Inquiry Form Modal Sheet */}
      {inquiryProgram && (
        <View className="bg-slate-800 p-5 rounded-2xl border-2 border-indigo-500 space-y-3 mb-8">
          <View className="flex-row justify-between items-center border-b border-slate-700 pb-2">
            <View>
              <Text className="text-xs font-black text-indigo-400 uppercase">Admissions Lead Form</Text>
              <Text className="text-xs font-bold text-white">{inquiryProgram}</Text>
            </View>
            <TouchableOpacity onPress={() => setInquiryProgram(null)}>
              <Text className="text-slate-400 font-bold">✕</Text>
            </TouchableOpacity>
          </View>

          <TextInput
            placeholder="Full Name *"
            placeholderTextColor="#94a3b8"
            value={name}
            onChangeText={setName}
            className="border border-slate-600 rounded-xl p-2.5 text-xs bg-slate-900 text-white"
          />

          <TextInput
            placeholder="Email Address *"
            placeholderTextColor="#94a3b8"
            value={email}
            onChangeText={setEmail}
            keyboardType="email-address"
            className="border border-slate-600 rounded-xl p-2.5 text-xs bg-slate-900 text-white"
          />

          <TextInput
            placeholder="Phone Number *"
            placeholderTextColor="#94a3b8"
            value={phone}
            onChangeText={setPhone}
            keyboardType="phone-pad"
            className="border border-slate-600 rounded-xl p-2.5 text-xs bg-slate-900 text-white"
          />

          <TextInput
            placeholder="Questions regarding scholarships, hostel, etc."
            placeholderTextColor="#94a3b8"
            value={message}
            onChangeText={setMessage}
            multiline
            numberOfLines={3}
            className="border border-slate-600 rounded-xl p-2.5 text-xs bg-slate-900 text-white min-h-[60px]"
          />

          <TouchableOpacity
            onPress={handleSendInquiry}
            className="bg-emerald-600 p-3 rounded-xl items-center shadow-md"
          >
            <Text className="text-white text-xs font-bold">Submit Admissions Inquiry</Text>
          </TouchableOpacity>
        </View>
      )}
    </ScrollView>
  );
}
