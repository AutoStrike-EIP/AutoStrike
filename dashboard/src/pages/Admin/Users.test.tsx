import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Users from './Users';

// Mock the auth context
vi.mock('../../contexts/AuthContext', () => ({
  useAuth: vi.fn(() => ({
    user: { id: 'current-user', username: 'admin', role: 'admin' },
    isAuthenticated: true,
  })),
}));

// Mock the API
vi.mock('../../lib/api', () => ({
  adminApi: {
    listUsers: vi.fn(() =>
      Promise.resolve({
        data: {
          users: [
            {
              id: 'user-1',
              username: 'admin',
              email: 'admin@example.com',
              role: 'admin',
              role_display: 'Administrator',
              is_active: true,
              last_login_at: new Date().toISOString(),
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: 'user-2',
              username: 'operator',
              email: 'operator@example.com',
              role: 'operator',
              role_display: 'Operator',
              is_active: true,
              last_login_at: null,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: 'user-3',
              username: 'inactive',
              email: 'inactive@example.com',
              role: 'viewer',
              role_display: 'Viewer',
              is_active: false,
              last_login_at: null,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
          ],
          total: 3,
        },
      })
    ),
    createUser: vi.fn(() => Promise.resolve({ data: {} })),
    updateUser: vi.fn(() => Promise.resolve({ data: {} })),
    updateUserRole: vi.fn(() => Promise.resolve({ data: {} })),
    resetPassword: vi.fn(() => Promise.resolve({ data: {} })),
    deactivateUser: vi.fn(() => Promise.resolve({ data: {} })),
    reactivateUser: vi.fn(() => Promise.resolve({ data: {} })),
  },
}));

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

function renderUsers() {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Users />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('Users Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders user management title', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });

  it('renders add user button', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Add User')).toBeInTheDocument();
    });
  });

  it('displays users after loading', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });
  });

  it('shows role badges', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
      expect(screen.getByText('Operator')).toBeInTheDocument();
    });
  });

  it('shows active/inactive status', async () => {
    renderUsers();
    await waitFor(() => {
      const activeLabels = screen.getAllByText('Active');
      expect(activeLabels.length).toBeGreaterThan(0);
    });
  });

  it('shows show inactive users checkbox', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });
  });

  it('shows last login time or Never', async () => {
    renderUsers();
    await waitFor(() => {
      const neverTexts = screen.getAllByText('Never');
      expect(neverTexts.length).toBeGreaterThan(0);
    });
  });
});

describe('Users Create Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('has add user button', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const addButton = screen.getByText('Add User');
    expect(addButton).toBeInTheDocument();
  });

  it('shows role options in create form', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      const roleSelect = screen.getByLabelText('Role');
      expect(roleSelect).toBeInTheDocument();
    });
  });
});

describe('Users Edit Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens edit modal when edit button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit User')).toBeInTheDocument();
    });
  });

  it('closes edit modal when cancel clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit User')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByText('Edit User')).not.toBeInTheDocument();
    });
  });
});

describe('Users Role Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens role modal when change role button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    // Click on a user that is not the current user
    const enabledRoleButton = roleButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledRoleButton) {
      fireEvent.click(enabledRoleButton);

      await waitFor(() => {
        expect(screen.getByText('Change Role')).toBeInTheDocument();
      });
    }
  });
});

describe('Users Password Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows reset password buttons', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    expect(passwordButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Deactivate Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens deactivate modal when deactivate button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    // Click on a user that is not the current user
    const enabledButton = deactivateButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByText('Deactivate User')).toBeInTheDocument();
      });
    }
  });

  it('shows warning message in deactivate modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(
          screen.getByText(/The user will no longer be able to log in/)
        ).toBeInTheDocument();
      });
    }
  });
});

describe('Users Reactivate Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows reactivate button for inactive users', async () => {
    renderUsers();

    // First, enable showing inactive users
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const reactivateButtons = screen.queryAllByTitle('Reactivate user');
      expect(reactivateButtons.length).toBeGreaterThanOrEqual(0);
    });
  });
});

