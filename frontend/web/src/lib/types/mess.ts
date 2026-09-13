/**
 * BLOCK_TYPES_MESS_001
 * Subsystem: Rank 8 - Mess Management System (mess)
 * Purpose:   TypeScript domain models for campus dining halls, weekly menus, QR tokens, and rebates.
 */

export type MealType = 'BREAKFAST' | 'LUNCH' | 'SNACKS' | 'DINNER';
export type DietCategory = 'VEG' | 'NON_VEG' | 'JAIN' | 'VEGAN';
export type MessPlanType = 'SEMESTER' | 'MONTHLY' | 'PER_MEAL' | 'FLEXI';
export type MessSubscriptionStatus = 'ACTIVE' | 'PAUSED' | 'EXPIRED' | 'CANCELLED';
export type DiningTokenStatus = 'GENERATED' | 'REDEEMED' | 'EXPIRED' | 'CANCELLED';
export type MessRebateStatus = 'PENDING' | 'APPROVED' | 'REJECTED' | 'CREDITED';
export type DayOfWeek = 'MONDAY' | 'TUESDAY' | 'WEDNESDAY' | 'THURSDAY' | 'FRIDAY' | 'SATURDAY' | 'SUNDAY';

export interface MessHall {
  id: string;
  tenant_id: string;
  name: string;
  code: string;
  capacity: number;
  location: string;
  caterer_name: string;
  manager_id?: string;
  is_active: boolean;
  created_at: string;
}

export interface MessMenuItem {
  id: string;
  tenant_id: string;
  mess_hall_id: string;
  day_of_week: DayOfWeek;
  meal_type: MealType;
  diet_category: DietCategory;
  title: string;
  description?: string;
  calories: number;
}

export interface MessSubscription {
  id: string;
  tenant_id: string;
  student_id: string;
  mess_hall_id: string;
  plan_type: MessPlanType;
  diet_preference: DietCategory;
  start_date: string;
  end_date: string;
  status: MessSubscriptionStatus;
  monthly_fee: number;
}

export interface MessDiningToken {
  id: string;
  tenant_id: string;
  student_id: string;
  mess_hall_id: string;
  meal_type: MealType;
  token_code: string;
  meal_date: string;
  status: DiningTokenStatus;
  redeemed_at?: string;
  device_reader_id?: string;
  created_at: string;
}

export interface MessRebateApplication {
  id: string;
  tenant_id: string;
  student_id: string;
  student_name?: string;
  roll_number?: string;
  subscription_id: string;
  from_date: string;
  to_date: string;
  total_days: number;
  daily_rate: number;
  total_rebate_amount: number;
  reason: string;
  gate_pass_id?: string;
  status: MessRebateStatus;
  approved_by_id?: string;
  created_at: string;
}

export interface MessMealFeedback {
  id: string;
  tenant_id: string;
  student_id: string;
  mess_hall_id: string;
  meal_date: string;
  meal_type: MealType;
  food_rating: number;
  hygiene_rating: number;
  comments?: string;
  created_at: string;
}
