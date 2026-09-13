/**
 * BLOCK_MOBILE_BILLING_SCREEN_001
 * Subsystem: Rank 5 - Central Payment & Billing System (billing)
 * Purpose:   Mobile student fee payment portal, itemized invoice breakdown, 1-tap UPI checkout, and receipts.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';

export default function MobileBillingScreen() {
  const [activeTab, setActiveTab] = useState<'INVOICES' | 'HISTORY'>('INVOICES');
  const [isPaying, setIsPaying] = useState(false);
  const [receiptNumber, setReceiptNumber] = useState<string | null>(null);

  const invoices = [
    {
      id: 'inv_101',
      number: 'INV-2026-2027-000101',
      name: 'Semester 1 - Tuition & Lab Fee',
      dueDate: '2026-10-15',
      total: 50000,
      paid: 25000,
      balance: 25000,
      status: 'PARTIALLY_PAID',
    },
    {
      id: 'inv_102',
      number: 'INV-2026-2027-000102',
      name: 'Hostel & Mess Boarding Fee',
      dueDate: '2026-10-31',
      total: 45000,
      paid: 0,
      balance: 45000,
      status: 'ISSUED',
    },
  ];

  const totalPending = invoices.reduce((acc, inv) => acc + inv.balance, 0);

  const handlePay = (invNumber: string, amount: number) => {
    setIsPaying(true);
    setTimeout(() => {
      setIsPaying(false);
      const rec = `REC-${new Date().getFullYear()}-${Math.floor(100000 + Math.random() * 900000)}`;
      setReceiptNumber(rec);
      Alert.alert('Payment Successful', `Settled ₹${amount.toLocaleString('en-IN')} for ${invNumber}.\nReceipt: ${rec}`);
    }, 800);
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Fee Payments & Invoices</Text>
        <Text className="text-xs text-slate-300">
          Rank 5 Billing · Double-Entry Ledger & Fast-Track Digital Settlement
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('INVOICES')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'INVOICES' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'INVOICES' ? 'text-slate-900' : 'text-slate-600'}`}>
            Pending Fees ({invoices.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('HISTORY')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'HISTORY' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'HISTORY' ? 'text-slate-900' : 'text-slate-600'}`}>
            Receipts & History
          </Text>
        </TouchableOpacity>
      </View>

      {activeTab === 'INVOICES' ? (
        <View className="space-y-4">
          {/* Outstanding Balance Banner */}
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm">
            <Text className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
              Total Outstanding Fee Balance
            </Text>
            <Text className="text-3xl font-extrabold text-slate-900 mb-1">
              ₹{totalPending.toLocaleString('en-IN')}
            </Text>
            <Text className="text-[11px] text-slate-400">
              Clear fees before due date to avoid automated late fine accruals.
            </Text>
          </View>

          {receiptNumber && (
            <View className="p-4 bg-emerald-50 border border-emerald-200 rounded-2xl">
              <Text className="text-xs font-bold text-emerald-900">✓ Latest Payment Settled</Text>
              <Text className="text-xs text-emerald-700">Receipt No: {receiptNumber}</Text>
            </View>
          )}

          {/* Invoice Cards */}
          <View className="space-y-3">
            <Text className="text-sm font-bold text-slate-900 px-1">Active Semester Invoices</Text>
            {invoices.map((inv) => (
              <View key={inv.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-3">
                <View className="flex-row justify-between items-start">
                  <View className="flex-1 pr-2">
                    <Text className="text-xs font-bold text-slate-900">{inv.name}</Text>
                    <Text className="text-[11px] font-mono text-slate-400 mt-0.5">{inv.number}</Text>
                  </View>
                  <View className={`px-2 py-0.5 rounded-full ${inv.status === 'PAID' ? 'bg-emerald-100' : 'bg-amber-100'}`}>
                    <Text className={`text-[10px] font-bold ${inv.status === 'PAID' ? 'text-emerald-800' : 'text-amber-800'}`}>
                      {inv.status}
                    </Text>
                  </View>
                </View>

                <View className="flex-row justify-between py-2 border-t border-b border-slate-100 text-xs">
                  <View>
                    <Text className="text-slate-400 text-[11px]">Due Date</Text>
                    <Text className="font-medium text-slate-700">{inv.dueDate}</Text>
                  </View>
                  <View>
                    <Text className="text-slate-400 text-[11px]">Total Fee</Text>
                    <Text className="font-medium text-slate-700">₹{inv.total.toLocaleString('en-IN')}</Text>
                  </View>
                  <View className="items-end">
                    <Text className="text-slate-400 text-[11px]">Pending Balance</Text>
                    <Text className="font-bold text-slate-900">₹{inv.balance.toLocaleString('en-IN')}</Text>
                  </View>
                </View>

                <TouchableOpacity
                  onPress={() => handlePay(inv.number, inv.balance)}
                  disabled={isPaying}
                  className="w-full py-3 bg-slate-900 rounded-xl items-center active:bg-slate-800"
                >
                  <Text className="text-xs font-bold text-white">
                    {isPaying ? 'Processing UPI Gateway...' : `Pay ₹${inv.balance.toLocaleString('en-IN')} via UPI`}
                  </Text>
                </TouchableOpacity>
              </View>
            ))}
          </View>
        </View>
      ) : (
        /* History & Receipts */
        <View className="p-5 bg-white rounded-2xl border border-slate-200 space-y-4">
          <Text className="text-sm font-bold text-slate-900">Verified Payment Receipts</Text>
          <View className="p-4 bg-slate-50 rounded-xl border border-slate-200">
            <View className="flex-row justify-between items-center mb-1">
              <Text className="font-mono text-xs font-bold text-slate-900">REC-2026-000001</Text>
              <Text className="text-[10px] font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full">SETTLED</Text>
            </View>
            <Text className="text-xs text-slate-600">Semester 1 Advance Installment</Text>
            <View className="flex-row justify-between items-baseline mt-2">
              <Text className="text-xs text-slate-400">10 Sep 2026 · Razorpay</Text>
              <Text className="font-bold text-sm text-slate-900">₹25,000</Text>
            </View>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
