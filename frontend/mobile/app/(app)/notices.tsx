/**
 * BLOCK_MOBILE_NOTICES_SCREEN_001
 * Subsystem: Rank 6 - Notice & Announcement System (notices)
 * Purpose:   Mobile campus circular feed, emergency bulletins, category filtering, and mandatory circular acknowledgement.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';

interface NoticeItem {
  id: string;
  title: string;
  category: 'ACADEMIC' | 'EXAMINATION' | 'EMERGENCY' | 'HOSTEL' | 'GENERAL' | 'PLACEMENT';
  priority: 'LOW' | 'STANDARD' | 'HIGH' | 'URGENT';
  targetAudience: 'ALL' | 'STUDENTS' | 'FACULTY' | 'STAFF';
  publishedAt: string;
  authorName: string;
  authorRole: string;
  content: string;
  isPinned: boolean;
  requiresAck: boolean;
  acknowledged: boolean;
}

export default function MobileNoticesScreen() {
  const [selectedCategory, setSelectedCategory] = useState<string>('ALL');
  const [notices, setNotices] = useState<NoticeItem[]>([
    {
      id: 'not_001',
      title: 'Mid-Semester Examination Schedule - Autumn 2026',
      category: 'EXAMINATION',
      priority: 'URGENT',
      targetAudience: 'STUDENTS',
      publishedAt: '2026-09-12 09:30',
      authorName: 'Dr. S. K. Raman',
      authorRole: 'Controller of Examinations',
      content: 'The Mid-Semester Examinations for all B.Tech / M.Tech batches commence from Oct 05, 2026. Hall tickets will be released on the portal 3 days prior.',
      isPinned: true,
      requiresAck: true,
      acknowledged: false,
    },
    {
      id: 'not_002',
      title: 'Hostel Curfew & Quiet Hours Notification',
      category: 'HOSTEL',
      priority: 'HIGH',
      targetAudience: 'STUDENTS',
      publishedAt: '2026-09-10 18:00',
      authorName: 'Prof. Ananya Roy',
      authorRole: 'Chief Warden',
      content: 'All hostel residents are advised that strict quiet hours will be in effect starting 10:30 PM due to upcoming examination preparations.',
      isPinned: true,
      requiresAck: false,
      acknowledged: true,
    },
    {
      id: 'not_003',
      title: 'Campus Placement Drive: Top Tier Tech Companies',
      category: 'PLACEMENT',
      priority: 'STANDARD',
      targetAudience: 'STUDENTS',
      publishedAt: '2026-09-08 14:15',
      authorName: 'Mr. Vikram Verma',
      authorRole: 'Head of Placements',
      content: 'Registrations open for upcoming Day-1 recruitment drives. Eligible 7th-semester students must upload updated resumes by Friday 5 PM.',
      isPinned: false,
      requiresAck: true,
      acknowledged: true,
    },
    {
      id: 'not_004',
      title: 'Annual Tech-Cultural Symposium Call for Coordinators',
      category: 'GENERAL',
      priority: 'LOW',
      targetAudience: 'ALL',
      publishedAt: '2026-09-05 11:00',
      authorName: 'Student Affairs Council',
      authorRole: 'Student Council',
      content: 'Applications are invited from enthusiastic student leaders for executive coordinator positions for the annual flagship symposium.',
      isPinned: false,
      requiresAck: false,
      acknowledged: false,
    },
  ]);

  const categories = ['ALL', 'EXAMINATION', 'ACADEMIC', 'HOSTEL', 'PLACEMENT', 'GENERAL'];

  const filteredNotices = notices.filter(
    (n) => selectedCategory === 'ALL' || n.category === selectedCategory
  );

  const handleAcknowledge = (id: string, title: string) => {
    setNotices((prev) =>
      prev.map((n) => (n.id === id ? { ...n, acknowledged: true } : n))
    );
    Alert.alert('Notice Acknowledged', `Your formal digital acknowledgement for "${title}" has been recorded.`);
  };

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Notice & Circular Feed</Text>
        <Text className="text-xs text-slate-300">
          Rank 6 Notices · Official Campus Bulletins, Urgent Alerts & Circular Tracking
        </Text>
      </View>

      {/* Category Horizontal Filter */}
      <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mb-4">
        <View className="flex-row space-x-2">
          {categories.map((cat) => (
            <TouchableOpacity
              key={cat}
              onPress={() => setSelectedCategory(cat)}
              className={`px-3.5 py-2 rounded-xl mr-2 ${
                selectedCategory === cat ? 'bg-slate-900' : 'bg-white border border-slate-200'
              }`}
            >
              <Text
                className={`text-xs font-semibold ${
                  selectedCategory === cat ? 'text-white' : 'text-slate-700'
                }`}
              >
                {cat}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>

      {/* Notices List */}
      <View className="space-y-4">
        {filteredNotices.map((notice) => (
          <View
            key={notice.id}
            className={`p-5 rounded-2xl bg-white border ${
              notice.isPinned
                ? 'border-amber-300 bg-amber-50/20'
                : notice.priority === 'URGENT'
                ? 'border-red-300 bg-red-50/20'
                : 'border-slate-200'
            } shadow-sm`}
          >
            {/* Header Tags */}
            <View className="flex-row items-center justify-between mb-2">
              <View className="flex-row items-center space-x-2">
                {notice.isPinned && (
                  <View className="bg-amber-100 px-2 py-0.5 rounded-md mr-1.5">
                    <Text className="text-[10px] font-bold text-amber-800">📌 PINNED</Text>
                  </View>
                )}
                <View
                  className={`px-2 py-0.5 rounded-md ${
                    notice.priority === 'URGENT'
                      ? 'bg-red-100'
                      : notice.priority === 'HIGH'
                      ? 'bg-orange-100'
                      : 'bg-slate-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      notice.priority === 'URGENT'
                        ? 'text-red-800'
                        : notice.priority === 'HIGH'
                        ? 'text-orange-800'
                        : 'text-slate-700'
                    }`}
                  >
                    {notice.priority}
                  </Text>
                </View>
                <View className="bg-slate-100 px-2 py-0.5 rounded-md ml-1.5">
                  <Text className="text-[10px] font-semibold text-slate-600">{notice.category}</Text>
                </View>
              </View>
              <Text className="text-[11px] text-slate-400">{notice.publishedAt}</Text>
            </View>

            {/* Title */}
            <Text className="text-base font-bold text-slate-900 mb-1">{notice.title}</Text>

            {/* Author */}
            <Text className="text-xs text-slate-500 mb-3">
              By <Text className="font-semibold text-slate-700">{notice.authorName}</Text> ({notice.authorRole})
            </Text>

            {/* Body */}
            <Text className="text-xs text-slate-600 leading-relaxed mb-4">{notice.content}</Text>

            {/* Action / Ack footer */}
            <View className="flex-row items-center justify-between pt-3 border-t border-slate-100">
              <Text className="text-[11px] text-slate-400">Audience: {notice.targetAudience}</Text>
              {notice.requiresAck ? (
                notice.acknowledged ? (
                  <View className="bg-emerald-50 px-3 py-1.5 rounded-lg border border-emerald-200">
                    <Text className="text-[11px] font-bold text-emerald-800">✓ Acknowledged</Text>
                  </View>
                ) : (
                  <TouchableOpacity
                    onPress={() => handleAcknowledge(notice.id, notice.title)}
                    className="bg-slate-900 px-3 py-1.5 rounded-lg"
                  >
                    <Text className="text-[11px] font-bold text-white">Acknowledge</Text>
                  </TouchableOpacity>
                )
              ) : null}
            </View>
          </View>
        ))}
      </View>
    </ScrollView>
  );
}
