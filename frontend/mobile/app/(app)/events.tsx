/**
 * BLOCK_MOBILE_EVENTS_SCREEN_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   Mobile student campus events discovery, 1-tap ticket booking, and digital QR gate pass wallet.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface EventItem {
  id: string;
  title: string;
  category: string;
  venue: string;
  dateTime: string;
  price: string;
  seatsLeft: number;
}

interface MyTicket {
  id: string;
  eventTitle: string;
  venue: string;
  dateTime: string;
  ticketCode: string;
  status: 'CONFIRMED' | 'CHECKED_IN';
}

export default function MobileEventsScreen() {
  const [activeTab, setActiveTab] = useState<'DISCOVER' | 'MY_TICKETS'>('DISCOVER');
  const [searchQuery, setSearchQuery] = useState('');

  const [events, setEvents] = useState<EventItem[]>([
    {
      id: 'e-1',
      title: 'Annual Tech Symposium 2026',
      category: 'TECHNICAL',
      venue: 'Main Auditorium',
      dateTime: '25 Sep · 09:00 AM',
      price: 'FREE',
      seatsLeft: 158,
    },
    {
      id: 'e-2',
      title: 'Tarang Cultural Night 2026',
      category: 'CULTURAL',
      venue: 'Open Amphitheatre',
      dateTime: '02 Oct · 06:00 PM',
      price: '₹150',
      seatsLeft: 190,
    },
    {
      id: 'e-3',
      title: 'Next-Gen Quantum Keynote',
      category: 'GUEST_LECTURE',
      venue: 'Seminar Hall 3B',
      dateTime: '28 Sep · 02:00 PM',
      price: 'FREE',
      seatsLeft: 32,
    },
  ]);

  const [tickets, setTickets] = useState<MyTicket[]>([
    {
      id: 't-1',
      eventTitle: 'Annual Tech Symposium 2026',
      venue: 'Main Auditorium',
      dateTime: '25 Sep · 09:00 AM',
      ticketCode: 'TCK-8f921a',
      status: 'CONFIRMED',
    },
  ]);

  const handleBookTicket = (event: EventItem) => {
    if (event.seatsLeft <= 0) {
      Alert.alert('Sold Out', 'No seats available for this event.');
      return;
    }
    const hash = Math.random().toString(16).substring(2, 8);
    const newTck: MyTicket = {
      id: `t-${Date.now()}`,
      eventTitle: event.title,
      venue: event.venue,
      dateTime: event.dateTime,
      ticketCode: `TCK-${hash}`,
      status: 'CONFIRMED',
    };
    setTickets([newTck, ...tickets]);
    setEvents((prev) =>
      prev.map((e) => (e.id === event.id ? { ...e, seatsLeft: e.seatsLeft - 1 } : e))
    );
    Alert.alert('Pass Confirmed', `Your QR pass for "${event.title}" is ready in your tickets wallet.`);
  };

  const handleSimulateCheckIn = (ticketCode: string) => {
    setTickets((prev) =>
      prev.map((t) => (t.ticketCode === ticketCode ? { ...t, status: 'CHECKED_IN' } : t))
    );
    Alert.alert('Gate Check-In Success', `Ticket ${ticketCode} verified at turnstile scanner.`);
  };

  const filteredEvents = events.filter(
    (e) =>
      e.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      e.venue.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">Campus Events Hub</Text>
        <Text className="text-xs text-slate-300">
          Rank 12 Events · Hackathons, Festivals, RSVPs & Digital QR Passes
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('DISCOVER')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'DISCOVER' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'DISCOVER' ? 'text-slate-900' : 'text-slate-600'}`}>
            Discover ({events.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('MY_TICKETS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'MY_TICKETS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'MY_TICKETS' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Passes ({tickets.length})
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: Discover */}
      {activeTab === 'DISCOVER' && (
        <View className="space-y-3">
          <View className="bg-white p-3 rounded-xl border border-slate-200">
            <TextInput
              placeholder="Search by event title or venue..."
              value={searchQuery}
              onChangeText={setSearchQuery}
              className="text-xs text-slate-800"
            />
          </View>

          {filteredEvents.map((evnt) => (
            <View key={evnt.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-3 shadow-sm">
              <View className="flex-row justify-between items-start">
                <View className="flex-1 mr-2">
                  <Text className="text-sm font-bold text-slate-900">{evnt.title}</Text>
                  <Text className="text-xs text-slate-500">📍 {evnt.venue}</Text>
                </View>
                <View className="bg-indigo-50 px-2 py-0.5 rounded">
                  <Text className="text-[10px] font-bold text-indigo-700">{evnt.category}</Text>
                </View>
              </View>

              <View className="flex-row justify-between bg-slate-50 p-2.5 rounded-xl">
                <View>
                  <Text className="text-[10px] text-slate-400">Schedule</Text>
                  <Text className="text-xs font-semibold text-slate-700">{evnt.dateTime}</Text>
                </View>
                <View className="items-end">
                  <Text className="text-[10px] text-slate-400">Entry</Text>
                  <Text className="text-xs font-bold text-slate-900">{evnt.price}</Text>
                </View>
              </View>

              <View className="flex-row items-center justify-between pt-1 border-t border-slate-100">
                <Text className="text-[11px] text-emerald-700 font-semibold">{evnt.seatsLeft} seats left</Text>
                <TouchableOpacity
                  onPress={() => handleBookTicket(evnt)}
                  className="px-4 py-2 bg-slate-900 rounded-xl"
                >
                  <Text className="text-xs font-bold text-white">🎟️ Book Pass</Text>
                </TouchableOpacity>
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Tab: My Tickets */}
      {activeTab === 'MY_TICKETS' && (
        <View className="space-y-4">
          {tickets.map((tck) => (
            <View key={tck.id} className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm space-y-4 items-center">
              <View className="flex-row justify-between items-center w-full">
                <Text className="text-sm font-bold text-slate-900 flex-1 mr-2">{tck.eventTitle}</Text>
                <View
                  className={`px-2 py-0.5 rounded ${
                    tck.status === 'CHECKED_IN' ? 'bg-emerald-100' : 'bg-blue-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      tck.status === 'CHECKED_IN' ? 'text-emerald-800' : 'text-blue-800'
                    }`}
                  >
                    {tck.status === 'CHECKED_IN' ? 'ADMITTED' : 'CONFIRMED PASS'}
                  </Text>
                </View>
              </View>

              {/* QR Mockup */}
              <View className="w-44 h-44 bg-slate-900 rounded-2xl p-4 items-center justify-center space-y-2">
                <View className="w-24 h-24 bg-white rounded-xl items-center justify-center">
                  <Text className="text-3xl">🎟️</Text>
                </View>
                <Text className="font-mono text-xs font-bold text-white tracking-widest">
                  {tck.ticketCode}
                </Text>
              </View>

              <View className="items-center">
                <Text className="text-xs font-semibold text-slate-700">📍 {tck.venue}</Text>
                <Text className="text-[11px] text-slate-400 mt-0.5">⏰ {tck.dateTime}</Text>
              </View>

              {tck.status === 'CONFIRMED' ? (
                <TouchableOpacity
                  onPress={() => handleSimulateCheckIn(tck.ticketCode)}
                  className="p-3 bg-emerald-700 rounded-xl w-full items-center"
                >
                  <Text className="text-xs font-bold text-white">Simulate Gate Scanner Punch</Text>
                </TouchableOpacity>
              ) : (
                <View className="p-2.5 bg-emerald-50 rounded-xl border border-emerald-100 w-full items-center">
                  <Text className="text-xs font-bold text-emerald-800">✓ Checked-In & Admitted</Text>
                </View>
              )}
            </View>
          ))}
        </View>
      )}
    </ScrollView>
  );
}
