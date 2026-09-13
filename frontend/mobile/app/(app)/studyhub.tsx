/**
 * BLOCK_MOBILE_STUDYHUB_SCREEN_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Mobile student portal for lecture notes, syllabus tracking, assignment submissions, and peer grading.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface CourseInfo {
  id: string;
  code: string;
  title: string;
  department: string;
  credits: number;
}

interface MaterialItem {
  id: string;
  title: string;
  type: string;
  unit: number;
  fileSize: string;
}

interface AssignmentItem {
  id: string;
  title: string;
  dueDate: string;
  maxMarks: number;
  status: 'PENDING' | 'SUBMITTED' | 'GRADED';
  marksObtained?: number;
  feedback?: string;
}

export default function MobileStudyHubScreen() {
  const [activeTab, setActiveTab] = useState<'MATERIALS' | 'ASSIGNMENTS' | 'PEER_REVIEW'>('MATERIALS');
  const [selectedCourse, setSelectedCourse] = useState<string>('CS-301');

  // Submit Modal state
  const [submittingAsgnId, setSubmittingAsgnId] = useState<string | null>(null);
  const [submissionUrl, setSubmissionUrl] = useState<string>('');
  const [submissionNotes, setSubmissionNotes] = useState<string>('');

  const courses: CourseInfo[] = [
    { id: 'crs-1', code: 'CS-301', title: 'Operating Systems', department: 'CSE', credits: 4 },
    { id: 'crs-2', code: 'CS-302', title: 'Database Engines', department: 'CSE', credits: 4 },
    { id: 'crs-3', code: 'CS-303', title: 'Computer Networks', department: 'CSE', credits: 3 },
  ];

  const materials: MaterialItem[] = [
    { id: 'm-1', title: 'Unit 1: Kernel Architecture & System Calls', type: 'SLIDE', unit: 1, fileSize: '2.4 MB' },
    { id: 'm-2', title: 'Lab 1: ProcFS Kernel Module in C', type: 'LAB_MANUAL', unit: 1, fileSize: '1.8 MB' },
    { id: 'm-3', title: 'Unit 2: Virtual Memory & 4-Level Paging', type: 'NOTE', unit: 2, fileSize: '3.1 MB' },
  ];

  const [assignments, setAssignments] = useState<AssignmentItem[]>([
    {
      id: 'a-1',
      title: 'Lab 01: Kernel Module & ProcFS Debugger',
      dueDate: '2026-09-18 23:59',
      maxMarks: 100,
      status: 'GRADED',
      marksObtained: 96,
      feedback: 'Clean kmemleak run and accurate /proc output.',
    },
    {
      id: 'a-2',
      title: 'Assignment 2: Lock-Free Multi-Producer Queue',
      dueDate: '2026-09-25 23:59',
      maxMarks: 50,
      status: 'PENDING',
    },
  ]);

  const handleSubmit = (asgnId: string) => {
    if (!submissionUrl) {
      Alert.alert('Missing Code Link', 'Please provide a Git repository or cloud storage archive link.');
      return;
    }
    setAssignments((prev) =>
      prev.map((a) => (a.id === asgnId ? { ...a, status: 'SUBMITTED' } : a))
    );
    setSubmittingAsgnId(null);
    setSubmissionUrl('');
    setSubmissionNotes('');
    Alert.alert('Submission Received', 'Your code has been uploaded and timestamped for grading.');
  };

  const handleDownload = (mat: MaterialItem) => {
    Alert.alert('Downloading Resource', `Downloading "${mat.title}" (${mat.fileSize}) to offline storage.`);
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Study Hub</Text>
        <Text className="text-xs text-slate-300">
          Rank 10 Study Hub · Syllabus, Lecture Notes, Assignments & Peer Grading
        </Text>
      </View>

      {/* Course Selector Pills */}
      <ScrollView horizontal showsHorizontalScrollIndicator={false} className="flex-row mb-4 space-x-2">
        {courses.map((c) => (
          <TouchableOpacity
            key={c.id}
            onPress={() => setSelectedCourse(c.code)}
            className={`px-3.5 py-2 rounded-xl border ${
              selectedCourse === c.code
                ? 'bg-slate-900 border-slate-900'
                : 'bg-white border-slate-200'
            }`}
          >
            <Text
              className={`text-xs font-bold ${
                selectedCourse === c.code ? 'text-white' : 'text-slate-700'
              }`}
            >
              {c.code} · {c.title}
            </Text>
          </TouchableOpacity>
        ))}
      </ScrollView>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('MATERIALS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'MATERIALS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'MATERIALS' ? 'text-slate-900' : 'text-slate-600'}`}>
            Notes ({materials.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('ASSIGNMENTS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'ASSIGNMENTS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'ASSIGNMENTS' ? 'text-slate-900' : 'text-slate-600'}`}>
            Assignments ({assignments.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('PEER_REVIEW')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'PEER_REVIEW' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'PEER_REVIEW' ? 'text-slate-900' : 'text-slate-600'}`}>
            Peer Review
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: Materials */}
      {activeTab === 'MATERIALS' && (
        <View className="space-y-3">
          {materials.map((mat) => (
            <View key={mat.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
              <View className="flex-row justify-between items-start">
                <Text className="text-sm font-bold text-slate-900 flex-1 mr-2">{mat.title}</Text>
                <View className="bg-indigo-50 px-2 py-0.5 rounded">
                  <Text className="text-[10px] font-bold text-indigo-700">{mat.type}</Text>
                </View>
              </View>

              <View className="flex-row justify-between items-center pt-1 border-t border-slate-100">
                <Text className="text-xs text-slate-400">Unit {mat.unit} · {mat.fileSize}</Text>
                <TouchableOpacity
                  onPress={() => handleDownload(mat)}
                  className="px-3 py-1.5 bg-slate-900 rounded-lg"
                >
                  <Text className="text-[11px] font-bold text-white">📥 Download</Text>
                </TouchableOpacity>
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Tab: Assignments */}
      {activeTab === 'ASSIGNMENTS' && (
        <View className="space-y-3">
          {assignments.map((asgn) => (
            <View key={asgn.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-3">
              <View className="flex-row justify-between items-start">
                <Text className="text-sm font-bold text-slate-900 flex-1 mr-2">{asgn.title}</Text>
                <View
                  className={`px-2 py-0.5 rounded ${
                    asgn.status === 'GRADED'
                      ? 'bg-emerald-100'
                      : asgn.status === 'SUBMITTED'
                      ? 'bg-blue-100'
                      : 'bg-amber-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      asgn.status === 'GRADED'
                        ? 'text-emerald-800'
                        : asgn.status === 'SUBMITTED'
                        ? 'text-blue-800'
                        : 'text-amber-800'
                    }`}
                  >
                    {asgn.status}
                  </Text>
                </View>
              </View>

              <View className="flex-row justify-between bg-slate-50 p-2.5 rounded-xl">
                <View>
                  <Text className="text-[10px] text-slate-400">Deadline</Text>
                  <Text className="text-xs font-semibold text-slate-700">{asgn.dueDate}</Text>
                </View>
                <View className="items-end">
                  <Text className="text-[10px] text-slate-400">Marks</Text>
                  <Text className="text-xs font-bold text-slate-900">
                    {asgn.marksObtained !== undefined ? `${asgn.marksObtained} / ` : ''}
                    {asgn.maxMarks} pts
                  </Text>
                </View>
              </View>

              {asgn.feedback && (
                <View className="p-2.5 bg-emerald-50 rounded-xl border border-emerald-100">
                  <Text className="text-[11px] font-bold text-emerald-900">Faculty Feedback:</Text>
                  <Text className="text-xs text-emerald-800 mt-0.5">{asgn.feedback}</Text>
                </View>
              )}

              {asgn.status === 'PENDING' && (
                <TouchableOpacity
                  onPress={() => setSubmittingAsgnId(asgn.id)}
                  className="py-2.5 bg-slate-900 rounded-xl items-center"
                >
                  <Text className="text-xs font-bold text-white">📤 Turn in Submission</Text>
                </TouchableOpacity>
              )}
            </View>
          ))}
        </View>
      )}

      {/* Tab: Peer Review */}
      {activeTab === 'PEER_REVIEW' && (
        <View className="space-y-3">
          <View className="p-4 bg-indigo-50 border border-indigo-200 rounded-2xl">
            <Text className="text-sm font-bold text-indigo-900">Double-Blind Peer Grading</Text>
            <Text className="text-xs text-indigo-700 mt-1">
              Review code submissions assigned anonymously from other students in your cohort.
            </Text>
          </View>

          <View className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
            <Text className="text-sm font-bold text-slate-900">Assigned Submission #4082</Text>
            <Text className="text-xs text-slate-500">Assignment 1: B-Tree Index Implementation</Text>
            <TouchableOpacity
              onPress={() => Alert.alert('Peer Evaluation', 'Please open web portal for full code diff viewer & rubric scoring.')}
              className="py-2.5 bg-indigo-600 rounded-xl items-center mt-2"
            >
              <Text className="text-xs font-bold text-white">Review & Score Code</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      {/* Submit Assignment Modal */}
      {submittingAsgnId && (
        <View className="mt-4 p-5 bg-white border border-slate-300 rounded-2xl shadow-lg space-y-3">
          <Text className="text-sm font-bold text-slate-900">Submit Code Archive</Text>
          <TextInput
            placeholder="Archive URL (https://storage.campus.internal/...)"
            value={submissionUrl}
            onChangeText={setSubmissionUrl}
            className="border border-slate-200 rounded-lg p-2 text-xs bg-slate-50"
          />
          <TextInput
            placeholder="Notes or build instructions..."
            value={submissionNotes}
            onChangeText={setSubmissionNotes}
            multiline
            numberOfLines={2}
            className="border border-slate-200 rounded-lg p-2 text-xs bg-slate-50"
          />
          <View className="flex-row space-x-2">
            <TouchableOpacity
              onPress={() => setSubmittingAsgnId(null)}
              className="flex-1 py-2 bg-slate-200 rounded-lg items-center"
            >
              <Text className="text-xs font-bold text-slate-700">Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity
              onPress={() => handleSubmit(submittingAsgnId)}
              className="flex-1 py-2 bg-slate-900 rounded-lg items-center"
            >
              <Text className="text-xs font-bold text-white">Confirm</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}
    </ScrollView>
  );
}
