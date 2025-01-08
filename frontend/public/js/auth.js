document.addEventListener('DOMContentLoaded', function() {
    // 登录表单处理
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // 注册表单处理
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
});

// 登录处理
async function handleLogin(e) {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const remember = document.getElementById('remember').checked;

    try {
        // 表单验证
        if (!username || !password) {
            throw new Error('请填写用户名和密码');
        }

        const response = await Api.login(username, password);
        
        // 保存token
        if (remember) {
            localStorage.setItem('token', response.token);
        } else {
            sessionStorage.setItem('token', response.token);
        }

        // 保存用户信息
        localStorage.setItem('user', JSON.stringify(response.user));

        showSuccess('登录成功，正在跳转...');
        
        // 延迟跳转
        setTimeout(() => {
            window.location.href = './pages/dashboard.html';
        }, 1500);
    } catch (error) {
        showError(error.message);
    }
}

// 注册处理
async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    try {
        // 表单验证
        if (!username || !email || !password || !confirmPassword) {
            throw new Error('请填写所有必填字段');
        }

        if (password !== confirmPassword) {
            throw new Error('两次输入的密码不一致');
        }

        if (password.length < 6 || password.length > 32) {
            throw new Error('密码长度必须在6-32个字符之间');
        }

        if (username.length < 3 || username.length > 50) {
            throw new Error('用户名长度必须在3-50个字符之间');
        }

        if (!isValidEmail(email)) {
            throw new Error('请输入有效的邮箱地址');
        }

        const response = await Api.register(username, email, password);
        
        showSuccess('注册成功，正在跳转到登录页面...');
        
        // 延迟跳转
        setTimeout(() => {
            window.location.href = '../index.html';
        }, 1500);
    } catch (error) {
        showError(error.message);
    }
}

// 显示错误信息
function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    
    removeExistingMessages();
    
    const form = document.querySelector('.auth-form');
    form.insertBefore(errorDiv, form.firstChild);
    
    setTimeout(() => {
        errorDiv.remove();
    }, 3000);
}

// 显示成功信息
function showSuccess(message) {
    const successDiv = document.createElement('div');
    successDiv.className = 'success-message';
    successDiv.textContent = message;
    
    removeExistingMessages();
    
    const form = document.querySelector('.auth-form');
    form.insertBefore(successDiv, form.firstChild);
    
    setTimeout(() => {
        successDiv.remove();
    }, 3000);
}

// 移除已存在的消息
function removeExistingMessages() {
    const existingMessages = document.querySelectorAll('.error-message, .success-message');
    existingMessages.forEach(msg => msg.remove());
}

// 邮箱验证
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// 检查登录状态
function checkAuth() {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');
    if (!token) {
        window.location.href = '/index.html';
    }
    return token;
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    sessionStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/index.html';
} 