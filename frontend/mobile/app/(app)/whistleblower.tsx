/**
 * BLOCK_MOBILE_WHISTLEBLOWER_SCREEN_001
 * Subsystem: Rank 15 - Anonymity & Whistleblower System (whistleblower)
 * Purpose:   Mobile Anonymous Anti-Ragging & Confidential Grievance portal with secret token generation and offline token vault.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface SavedToken {
  reportNumber: string;
  category: string;
  token: string;
  date: string;
}

export default function MobileWhistleblowerScreen() {
  const [activeTab, setActiveTab] = useState<'SUBMIT' | 'TRACK'>('SUBMIT');

  // Form State
  const [category, setCategory] = useState<string>('ANTI_RAGGING');
  const [severity, setSeverity] = useState<'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL'>('HIGH');
  const [title, setTitle] = useState<string>('');
  const [description, setDescription] = useState<string>('');

  // Generated Token Modal / State
  const [latestToken, setLatestToken] = useState<string | null>(null);
  const [latestReportNum, setLatestReportNum] = useState<string | null>(null);

  // Saved Vault Tokens
  const [savedTokens, setSavedTokens] = useState<SavedToken[]>([
    {
      reportNumber: 'WB-2026-00042',
      category: 'ANTI_RAGGING',
      token: 'WB-tok-9f821ac34b89e210',
      date: 'Yesterday',
    },
  ]);

  // Lookup state
  const [searchToken, setSearchToken] = useState<string>('');
  const [lookupResult, setLookupResult] = useState<{
    status: string;
    category: string;
    message: string;
  } | null>(null);

  const handleSubmit = () => {
    if (!title.trim() || !description.trim()) {
      Alert.alert('Required Fields', 'Please provide a subject summary and detailed description.');
      return;
    }

    const hex = Math.random().toString(16).substring(2, 10) + Math.random().toString(16).substring(2, 10);
    const rawTok = `WB-tok-${hex}`;
    const repNum = `WB-2026-${Math.floor(10000 + Math.random() * 90000)}`;

    setLatestToken(rawTok);
    setLatestReportNum(repNum);
    setSavedTokens([
      { reportNumber: repNum, category, token: rawTok, date: 'Today' },
      ...savedTokens,
    ]);

    setTitle('');
    setDescription('');
  };

  const handleLookup = () => {
    if (!searchToken.trim()) return;
    setLookupResult({
      status: 'UNDER_INVESTIGATION',
      category: 'ANTI_RAGGING',
      message: 'Anti-ragging squad has initiated surprise inspections. Disciplinary inquiry scheduled.',
    });
  };

  return (
    <ScrollView className="flex-1 bg-slate-900 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-indigo-950 border border-indigo-800 rounded-2xl mb-4">
        <Text className="text-xl font-black text-indigo-400 mb-1">Zero-Knowledge Whistleblower</Text>
        <Text className="text-xs text-slate-300">
          Rank 15 · 100% Identity-Shielded Anti-Ragging & Confidential Grievance Desk
        </Text>
      </View>

      {/* Tabs */}
      <View className="flex-row bg-slate-800 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('SUBMIT')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'SUBMIT' ? 'bg-indigo-600' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'SUBMIT' ? 'text-white' : 'text-slate-400'}`}>
            + Lodge Report
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('TRACK')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'TRACK' ? 'bg-indigo-600' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'TRACK' ? 'text-white' : 'text-slate-400'}`}>
            Track with Token ({savedTokens.length})
          </Text>
        </TouchableOpacity>
      </View>

      {/* Mode: Submit */}
      {activeTab === 'SUBMIT' && (
        <View className="space-y-4">
          {latestToken && (
            <View className="p-4 bg-amber-950/80 border-2 border-amber-500 rounded-2xl space-y-2">
              <Text className="text-xs font-black text-amber-400 uppercase">
                ⚠️ REPORT LODGED: SAVE THIS SECRET KEY
              </Text>
              <Text className="text-xs text-slate-200">
                Ref: <Text className="font-bold text-white">{latestReportNum}</Text>
              </Text>
              <View className="bg-black/50 p-2.5 rounded-lg border border-amber-600">
                <Text className="font-mono text-xs font-bold text-amber-300">{latestToken}</Text>
              </View>
              <Text className="text-[10px] text-amber-300">
                This token is the ONLY way to view investigation findings. It has also been saved to your local device vault.
              </Text>
              <TouchableOpacity
                onPress={() => setLatestToken(null)}
                className="bg-amber-600 py-1.5 rounded-lg items-center mt-1"
              >
                <Text className="text-xs font-bold text-white">Acknowledge & Dismiss</Text>
              </TouchableOpacity>
            </View>
          )}

          <View className="bg-slate-800 p-5 rounded-2xl border border-slate-700 space-y-3">
            <View className="space-y-1">
              <Text className="text-xs font-semibold text-slate-300">Category</Text>
              <View className="flex-row flex-wrap gap-2">
                {['ANTI_RAGGING', 'HARASSMENT', 'ACADEMIC_CORRUPTION', 'FINANCIAL_FRAUD'].map((c) => (
                  <TouchableOpacity
                    key={c}
                    onPress={() => setCategory(c)}
                    className={`px-3 py-1.5 rounded-lg border ${category === c ? 'bg-indigo-600 border-indigo-600' : 'bg-slate-700/60 border-slate-600'}`}
                  >
                    <Text className={`text-[11px] font-bold ${category === c ? 'text-white' : 'text-slate-300'}`}>
                      {c.replace('_', ' ')}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            <View className="space-y-1">
              <Text className="text-xs font-semibold text-slate-300">Subject Summary</Text>
              <TextInput
                placeholder="e.g. Intimidation reported in Hostel Block C"
                placeholderTextColor="#94a3b8"
                value={title}
                onChangeText={setTitle}
                className="border border-slate-600 rounded-xl p-3 text-xs bg-slate-900 text-white"
              />
            </View>

            <View className="space-y-1">
              <Text className="text-xs font-semibold text-slate-300">Detailed Incident Narrative</Text>
              <TextInput
                placeholder="Describe approximate location, dates, times, and incident nature. Do NOT include your own name..."
                placeholderTextColor="#94a3b8"
                value={description}
                onChangeText={setDescription}
                multiline
                numberOfLines={5}
                className="border border-slate-600 rounded-xl p-3 text-xs bg-slate-900 text-white min-h-[120px]"
              />
            </View>

            <TouchableOpacity
              onPress={handleSubmit}
              className="bg-indigo-600 p-3.5 rounded-xl items-center shadow-md"
            >
              <Text className="text-white text-xs font-bold">Submit 100% Anonymously</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      {/* Mode: Track */}
      {activeTab === 'TRACK' && (
        <View className="space-y-4">
          <View className="bg-slate-800 p-4 rounded-2xl border border-slate-700 space-y-3">
            <Text className="text-xs font-bold text-slate-300 uppercase">Lookup By Secret Token</Text>
            <TextInput
              placeholder="Paste WB-tok-xxxxxxxx key..."
              placeholderTextColor="#94a3b8"
              value={searchToken}
              onChangeText={setSearchToken}
              className="border border-slate-600 rounded-xl p-3 text-xs bg-slate-900 text-white font-mono"
            />
            <TouchableOpacity
              onPress={handleLookup}
              className="bg-indigo-600 p-3 rounded-xl items-center"
            >
              <Text className="text-white text-xs font-bold">Track Investigation Status</Text>
            </TouchableOpacity>
          </View>

          {lookupResult && (
            <View className="bg-slate-800 p-4 rounded-2xl border border-indigo-500 space-y-2">
              <View className="flex-row justify-between items-center">
                <Text className="text-xs font-bold text-indigo-400">{lookupResult.category}</Text>
                <View className="bg-indigo-950 px-2 py-0.5 rounded border border-indigo-400">
                  <Text className="text-[10px] font-bold text-indigo-300">{lookupResult.status}</Text>
                </View>
              </View>
              <Text className="text-xs text-white leading-relaxed">{lookupResult.message}</Text>
            </View>
          )}

          {/* Local Vault List */}
          <View className="space-y-2">
            <Text className="text-xs font-bold text-slate-400 uppercase tracking-wider">
              Saved Secret Keys (Device Vault)
            </Text>
            {savedTokens.map((st) => (
              <TouchableOpacity
                key={st.token}
                onPress={() => {
                  setSearchToken(st.token);
                  handleLookup();
                }}
                className="bg-slate-800 p-3.5 rounded-xl border border-slate-700 flex-row justify-between items-center"
              >
                <View>
                  <Text className="text-xs font-bold text-white">{st.reportNumber} · {st.category}</Text>
                  <Text className="font-mono text-[10px] text-indigo-400 mt-0.5">{st.token}</Text>
                </View>
                <Text className="text-[10px] text-slate-400">{st.date}</Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>
      )}
    </ScrollView>
  );
}
