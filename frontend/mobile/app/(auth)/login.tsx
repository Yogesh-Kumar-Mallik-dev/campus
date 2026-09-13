/**
 * BLOCK_MOBILE_LOGIN_SCREEN_001
 * Purpose: Authentication screen routing on successful login.
 */

import React from 'react';
import { View } from 'react-native';
import { useRouter } from 'expo-router';
import { MobileLoginForm } from '../../src/components/LoginForm';

export default function LoginScreen() {
  const router = useRouter();

  const handleSuccess = () => {
    router.replace('/(app)');
  };

  return (
    <View className="flex-1 justify-center px-6 bg-slate-100">
      <MobileLoginForm onSuccess={handleSuccess} />
    </View>
  );
}
