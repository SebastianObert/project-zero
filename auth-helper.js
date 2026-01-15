// auth-helper.js - Helper functions untuk authentication
const API_BASE_URL = 'http://localhost:8080';

// Cek apakah user sudah login
function isLoggedIn() {
    return !!localStorage.getItem('token');
}

// Get current user
function getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

// Get auth headers dengan JWT token
function getAuthHeaders(includeContentType = true) {
    const token = localStorage.getItem('token');
    const headers = {
        'Authorization': `Bearer ${token}`
    };
    
    if (includeContentType) {
        headers['Content-Type'] = 'application/json';
    }
    
    return headers;
}

// Fetch dengan auth
async function authFetch(url, options = {}) {
    if (!isLoggedIn()) {
        window.location.href = '/login.html';
        throw new Error('Not authenticated');
    }
    
    const headers = options.headers || {};
    const token = localStorage.getItem('token');
    
    // Jangan tambah Content-Type jika body adalah FormData
    if (!(options.body instanceof FormData)) {
        headers['Content-Type'] = 'application/json';
    }
    
    headers['Authorization'] = `Bearer ${token}`;
    
    const response = await fetch(url, {
        ...options,
        headers
    });
    
    // Jika unauthorized, redirect ke login
    if (response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login.html';
        throw new Error('Unauthorized');
    }
    
    return response;
}

// Logout
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login.html';
}

// Check login on page load
function requireAuth() {
    if (!isLoggedIn()) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// Update UI dengan info user
function updateUserUI() {
    const user = getCurrentUser();
    if (user && user.name) {
        // Update semua element yang menampilkan nama user
        document.querySelectorAll('.user-name').forEach(el => {
            el.textContent = user.name;
        });
        
        // Update initial
        document.querySelectorAll('.user-initial').forEach(el => {
            const initials = user.name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2);
            el.textContent = initials;
        });
    }
}
