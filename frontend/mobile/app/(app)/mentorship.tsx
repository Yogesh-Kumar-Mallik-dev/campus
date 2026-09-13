/**
 * BLOCK_MOBILE_MENTORSHIP_SCREEN_001
 * Subsystem: Rank 11 - Progress Tracker & Mentorship System (mentorship)
 * Purpose:   Mobile student portal for mentor contact, counseling session logs, and GPA/attendance progress tracker.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface MentorInfo {
  name: string;
  designation: string;
  department: string;
  email: string;
  officeLocation: string;
}

interface SessionLog {
  id: string;
  date: string;
  meetingType: string;
  status: 'SCHEDULED' | 'COMPLETED';
  summary?: string;
  actionItems?: string;
}

interface ProgressSummary {
  semester: number;
  sgpa: number;
  cgpa: number;
  attendance: number;
  atRiskStatus: 'NORMAL' | 'WATCHLIST' | 'CRITICAL_INTERVENTION';
}

export default function MobileMentorshipScreen() {
  const [activeTab, setActiveTab] = useState<'MENTOR' | 'SESSIONS' | 'PROGRESS'>('MENTOR');

  const mentor: MentorInfo = {
    name: 'Dr. Suresh Rao',
    designation: 'Professor & Faculty Mentor',
    department: 'Computer Science & Engineering',
    email: 'suresh.rao@campus.edu',
    officeLocation: 'Room 204, Academic Block A',
  };

  const [sessions, setSessions] = useState<SessionLog[]>([
    {
      id: 'sess-1',
      date: '2026-09-12 14:30',
      meetingType: 'Academic Review',
      status: 'COMPLETED',
      summary: 'Reviewed semester 4 mid-term. On track for honors thesis.',
      actionItems: 'Submit internship applications before end of month.',
    },
    {
      id: 'sess-2',
      date: '2026-09-20 11:00',
      meetingType: '1-on-1 Counseling',
      status: 'SCHEDULED',
    },
  ]);

  const progress: ProgressSummary = {
    semester: 4,
    sgpa: 8.6,
    cgpa: 8.4,
    attendance: 92.5,
    atRiskStatus: 'NORMAL',
  };

  // Request Meeting State
  const [isRequesting, setIsRequesting] = useState(false);
  const [requestReason, setRequestReason] = useState('');

  const handleBookMeeting = () => {
    if (!requestReason) {
      Alert.alert('Missing Reason', 'Please describe the purpose of the counseling session.');
      return;
    }
    const newSess: SessionLog = {
      id: `sess-${Date.now()}`,
      date: 'Requested (Pending Confirmation)',
      meetingType: 'Ad-hoc Counseling',
      status: 'SCHEDULED',
      summary: requestReason,
    };
    setSessions([newSess, ...sessions]);
    setIsRequesting(false);
    setRequestReason('');
    Alert.alert('Meeting Requested', 'Your request has been sent to Dr. Suresh Rao.');
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Mentorship & Progress</Text>
        <Text className="text-xs text-slate-300">
          Rank 11 Mentorship · Faculty Guidance, Counseling Logs & Academic Tracking
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('MENTOR')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'MENTOR' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'MENTOR' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Mentor
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('SESSIONS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'SESSIONS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'SESSIONS' ? 'text-slate-900' : 'text-slate-600'}`}>
            Sessions ({sessions.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('PROGRESS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'PROGRESS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'PROGRESS' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Progress
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: My Mentor */}
      {activeTab === 'MENTOR' && (
        <View className="space-y-4">
          <View className="p-6 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-4">
            <View className="flex-row items-center space-x-3">
              <View className="w-14 h-14 bg-indigo-100 rounded-2xl items-center justify-center">
                <Text className="text-2xl">👨‍🏫</Text>
              </View>
              <View className="flex-1">
                <Text className="text-base font-bold text-slate-900">{mentor.name}</Text>
                <Text className="text-xs text-slate-500">{mentor.designation}</Text>
                <Text className="text-[11px] text-slate-400">{mentor.department}</Text>
              </View>
            </View>

            <View className="bg-slate-50 p-3.5 rounded-xl space-y-1.5 border border-slate-100">
              <View className="flex-row justify-between">
                <Text className="text-xs text-slate-400">Office</Text>
                <Text className="text-xs font-semibold text-slate-700">{mentor.officeLocation}</Text>
              </View>
              <View className="flex-row justify-between">
                <Text className="text-xs text-slate-400">Email</Text>
                <Text className="text-xs font-semibold text-indigo-600">{mentor.email}</Text>
              </View>
            </View>

            {!isRequesting ? (
              <TouchableOpacity
                onPress={() => setIsRequesting(true)}
                className="p-3.5 bg-slate-900 rounded-xl items-center"
              >
                <Text className="text-xs font-bold text-white">📅 Request 1-on-1 Counseling Meeting</Text>
              </TouchableOpacity>
            ) : (
              <View className="space-y-2 pt-2 border-t border-slate-100">
                <Text className="text-xs font-bold text-slate-800">Meeting Topic / Questions</Text>
                <TextInput
                  placeholder="e.g. Discuss elective subject selection and exam preparation..."
                  value={requestReason}
                  onChangeText={setRequestReason}
                  multiline
                  numberOfLines={2}
                  className="border border-slate-200 rounded-lg p-2 text-xs bg-slate-50"
                />
                <View className="flex-row space-x-2 pt-1">
                  <TouchableOpacity
                    onPress={() => setIsRequesting(false)}
                    className="flex-1 py-2 bg-slate-200 rounded-lg items-center"
                  >
                    <Text className="text-xs font-bold text-slate-700">Cancel</Text>
                  </TouchableOpacity>
                  <TouchableOpacity
                    onPress={handleBookMeeting}
                    className="flex-1 py-2 bg-slate-900 rounded-lg items-center"
                  >
                    <Text className="text-xs font-bold text-white">Send Request</Text>
                  </TouchableOpacity>
                </View>
              </View>
            )}
          </View>
        </View>
      )}

      {/* Tab: Sessions */}
      {activeTab === 'SESSIONS' && (
        <View className="space-y-3">
          {sessions.map((sess) => (
            <View key={sess.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
              <View className="flex-row justify-between items-start">
                <Text className="text-sm font-bold text-slate-900">{sess.meetingType}</Text>
                <View
                  className={`px-2 py-0.5 rounded ${
                    sess.status === 'COMPLETED' ? 'bg-emerald-100' : 'bg-blue-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      sess.status === 'COMPLETED' ? 'text-emerald-800' : 'text-blue-800'
                    }`}
                  >
                    {sess.status}
                  </Text>
                </View>
              </View>

              <Text className="text-xs text-slate-400">⏰ {sess.date}</Text>

              {sess.summary && (
                <View className="bg-slate-50 p-2.5 rounded-xl border border-slate-100 space-y-1 mt-1">
                  <Text className="text-[11px] font-bold text-slate-700">Notes:</Text>
                  <Text className="text-xs text-slate-600">{sess.summary}</Text>
                  {sess.actionItems && (
                    <Text className="text-[11px] text-slate-500 pt-1">
                      <Text className="font-bold text-slate-700">Action Items: </Text>
                      {sess.actionItems}
                    </Text>
                  )}
                </View>
              )}
            </View>
          ))}
        </View>
      )}

      {/* Tab: Progress */}
      {activeTab === 'PROGRESS' && (
        <View className="space-y-4">
          <View className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-4">
            <View className="flex-row justify-between items-center border-b border-slate-100 pb-3">
              <Text className="text-sm font-bold text-slate-900">Semester {progress.semester} Standing</Text>
              <View className="bg-emerald-100 px-2.5 py-0.5 rounded-full">
                <Text className="text-[10px] font-bold text-emerald-800">ACADEMIC STANDING: GOOD</Text>
              </View>
            </View>

            <View className="grid grid-cols-2 gap-3 flex-row">
              <View className="flex-1 bg-slate-50 p-3 rounded-xl border border-slate-100">
                <Text className="text-[10px] text-slate-400">Semester SGPA</Text>
                <Text className="text-xl font-extrabold text-slate-900 mt-0.5">{progress.sgpa.toFixed(2)}</Text>
              </View>
              <View className="flex-1 bg-slate-50 p-3 rounded-xl border border-slate-100">
                <Text className="text-[10px] text-slate-400">Cumulative CGPA</Text>
                <Text className="text-xl font-extrabold text-slate-900 mt-0.5">{progress.cgpa.toFixed(2)}</Text>
              </View>
            </View>

            <View className="bg-emerald-50 p-3.5 rounded-xl border border-emerald-100 flex-row justify-between items-center">
              <View>
                <Text className="text-xs font-bold text-emerald-900">Biometric Attendance</Text>
                <Text className="text-[10px] text-emerald-700">Above mandatory 75% cutoff</Text>
              </View>
              <Text className="text-xl font-black text-emerald-900">{progress.attendance.toFixed(1)}%</Text>
            </View>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
