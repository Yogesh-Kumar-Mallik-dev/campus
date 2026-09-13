/**
 * BLOCK_MOBILE_COMPLIANCE_SCREEN_001
 * Purpose: Mobile screen displaying accreditation summaries with framework selector.
 */

import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, ActivityIndicator } from 'react-native';
import { apiClient } from '@campus/api-client';
import { ComplianceReportCard, type MobileComplianceReport } from '../../src/components/audit/ComplianceReportCard';

const frameworks = ['NAAC', 'NIRF', 'ISO_27001', 'ABET'];

export default function ComplianceScreen() {
  const [selectedFramework, setSelectedFramework] = useState('NAAC');
  const [report, setReport] = useState<MobileComplianceReport | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchReport = async (framework: string) => {
    setLoading(true);
    try {
      const res = await apiClient.get<any>(`/audit/compliance-reports?tenant_id=tenant_default&framework=${framework}`);
      setReport(res.data || null);
    } catch (err) {
      console.error('Failed to load compliance report:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReport(selectedFramework);
  }, [selectedFramework]);

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Framework Selector */}
      <View className="flex-row gap-2 mb-4">
        {frameworks.map((fw) => (
          <TouchableOpacity
            key={fw}
            onPress={() => setSelectedFramework(fw)}
            className={`px-3 py-2 rounded-xl border ${
              selectedFramework === fw
                ? 'bg-slate-900 border-slate-900'
                : 'bg-white border-slate-200'
            }`}
          >
            <Text
              className={`text-xs font-bold ${
                selectedFramework === fw ? 'text-white' : 'text-slate-700'
              }`}
            >
              {fw}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {loading ? (
        <View className="p-8 items-center justify-center">
          <ActivityIndicator color="#0f172a" />
          <Text className="text-xs text-slate-500 mt-2">Compiling accreditation metrics...</Text>
        </View>
      ) : report ? (
        <ComplianceReportCard report={report} />
      ) : (
        <View className="p-8 items-center justify-center">
          <Text className="text-sm text-slate-400">No data available for this framework.</Text>
        </View>
      )}
    </ScrollView>
  );
}
