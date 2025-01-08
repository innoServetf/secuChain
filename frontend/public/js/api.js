const API_BASE_URL = 'http://localhost:8080/api/v1';

class Api {
    static async request(url, options = {}) {
        const token = localStorage.getItem('token') || sessionStorage.getItem('token');
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(`${API_BASE_URL}${url}`, {
                ...options,
                headers,
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || '请求失败');
            }

            return data;
        } catch (error) {
            if (error.message === 'Failed to fetch') {
                throw new Error('网络连接失败，请检查网络设置');
            }
            throw error;
        }
    }

    static async login(username, password) {
        return await this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }

    static async register(username, email, password) {
        return await this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
        });
    }

    static async getUserInfo() {
        return await this.request('/user/info');
    }

    static async updatePassword(oldPassword, newPassword) {
        return await this.request('/user/password', {
            method: 'PUT',
            body: JSON.stringify({ 
                old_password: oldPassword,
                new_password: newPassword 
            }),
        });
    }
} 