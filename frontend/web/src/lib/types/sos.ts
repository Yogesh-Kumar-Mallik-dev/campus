/**
 * BLOCK_TYPES_SOS_001
 * Subsystem: Rank 14 - SOS & Emergency Response (sos)
 * Purpose:   TypeScript domain definitions for emergency incidents, dispatch responders, and telemetry broadcast.
 */

export type SOSEmergencyType =
  | 'MEDICAL'
  | 'FIRE'
  | 'SECURITY_THREAT'
  | 'NATURAL_HAZARD'
  | 'HARASSMENT_RAGGING'
  | 'OTHER';

export type SOSIncidentStatus =
  | 'TRIGGERED'
  | 'ACKNOWLEDGED'
  | 'DISPATCHED'
  | 'ON_SCENE'
  | 'RESOLVED'
  | 'FALSE_ALARM';

export type SOSResponderRole =
  | 'CAMPUS_SECURITY'
  | 'PARAMEDIC'
  | 'HOSTEL_WARDEN'
  | 'FIRE_SAFETY_OFFICER'
  | 'POLICE_LIAISON';

export type SOSResponderStatus = 'ASSIGNED' | 'EN_ROUTE' | 'ON_SCENE' | 'COMPLETED';

export type SOSTelemetryChannel = 'SMS_GATEWAY' | 'PUSH_NOTIFICATION' | 'CAMPUS_PA' | 'PUBLIC_SIREN';

export interface SOSIncident {
  id: string;
  tenantId: string;
  alertNumber: string;
  userId: string;
  userName?: string;
  userPhone?: string;
  emergencyType: SOSEmergencyType;
  latitude: number;
  longitude: number;
  locationDescription: string;
  status: SOSIncidentStatus;
  triggeredAt: string;
  acknowledgedAt?: string;
  resolvedAt?: string;
  resolvedById?: string;
  resolvedByName?: string;
  resolutionNotes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface SOSDispatchResponder {
  id: string;
  tenantId: string;
  incidentId: string;
  responderId: string;
  responderName?: string;
  role: SOSResponderRole;
  status: SOSResponderStatus;
  dispatchedAt: string;
  arrivedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface SOSTelemetryBroadcast {
  id: string;
  tenantId: string;
  incidentId: string;
  channel: SOSTelemetryChannel;
  payload: string;
  deliveredCount: number;
  broadcastAt: string;
  createdAt: string;
}
