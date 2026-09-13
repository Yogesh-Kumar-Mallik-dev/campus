/**
 * BLOCK_MOBILE_AUDIT_SCREEN_001
 * Purpose: Mobile Audit Screen integrating AuditLogList with live API fetching and chain verification.
 */

import React, { useState, useEffect } from 'react';
import { View, Alert } from 'react-native';
import { apiClient } from '@campus/api-client';
import { AuditLogList, type MobileAuditLog } from '../../src/components/audit/AuditLogList';

export default function AuditScreen() {
  const [logs, setLogs] = useState<MobileAuditLog[]>([]);
  const [isChainIntact, setIsChainIntact] = useState(true);

  const fetchLogs = async () => {
    try {
      const res = await apiClient.get<any>('/audit/logs?tenant_id=tenant_default&limit=30');
      setLogs(res.data || []);
    } catch (err: any) {
      console.error('Failed to load audit logs:', err);
    }
  };

  const handleSelectLog = (log: MobileAuditLog) => {
    Alert.alert(
      log.action,
      `ID: ${log.id}\nActor: ${log.actor_id || 'System'}\nStatus: ${log.status}\nHash: ${log.hash}\nTime: ${new Date(log.created_at).toLocaleString()}`
    );
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  return (
    <View className="flex-1">
      <AuditLogList
        logs={logs}
        isChainIntact={isChainIntact}
        onRefresh={fetchLogs}
        onSelectLog={handleSelectLog}
      />
    </View>
  );
}
