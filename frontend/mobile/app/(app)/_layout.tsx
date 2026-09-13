/**
 * BLOCK_MOBILE_APP_LAYOUT_001
 * Purpose: Navigation stack for authenticated workspace views.
 */

import React from 'react';
import { Stack } from 'expo-router';

export default function AppLayout() {
  return (
    <Stack screenOptions={{ headerStyle: { backgroundColor: '#ffffff' }, headerTintColor: '#0f172a' }}>
      <Stack.Screen name="index" options={{ title: 'The Hub' }} />
      <Stack.Screen name="audit" options={{ title: 'Audit Ledger' }} />
      <Stack.Screen name="compliance" options={{ title: 'Accreditation Reports' }} />
    </Stack>
  );
}
