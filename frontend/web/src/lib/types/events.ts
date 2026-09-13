/**
 * BLOCK_TYPES_EVENTS_001
 * Subsystem: Rank 12 - Event Organisation System (events)
 * Purpose:   TypeScript domain models for campus events, venue bookings, ticketing, and gate check-ins.
 */

export type EventCategory =
  | 'CULTURAL'
  | 'TECHNICAL'
  | 'SPORTS'
  | 'ACADEMIC_WORKSHOP'
  | 'GUEST_LECTURE'
  | 'CAREER_FAIR';

export type EventStatus = 'DRAFT' | 'PUBLISHED' | 'CANCELLED' | 'COMPLETED';

export type TicketStatus = 'CONFIRMED' | 'CHECKED_IN' | 'CANCELLED';

export type VenueBookingStatus = 'CONFIRMED' | 'CANCELLED';

export interface CampusEvent {
  id: string;
  tenantId: string;
  title: string;
  description: string;
  category: EventCategory;
  venueName: string;
  venueCapacity: number;
  startTime: string;
  endTime: string;
  registrationDeadline: string;
  organizerId: string;
  organizerName?: string;
  status: EventStatus;
  isTicketed: boolean;
  ticketPrice: number;
  maxTickets: number;
  ticketsSold: number;
  createdAt: string;
  updatedAt: string;
}

export interface EventTicket {
  id: string;
  tenantId: string;
  eventId: string;
  eventTitle?: string;
  venueName?: string;
  startTime?: string;
  studentId: string;
  studentName?: string;
  rollNumber?: string;
  ticketCode: string;
  status: TicketStatus;
  checkedInAt?: string;
  checkedInById?: string;
  createdAt: string;
  updatedAt: string;
}

export interface EventVenueBooking {
  id: string;
  tenantId: string;
  venueName: string;
  eventId: string;
  eventTitle?: string;
  bookedStartTime: string;
  bookedEndTime: string;
  bookedById: string;
  status: VenueBookingStatus;
  createdAt: string;
  updatedAt: string;
}
