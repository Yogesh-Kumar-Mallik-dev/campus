/**
 * BLOCK_MOBILE_COMPLIANCE_CARD_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Accreditation & compliance metrics cards for mobile screens.
 */

import React from 'react';
import { View, Text } from 'react-native';

export interface MobileComplianceReport {
  framework: string;
  total_events: number;
  successful_events: number;
  failed_events: number;
  security_incidents: number;
  chain_integrity: boolean;
  top_actions: { action: string; count: number }[];
}

interface ComplianceReportCardProps {
  report: MobileComplianceReport;
}

export function ComplianceReportCard({ report }: ComplianceReportCardProps) {
  return (
    <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-4">
      <View className="flex-row items-center justify-between">
        <Text className="text-base font-bold text-slate-900">{report.framework} Accreditation</Text>
        <View
          className={`px-2.5 py-0.5 rounded-full ${
            report.chain_integrity ? 'bg-emerald-100' : 'bg-red-100'
          }`}
        >
          <Text
            className={`text-xs font-bold ${
              report.chain_integrity ? 'text-emerald-800' : 'text-red-800'
            }`}
          >
            {report.chain_integrity ? '✓ INTACT' : '⚠ TAMPERED'}
          </Text>
        </View>
      </View>

      <View className="grid grid-cols-2 gap-3">
        <View className="p-3 bg-slate-50 rounded-xl">
          <Text className="text-[11px] font-semibold text-slate-500 uppercase">Total Events</Text>
          <Text className="text-xl font-bold text-slate-900 mt-1">{report.total_events.toLocaleString()}</Text>
        </View>
        <View className="p-3 bg-slate-50 rounded-xl">
          <Text className="text-[11px] font-semibold text-slate-500 uppercase">Security Alerts</Text>
          <Text className={`text-xl font-bold mt-1 ${report.security_incidents > 0 ? 'text-amber-600' : 'text-slate-900'}`}>
            {report.security_incidents.toLocaleString()}
          </Text>
        </View>
      </View>

      <View className="space-y-2 pt-2 border-t border-slate-100">
        <Text className="text-xs font-bold uppercase text-slate-600">Top Operations</Text>
        {report.top_actions.slice(0, 4).map((act, i) => (
          <View key={i} className="flex-row items-center justify-between py-1">
            <Text className="text-xs font-mono text-slate-800" numberOfLines={1}>{act.action}</Text>
            <Text className="text-xs font-bold text-slate-600">{act.count}</Text>
          </View>
        ))}
      </View>
    </View>
  );
}
