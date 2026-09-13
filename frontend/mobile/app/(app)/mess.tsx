/**
 * BLOCK_MOBILE_MESS_SCREEN_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   Mobile student dining portal, 1-tap QR meal punch, weekly nutrition menu, and leave fee rebates.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

export default function MobileMessScreen() {
  const [activeTab, setActiveTab] = useState<'TOKEN' | 'MENU' | 'REBATE'>('TOKEN');
  const [selectedMeal, setSelectedMeal] = useState<'BREAKFAST' | 'LUNCH' | 'SNACKS' | 'DINNER'>('LUNCH');
  const [tokenGenerated, setTokenGenerated] = useState<boolean>(false);
  const [tokenCode, setTokenCode] = useState<string>('TOK-9d7a2f10');
  const [tokenRedeemed, setTokenRedeemed] = useState<boolean>(false);

  // Rebate Form
  const [rebateDays, setRebateDays] = useState<string>('4');
  const [rebateReason, setRebateReason] = useState<string>('Sports Tournament Trip');

  const handleGenerateToken = () => {
    const hash = Math.random().toString(16).substring(2, 10);
    setTokenCode(`TOK-${hash}`);
    setTokenGenerated(true);
    setTokenRedeemed(false);
  };

  const handleSimulateScan = () => {
    setTokenRedeemed(true);
    Alert.alert('Meal Punch Success', `Your ${selectedMeal} meal token has been validated at North Mess entrance.`);
  };

  const handleApplyRebate = () => {
    const days = parseInt(rebateDays, 10) || 0;
    if (days < 3) {
      Alert.alert('Rebate Ineligible', 'Mess fee rebate requires minimum 3 continuous days of leave.');
      return;
    }
    const amount = days * 150;
    Alert.alert('Rebate Submitted', `Your request for ₹${amount} (${days} days @ ₹150/day) has been sent to the Mess Manager.`);
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Mess & Dining Hub</Text>
        <Text className="text-xs text-slate-300">
          Rank 8 Mess · 1-Tap QR Meal Token, Nutrition Menu & Leave Rebates
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('TOKEN')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'TOKEN' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'TOKEN' ? 'text-slate-900' : 'text-slate-600'}`}>
            QR Token
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('MENU')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'MENU' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'MENU' ? 'text-slate-900' : 'text-slate-600'}`}>
            Today's Menu
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('REBATE')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'REBATE' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'REBATE' ? 'text-slate-900' : 'text-slate-600'}`}>
            Fee Rebates
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: QR Token */}
      {activeTab === 'TOKEN' && (
        <View className="space-y-4">
          <!-- Meal Selector Pills -->
          <View className="flex-row justify-between bg-white p-2 rounded-xl border border-slate-200">
            {(['BREAKFAST', 'LUNCH', 'SNACKS', 'DINNER'] as const).map((meal) => (
              <TouchableOpacity
                key={meal}
                onPress={() => {
                  setSelectedMeal(meal);
                  setTokenGenerated(false);
                  setTokenRedeemed(false);
                }}
                className={`px-3 py-1.5 rounded-lg ${
                  selectedMeal === meal ? 'bg-slate-900' : 'bg-transparent'
                }`}
              >
                <Text
                  className={`text-[10px] font-bold ${
                    selectedMeal === meal ? 'text-white' : 'text-slate-600'
                  }`}
                >
                  {meal}
                </Text>
              </TouchableOpacity>
            ))}
          </View>

          {/* Token Card */}
          <View className="p-6 bg-white rounded-2xl border border-slate-200 shadow-sm items-center space-y-4">
            <View className="flex-row items-center justify-between w-full">
              <View className="bg-emerald-100 px-2.5 py-0.5 rounded-md">
                <Text className="text-[10px] font-bold text-emerald-800">ACTIVE SUBSCRIPTION</Text>
              </View>
              <Text className="text-xs text-slate-500 font-semibold">Central North Mess</Text>
            </View>

            {/* QR Mockup */}
            <View className="w-48 h-48 bg-slate-900 rounded-2xl p-4 items-center justify-center space-y-2">
              <View className="w-24 h-24 bg-white rounded-xl items-center justify-center">
                <Text className="text-3xl">🍱</Text>
              </View>
              <Text className="font-mono text-xs font-bold text-white tracking-widest">
                {tokenGenerated ? tokenCode : 'SCAN READY'}
              </Text>
            </View>

            <View className="items-center">
              <Text className="text-base font-bold text-slate-900">{selectedMeal} DINING PASS</Text>
              <Text className="text-xs text-slate-400">Valid for Today · 1 Serving</Text>
            </View>

            {tokenGenerated ? (
              tokenRedeemed ? (
                <View className="p-3 bg-emerald-50 border border-emerald-200 rounded-xl w-full items-center">
                  <Text className="text-xs font-bold text-emerald-900">✓ Token Redeemed & Verified</Text>
                  <Text className="text-[11px] text-emerald-700">Meal serving authorized at counter</Text>
                </View>
              ) : (
                <TouchableOpacity
                  onPress={handleSimulateScan}
                  className="p-3.5 bg-emerald-700 rounded-xl w-full items-center"
                >
                  <Text className="text-xs font-bold text-white">Simulate Gate Scanner Punch</Text>
                </TouchableOpacity>
              )
            ) : (
              <TouchableOpacity
                onPress={handleGenerateToken}
                className="p-3.5 bg-slate-900 rounded-xl w-full items-center"
              >
                <Text className="text-xs font-bold text-white">Generate {selectedMeal} Token</Text>
              </TouchableOpacity>
            )}
          </View>
        </View>
      )}

      {/* Tab: Menu */}
      {activeTab === 'MENU' && (
        <View className="space-y-4">
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-3">
            <View className="flex-row justify-between items-center border-b border-slate-100 pb-2">
              <Text className="text-sm font-bold text-slate-900">Breakfast (07:30 - 09:30)</Text>
              <Text className="text-xs font-semibold text-emerald-600">450 kcal</Text>
            </View>
            <Text className="text-xs text-slate-700">Masala Idli, Medu Vada, Sambar, Coconut Chutney, Tea/Coffee</Text>
          </View>

          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-3">
            <View className="flex-row justify-between items-center border-b border-slate-100 pb-2">
              <Text className="text-sm font-bold text-slate-900">Lunch (12:30 - 14:30)</Text>
              <Text className="text-xs font-semibold text-emerald-600">780 kcal</Text>
            </View>
            <Text className="text-xs text-slate-700">Paneer Butter Masala, Dal Makhani, Jeera Rice, Tawa Roti, Gulab Jamun</Text>
          </View>

          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-3">
            <View className="flex-row justify-between items-center border-b border-slate-100 pb-2">
              <Text className="text-sm font-bold text-slate-900">Dinner (19:30 - 21:30)</Text>
              <Text className="text-xs font-semibold text-emerald-600">710 kcal</Text>
            </View>
            <Text className="text-xs text-slate-700">Mix Veg Curry, Aloo Gobi, Steamed Basmati, Phulka, Tomato Soup, Kheer</Text>
          </View>
        </View>
      )}

      {/* Tab: Rebates */}
      {activeTab === 'REBATE' && (
        <View className="space-y-4">
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-3">
            <Text className="text-sm font-bold text-slate-900">Apply for Dining Fee Rebate</Text>
            <Text className="text-xs text-slate-500">
              Absence of ≥ 3 days entitles a ₹150/day deduction directly credited to your student fee invoice.
            </Text>

            <View className="space-y-1">
              <Text className="text-[11px] font-semibold text-slate-600">Number of Absent Days</Text>
              <TextInput
                keyboardType="numeric"
                value={rebateDays}
                onChangeText={setRebateDays}
                className="border border-slate-200 rounded-lg p-2 text-xs bg-slate-50"
              />
            </View>

            <View className="space-y-1">
              <Text className="text-[11px] font-semibold text-slate-600">Reason / Gate Pass Reference</Text>
              <TextInput
                value={rebateReason}
                onChangeText={setRebateReason}
                className="border border-slate-200 rounded-lg p-2 text-xs bg-slate-50"
              />
            </View>

            <View className="p-3 bg-emerald-50 rounded-xl border border-emerald-200 flex-row justify-between items-center">
              <Text className="text-xs font-bold text-emerald-900">Estimated Rebate Credit:</Text>
              <Text className="text-base font-extrabold text-emerald-900">
                ₹{(parseInt(rebateDays, 10) || 0) * 150}
              </Text>
            </View>

            <TouchableOpacity
              onPress={handleApplyRebate}
              className="p-3.5 bg-slate-900 rounded-xl items-center mt-2"
            >
              <Text className="text-xs font-bold text-white">Submit Rebate Request</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
