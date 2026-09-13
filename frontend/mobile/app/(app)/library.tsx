/**
 * BLOCK_MOBILE_LIBRARY_SCREEN_001
 * Subsystem: Rank 9 - E-Library Management System (library)
 * Purpose:   Mobile campus library portal: catalog discovery, book reservations, loan renewal, and digital e-resource access.
 */

import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, TextInput } from 'react-native';

interface BookItem {
  id: string;
  title: string;
  author: string;
  isbn: string;
  category: string;
  availableCopies: number;
  isDigital: boolean;
}

interface ActiveLoan {
  id: string;
  title: string;
  borrowedAt: string;
  dueDate: string;
  renewalCount: number;
  overdueFine: number;
}

export default function MobileLibraryScreen() {
  const [activeTab, setActiveTab] = useState<'CATALOG' | 'LOANS' | 'EBOOKS'>('CATALOG');
  const [searchQuery, setSearchQuery] = useState<string>('');

  const [books] = useState<BookItem[]>([
    {
      id: 'book-01',
      title: 'Introduction to Algorithms (4th Edition)',
      author: 'Thomas H. Cormen, Charles E. Leiserson',
      isbn: '978-0262046305',
      category: 'Computer Science',
      availableCopies: 3,
      isDigital: true,
    },
    {
      id: 'book-02',
      title: 'Operating System Concepts',
      author: 'Abraham Silberschatz, Peter B. Galvin',
      isbn: '978-1119800361',
      category: 'Operating Systems',
      availableCopies: 0,
      isDigital: false,
    },
    {
      id: 'book-03',
      title: 'Database System Concepts',
      author: 'Avi Silberschatz, Henry F. Korth',
      isbn: '978-0078022159',
      category: 'Databases',
      availableCopies: 2,
      isDigital: true,
    },
  ]);

  const [loans, setLoans] = useState<ActiveLoan[]>([
    {
      id: 'loan-01',
      title: 'Introduction to Algorithms (4th Edition)',
      borrowedAt: '2026-09-01',
      dueDate: '2026-09-15',
      renewalCount: 0,
      overdueFine: 0,
    },
    {
      id: 'loan-02',
      title: 'Clean Architecture in Go',
      borrowedAt: '2026-08-20',
      dueDate: '2026-09-03',
      renewalCount: 1,
      overdueFine: 50, // 10 days overdue * ₹5/day
    },
  ]);

  const handleRenew = (loanId: string) => {
    setLoans((prev) =>
      prev.map((l) => {
        if (l.id === loanId) {
          if (l.renewalCount >= 2) {
            Alert.alert('Renewal Limit Reached', 'Books can only be renewed up to 2 times.');
            return l;
          }
          Alert.alert('Loan Renewed', `Due date extended by 14 days for ${l.title}.`);
          return {
            ...l,
            renewalCount: l.renewalCount + 1,
            dueDate: '2026-09-29',
          };
        }
        return l;
      })
    );
  };

  const handleReserve = (book: BookItem) => {
    Alert.alert(
      'Hold Placed',
      `You will be notified when a copy of "${book.title}" is returned to the circulation desk.`
    );
  };

  const handleOpenEbook = (book: BookItem) => {
    Alert.alert('Digital Access Granted', `Opening secure PDF reader for "${book.title}". DRM active.`);
  };

  const filteredBooks = books.filter(
    (b) =>
      b.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      b.author.toLowerCase().includes(searchQuery.toLowerCase()) ||
      b.isbn.includes(searchQuery)
  );

  return (
    <ScrollView className="flex-1 bg-slate-50 p-4">
      {/* Header Banner */}
      <View className="p-5 bg-slate-900 rounded-2xl mb-4">
        <Text className="text-xl font-bold text-white mb-1">E-Library Portal</Text>
        <Text className="text-xs text-slate-300">
          Rank 9 E-Library · Catalog, Borrow Loans, Fines & Digital Resources
        </Text>
      </View>

      {/* Mode Switcher */}
      <View className="flex-row bg-slate-200 p-1 rounded-xl mb-4">
        <TouchableOpacity
          onPress={() => setActiveTab('CATALOG')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'CATALOG' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'CATALOG' ? 'text-slate-900' : 'text-slate-600'}`}>
            Catalog ({books.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('LOANS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'LOANS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'LOANS' ? 'text-slate-900' : 'text-slate-600'}`}>
            My Loans ({loans.length})
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => setActiveTab('EBOOKS')}
          className={`flex-1 py-2 rounded-lg items-center ${activeTab === 'EBOOKS' ? 'bg-white shadow-sm' : ''}`}
        >
          <Text className={`text-xs font-bold ${activeTab === 'EBOOKS' ? 'text-slate-900' : 'text-slate-600'}`}>
            E-Books
          </Text>
        </TouchableOpacity>
      </View>

      {/* Tab: Catalog */}
      {activeTab === 'CATALOG' && (
        <View className="space-y-3">
          <View className="bg-white p-3 rounded-xl border border-slate-200">
            <TextInput
              placeholder="Search by Title, Author, or ISBN..."
              value={searchQuery}
              onChangeText={setSearchQuery}
              className="text-xs text-slate-800"
            />
          </View>

          {filteredBooks.map((book) => (
            <View key={book.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
              <View className="flex-row justify-between items-start">
                <View className="flex-1 mr-2">
                  <Text className="text-sm font-bold text-slate-900">{book.title}</Text>
                  <Text className="text-xs text-slate-500">{book.author}</Text>
                </View>
                <View
                  className={`px-2 py-0.5 rounded ${
                    book.availableCopies > 0 ? 'bg-emerald-100' : 'bg-rose-100'
                  }`}
                >
                  <Text
                    className={`text-[10px] font-bold ${
                      book.availableCopies > 0 ? 'text-emerald-800' : 'text-rose-800'
                    }`}
                  >
                    {book.availableCopies > 0 ? `${book.availableCopies} Available` : 'Reserved'}
                  </Text>
                </View>
              </View>

              <View className="flex-row items-center space-x-2 pt-1">
                <Text className="text-[10px] bg-slate-100 text-slate-600 px-2 py-0.5 rounded font-mono">
                  {book.isbn}
                </Text>
                <Text className="text-[10px] bg-indigo-50 text-indigo-700 px-2 py-0.5 rounded font-semibold">
                  {book.category}
                </Text>
                {book.isDigital && (
                  <Text className="text-[10px] bg-blue-50 text-blue-700 px-2 py-0.5 rounded font-semibold">
                    ⚡ PDF E-Book
                  </Text>
                )}
              </View>

              <View className="flex-row space-x-2 pt-2 border-t border-slate-100">
                {book.availableCopies > 0 ? (
                  <TouchableOpacity
                    onPress={() => Alert.alert('Borrow Request', 'Please present your Student ID barcode at the Circulation Desk.')}
                    className="flex-1 py-2 bg-slate-900 rounded-lg items-center"
                  >
                    <Text className="text-xs font-bold text-white">Borrow Physical</Text>
                  </TouchableOpacity>
                ) : (
                  <TouchableOpacity
                    onPress={() => handleReserve(book)}
                    className="flex-1 py-2 bg-amber-600 rounded-lg items-center"
                  >
                    <Text className="text-xs font-bold text-white">Place Hold / Reserve</Text>
                  </TouchableOpacity>
                )}
                {book.isDigital && (
                  <TouchableOpacity
                    onPress={() => handleOpenEbook(book)}
                    className="flex-1 py-2 bg-blue-600 rounded-lg items-center"
                  >
                    <Text className="text-xs font-bold text-white">Read Digital</Text>
                  </TouchableOpacity>
                )}
              </View>
            </View>
          ))}
        </View>
      )}

      {/* Tab: Loans */}
      {activeTab === 'LOANS' && (
        <View className="space-y-3">
          {loans.map((loan) => (
            <View key={loan.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
              <View className="flex-row justify-between items-start">
                <Text className="text-sm font-bold text-slate-900 flex-1">{loan.title}</Text>
                {loan.overdueFine > 0 && (
                  <View className="bg-rose-100 px-2 py-0.5 rounded">
                    <Text className="text-[10px] font-bold text-rose-800">
                      Fine: ₹{loan.overdueFine}
                    </Text>
                  </View>
                )}
              </View>

              <View className="flex-row justify-between bg-slate-50 p-2 rounded-lg">
                <View>
                  <Text className="text-[10px] text-slate-400">Borrowed On</Text>
                  <Text className="text-xs font-semibold text-slate-700">{loan.borrowedAt}</Text>
                </View>
                <View>
                  <Text className="text-[10px] text-slate-400">Due Date</Text>
                  <Text
                    className={`text-xs font-bold ${
                      loan.overdueFine > 0 ? 'text-rose-600' : 'text-slate-900'
                    }`}
                  >
                    {loan.dueDate}
                  </Text>
                </View>
                <View>
                  <Text className="text-[10px] text-slate-400">Renewals</Text>
                  <Text className="text-xs font-semibold text-slate-700">{loan.renewalCount}/2</Text>
                </View>
              </View>

              <TouchableOpacity
                onPress={() => handleRenew(loan.id)}
                disabled={loan.renewalCount >= 2}
                className={`py-2 rounded-lg items-center ${
                  loan.renewalCount >= 2 ? 'bg-slate-200' : 'bg-slate-900'
                }`}
              >
                <Text
                  className={`text-xs font-bold ${
                    loan.renewalCount >= 2 ? 'text-slate-400' : 'text-white'
                  }`}
                >
                  {loan.renewalCount >= 2 ? 'Max Renewals Reached' : 'Renew Loan (+14 Days)'}
                </Text>
              </TouchableOpacity>
            </View>
          ))}
        </View>
      )}

      {/* Tab: E-Books */}
      {activeTab === 'EBOOKS' && (
        <View className="space-y-3">
          <View className="p-4 bg-blue-50 border border-blue-200 rounded-2xl">
            <Text className="text-sm font-bold text-blue-900">Campus E-Book Vault</Text>
            <Text className="text-xs text-blue-700 mt-1">
              Unlimited concurrent access to course textbooks and IEEE research archives for enrolled students.
            </Text>
          </View>

          {books
            .filter((b) => b.isDigital)
            .map((ebook) => (
              <View key={ebook.id} className="p-4 bg-white rounded-2xl border border-slate-200 space-y-2">
                <Text className="text-sm font-bold text-slate-900">{ebook.title}</Text>
                <Text className="text-xs text-slate-500">{ebook.author}</Text>
                <TouchableOpacity
                  onPress={() => handleOpenEbook(ebook)}
                  className="py-2.5 bg-blue-600 rounded-xl items-center mt-2"
                >
                  <Text className="text-xs font-bold text-white">📖 Open DRM E-Book</Text>
                </TouchableOpacity>
              </View>
            ))}
        </View>
      )}
    </ScrollView>
  );
}
