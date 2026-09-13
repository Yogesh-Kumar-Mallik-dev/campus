/**
 * BLOCK_MOBILE_HOSTEL_SCREEN_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   Mobile residential accommodation hub, bed allotment card, digital gate pass application, and curfew alerts.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface MobileGatePass {
  id: string;
  reason: string;
  destination: string;
  expectedOutAt: string;
  expectedInAt: string;
  status: 'PENDING' | 'APPROVED' | 'OUT_CAMPUS' | 'RETURNED' | 'REJECTED';
}

export default function MobileHostelScreen() {
  const [activeTab, setActiveTab] = useState<'ROOM' | 'GATE_PASS' | 'RULES'>('ROOM');
  const [gatePasses, setGatePasses] = useState<MobileGatePass[]>([
    {
      id: 'gp_001',
      reason: 'Weekend Home Visit',
      destination: 'New Delhi (Home)',
      expectedOutAt: '2026-09-15 16:00',
      expectedInAt: '2026-09-17 20:00',
      status: 'APPROVED',
    },
    {
      id: 'gp_002',
      reason: 'Medical Consultation',
      destination: 'Apollo Hospital',
      expectedOutAt: '2026-09-10 11:00',
      expectedInAt: '2026-09-10 15:00',
      status: 'RETURNED',
    },
  ]);

  const [isApplying, setIsApplying] = useState(false);
  const [reason, setReason] = useState('');
  const [destination, setDestination] = useState('');

  const handleApply = () => {
    if (!reason || !destination) {
      Alert.alert('Incomplete Form', 'Please enter reason and destination.');
      return;
    }

    const newPass: MobileGatePass = {
      id: `gp_${Date.now()}`,
      reason,
      destination,
      expectedOutAt: '2026-09-18 10:00',
      expectedInAt: '2026-09-18 18:00',
      status: 'PENDING',
    };

    setGatePasses([newPass, ...gatePasses]);
    setIsApplying(false);
    setReason('');
    setDestination('');
    Alert.alert('Gate Pass Submitted', 'Your outing request has been forwarded to the Chief Warden for review.');
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Hostel & Residence Portal</Text>
        <Text className="text-xs text-slate-300">
          Rank 7 Hostel · Bed Allocation, Digital Gate Passes & Curfew Enforcement
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('ROOM')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'ROOM' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'ROOM' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Allotment
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('GATE_PASS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'GATE_PASS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'GATE_PASS' ? 'text-slate-900' : 'text-slate-600'}`}>
            Gate Passes ({gatePasses.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('RULES')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'RULES' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'RULES' ? 'text-slate-900' : 'text-slate-600'}`}>
            Curfew Rules
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: My Allotment */}
      {activeTab === 'ROOM' && (
        <View className="space-y-4">
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm">
            <View className="flex-row items-center justify-between mb-3">
              <View className="bg-blue-100 px-2.5 py-1 rounded-md">
                <Text className="text-[10px] font-bold text-blue-900">ACTIVE ALLOCATION</Text>
              </View>
              <Text className="text-xs font-semibold text-slate-500">Academic Year 2026-27</Text>
            </View>

            <Text className="text-xl font-bold text-slate-900 mb-1">Aryabhatta Block A (BH-A)</Text>
            <Text className="text-xs text-slate-500 mb-4">Floor 1 · Room 101 (Double Sharing - AC)</Text>

            <View className="bg-slate-50 p-3 rounded-xl border border-slate-100 space-y-2 mb-4">
              <View className="flex-row justify-between">
                <Text className="text-xs text-slate-500">Bed Number:</Text>
                <Text className="text-xs font-bold text-slate-900">A-101-1</Text>
              </View>
              <View className="flex-row justify-between">
                <Text className="text-xs text-slate-500">Roommate:</Text>
                <Text className="text-xs font-bold text-slate-900">Aditya Verma (2026CSE002)</Text>
              </View>
              <View className="flex-row justify-between">
                <Text className="text-xs text-slate-500">Hostel Fee per Semester:</Text>
                <Text className="text-xs font-bold text-slate-900">₹32,000 (Paid)</Text>
              </View>
            </View>

            <View className="p-3 bg-emerald-50 rounded-xl border border-emerald-200">
              <Text className="text-xs font-bold text-emerald-900">Chief Warden: Prof. Ananya Roy</Text>
              <Text className="text-[11px] text-emerald-700">Warden Office: Ground Floor Admin Block</Text>
            </View>
          </View>
        </View>
      )}

      {/* Tab: Gate Passes */}
      {activeTab === 'GATE_PASS' && (
        <View className="space-y-4">
          <TouchableOpacity
            onPress={() => setIsApplying(!isApplying)}
            className="p-3.5 bg-slate-900 rounded-xl items-center"
          >
            <Text className="text-xs font-bold text-white">
              {isApplying ? '✕ Cancel Application' : '+ Apply for Gate Pass'}
            </Text>
          </TouchableOpacity>

          {isApplying && (
            <View className="p-5 bg-white rounded-2xl border border-slate-300 shadow-sm space-y-3">
              <Text className="text-sm font-bold text-slate-900">New Gate Pass Request</Text>
              <View>
                <Text className="text-[11px] font-semibold text-slate-500 mb-1">Reason for Leave</Text>
                <TextInput
                  placeholder="e.g. Weekend visit home"
                  value={reason}
                  onChangeText={setReason}
                  className="border border-slate-200 rounded-lg p-2.5 text-xs bg-slate-50"
                />
              </View>
              <View>
                <Text className="text-[11px] font-semibold text-slate-500 mb-1">Destination Address</Text>
                <TextInput
                  placeholder="e.g. Delhi / Gurugram"
                  value={destination}
                  onChangeText={setDestination}
                  className="border border-slate-200 rounded-lg p-2.5 text-xs bg-slate-50"
                />
              </View>
              <TouchableOpacity
                onPress={handleApply}
                className="p-3 bg-emerald-700 rounded-lg items-center mt-2"
              >
                <Text className="text-xs font-bold text-white">Submit to Warden</Text>
              </TouchableOpacity>
            </View>
          )}

          {gatePasses.map((gp) => (
            <View key={gp.id} className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm">
              <View className="flex-row items-center justify-between mb-2">
                <Text className="text-sm font-bold text-slate-900">{gp.reason}</Text>
                <View
                  className={`px-2 py-0.5 rounded-md ${
                    gp.status === 'APPROVED'
                      ? 'bg-blue-100'
                      : gp.status === 'OUT_CAMPUS'
                      ? 'bg-purple-100'
                      : gp.status === 'RETURNED'
                      ? 'bg-emerald-100'
                      : 'bg-amber-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      gp.status === 'APPROVED'
                        ? 'text-blue-800'
                        : gp.status === 'OUT_CAMPUS'
                        ? 'text-purple-800'
                        : gp.status === 'RETURNED'
                        ? 'text-emerald-800'
                        : 'text-amber-800'
                    }`}
                  >
                    {gp.status}
                  </Text>
                </View>
              </View>

              <Text className="text-xs text-slate-500 mb-3">Destination: {gp.destination}</Text>

              <View className="border-t border-slate-100 pt-2 flex-row justify-between">
                <Text className="text-[11px] text-slate-400">Out: {gp.expectedOutAt}</Text>
                <Text className="text-[11px] text-slate-400">In: {gp.expectedInAt}</Text>
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Tab: Curfew Rules */}
      {activeTab === 'RULES' && (
        <View className="space-y-4">
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-3">
            <Text className="text-sm font-bold text-slate-900">Hostel Rules & Curfew Timings</Text>
            <View className="space-y-2 text-xs text-slate-600">
              <Text>• <Text className="font-bold text-slate-800">Main Gate Curfew:</Text> Strict entry closure at 10:30 PM.</Text>
              <Text>• <Text className="font-bold text-slate-800">Quiet Hours:</Text> 11:00 PM to 06:00 AM daily.</Text>
              <Text>• <Text className="font-bold text-slate-800">Outing Passes:</Text> Overnight leave requires warden authorization at least 24 hours in advance.</Text>
              <Text>• <Text className="font-bold text-slate-800">Late Return Fine:</Text> Automatic fine of ₹100 is logged upon curfew overrun.</Text>
            </View>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
