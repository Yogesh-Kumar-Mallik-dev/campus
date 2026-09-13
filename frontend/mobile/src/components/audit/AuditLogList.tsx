/**
 * BLOCK_MOBILE_AUDIT_LOG_LIST_001
 * Subsystem: Rank 2 - Central Audit & Compliance System (audit)
 * Purpose:   Native mobile audit trail viewer with tamper-evident chain verification status.
 * Viewports: iOS and Android handheld viewports with smooth scrolling and responsive touch targets.
 */

import React from 'react';
import { View, Text, FlatList, TouchableOpacity } from 'react-native';

export interface MobileAuditLog {
  id: string;
  tenant_id: string;
  actor_id?: string;
  actor_role?: string;
  action: string;
  resource_type: string;
  status: 'SUCCESS' | 'FAILURE' | 'ATTEMPTED';
  hash: string;
  created_at: string;
}

interface AuditLogListProps {
  logs: MobileAuditLog[];
  isChainIntact?: boolean;
  onRefresh?: () => void;
  onSelectLog?: (log: MobileAuditLog) => void;
}

export function AuditLogList({
  logs,
  isChainIntact = true,
  onRefresh,
  onSelectLog,
}: AuditLogListProps) {
  const renderItem = ({ item }: { item: MobileAuditLog }) => {
    const isSuccess = item.status === 'SUCCESS';

    return (
      <TouchableOpacity
        onPress={() => onSelectLog?.(item)}
        activeOpacity={0.7}
        className="p-4 bg-white rounded-xl mb-2.5 border border-slate-200"
      >
        <View className="flex-row items-center justify-between mb-1.5">
          <Text className="text-xs font-mono font-bold text-slate-900 flex-1 mr-2" numberOfLines={1}>
            {item.action}
          </Text>
          <View
            className={`px-2 py-0.5 rounded-full ${
              isSuccess ? 'bg-emerald-100' : 'bg-red-100'
            }`}
          >
            <Text
              className={`text-[10px] font-bold ${
                isSuccess ? 'text-emerald-800' : 'text-red-800'
              }`}
            >
              {item.status}
            </Text>
          </View>
        </View>

        <View className="flex-row items-center justify-between text-xs text-slate-500 mb-2">
          <Text className="text-xs text-slate-600">
            Actor: <Text className="font-semibold">{item.actor_id || 'System'}</Text>
            {item.actor_role ? ` (${item.actor_role})` : ''}
          </Text>
          <Text className="text-[11px] text-slate-400">
            {new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
          </Text>
        </View>

        <View className="pt-2 border-t border-slate-100 flex-row items-center justify-between">
          <Text className="text-[10px] font-mono text-slate-400">
            SHA256: {item.hash.substring(0, 8)}...
          </Text>
          <Text className="text-xs font-medium text-blue-600">Inspect →</Text>
        </View>
      </TouchableOpacity>
    );
  };

  return (
    <View className="flex-1 bg-slate-50 px-4 pt-3">
      {/* Chain Status Banner */}
      <View
        className={`p-3 rounded-xl mb-3 flex-row items-center justify-between ${
          isChainIntact ? 'bg-emerald-50 border border-emerald-200' : 'bg-red-50 border border-red-200'
        }`}
      >
        <View className="flex-row items-center">
          <Text className="text-sm mr-2">{isChainIntact ? '🛡️' : '⚠️'}</Text>
          <Text
            className={`text-xs font-bold ${
              isChainIntact ? 'text-emerald-900' : 'text-red-900'
            }`}
          >
            {isChainIntact ? 'Ledger Hash Chain Intact' : 'Chain Tamper Alert Detected'}
          </Text>
        </View>
        <Text className="text-[10px] font-mono text-slate-500">SHA-256</Text>
      </View>

      <FlatList
        data={logs}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        onRefresh={onRefresh}
        refreshing={false}
        contentContainerStyle={{ paddingBottom: 20 }}
        ListEmptyComponent={
          <View className="p-8 items-center justify-center">
            <Text className="text-sm text-slate-400">No audit events found.</Text>
          </View>
        }
      />
    </View>
  );
}
