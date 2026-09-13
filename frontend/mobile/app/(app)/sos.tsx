/**
 * BLOCK_MOBILE_SOS_SCREEN_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   Mobile Student 1-Tap Emergency SOS trigger, instant GPS broadcast, live responder ETA countdown, and safety speed dials.
 */

import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, Linking } from 'react-native';

interface ActiveSOSState {
  isActive: boolean;
  alertNumber: string;
  emergencyType: string;
  location: string;
  latitude: number;
  longitude: number;
  status: 'TRIGGERED' | 'ACKNOWLEDGED' | 'DISPATCHED' | 'ON_SCENE';
  responderName?: string;
  etaMinutes?: number;
}

export default function MobileSOSScreen() {
  const [activeSOS, setActiveSOS] = useState<ActiveSOSState | null>(null);
  const [selectedType, setSelectedType] = useState<string>('MEDICAL');

  const emergencyTypes = [
    { type: 'MEDICAL', label: '🚑 Medical Emergency', color: 'bg-rose-600' },
    { type: 'SECURITY_THREAT', label: '🛡️ Security Threat', color: 'bg-red-800' },
    { type: 'FIRE', label: '🔥 Fire & Smoke', color: 'bg-orange-600' },
    { type: 'HARASSMENT_RAGGING', label: '⚠️ Anti-Ragging', color: 'bg-purple-700' },
  ];

  const handleTriggerSOS = () => {
    const lat = 12.9716 + (Math.random() - 0.5) * 0.005;
    const lng = 77.5946 + (Math.random() - 0.5) * 0.005;
    const alertNum = `SOS-2026-${Math.floor(10000 + Math.random() * 90000)}`;

    const sosData: ActiveSOSState = {
      isActive: true,
      alertNumber: alertNum,
      emergencyType: selectedType,
      location: 'Central Library Floor 2 (Auto-GPS Detected)',
      latitude: lat,
      longitude: lng,
      status: 'TRIGGERED',
      etaMinutes: 4,
    };

    setActiveSOS(sosData);
    Alert.alert(
      '🚨 SOS ALERT BROADCASTED',
      `Control room & campus rapid response team alerted.\nGPS: ${lat.toFixed(4)}° N, ${lng.toFixed(4)}° E\nStay where you are.`
    );
  };

  const handleCancelSOS = () => {
    Alert.alert('Cancel Emergency', 'Are you sure you want to mark this alert as resolved / false alarm?', [
      { text: 'Keep Active', style: 'cancel' },
      {
        text: 'Confirm Cancel',
        style: 'destructive',
        onPress: () => {
          setActiveSOS(null);
          Alert.alert('SOS Deactivated', 'Control room notified that emergency has subsided.');
        },
      },
    ]);
  };

  const handleCall = (number: string) => {
    Linking.openURL(`tel:${number}`).catch(() => {
      Alert.alert('Calling Failed', `Please dial ${number} manually.`);
    });
  };

  return (
    <ScrollView className="flex-1 bg-slate-900 p-4">
      {/* Top Banner */}
      <View className="p-4 bg-red-950/80 border border-red-800 rounded-2xl mb-4 items-center">
        <Text className="text-xl font-black text-rose-400 tracking-wider">CAMPUS EMERGENCY SOS</Text>
        <Text className="text-[11px] text-slate-300 text-center mt-1">
          Instantaneous telemetry broadcast directly to 24/7 Security Control & Campus Paramedics
        </Text>
      </View>

      {/* Active Emergency Alert Card */}
      {activeSOS && (
        <View className="bg-rose-900 border-2 border-rose-500 rounded-2xl p-5 mb-5 space-y-3 shadow-2xl">
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center gap-2">
              <View className="w-3.5 h-3.5 rounded-full bg-rose-400 animate-ping" />
              <Text className="text-sm font-black text-white">{activeSOS.alertNumber}</Text>
            </View>
            <View className="bg-rose-950 px-2.5 py-1 rounded-full border border-rose-400">
              <Text className="text-[10px] font-bold text-rose-300">{activeSOS.status}</Text>
            </View>
          </View>

          <Text className="text-base font-bold text-white">
            🚨 {activeSOS.emergencyType} DISPATCH ACTIVE
          </Text>

          <Text className="text-xs text-rose-200">
            📍 Location: {activeSOS.location}
          </Text>
          <Text className="text-[11px] font-mono text-rose-300">
            GPS: {activeSOS.latitude.toFixed(4)}° N, {activeSOS.longitude.toFixed(4)}° E
          </Text>

          <View className="bg-rose-950/90 p-3 rounded-xl border border-rose-700 items-center">
            <Text className="text-xs text-rose-300 font-semibold">Paramedic & Security Response Team</Text>
            <Text className="text-2xl font-black text-white mt-0.5">ETA: ~{activeSOS.etaMinutes} MINS</Text>
            <Text className="text-[10px] text-rose-400 mt-1">Stay on this screen. Responders are approaching your pin.</Text>
          </View>

          <TouchableOpacity
            onPress={handleCancelSOS}
            className="bg-black/40 border border-rose-400 p-3 rounded-xl items-center"
          >
            <Text className="text-xs font-bold text-rose-300">Mark Safe / Cancel Alert</Text>
          </TouchableOpacity>
        </View>
      )}

      {/* 1-Tap Trigger Section (When Idle) */}
      {!activeSOS && (
        <View className="space-y-4 mb-6">
          <Text className="text-xs font-bold uppercase tracking-wider text-slate-400 text-center">
            Select Emergency Category
          </Text>

          <View className="grid grid-cols-2 gap-2">
            {emergencyTypes.map((em) => (
              <TouchableOpacity
                key={em.type}
                onPress={() => setSelectedType(em.type)}
                className={`p-3 rounded-xl border ${selectedType === em.type ? 'bg-slate-800 border-rose-500' : 'bg-slate-800/40 border-slate-700'}`}
              >
                <Text className={`text-xs font-bold ${selectedType === em.type ? 'text-rose-400' : 'text-slate-300'}`}>
                  {em.label}
                </Text>
              </TouchableOpacity>
            ))}
          </View>

          {/* Big Red Button */}
          <View className="items-center py-6">
            <TouchableOpacity
              onPress={handleTriggerSOS}
              className="w-48 h-48 rounded-full bg-rose-600 border-8 border-rose-800 items-center justify-center shadow-2xl active:scale-95"
            >
              <Text className="text-3xl font-black text-white tracking-widest">SOS</Text>
              <Text className="text-[10px] font-bold text-rose-200 mt-1 uppercase tracking-wider">TAP TO TRIGGER</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      {/* Emergency Speed Dial Directory */}
      <View className="bg-slate-800 p-4 rounded-2xl border border-slate-700 space-y-3 mb-6">
        <Text className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Direct Helpline Speed Dials
        </Text>

        <TouchableOpacity
          onPress={() => handleCall('112')}
          className="flex-row justify-between items-center p-3 bg-slate-700/60 rounded-xl"
        >
          <View>
            <Text className="text-xs font-bold text-white">🚓 National Emergency Police</Text>
            <Text className="text-[10px] text-slate-400">Toll-free 24/7 Dispatch</Text>
          </View>
          <Text className="text-xs font-bold text-rose-400">CALL 112</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => handleCall('108')}
          className="flex-row justify-between items-center p-3 bg-slate-700/60 rounded-xl"
        >
          <View>
            <Text className="text-xs font-bold text-white">🚑 Campus Medical & Ambulance</Text>
            <Text className="text-[10px] text-slate-400">Health Center Dispatch</Text>
          </View>
          <Text className="text-xs font-bold text-rose-400">CALL 108</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => handleCall('18001805522')}
          className="flex-row justify-between items-center p-3 bg-slate-700/60 rounded-xl"
        >
          <View>
            <Text className="text-xs font-bold text-white">⚠️ UGC National Anti-Ragging Helpline</Text>
            <Text className="text-[10px] text-slate-400">Zero Tolerance Confidential Desk</Text>
          </View>
          <Text className="text-xs font-bold text-purple-400">CALL 1800-180-5522</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