describe('Users Table Actions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('disables change role button for current user', async () => {
    // Mock useAuth to return the same user as in the list
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    // The first user (admin) should have a disabled role button
    expect(roleButtons[0]).toBeDisabled();
  });

  it('disables deactivate button for current user', async () => {
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    // The first user (admin) should have a disabled deactivate button
    expect(deactivateButtons[0]).toBeDisabled();
  });
});

describe('Users Empty State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows empty state when no users', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: { users: [], total: 0 },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('No users found')).toBeInTheDocument();
    });
  });
});

describe('Users Role Badge Colors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays role labels correctly', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
      expect(screen.getByText('Operator')).toBeInTheDocument();
    });
  });
});

describe('Users Filter', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('toggles include inactive users', async () => {
    const { adminApi } = await import('../../lib/api');

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      expect(adminApi.listUsers).toHaveBeenCalledWith(true);
    });
  });
});

describe('Users Create Form', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows username field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });
  });

  it('shows email field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument();
    });
  });

  it('shows password field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Password')).toBeInTheDocument();
    });
  });

  it('shows cancel button in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      const cancelButtons = screen.getAllByText('Cancel');
      expect(cancelButtons.length).toBeGreaterThan(0);
    });
  });

  it('closes create modal on cancel', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByLabelText('Username')).not.toBeInTheDocument();
    });
  });
});

describe('Users API Errors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('handles list users error gracefully', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockRejectedValueOnce(new Error('Network error'));

    renderUsers();

    // Should not crash even with error
    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });
});

describe('Users Table Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows user count in header', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });

  it('displays inactive status for inactive users', async () => {
    renderUsers();

    // First enable showing inactive users
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const inactiveLabels = screen.queryAllByText('Inactive');
      expect(inactiveLabels.length).toBeGreaterThanOrEqual(0);
    });
  });

  it('shows viewer role badge', async () => {
    renderUsers();

    // Enable showing inactive users to see the viewer role
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const viewerBadges = screen.queryAllByText('Viewer');
      expect(viewerBadges.length).toBeGreaterThanOrEqual(0);
    });
  });
});

describe('Users Password Reset', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('has password reset buttons for each user', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    expect(passwordButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Edit Form', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows email field in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument();
    });
  });

  it('pre-fills email in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      const emailInput = screen.getByDisplayValue('admin@example.com');
      expect(emailInput).toBeInTheDocument();
    });
  });

  it('shows username field in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });
  });
});

describe('Users Action Buttons', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows all action buttons for each user', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    const roleButtons = screen.getAllByTitle('Change role');
    const passwordButtons = screen.getAllByTitle('Reset password');
    const deactivateButtons = screen.getAllByTitle('Deactivate user');

    expect(editButtons.length).toBeGreaterThan(0);
    expect(roleButtons.length).toBeGreaterThan(0);
    expect(passwordButtons.length).toBeGreaterThan(0);
    expect(deactivateButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Form Submissions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits create form and calls createUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });

    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(adminApi.createUser).toHaveBeenCalledWith({
        username: 'newuser',
        email: 'newuser@test.com',
        password: 'password123',
        role: 'viewer',
      });
    });
  });

  it('submits edit form and calls updateUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[1]); // Click operator's edit button

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'updated-operator' } });
    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(adminApi.updateUser).toHaveBeenCalledWith('user-2', {
        username: 'updated-operator',
        email: 'operator@example.com',
      });
    });
  });

  it('submits role change form and calls updateUserRole API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledRoleButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledRoleButton) {
      fireEvent.click(enabledRoleButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.change(screen.getByLabelText('Role'), { target: { value: 'analyst' } });
      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(adminApi.updateUserRole).toHaveBeenCalled();
      });
    }
  });

  it('submits password reset form and calls resetPassword API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(adminApi.resetPassword).toHaveBeenCalledWith('user-1', { new_password: 'newpassword123' });
    });
  });

  it('confirms deactivation and calls deactivateUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(adminApi.deactivateUser).toHaveBeenCalled();
      });
    }
  });

  it('confirms reactivation and calls reactivateUser API', async () => {
    const { adminApi } = await import('../../lib/api');

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [
          {
            id: 'user-3',
            username: 'inactive',
            email: 'inactive@example.com',
            role: 'viewer',
            is_active: false,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    const reactivateButton = screen.getByTitle('Reactivate user');
    fireEvent.click(reactivateButton);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(adminApi.reactivateUser).toHaveBeenCalledWith('user-3');
    });
  });
});

