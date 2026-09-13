/**
 * BLOCK_MOBILE_LOGIN_FORM_001
 * Purpose: Mobile viewport authentication screen using React Native Reusables styling.
 * Viewports: Handheld mobile devices with safe-area insets and 44px touch targets.
 */

import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, ActivityIndicator, Alert } from 'react-native';
import { apiClient } from '@campus/api-client';

interface LoginFormProps {
  onSuccess: (tokens: { access_token: string; refresh_token: string }) => void;
}

export function MobileLoginForm({ onSuccess }: LoginFormProps) {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [mfaCode, setMfaCode] = useState('');
  const [isMfaStep, setIsMfaStep] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    if (!identifier || !password) {
      Alert.alert('Validation Error', 'Please enter both your identifier and password.');
      return;
    }

    setLoading(true);
    try {
      const res = await apiClient.post('/auth/login', {
        tenant_id: 'tenant_default',
        identifier,
        password,
        mfa_code: isMfaStep ? mfaCode : undefined,
      });

      if (res.requires_mfa) {
        setIsMfaStep(true);
      } else if (res.tokens) {
        onSuccess(res.tokens);
      }
    } catch (err: any) {
      Alert.alert('Login Failed', err.detail || err.message || 'Authentication error.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <View className="w-full px-6 py-8 bg-white rounded-2xl border border-slate-200">
      <Text className="text-2xl font-bold text-slate-900 text-center mb-1">
        {isMfaStep ? 'Two-Factor Verification' : 'Campus Sign In'}
      </Text>
      <Text className="text-sm text-slate-500 text-center mb-6">
        {isMfaStep ? 'Enter your 6-digit TOTP code' : 'Access your courses, hostel & services'}
      </Text>

      {!isMfaStep ? (
        <>
          <View className="mb-4">
            <Text className="text-xs font-semibold uppercase text-slate-700 mb-1.5">Identifier</Text>
            <TextInput
              className="h-12 px-4 rounded-xl border border-slate-300 text-slate-900 text-base"
              placeholder="Email or Username"
              autoCapitalize="none"
              value={identifier}
              onChangeText={setIdentifier}
            />
          </View>

          <View className="mb-6">
            <Text className="text-xs font-semibold uppercase text-slate-700 mb-1.5">Password</Text>
            <TextInput
              className="h-12 px-4 rounded-xl border border-slate-300 text-slate-900 text-base"
              placeholder="••••••••••••"
              secureTextEntry
              value={password}
              onChangeText={setPassword}
            />
          </View>
        </>
      ) : (
        <View className="mb-6">
          <Text className="text-xs font-semibold uppercase text-slate-700 mb-1.5">TOTP Code</Text>
          <TextInput
            className="h-14 text-center font-mono text-2xl tracking-widest rounded-xl border border-slate-300 text-slate-900"
            placeholder="123456"
            keyboardType="number-pad"
            maxLength={6}
            value={mfaCode}
            onChangeText={setMfaCode}
          />
        </View>
      )}

      <TouchableOpacity
        className="h-12 bg-slate-900 rounded-xl items-center justify-center min-h-[44px]"
        onPress={handleLogin}
        disabled={loading}
      >
        {loading ? (
          <ActivityIndicator color="#ffffff" />
        ) : (
          <Text className="text-white font-semibold text-base">
            {isMfaStep ? 'Verify Code' : 'Sign In'}
          </Text>
        )}
      </TouchableOpacity>
    </View>
  );
}
