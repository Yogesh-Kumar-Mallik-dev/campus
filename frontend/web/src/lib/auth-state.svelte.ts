/**
 * BLOCK_WEB_AUTH_STATE_001
 * Purpose: Svelte 5 Rune-based reactive authentication & workspace state.
 * Domain:  Web Presentation Layer
 */

import { apiClient, ApiError } from '@campus/api-client';

export interface UserProfile {
  id: string;
  tenant_id: string;
  tier: 'SUPER_ADMIN' | 'ADMIN' | 'USER';
  username: string;
  email: string;
  roles: Array<{ role_key: string; role_name: string; scope_type?: string; scope_id?: string }>;
  permissions: string[];
  allowed_workspaces: string[];
}

class AuthState {
  // Svelte 5 Runes for fine-grained reactivity
  user = $state<UserProfile | null>(null);
  token = $state<string | null>(null);
  refreshToken = $state<string | null>(null);
  activeWorkspace = $state<string>('student_hub');
  isLoading = $state<boolean>(false);
  errorMessage = $state<string | null>(null);

  get isAuthenticated(): boolean {
    return this.user !== null && this.token !== null;
  }

  get availableWorkspaces(): string[] {
    return this.user?.allowed_workspaces || [];
  }

  async login(tenantId: string, identifier: string, password: string, mfaCode?: string) {
    this.isLoading = true;
    this.errorMessage = null;

    try {
      const res = await apiClient.post('/auth/login', {
        tenant_id: tenantId,
        identifier,
        password,
        mfa_code: mfaCode,
      });

      if (res.requires_mfa) {
        this.isLoading = false;
        return { requiresMFA: true, mfaTicket: res.mfa_ticket };
      }

      if (res.tokens) {
        this.token = res.tokens.access_token;
        this.refreshToken = res.tokens.refresh_token;
        apiClient.setToken(this.token);
      }

      if (res.user) {
        this.user = res.user;
        if (res.user.allowed_workspaces && res.user.allowed_workspaces.length > 0) {
          this.activeWorkspace = res.user.allowed_workspaces[0];
        }
      }

      return { requiresMFA: false };
    } catch (err) {
      if (err instanceof ApiError) {
        this.errorMessage = err.detail || err.title;
      } else {
        this.errorMessage = 'Failed to connect to authentication gateway';
      }
      throw err;
    } finally {
      this.isLoading = false;
    }
  }

  async logout() {
    try {
      if (this.token) {
        await apiClient.post('/auth/logout', {});
      }
    } catch {
      // Ignore network failures on logout
    } finally {
      this.user = null;
      this.token = null;
      this.refreshToken = null;
      apiClient.setToken(null);
    }
  }

  switchWorkspace(workspaceKey: string) {
    if (this.availableWorkspaces.includes(workspaceKey)) {
      this.activeWorkspace = workspaceKey;
    }
  }

  hasCapability(permissionKey: string): boolean {
    if (!this.user) return false;
    if (this.user.tier === 'SUPER_ADMIN') return true;
    return this.user.permissions?.includes(permissionKey) || false;
  }
}

export const authState = new AuthState();
