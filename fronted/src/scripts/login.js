document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const messageEl = document.getElementById('message');
    const togglePasswordBtn = document.querySelector('.toggle-password');
    const passwordInput = document.getElementById('password');
    const body = document.body;

    // 密码显示切换
    if (togglePasswordBtn) {
        togglePasswordBtn.addEventListener('click', () => {
            const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
            passwordInput.setAttribute('type', type);
            togglePasswordBtn.querySelector('i').classList.toggle('fa-eye');
            togglePasswordBtn.querySelector('i').classList.toggle('fa-eye-slash');

            if (type === 'text') {
                enableDarkMode();
            } else {
                disableDarkMode();
            }
        });
    }

    // 表单提交处理
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const remember = document.getElementById('remember').checked;
        
        // 基本验证
        if (!username || !password) {
            showMessage('请填写用户名和密码', 'error');
            return;
        }
        
        // 显示加载状态
        const submitBtn = e.target.querySelector('button[type="submit"]');
        setLoadingState(submitBtn, true);
        
        try {
            const response = await fetch('/api/v1/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password, remember })
            });
            
            const data = await response.json();
            
            handleResponse(response.status, data);
        } catch (error) {
            console.error('登录错误:', error);
            showMessage('网络错误，请检查连接', 'error');
        } finally {
            setLoadingState(submitBtn, false);
        }
    });

    // 响应处理函数
    function handleResponse(status, data) {
        switch (status) {
            case 200:
                showMessage('登录成功！正在跳转...', 'success');
                // 存储token和用户信息
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));
                
                // 延迟跳转
                setTimeout(() => {
                    window.location.href = '/dashboard.html';
                }, 1000);
                break;
                
            case 401:
                showMessage('用户名或密码错误', 'error');
                break;
                
            case 400:
                showMessage('请输入有效的用户名和密码', 'error');
                break;
                
            case 500:
                showMessage('服务器错误，请稍后重试', 'error');
                break;
                
            default:
                showMessage(data.message || '登录失败，请重试', 'error');
        }
    }

    // 显示消息提示
    function showMessage(text, type = 'error') {
        messageEl.textContent = text;
        messageEl.className = `message ${type}`;
        
        // 如果是成功消息，自动清除
        if (type === 'success') {
            setTimeout(() => {
                messageEl.textContent = '';
                messageEl.className = 'message';
            }, 3000);
        }
    }

    // 设置按钮加载状态
    function setLoadingState(button, isLoading) {
        if (isLoading) {
            button.disabled = true;
            button.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 登录中...';
        } else {
            button.disabled = false;
            button.innerHTML = '<span>登录</span><i class="fas fa-arrow-right"></i>';
        }
    }

    // 输入时清除错误消息
    document.querySelectorAll('input').forEach(input => {
        input.addEventListener('input', () => {
            messageEl.textContent = '';
            messageEl.className = 'message';
        });
    });

    function enableDarkMode() {
        body.classList.add('dark-mode');
    }

    function disableDarkMode() {
        body.classList.remove('dark-mode');
    }
});