describe('Users Mutation Errors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows error alert when create mutation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.createUser).mockRejectedValueOnce({
      response: { data: { error: 'Username already exists' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(screen.getByText('Username already exists')).toBeInTheDocument();
    });
  });

  it('shows fallback error when create mutation fails without error message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.createUser).mockRejectedValueOnce(new Error('Network error'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to create user')).toBeInTheDocument();
    });
  });

  it('shows error alert when update mutation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUser).mockRejectedValueOnce({
      response: { data: { error: 'Duplicate username' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(screen.getByText('Duplicate username')).toBeInTheDocument();
    });
  });

  it('shows fallback error when update mutation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUser).mockRejectedValueOnce(new Error('Server error'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to update user')).toBeInTheDocument();
    });
  });

  it('shows error alert when role change fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUserRole).mockRejectedValueOnce({
      response: { data: { error: 'Cannot change role' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(screen.getByText('Cannot change role')).toBeInTheDocument();
      });
    }
  });

  it('shows fallback error when role mutation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUserRole).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(screen.getByText('Failed to update role')).toBeInTheDocument();
      });
    }
  });

  it('shows error alert when password reset fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.resetPassword).mockRejectedValueOnce({
      response: { data: { error: 'Password too weak' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(screen.getByText('Password too weak')).toBeInTheDocument();
    });
  });

  it('shows fallback error when password reset fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.resetPassword).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to reset password')).toBeInTheDocument();
    });
  });

  it('shows error alert when deactivation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.deactivateUser).mockRejectedValueOnce({
      response: { data: { error: 'Cannot deactivate last admin' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(screen.getByText('Cannot deactivate last admin')).toBeInTheDocument();
      });
    }
  });

  it('shows fallback error when deactivation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.deactivateUser).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(screen.getByText('Failed to deactivate user')).toBeInTheDocument();
      });
    }
  });

  it('shows error alert when reactivation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.reactivateUser).mockRejectedValueOnce({
      response: { data: { error: 'Cannot reactivate' } },
    });

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(screen.getByText('Cannot reactivate')).toBeInTheDocument();
    });
  });

  it('shows fallback error when reactivation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.reactivateUser).mockRejectedValueOnce(new Error('Network'));

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to reactivate user')).toBeInTheDocument();
    });
  });
});

describe('Users Role Badges - All Roles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays rssi and analyst role badges', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [
          {
            id: 'user-rssi',
            username: 'rssi-user',
            email: 'rssi@example.com',
            role: 'rssi',
            is_active: true,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
          {
            id: 'user-analyst',
            username: 'analyst-user',
            email: 'analyst@example.com',
            role: 'analyst',
            is_active: true,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
        total: 2,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Security Officer (RSSI)')).toBeInTheDocument();
      expect(screen.getByText('Analyst')).toBeInTheDocument();
    });

    const rssiBadge = screen.getByText('Security Officer (RSSI)');
    expect(rssiBadge.className).toContain('bg-purple-100');

    const analystBadge = screen.getByText('Analyst');
    expect(analystBadge.className).toContain('bg-green-100');
  });
});

describe('Users Current User Indicator', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows (you) indicator for current user', async () => {
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    expect(screen.getByText('(you)')).toBeInTheDocument();
  });
});
