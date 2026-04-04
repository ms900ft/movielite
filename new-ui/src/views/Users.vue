<template>
  <div class="users-container">
    <h1>Users</h1>

    <div v-if="loading" class="loading">Loading users...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <button @click="showCreateDialog" class="add-button">Add User</button>

      <table class="users-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Admin</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.ID">
            <td>{{ user.ID }}</td>
            <td>{{ user.UserName }}</td>
            <td>{{ user.IsAdmin ? 'Yes' : 'No' }}</td>
            <td>
              <button @click="showEditDialog(user)" class="edit-btn">Edit</button>
              <button @click="confirmDelete(user)" class="delete-btn">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Dialog v-model:visible="dialogVisible" modal :header="isEditing ? 'Edit User' : 'Add User'" :style="{ width: '30rem' }">
      <div class="form-group">
        <label for="username">Username:</label>
        <input id="username" v-model="formData.username" type="text" placeholder="Username" />
      </div>
      <div class="form-group">
        <label for="password">Password:</label>
        <input id="password" v-model="formData.password" type="password" :placeholder="isEditing ? 'Leave blank to keep current' : 'Password'" />
      </div>
      <div class="form-group">
        <label>
          <input v-model="formData.isAdmin" type="checkbox" />
          Admin
        </label>
      </div>
      <div class="dialog-actions">
        <button @click="dialogVisible = false" class="cancel-btn">Cancel</button>
        <button @click="saveUser" class="save-btn">{{ isEditing ? 'Update' : 'Create' }}</button>
      </div>
    </Dialog>

    <Dialog v-model:visible="deleteDialogVisible" modal header="Confirm Delete" :style="{ width: '30rem' }">
      <p>Are you sure you want to delete user "{{ selectedUser?.UserName }}"?</p>
      <div class="dialog-actions">
        <button @click="deleteDialogVisible = false" class="cancel-btn">Cancel</button>
        <button @click="deleteUser" class="delete-btn">Delete</button>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Dialog from 'primevue/dialog';
import { moviesService } from '../services/movies.js';

const users = ref([]);
const loading = ref(true);
const error = ref(null);
const dialogVisible = ref(false);
const deleteDialogVisible = ref(false);
const isEditing = ref(false);
const selectedUser = ref(null);
const formData = ref({
  username: '',
  password: '',
  isAdmin: false
});

const fetchUsers = async () => {
  try {
    loading.value = true;
    const response = await moviesService.getUsers();
    users.value = response || [];
  } catch (err) {
    error.value = 'Failed to load users.';
    console.error('Error fetching users:', err);
  } finally {
    loading.value = false;
  }
};

const showCreateDialog = () => {
  isEditing.value = false;
  selectedUser.value = null;
  formData.value = { username: '', password: '', isAdmin: false };
  dialogVisible.value = true;
};

const showEditDialog = (user) => {
  isEditing.value = true;
  selectedUser.value = user;
  formData.value = { username: user.UserName, password: '', isAdmin: user.IsAdmin };
  dialogVisible.value = true;
};

const saveUser = async () => {
  if (!formData.value.username.trim()) {
    alert('Username is required');
    return;
  }
  if (!isEditing.value && !formData.value.password) {
    alert('Password is required for new users');
    return;
  }
  try {
    if (isEditing.value) {
      const payload = { username: formData.value.username };
      if (formData.value.password) {
        payload.password = formData.value.password;
      }
      payload.isadmin = formData.value.isAdmin;
      await moviesService.updateUser(selectedUser.value.id, payload);
    } else {
      await moviesService.createUser({
        username: formData.value.username,
        password: formData.value.password,
        isadmin: formData.value.isAdmin
      });
    }
    dialogVisible.value = false;
    fetchUsers();
  } catch (err) {
    console.error('Error saving user:', err);
    alert('Failed to save user.');
  }
};

const confirmDelete = (user) => {
  selectedUser.value = user;
  deleteDialogVisible.value = true;
};

const deleteUser = async () => {
  try {
    await moviesService.deleteUser(selectedUser.value.ID);
    deleteDialogVisible.value = false;
    fetchUsers();
  } catch (err) {
    console.error('Error deleting user:', err);
    alert('Failed to delete user.');
  }
};

onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.users-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  margin-bottom: 20px;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  font-size: 18px;
}

.error {
  color: red;
}

.add-button {
  padding: 8px 16px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-bottom: 16px;
}

.add-button:hover {
  background-color: #218838;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #dee2e6;
}

.users-table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.users-table tr:hover {
  background-color: #f1f3f5;
}

.edit-btn {
  padding: 4px 12px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 8px;
  font-size: 13px;
}

.edit-btn:hover {
  background-color: #0056b3;
}

.delete-btn {
  padding: 4px 12px;
  background-color: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.delete-btn:hover {
  background-color: #c82333;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 4px;
  font-weight: 600;
}

.form-group input[type="text"],
.form-group input[type="password"] {
  width: 100%;
  padding: 8px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.cancel-btn {
  padding: 8px 16px;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-btn:hover {
  background-color: #5a6268;
}

.save-btn {
  padding: 8px 16px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.save-btn:hover {
  background-color: #218838;
}
</style>
