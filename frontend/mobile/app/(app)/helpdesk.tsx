/**
 * BLOCK_MOBILE_HELPDESK_SCREEN_001
 * Subsystem: Rank 13 - Application & Query Helpdesk (helpdesk)
 * Purpose:   Mobile student support desk: raise query tickets, track SLA countdowns, live response messaging, and 1-tap urgent escalation.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput, Modal } from 'react-native';

interface MobileTicket {
  id: string;
  ticketNumber: string;
  category: string;
  title: string;
  description: string;
  priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
  status: 'OPEN' | 'IN_PROGRESS' | 'WAITING_FOR_APPLICANT' | 'RESOLVED' | 'CLOSED';
  slaHoursLeft: number;
  lastReply?: string;
  rating?: number;
}

export default function MobileHelpdeskScreen() {
  const [activeTab, setActiveTab] = useState<'MY_TICKETS' | 'RAISE_QUERY'>('MY_TICKETS');
  const [searchQuery, setSearchQuery] = useState('');

  const [tickets, setTickets] = useState<MobileTicket[]>([
    {
      id: 't-1',
      ticketNumber: 'HD-2026-00042',
      category: 'Academic Affairs',
      title: 'End-sem Exam slot overlap with Lab Exam',
      description: 'My CS302 exam conflicts with CS392 Compiler Lab at 10:00 AM on 18 Oct.',
      priority: 'URGENT',
      status: 'IN_PROGRESS',
      slaHoursLeft: 18,
      lastReply: 'Exam Office: Special buffer slot arranged on 19 Oct 09:00 AM.',
    },
    {
      id: 't-2',
      ticketNumber: 'HD-2026-00043',
      category: 'Hostel Facilities',
      title: 'Water filter motor faulty in Block C',
      description: 'Drinking water dispenser unit is not chilling properly.',
      priority: 'MEDIUM',
      status: 'WAITING_FOR_APPLICANT',
      slaHoursLeft: 44,
      lastReply: 'Warden Desk: Technician visit scheduled for 03:00 PM today.',
    },
    {
      id: 't-3',
      ticketNumber: 'HD-2026-00041',
      category: 'IT & LMS Support',
      title: 'Portal MFA Authenticator reset',
      description: 'Device changed, need TOTP reset.',
      priority: 'MEDIUM',
      status: 'RESOLVED',
      slaHoursLeft: 0,
      rating: 5,
      lastReply: 'IT Desk: MFA reset link sent to your registered institutional email.',
    },
  ]);

  // Form state
  const [newTitle, setNewTitle] = useState('');
  const [newCategory, setNewCategory] = useState('Academic Affairs');
  const [newPriority, setNewPriority] = useState<'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT'>('MEDIUM');
  const [newDescription, setNewDescription] = useState('');

  // Rating Modal
  const [ratingModalVisible, setRatingModalVisible] = useState(false);
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null);

  const handleCreateTicket = () => {
    if (!newTitle.trim() || !newDescription.trim()) {
      Alert.alert('Required Fields', 'Please enter a title and detailed query description.');
      return;
    }

    const newTicket: MobileTicket = {
      id: `t-${Date.now()}`,
      ticketNumber: `HD-2026-${String(tickets.length + 50).padStart(5, '0')}`,
      category: newCategory,
      title: newTitle.trim(),
      description: newDescription.trim(),
      priority: newPriority,
      status: 'OPEN',
      slaHoursLeft: 48,
    };

    setTickets([newTicket, ...tickets]);
    setNewTitle('');
    setNewDescription('');
    setActiveTab('MY_TICKETS');
    Alert.alert('Ticket Submitted', `Your ticket ${newTicket.ticketNumber} has been logged in the SLA radar.`);
  };

  const handleEscalate = (ticketId: string) => {
    setTickets((prev) =>
      prev.map((t) => (t.id === ticketId ? { ...t, priority: 'URGENT' } : t))
    );
    Alert.alert('Escalation Alert', 'Priority bumped to URGENT. Department supervisor notified.');
  };

  const handleRate = (stars: number) => {
    if (selectedTicketId) {
      setTickets((prev) =>
        prev.map((t) => (t.id === selectedTicketId ? { ...t, rating: stars } : t))
      );
    }
    setRatingModalVisible(false);
    setSelectedTicketId(null);
    Alert.alert('Feedback Recorded', `Thank you for rating our support team ${stars} stars!`);
  };

  const filteredTickets = tickets.filter(
    (t) =>
      t.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.ticketNumber.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.category.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Helpdesk & Grievance Portal</Text>
        <Text className="text-xs text-slate-300">
          Rank 13 Helpdesk · 24/7 SLA Monitored Institutional Query Desk
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('MY_TICKETS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'MY_TICKETS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'MY_TICKETS' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Tickets ({tickets.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('RAISE_QUERY')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'RAISE_QUERY' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'RAISE_QUERY' ? 'text-slate-900' : 'text-slate-600'}`}>
            + Raise Query
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: My Tickets */}
      {activeTab === 'MY_TICKETS' && (
        <View className="space-y-3">
          <View className="bg-white p-3 rounded-xl border border-slate-200">
            <TextInput
              placeholder="Search tickets by subject, #, or category..."
              value={searchQuery}
              onChangeText={setSearchQuery}
              className="text-xs text-slate-800"
            />
          </View>

          {filteredTickets.map((t) => (
            <View key={t.id} className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm space-y-2">
              <View className="flex-row justify-between items-center">
                <View className="flex-row items-center gap-1.5">
                  <Text className="font-mono text-xs font-bold text-blue-600 bg-blue-50 px-2 py-0.5 rounded">
                    {t.ticketNumber}
                  </Text>
                  <Text
                    className={`text-[10px] font-bold px-2 py-0.5 rounded-full text-white ${
                      t.priority === 'URGENT'
                        ? 'bg-rose-500'
                        : t.priority === 'HIGH'
                        ? 'bg-amber-500'
                        : 'bg-blue-500'
                    }`}
                  >
                    {t.priority}
                  </Text>
                </View>

                <View
                  className={`px-2 py-0.5 rounded text-[10px] font-semibold ${
                    t.status === 'RESOLVED'
                      ? 'bg-emerald-100 text-emerald-800'
                      : t.status === 'WAITING_FOR_APPLICANT'
                      ? 'bg-amber-100 text-amber-800'
                      : 'bg-slate-100 text-slate-700'
                  }`}
                >
                  <Text className="text-[10px] font-semibold text-slate-700">{t.status}</Text>
                </View>
              </View>

              <Text className="text-sm font-bold text-slate-900">{t.title}</Text>
              <Text className="text-xs text-slate-600">{t.description}</Text>

              {t.lastReply && (
                <View className="bg-slate-50 p-2.5 rounded-lg border border-slate-100">
                  <Text className="text-[11px] text-slate-700 italic">💬 {t.lastReply}</Text>
                </View>
              )}

              <View className="flex-row justify-between items-center pt-2 border-t border-slate-100">
                <Text className="text-[10px] text-slate-500 font-medium">
                  ⏱️ {t.slaHoursLeft > 0 ? `${t.slaHoursLeft}h SLA left` : 'SLA Completed'}
                </Text>

                <View className="flex-row gap-2">
                  {t.status === 'RESOLVED' && !t.rating && (
                    <TouchableOpacity
                      onPress={() => {
                        setSelectedTicketId(t.id);
                        setRatingModalVisible(true);
                      }}
                      className="bg-emerald-600 px-3 py-1.5 rounded-lg"
                    >
                      <Text className="text-white text-[10px] font-bold">★ Rate Resolution</Text>
                    </TouchableOpacity>
                  )}

                  {t.priority !== 'URGENT' && t.status !== 'RESOLVED' && (
                    <TouchableOpacity
                      onPress={() => handleEscalate(t.id)}
                      className="bg-rose-500 px-3 py-1.5 rounded-lg"
                    >
                      <Text className="text-white text-[10px] font-bold">⚡ Escalate</Text>
                    </TouchableOpacity>
                  )}
                </View>
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Tab: Raise Query */}
      {activeTab === 'RAISE_QUERY' && (
        <View className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm space-y-4">
          <Text className="text-base font-bold text-slate-900">New Query Submission</Text>

          <View className="space-y-1">
            <Text className="text-xs font-semibold text-slate-700">Category</Text>
            <View className="flex-row flex-wrap gap-2">
              {['Academic Affairs', 'Hostel Facilities', 'Fee & Billing', 'IT & LMS Support'].map((cat) => (
                <TouchableOpacity
                  key={cat}
                  onPress={() => setNewCategory(cat)}
                  className={`px-3 py-1.5 rounded-lg border ${newCategory === cat ? 'bg-blue-600 border-blue-600' : 'bg-slate-50 border-slate-200'}`}
                >
                  <Text className={`text-xs font-semibold ${newCategory === cat ? 'text-white' : 'text-slate-700'}`}>
                    {cat}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View className="space-y-1">
            <Text className="text-xs font-semibold text-slate-700">Priority</Text>
            <View className="flex-row gap-2">
              {(['LOW', 'MEDIUM', 'HIGH', 'URGENT'] as const).map((p) => (
                <TouchableOpacity
                  key={p}
                  onPress={() => setNewPriority(p)}
                  className={`flex-1 py-1.5 rounded-lg border items-center ${newPriority === p ? 'bg-slate-900 border-slate-900' : 'bg-slate-50 border-slate-200'}`}
                >
                  <Text className={`text-xs font-bold ${newPriority === p ? 'text-white' : 'text-slate-700'}`}>
                    {p}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View className="space-y-1">
            <Text className="text-xs font-semibold text-slate-700">Query Title</Text>
            <TextInput
              placeholder="e.g. Exam attendance threshold clarification"
              value={newTitle}
              onChangeText={setNewTitle}
              className="border border-slate-200 rounded-xl p-3 text-xs bg-slate-50"
            />
          </View>

          <View className="space-y-1">
            <Text className="text-xs font-semibold text-slate-700">Detailed Description</Text>
            <TextInput
              placeholder="Explain your inquiry or issue with exact details..."
              value={newDescription}
              onChangeText={setNewDescription}
              multiline
              numberOfLines={4}
              className="border border-slate-200 rounded-xl p-3 text-xs bg-slate-50 min-h-[100px]"
            />
          </View>

          <TouchableOpacity
            onPress={handleCreateTicket}
            className="bg-blue-600 p-3.5 rounded-xl items-center shadow-sm"
          >
            <Text className="text-white text-xs font-bold">Submit Query to Helpdesk</Text>
          </TouchableOpacity>
        </View>
      )}

      {/* CSAT Rating Modal */}
      <Modal visible={ratingModalVisible} transparent animationType="fade">
        <View className="flex-1 bg-black/60 justify-center items-center p-4">
          <View className="bg-white w-full max-w-xs rounded-2xl p-5 items-center space-y-4 shadow-xl">
            <Text className="text-base font-bold text-slate-900">Rate Ticket Resolution</Text>
            <Text className="text-xs text-slate-600 text-center">
              How satisfied are you with the resolution provided by our support officer?
            </Text>

            <View className="flex-row gap-2">
              {[1, 2, 3, 4, 5].map((star) => (
                <TouchableOpacity key={star} onPress={() => handleRate(star)} className="p-1">
                  <Text className="text-2xl text-amber-500">★</Text>
                </TouchableOpacity>
              ))}
            </View>

            <TouchableOpacity
              onPress={() => setRatingModalVisible(false)}
              className="bg-slate-100 px-4 py-2 rounded-lg"
            >
              <Text className="text-xs font-semibold text-slate-700">Cancel</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </ScrollView>
  );
}
